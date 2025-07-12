package token

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ed25519"
)

func createKeys(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create key directory: %w", err)
	}

	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return fmt.Errorf("failed to generate keys: %w", err)
	}

	pvtPath := filepath.Join(path, "private.key")
	if err := os.WriteFile(pvtPath, privateKey, 0600); err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}

	publicPath := filepath.Join(path, "public.key")
	if err := os.WriteFile(publicPath, publicKey, 0644); err != nil {
		return fmt.Errorf("failed to write public key: %w", err)
	}

	return nil
}
