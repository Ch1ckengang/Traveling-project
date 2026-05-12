package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"travel-backend/domain"
	"travel-backend/internal/shared"

	"github.com/golang-jwt/jwt/v5"
)

type tokenClaims struct {
	UserID uint   `json:"uid"`
	Email  string `json:"email"`
	Type   string `json:"typ"`
	JTI    string `json:"jti"`
	jwt.RegisteredClaims
}

const (
	defaultJWTIssuer        = "traveling-backend"
	defaultAccessTTLMinutes = 15
	defaultRefreshTTLHours  = 168
)

func issueTokenPair(user *domain.User) (*domain.TokenPair, error) {
	accessTTL := time.Duration(getEnvAsInt("JWT_ACCESS_TTL_MINUTES", defaultAccessTTLMinutes)) * time.Minute
	refreshTTL := time.Duration(getEnvAsInt("JWT_REFRESH_TTL_HOURS", defaultRefreshTTLHours)) * time.Hour

	now := time.Now()
	accessExp := now.Add(accessTTL)
	refreshExp := now.Add(refreshTTL)

	accessJTI, err := generateJTI()
	if err != nil {
		return nil, err
	}

	refreshJTI, err := generateJTI()
	if err != nil {
		return nil, err
	}

	accessToken, err := signToken(user, "access", accessJTI, now, accessExp)
	if err != nil {
		return nil, err
	}

	refreshToken, err := signToken(user, "refresh", refreshJTI, now, refreshExp)
	if err != nil {
		return nil, err
	}

	log.Printf("[TOKEN][ISSUE] uid=%d type=access jti=%s exp=%s", user.ID, accessJTI, accessExp.Format(time.RFC3339))
	log.Printf("[TOKEN][ISSUE] uid=%d type=refresh jti=%s exp=%s", user.ID, refreshJTI, refreshExp.Format(time.RFC3339))

	// FIXED: Use snake_case to match DTO definition in auth_dto.go
	return &domain.TokenPair{
		TokenType:        "Bearer",
		AccessToken:      accessToken,      // JSON will be "access_token" (from struct tag)
		RefreshToken:     refreshToken,     // JSON will be "refresh_token" (from struct tag)
		ExpiresIn:        int64(accessTTL.Seconds()),
		RefreshExpiresIn: int64(refreshTTL.Seconds()),
	}, nil
}

func refreshTokenPair(rawRefreshToken string) (*domain.TokenPair, error) {
	claims, err := parseAndValidateToken(rawRefreshToken, "refresh")
	if err != nil {
		return nil, shared.ErrInvalidRefreshToken
	}

	user := &domain.User{
		ID:    claims.UserID,
		Email: claims.Email,
	}

	pair, err := issueTokenPair(user)
	if err != nil {
		return nil, shared.ErrTokenIssueFailed
	}

	log.Printf("[TOKEN][REFRESH] uid=%d old_jti=%s", claims.UserID, claims.JTI)

	return pair, nil
}

func parseAccessToken(rawAccessToken string) (*tokenClaims, error) {
	claims, err := parseAndValidateToken(rawAccessToken, "access")
	if err != nil {
		return nil, shared.ErrInvalidAccessToken
	}

	return claims, nil
}

func parseAndValidateToken(tokenString, expectedType string) (*tokenClaims, error) {
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return nil, errors.New("empty token")
	}

	claims := &tokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		if expectedType == "refresh" {
			return []byte(getEnvOrDefault("JWT_REFRESH_SECRET", "change-me-refresh-secret")), nil
		}
		return []byte(getEnvOrDefault("JWT_ACCESS_SECRET", "change-me-access-secret")), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.Type != expectedType {
		return nil, errors.New("token type mismatch")
	}

	if claims.UserID == 0 {
		return nil, errors.New("invalid user id claim")
	}

	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func signToken(user *domain.User, tokenType, jti string, issuedAt, expiresAt time.Time) (string, error) {
	issuer := getEnvOrDefault("JWT_ISSUER", defaultJWTIssuer)

	claims := tokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		Type:   tokenType,
		JTI:    jti,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   strconv.FormatUint(uint64(user.ID), 10),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			ID:        jti,
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := getEnvOrDefault("JWT_ACCESS_SECRET", "change-me-access-secret")
	if tokenType == "refresh" {
		secret = getEnvOrDefault("JWT_REFRESH_SECRET", "change-me-refresh-secret")
	}

	return jwtToken.SignedString([]byte(secret))
}

func generateJTI() (string, error) {
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}

	return hex.EncodeToString(raw), nil
}

func getEnvOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}

func getEnvAsInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
}
