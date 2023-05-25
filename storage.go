package main

import "io"

type PathTransformFunc func(string) string

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}

type Store struct {
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathName := key

}
