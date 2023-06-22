package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
		Root:              "mar",
	}
	return NewStore(opts)
}

func tearDown(t *testing.T, s *Store) {
	s.Clear()
}

func TestStoreDeleteKey(t *testing.T) {
	s := newStore()
	defer tearDown(t, s)

	key := "mario"
	data := []byte("mp3bytes")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

func TestStore(t *testing.T) {
	s := newStore()
	defer tearDown(t, s)

	key := "mario"
	data := []byte("mp3bytes")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if ok := s.Has(key); !ok {
		t.Errorf("expected to have key %s", key)
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
