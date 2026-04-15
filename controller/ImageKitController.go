package controller

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type imagekitAuthResponse struct {
	Token     string `json:"token"`
	Expire    int64  `json:"expire"`
	Signature string `json:"signature"`
	PublicKey string `json:"publicKey"`
}

// generateToken produces a cryptographically random 32-char hex token.
func generateToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// GetImageKitAuth returns a short-lived signature bundle for client-side uploads.
// The private key never leaves the server.
// GET /api/imagekit/auth/
func GetImageKitAuth(w http.ResponseWriter, r *http.Request) {
	privateKey := os.Getenv("IMAGEKIT_PRIVATE_KEY")
	publicKey := os.Getenv("IMAGEKIT_PUBLIC_KEY")
	if privateKey == "" || publicKey == "" {
		http.Error(w, "ImageKit not configured", http.StatusServiceUnavailable)
		return
	}

	token, err := generateToken()
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Expire 40 minutes from now (ImageKit requires < 1 hour in the future).
	expire := time.Now().Unix() + 2400

	// Signature = HMAC-SHA1(token + strconv.Itoa(expire), privateKey)
	mac := hmac.New(sha1.New, []byte(privateKey))
	mac.Write([]byte(fmt.Sprintf("%s%d", token, expire)))
	signature := hex.EncodeToString(mac.Sum(nil))

	w.Header().Set("Content-Type", "application/json")
	// Allow the frontend origin to read this response.
	w.Header().Set("Cache-Control", "no-store")
	json.NewEncoder(w).Encode(imagekitAuthResponse{
		Token:     token,
		Expire:    expire,
		Signature: signature,
		PublicKey: publicKey,
	})
}
