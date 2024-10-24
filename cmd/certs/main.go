package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"metrix/pkg/crypto"
	"metrix/pkg/logger"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
)

func createPrivateKey(privateKeyPath string) (*rsa.PrivateKey, error) {
	f, err := os.Create(privateKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create private key")
	}
	defer func() {
		err = f.Close()
	}()

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, errors.Wrap(err, "failed with an error")
	}
	if err := pem.Encode(
		f,
		&pem.Block{
			Type:  crypto.PrivateKeyTitle,
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	); err != nil {
		return nil, errors.Wrap(err, "failed with an error")
	}

	return privateKey, nil
}

func createPublicKey(publicKeyPath string, privateKey *rsa.PrivateKey) (*rsa.PublicKey, error) {
	f, err := os.Create(publicKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed with an error")
	}
	defer func() {
		err = f.Close()
	}()
	err = pem.Encode(f, &pem.Block{
		Type:  crypto.PublicKeyTitle,
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed with an error")
	}
	return &privateKey.PublicKey, nil
}

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	privateKey, err := createPrivateKey("private.key")
	if err != nil {
		logger.Fatal(ctx, "failed to create private key", err)
	}

	_, err = createPublicKey("public.key", privateKey)
	if err != nil {
		logger.Fatal(ctx, "failed to create public key", err)
	}
}
