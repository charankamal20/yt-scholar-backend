package token

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/ed25519"
)

type Options struct {
	TokenLifetime time.Duration
	KeyDir        string
	Issuer        string
	Audience      string
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

// Standard PASETO claims following the specification
type Payload struct {
	// Standard claims
	Issuer     string    `json:"iss,omitempty"`
	Subject    string    `json:"sub,omitempty"`
	Audience   string    `json:"aud,omitempty"`
	Expiration time.Time `json:"exp"`
	NotBefore  time.Time `json:"nbf,omitempty"`
	IssuedAt   time.Time `json:"iat"`
	JWTID      string    `json:"jti,omitempty"`

	// Custom claims
	UserID string   `json:"user_id,omitempty"`
	Email  string   `json:"email,omitempty"`
	Roles  []string `json:"roles,omitempty"`
}

func WithTokenLifetime(d time.Duration) Option {
	return func(o *Options) {
		o.TokenLifetime = d
	}
}

func WithKeyDir(dir string) Option {
	return func(o *Options) {
		o.KeyDir = dir
	}
}

func WithIssuer(issuer string) Option {
	return func(o *Options) {
		o.Issuer = issuer
	}
}

func WithAudience(audience string) Option {
	return func(o *Options) {
		o.Audience = audience
	}
}

func keysExist(keyDir string) bool {
	privPath := filepath.Join(keyDir, "private.key")
	pubPath := filepath.Join(keyDir, "public.key")

	_, err1 := os.Stat(privPath)
	_, err2 := os.Stat(pubPath)

	return err1 == nil && err2 == nil
}

func NewPasetoMaker(opts ...Option) (*PasetoMaker, error) {
	cfg := Options{
		TokenLifetime: 15 * time.Minute,
		KeyDir:        "pkg/keys",
		Issuer:        "youtube_scholar_api",
		Audience:      "youtube_scholar",
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	if !keysExist(cfg.KeyDir) {
		if err := createKeys(cfg.KeyDir); err != nil {
			return nil, fmt.Errorf("failed to create keys: %w", err)
		}
	}

	privPath := filepath.Join(cfg.KeyDir, "private.key")
	pubPath := filepath.Join(cfg.KeyDir, "public.key")

	privateKeyBytes, err := os.ReadFile(privPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	publicKeyBytes, err := os.ReadFile(pubPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}

	if len(privateKeyBytes) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid private key size: expected %d, got %d",
			ed25519.PrivateKeySize, len(privateKeyBytes))
	}

	if len(publicKeyBytes) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid public key size: expected %d, got %d",
			ed25519.PublicKeySize, len(publicKeyBytes))
	}

	return &PasetoMaker{
		paseto:     paseto.NewV2(),
		privateKey: ed25519.PrivateKey(privateKeyBytes),
		publicKey:  ed25519.PublicKey(publicKeyBytes),
		opts:       cfg,
	}, nil
}

func generateJTI() string {
	bytes := make([]byte, 16)

	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(bytes)
}

type TokenOptions struct {
	Subject   string
	NotBefore *time.Time
	JWTID     string
	Roles     []string
}

func (maker *PasetoMaker) CreateToken(userID, email string, opts *TokenOptions) (string, error) {
	// Validate inputs
	if userID == "" {
		return "", errors.New("userID cannot be empty")
	}
	if email == "" {
		return "", errors.New("email cannot be empty")
	}

	now := time.Now()
	jti := generateJTI()
	if opts != nil && opts.JWTID != "" {
		jti = opts.JWTID
	}

	payload := &Payload{
		Issuer:     maker.opts.Issuer,
		Subject:    userID,
		Audience:   maker.opts.Audience,
		Expiration: now.Add(maker.opts.TokenLifetime),
		IssuedAt:   now,
		JWTID:      jti,
		UserID:     userID,
		Email:      email,
	}

	// Optional claims
	if opts != nil {
		if opts.Subject != "" {
			payload.Subject = opts.Subject
		}
		if opts.NotBefore != nil {
			payload.NotBefore = *opts.NotBefore
		}
		if opts.Roles != nil {
			payload.Roles = opts.Roles
		}
	}

	token, err := maker.paseto.Sign(maker.privateKey, payload, nil)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return token, nil
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	if token == "" {
		return nil, errors.New("token cannot be empty")
	}

	var payload Payload
	err := maker.paseto.Verify(token, maker.publicKey, &payload, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	now := time.Now()

	if now.After(payload.Expiration) {
		return nil, errors.New("token expired")
	}

	if !payload.NotBefore.IsZero() && now.Before(payload.NotBefore) {
		return nil, errors.New("token not yet valid")
	}

	if payload.Issuer != "" && payload.Issuer != maker.opts.Issuer {
		return nil, errors.New("invalid issuer")
	}

	if payload.Audience != "" && payload.Audience != maker.opts.Audience {
		return nil, errors.New("invalid audience")
	}

	if payload.UserID == "" {
		return nil, errors.New("invalid token: missing user ID")
	}

	return &payload, nil
}

func (maker *PasetoMaker) PublicKey() []byte {
	return maker.publicKey
}
