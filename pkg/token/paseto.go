package token

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/ed25519"
)

type Options struct {
	TokenLifetime time.Duration
	KeyDir        string
}

type Option func(*Options)

var defaultOptions = Options{
	TokenLifetime: 24 * time.Hour,
}

type PasetoMaker struct {
	paseto     *paseto.V2
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
	opts       Options
}

type Payload struct {
	UserID string    `json:"user_id"`
	Email  string    `json:"email"`
	Issued time.Time `json:"issued_at"`
	Expiry time.Time `json:"expiry_at"`
}

func WithTokenLifetime(d time.Duration) Option {
	return func(o *Options) {
		o.TokenLifetime = d
	}
}

// WithKeyDir lets callers specify a custom directory for key files.
func WithKeyDir(dir string) Option {
	return func(o *Options) {
		o.KeyDir = dir
	}
}

func NewPasetoMaker(opts ...Option) *PasetoMaker {
	cfg := Options{
		TokenLifetime: 24 * time.Hour,
		KeyDir:        "pkg/keys",
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	if os.Getenv("ENVIRONMENT") == "PRODUCTION" {
		createKeys(cfg.KeyDir)
	}

	privPath := filepath.Join(cfg.KeyDir, "private.key")
	pubPath := filepath.Join(cfg.KeyDir, "public.key")

	privateKeyBytes, err := os.ReadFile(privPath)
	if err != nil {
		panic("private.key not found: " + err.Error())
	}
	publicKeyBytes, err := os.ReadFile(pubPath)
	if err != nil {
		panic("public.key not found: " + err.Error())
	}

	return &PasetoMaker{
		paseto:     paseto.NewV2(),
		privateKey: ed25519.PrivateKey(privateKeyBytes),
		publicKey:  ed25519.PublicKey(publicKeyBytes),
		opts:       cfg,
	}
}

func (maker *PasetoMaker) CreateToken(userID string, email string) (string, error) {
	payload := &Payload{
		UserID: userID,
		Email:  email,
		Issued: time.Now(),
		Expiry: time.Now().Add(maker.opts.TokenLifetime),
	}
	token, err := maker.paseto.Sign(maker.privateKey, payload, nil)
	return token, err
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	var payload Payload
	err := maker.paseto.Verify(token, maker.publicKey, &payload, nil)
	if err != nil {
		return nil, err
	}
	if time.Now().After(payload.Expiry) {
		return nil, errors.New("token expired")
	}
	return &payload, nil
}
