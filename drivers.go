package main

import (
	"io/ioutil"
)

type FileSaver interface {
	Put(data []byte) error
}

type Img struct {
	b []byte
	S Storage
}

type ImgOperator interface {
	Put(Img) error
	Get(name string) (Img, error)
}

type Video struct {
	b []byte
	S Storage
}

type VideoOperator interface {
	Put(Video) error
	Get(name string) (Video, error)
}

type Storage struct {
	Path string
}

func (s *Storage) Get(name string) ([]byte, error) {
	return ioutil.ReadFile(s.Path + name)
}

func (s *Storage) Put(name string, data []byte) error {
	return ioutil.WriteFile(s.Path+name, data, 0666)
}

func NewStorage(storagePath string) StorageOperator {
	return &Storage{
		Path: storagePath,
	}
}

type StorageOperator interface {
	Get(name string) ([]byte, error)
	Put(name string, data []byte) error
}
