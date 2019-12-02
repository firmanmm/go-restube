package service

import (
	"io/ioutil"
	"os"
)

type FileBasedStorageService struct {
	baseDir string
}

func (f *FileBasedStorageService) Create(name string) error {
	file, err := os.Create(f.baseDir + name)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

func (f *FileBasedStorageService) Rename(from string, to string) error {
	oldName := f.baseDir + from
	newName := f.baseDir + to
	return os.Rename(oldName, newName)
}

func (f *FileBasedStorageService) Read(name string) ([]byte, error) {
	name = f.baseDir + name
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func (f *FileBasedStorageService) Contain(name string) (bool, error) {
	name = f.baseDir + name
	_, err := os.Stat(name)
	return err != nil, err
}

func (f *FileBasedStorageService) Write(name string, body []byte) error {
	name = f.baseDir + name
	return ioutil.WriteFile(name, body, 0644)
}

func NewFileStorageService(baseDir string) *FileBasedStorageService {
	instance := new(FileBasedStorageService)
	instance.baseDir = baseDir + "/"
	os.MkdirAll(baseDir, 0644)
	return instance
}
