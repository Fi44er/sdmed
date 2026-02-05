package auth_adapters

// import (
// 	"encoding/base64"
// 	"fmt"
// 	"time"

// 	auth_entity "github.com/Fi44er/sdmed/internal/module/auth/entity"
// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/google/uuid"
// )

// type TokenService struct{}

// func NewTokenService() *TokenService {
// 	return &TokenService{}
// }

// func (s *TokenService) CreateToken(userID, deviceID string, ttl time.Duration, privateKey string) (*auth_entity.TokenDetails, error) {
// 	now := time.Now().UTC()
// 	td := &auth_entity.TokenDetails{
// 		ExpiresIn: new(int64),
// 		Token:     new(string),
// 		UserID:    userID,
// 		DeviceID:  deviceID,
// 	}
// 	*td.ExpiresIn = now.Add(ttl).Unix()
// 	td.TokenUUID = uuid.NewString()
// 	td.UserID = userID

// 	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
// 	if err != nil {
// 		return nil, fmt.Errorf("could not decode private key: %w", err)
// 	}
// 	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
// 	if err != nil {
// 		return nil, fmt.Errorf("parse private key: %w", err)
// 	}

// 	*td.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
// 		"sub":        userID,
// 		"token_uuid": td.TokenUUID,
// 		"device_id":  deviceID,
// 		"user_id":    userID,
// 		"exp":        td.ExpiresIn,
// 		"iat":        now.Unix(),
// 		"nbf":        now.Unix(),
// 	}).SignedString(key)
// 	if err != nil {
// 		return nil, fmt.Errorf("sign token: %w", err)
// 	}

// 	return td, nil
// }

// func (s *TokenService) ValidateToken(token string, publicKey string) (*auth_entity.TokenDetails, error) {
// 	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
// 	if err != nil {
// 		return nil, fmt.Errorf("could not decode public key: %w", err)
// 	}
// 	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
// 	if err != nil {
// 		return nil, fmt.Errorf("parse public key: %w", err)
// 	}

// 	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
// 		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %s", t.Header["alg"])
// 		}
// 		return key, nil
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("parse token: %w", err)
// 	}

// 	claims, ok := parsedToken.Claims.(jwt.MapClaims)
// 	if !ok || !parsedToken.Valid {
// 		return nil, fmt.Errorf("invalid token")
// 	}

// 	userID, ok := claims["sub"].(string)
// 	if !ok {
// 		return nil, fmt.Errorf("missing or invalid 'sub' claim")
// 	}

// 	// NEW: извлекаем device_id из claims
// 	deviceID, ok := claims["device_id"].(string)
// 	if !ok {
// 		return nil, fmt.Errorf("missing or invalid 'device_id' claim")
// 	}

// 	exp, ok := claims["exp"].(float64)
// 	if !ok {
// 		return nil, fmt.Errorf("missing or invalid 'exp' claim")
// 	}

// 	expInt := int64(exp)

// 	return &auth_entity.TokenDetails{
// 		Token:     &token,
// 		UserID:    userID,
// 		DeviceID:  deviceID, // NEW
// 		ExpiresIn: &expInt,
// 	}, nil
// }
