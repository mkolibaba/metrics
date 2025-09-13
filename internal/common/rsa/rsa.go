package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type Encryptor struct {
	publicKey *rsa.PublicKey
}

func NewEncryptor(publicKeyPath string) (*Encryptor, error) {
	op := "new rsa encryptor"

	content, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	block, _ := pem.Decode(content)
	if block == nil {
		return nil, fmt.Errorf("%s: no key found in file", op)
	}

	if block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("%s: unsupported key type %s", op, block.Type)
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Encryptor{publicKey: publicKey}, nil
}

func (e *Encryptor) Encrypt(data []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, e.publicKey, data)
}

type Decryptor struct {
	privateKey *rsa.PrivateKey
}

func NewDecryptor(privateKeyPath string) (*Decryptor, error) {
	op := "new rsa decryptor"

	content, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	block, _ := pem.Decode(content)
	if block == nil {
		return nil, fmt.Errorf("%s: no key found in file", op)
	}

	if block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("%s: unsupported key type %s", op, block.Type)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Decryptor{privateKey: privateKey}, nil
}

func (e *Decryptor) Decrypt(data []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, e.privateKey, data)
}
