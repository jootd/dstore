package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestStoreDeleteKey(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
		Root:              "mar",
	}

	s := NewStore(opts)
	key := "mario"
	data := []byte("jpg bytes")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
		Root:              "mar",
	}

	s := NewStore(opts)
	key := "mario"
	data := []byte("mp3 bytes")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := ioutil.ReadAll(r)
	fmt.Println(string(b))
	if string(b) != string(data) {
		t.Errorf("want %s  have %s", data, b)

	}
}

func TestPathTransform(t *testing.T) {

	key := "picnic"
	// expected := "8cb7f/0ef89/17f17/21d79/bcb5f/0abe1/4648e/dc6f4"

	pathName := CASPathTransformFunc(key)

	// if pathName != expected {
	// 	t.Errorf("have %s  want %s", pathName, expected)
	// }

	fmt.Println(pathName)

}
