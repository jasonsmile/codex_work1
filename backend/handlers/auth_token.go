package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

const defaultTokenSecret = "medical-info-platform-secret"

type tokenClaims struct {
	UUID        string `json:"uuid"`
	Username    string `json:"username"`
	AuthorityID uint64 `json:"authorityId"`
	ExpiresAt   int64  `json:"expiresAt"`
}

func createToken(userUUID, username string, authorityID uint64) (string, error) {
	claims := tokenClaims{
		UUID:        userUUID,
		Username:    username,
		AuthorityID: authorityID,
		ExpiresAt:   time.Now().Add(24 * time.Hour).Unix(),
	}

	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	signature := signPayload(encodedPayload)
	return encodedPayload + "." + signature, nil
}

func parseToken(token string) (tokenClaims, error) {
	var claims tokenClaims
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return claims, errors.New("token格式不正确")
	}

	expectedSignature := signPayload(parts[0])
	if !hmac.Equal([]byte(expectedSignature), []byte(parts[1])) {
		return claims, errors.New("token签名无效")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return claims, err
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return claims, err
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return claims, errors.New("token已过期")
	}
	if claims.UUID == "" || claims.Username == "" || claims.AuthorityID == 0 {
		return claims, errors.New("token内容不完整")
	}

	return claims, nil
}

func signPayload(payload string) string {
	mac := hmac.New(sha256.New, []byte(getTokenSecret()))
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func getTokenSecret() string {
	secret := os.Getenv("AUTH_SECRET")
	if secret == "" {
		return defaultTokenSecret
	}
	return secret
}

func roleName(authorityID uint64) string {
	return fmt.Sprintf("role_%d", authorityID)
}
