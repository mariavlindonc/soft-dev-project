package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"os"
	"strconv"
	"time"
)

type jwtClaims struct {
	Sub  uint   `json:"sub"`
	Role string `json:"role"`
	Exp  int64  `json:"exp"`
}

func GenerateJWT(userID uint, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-do-not-use-in-production"
	}

	expiryHours := 24
	if v := os.Getenv("JWT_EXPIRY_HOURS"); v != "" {
		if h, err := strconv.Atoi(v); err == nil && h > 0 {
			expiryHours = h
		}
	}

	claims := jwtClaims{
		Sub:  userID,
		Role: role,
		Exp:  time.Now().Add(time.Duration(expiryHours) * time.Hour).Unix(),
	}

	payload, _ := json.Marshal(claims)
	header := `{"alg":"HS256","typ":"JWT"}`

	headerB64 := base64.RawURLEncoding.EncodeToString([]byte(header))
	payloadB64 := base64.RawURLEncoding.EncodeToString(payload)

	sig := hmacSHA256([]byte(secret), headerB64+"."+payloadB64)
	sigB64 := base64.RawURLEncoding.EncodeToString(sig)

	return headerB64 + "." + payloadB64 + "." + sigB64, nil
}

func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}
