package service

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/rylio/ytdl"
)

type IFileStorage interface {
	Create(name string) error
	Write(name string, body []byte) error
	Read(name string) ([]byte, error)
	Rename(from string, to string) error
	Contain(name string) (bool, error)
}

type DownloadJob struct {
	Link string
	Mode int
	URL  string
}

type DownloaderService struct {
	workQueue  chan *DownloadJob
	storage    IFileStorage
	bufferPool *sync.Pool
}

func (d *DownloaderService) GetVideoQuality(link string) ([]string, error) {
	videoInfo, err := ytdl.GetVideoInfo(link)
	if err != nil {
		return nil, err
	}
	qualities := make([]string, 0, len(videoInfo.Formats))
	for _, v := range videoInfo.Formats {
		combined := fmt.Sprintf("%s %s %s/%s", v.Resolution, v.Extension, v.VideoEncoding, v.AudioEncoding)
		qualities = append(qualities, combined)
	}
	return qualities, nil
}

func (d *DownloaderService) Request(link string, mode int) (string, error) {
	videoInfo, err := ytdl.GetVideoInfo(link)
	if err != nil {
		log.Println(err.Error())
		return "", errors.New("Failed to get video info")
	}
	if len(videoInfo.Formats) <= mode {
		return "", errors.New("Mode not supported, please use /info API to find supported mode")
	}
	hashSource := fmt.Sprintf("[%d]%s", mode, videoInfo.Title)
	hashResult := sha512.Sum512([]byte(hashSource))
	encodedName := base64.URLEncoding.EncodeToString(hashResult[:])
	job := new(DownloadJob)
	job.Link = link
	job.Mode = mode
	job.URL = fmt.Sprintf("%s.%s", encodedName, videoInfo.Formats[mode].Extension)
	d.storage.Create("-" + job.URL)
	d.workQueue <- job
	return job.URL, nil
}

func (d *DownloaderService) Download(link string) ([]byte, error) {
	if ok, err := d.storage.Contain(link); ok {
		data, err := d.storage.Read(link)
		if err != nil {
			log.Println(err.Error())
			return nil, errors.New("Failed to request download")
		}
		return data, nil
	} else if err != nil {
		log.Println(err.Error())
		return nil, errors.New("Failed to request download")
	}
	if ok, err := d.storage.Contain("-" + link); ok {
		return nil, errors.New("File is still being processed")
	} else if err != nil {
		log.Println(err.Error())
		return nil, errors.New("Failed to request download")
	}
	return nil, errors.New("This file is not requested yet")
}

func (d *DownloaderService) downloaderRoutine() {
	for true {
		job := <-d.workQueue
		videoInfo, err := ytdl.GetVideoInfo(job.Link)
		if err != nil {
			log.Printf("Error when getting video info, %s", err.Error())
			continue
		}
		buffer := d.bufferPool.Get().(*bytes.Buffer)
		defer buffer.Reset()
		defer d.bufferPool.Put(buffer)
		if err := videoInfo.Download(videoInfo.Formats[job.Mode], buffer); err != nil {
			log.Printf("Error when downloading video, %s", err.Error())
			continue
		}
		d.storage.Write("-"+job.URL, buffer.Bytes())
		d.storage.Rename("-"+job.URL, job.URL)
	}
}

func NewDownloaderService(storage IFileStorage) *DownloaderService {
	instance := new(DownloaderService)
	instance.storage = storage
	instance.workQueue = make(chan *DownloadJob, 100)
	instance.bufferPool = new(sync.Pool)
	instance.bufferPool.New = func() interface{} {
		return bytes.NewBuffer(nil)
	}
	go instance.downloaderRoutine()
	return instance
}
