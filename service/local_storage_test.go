package main

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestLocalStorage_Get(t *testing.T) {
	testString := "test string"
	ls := NewLocalStorage("/tmp")
	input := strings.NewReader(testString)
	_, err := ls.Put("","test.txt", input)

	if err != nil {
		t.Fatalf("failed to Put file: %v", err.Error())
	}

	r, err := ls.Get("","test.txt")
	if err != nil {
		t.Fatalf("failed to Put file: %v", err.Error())
	}

	bytes, err := ioutil.ReadAll(r)
	if testString != string(bytes) {
		t.Fatalf("strings do not match:\nExpected: <%v>;\nActual: <%v>\n", testString, string(bytes))
	}

	err = r.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = ls.Delete("","test.txt")
	if err != nil {
		t.Fatal(err)
	}
}

func TestLocalStorage_Put(t *testing.T) {
	testString := "test string"
	ls := NewLocalStorage("/tmp")
	input := strings.NewReader(testString)
	f, err := ls.Put("", "test.txt", input)

	if err != nil {
		t.Fatalf("failed to Put file: %v", err.Error())
	}

	if !FileExists("/tmp/" + f, true) {
		t.Fatalf("file doesn't exist after Put")
	}

	bytes, err := ioutil.ReadFile("/tmp/" + f)

	if testString != string(bytes) {
		t.Fatalf("strings do not match:\nExpected: <%v>;\nActual: <%v>\n", testString, string(bytes))
	}

	err = ls.Delete("", "test.txt")
	if err != nil {
		t.Fatal(err)
	}
}

func TestLocalStorage_Delete(t *testing.T) {
	ls := NewLocalStorage("/tmp")
	f, err := ls.Put("", "test2.txt", strings.NewReader("some content"))

	if err != nil {
		t.Fatalf("failed to Put file: %v", err.Error())
	}

	err = ls.Delete("","test2.txt")
	if err != nil {
		t.Fatal(err)
	}

	if FileExists("/tmp/" + f, true) {
		t.Fatalf("file still exists after delete")
	}
}