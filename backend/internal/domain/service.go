package domain

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailTaken         = errors.New("email already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Service struct {
	repo       UserRepository
	privateKey *rsa.PrivateKey
}

func New(repo UserRepository, privateKeyPEM string) (*Service, error) {
	// Support \n-escaped PEM strings stored in env vars
	pemStr := strings.ReplaceAll(privateKeyPEM, `\n`, "\n")

	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block from JWT_PRIVATE_KEY")
	}

	var rsaKey *rsa.PrivateKey
	switch block.Type {
	case "RSA PRIVATE KEY":
		k, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parsing PKCS1 RSA private key: %w", err)
		}
		rsaKey = k
	case "PRIVATE KEY":
		k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parsing PKCS8 private key: %w", err)
		}
		var ok bool
		rsaKey, ok = k.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("JWT_PRIVATE_KEY must be an RSA key")
		}
	default:
		return nil, fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}

	return &Service{repo: repo, privateKey: rsaKey}, nil
}

func (s *Service) Register(ctx context.Context, email, name, password string) (*User, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("hashing password: %w", err)
	}

	user, err := s.repo.Create(ctx, email, name, string(hash))
	if err != nil {
		if isUniqueViolation(err) {
			return nil, "", ErrEmailTaken
		}
		return nil, "", fmt.Errorf("creating user: %w", err)
	}

	token, err := s.issueToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*User, string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	token, err := s.issueToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *Service) PublicKey() *rsa.PublicKey {
	return &s.privateKey.PublicKey
}

type claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func (s *Service) issueToken(user *User) (string, error) {
	now := time.Now()
	c := claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			Issuer:    "issue-tracker",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
	signed, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return signed, nil
}
