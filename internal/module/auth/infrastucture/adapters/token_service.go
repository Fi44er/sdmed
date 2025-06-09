package adapters

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/Fi44er/sdmed/internal/module/auth/entity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService struct{}

func NewTokenService() *TokenService {
	return &TokenService{}
}

func (s *TokenService) CreateToken(userID string, ttl time.Duration, privateKey string) (*entity.TokenDetails, error) {
	now := time.Now().UTC()
	td := &entity.TokenDetails{
		ExpiresIn: new(int64),
		Token:     new(string),
	}
	*td.ExpiresIn = now.Add(ttl).Unix()
	td.TokenUUID = uuid.NewString()
	td.UserID = userID

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode private key: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	*td.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":        userID,
		"token_uuid": td.TokenUUID,
		"exp":        td.ExpiresIn,
		"iat":        now.Unix(),
		"nbf":        now.Unix(),
	}).SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("sign token: %w", err)
	}

	return td, nil
}

func (s *TokenService) ValidateToken(token string, publicKey string) (*entity.TokenDetails, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode public key: %w", err)
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return &entity.TokenDetails{
		TokenUUID: fmt.Sprint(claims["token_uuid"]),
		UserID:    fmt.Sprint(claims["sub"]),
	}, nil
}
