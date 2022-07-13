package internal

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"log"
)

func readPriKey(privateKey string) (*rsa.PrivateKey, error) {
	p, _ := pem.Decode([]byte(privateKey))
	key, err := x509.ParsePKCS1PrivateKey(p.Bytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func RsaDecrypt(key string, cipherText string) (string, error) {
	rsaKey, err := readPriKey(key)
	if err != nil {
		return "", err
	}
	p, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		log.Println("decode base string err", err)
	}
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, rsaKey, p)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
