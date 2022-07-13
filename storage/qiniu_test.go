package storage

import (
	"os"
	"testing"
)

var err error
var key = "qiniu_test.go"

func TestNewQiNiu(t *testing.T) {
	FileStorage, err = NewQiNiu("N_DVuxqdvEMpqm2uVUyc48tOAbis1ZJvXGDLSsvI",
		"GDXTLnyv_aaJVwHvhz_jkYoNcVK5lmWd1T30CX4p",
		"ciprun-test")
	if err != nil {
		t.Error(err)
	}
}

func TestQiNiuOSSStorage_Put(t *testing.T) {
	f, _ := os.Open(key)
	err = FileStorage.Put(key, f)
	if err != nil {
		t.Error(err)
	}
}

func TestQiNiuOSSStorage_SignURL(t *testing.T) {
	privateUrl, err := FileStorage.SignURL(key, "", 3600)
	if err != nil {
		t.Error(err)
	}
	t.Log(privateUrl)
}

func TestQiNiuOSSStorage_IsExist(t *testing.T) {
	ok, err := FileStorage.IsExist(key)
	if err != nil {
		t.Error(err)
	}
	t.Log(ok)
}

func TestQiNiuOSSStorage_Delete(t *testing.T) {
	files, err := FileStorage.Delete(key)
	if err != nil {
		t.Error(err)
	}
	t.Log(files)
}
