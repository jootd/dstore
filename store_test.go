package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: DefaultPathTransformFunc,
	}

	s := NewStore(opts)
	data := bytes.NewReader([]byte("bytes"))

	if err := s.writeStream("mypic", data); err != nil {
		t.Error(err)
	}

}

func TestPathTransform(t *testing.T) {

	key := "picnic"
	expected := "8cb7f/0ef89/17f17/21d79/bcb5f/0abe1/4648e/dc6f4"

	pathName := CASPathTransformFunc(key)

	if pathName != expected {
		t.Errorf("have %s  want %s", pathName, expected)
	}

	fmt.Println(pathName)

}
