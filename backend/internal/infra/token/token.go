package token

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	app "github.com/umekikazuya/me/internal/app/identity"
	domain "github.com/umekikazuya/me/internal/domain/identity"
	"github.com/umekikazuya/me/pkg/errs"
)

var _ app.TokenService = (*JWTTokenService)(nil)

type JWTTokenService struct {
	secret   []byte
	atExpiry time.Duration
}

func NewJWTTokenService(secret string, atExpiry time.Duration) *JWTTokenService {
	return &JWTTokenService{
		secret:   []byte(secret),
		atExpiry: atExpiry,
	}
}

func (s *JWTTokenService) GenerateAT(ctx context.Context, identity domain.Identity) (string, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(time.Minute * 15)

	claims := jwt.RegisteredClaims{
		Subject:   identity.ID(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.secret)
	if err != nil {
		return "", app.ErrTokenInvalid
	}

	return signedToken, nil
}

func (s *JWTTokenService) GenerateRT(ctx context.Context) (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf(
			"GenerateRT: 乱数生成処理が失敗しました %w",
			errs.ErrInternal,
		)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (s *JWTTokenService) Hash(
	ctx context.Context,
	input string,
) (string, error) {
	if input == "" {
		return "", fmt.Errorf(
			"Hash: 入力文字が空です %w", errs.ErrInternal,
		)
	}
	hashed := sha256.Sum256(
		[]byte(input),
	)

	return hex.EncodeToString(hashed[:]), nil
}

func (s *JWTTokenService) Validate(ctx context.Context, token string) error {
	// トークンの長さ制限を追加
	if len(token) > 4096 {
		return app.ErrTokenInvalid
	}
	validatedToken, err := jwt.ParseWithClaims(
		token, &jwt.RegisteredClaims{},
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, app.ErrTokenInvalid
			}
			return s.secret, nil
		},
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return app.ErrTokenExpired
		}
		return app.ErrTokenInvalid
	}

	if !validatedToken.Valid {
		return app.ErrTokenInvalid
	}

	claims, ok := validatedToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return app.ErrTokenInvalid
	}

	_, err = uuid.Parse(claims.Subject)
	if err != nil {
		return app.ErrTokenInvalid
	}

	return nil
}
