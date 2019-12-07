package service

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

type FileBlobStorageService struct {
	baseDir   string
	container azblob.ContainerURL
}

func (f *FileBlobStorageService) Create(name string) error {
	blobUrl := f.container.NewBlockBlobURL(name)
	ctx := context.Background()
	reader := bytes.NewReader([]byte(""))
	_, err := blobUrl.Upload(ctx, reader, azblob.BlobHTTPHeaders{}, azblob.Metadata{}, azblob.BlobAccessConditions{})
	return err
}

func (f *FileBlobStorageService) Rename(from string, to string) error {
	newBlobUrl := f.container.NewBlockBlobURL(to)
	oldBlobUrl := f.container.NewBlockBlobURL(from)
	ctx := context.Background()
	_, err := newBlobUrl.StartCopyFromURL(ctx, oldBlobUrl.URL(), azblob.Metadata{}, azblob.ModifiedAccessConditions{}, azblob.BlobAccessConditions{})
	if err != nil {
		return err
	}
	for true {
		resp, err := newBlobUrl.GetProperties(ctx, azblob.BlobAccessConditions{})
		if err != nil {
			return err
		}
		if resp.CopyStatus() == azblob.CopyStatusSuccess {
			break
		}
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
	_, err = oldBlobUrl.Delete(ctx, azblob.DeleteSnapshotsOptionInclude, azblob.BlobAccessConditions{})
	if err != nil {
		return err
	}
	return nil
}

func (f *FileBlobStorageService) Read(name string) ([]byte, error) {
	blobUrl := f.container.NewBlockBlobURL(name)
	resp, err := blobUrl.Download(nil, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false)
	if err != nil {
		return nil, err
	}
	readCloser := resp.Body(azblob.RetryReaderOptions{
		MaxRetryRequests: 10,
	})
	defer readCloser.Close()
	return ioutil.ReadAll(readCloser)
}

func (f *FileBlobStorageService) Contain(name string) (bool, error) {
	blobUrl := f.container.NewBlockBlobURL(name)
	ctx := context.Background()
	resp, err := blobUrl.GetProperties(ctx, azblob.BlobAccessConditions{})
	if err != nil {
		return false, err
	}
	return resp.Response().StatusCode != http.StatusOK, nil
}

func (f *FileBlobStorageService) Write(name string, body []byte) error {
	blobUrl := f.container.NewBlockBlobURL(name)
	ctx := context.Background()
	reader := bytes.NewReader(body)
	_, err := blobUrl.Upload(ctx, reader, azblob.BlobHTTPHeaders{}, azblob.Metadata{}, azblob.BlobAccessConditions{})
	return err
}

func NewFileBlobStorageService(baseDir, accountName, accountKey string) *FileBlobStorageService {
	instance := new(FileBlobStorageService)
	instance.baseDir = baseDir + "/"
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatalln(err.Error())
	}
	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	targetUrl, err := url.Parse(
		fmt.Sprintf("http://%s.blob.core.windows.net/youtubedl", accountName),
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
	containerUrl := azblob.NewContainerURL(*targetUrl, pipeline)
	ctx := context.Background()
	if _, err := containerUrl.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone); err != nil {
		log.Fatalln(err.Error())
	}
	instance.container = containerUrl

	return instance
}
