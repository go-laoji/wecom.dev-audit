package storage

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestNewCos(t *testing.T) {
	err = godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	FileStorage, err = NewCos()
	if err != nil {
		t.Error(err)
	}
}

func TestCOSStorage_Put(t *testing.T) {
	f, _ := os.Open(key)
	err = FileStorage.Put(key, f)
	if err != nil {
		t.Error(err)
	}
}

func TestCOSStorage_SignURL(t *testing.T) {
	signUrl, err := FileStorage.SignURL(key, http.MethodGet, 3600)
	if err != nil {
		t.Error(err)
	}
	t.Log(signUrl)
}

func TestCOSStorage_IsExist(t *testing.T) {

}

func TestCOSStorage_Delete(t *testing.T) {

}
