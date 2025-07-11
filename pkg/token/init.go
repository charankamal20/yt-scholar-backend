package token

import (
	"os"
	"path/filepath"

	"golang.org/x/crypto/ed25519"
)

func createKeys(path string) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}

	pvtPath := filepath.Join(path, "private.key")
	err = os.WriteFile(pvtPath, privateKey, 0600)
	if err != nil {
		panic(err)
	}

	publicPath := filepath.Join(path, "public.key")
	err = os.WriteFile(publicPath, publicKey, 0644)
	if err != nil {
		panic(err)
	}
}
