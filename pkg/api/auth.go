package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type signinReq struct {
	Password string `json:"password"`
}

type signinResp struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func hashPassword(pass string) string {
	sum := sha256.Sum256([]byte(pass))
	return hex.EncodeToString(sum[:])
}

func jwtSecret(pass string) []byte {
	sum := sha256.Sum256([]byte("secret:" + pass))
	return sum[:]
}

func makeToken(pass string) (string, error) {
	claims := jwt.MapClaims{
		"ph":  hashPassword(pass),              // password hash
		"exp": time.Now().Add(8 * time.Hour).Unix(), // 8 часов
		"iat": time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(jwtSecret(pass))
}

func validateToken(tokenStr string, pass string) bool {
	parsed, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret(pass), nil
	})
	if err != nil || parsed == nil || !parsed.Valid {
		return false
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	ph, _ := claims["ph"].(string)
	return ph == hashPassword(pass)
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		// если пароль не задан, то логин не нужен
		writeJSON(w, http.StatusOK, signinResp{Token: ""})
		return
	}

	var req signinReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, signinResp{Error: err.Error()})
		return
	}

	if req.Password != pass {
		writeJSON(w, http.StatusUnauthorized, signinResp{Error: "Неверный пароль"})
		return
	}

tok, err := makeToken(pass)
if err != nil {
	writeJSON(w, http.StatusInternalServerError, signinResp{
		Error: err.Error(),
	})
	return
}

writeJSON(w, http.StatusOK, signinResp{
	Token: tok,
})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass == "" {
			next(w, r)
			return
		}

		c, err := r.Cookie("token")
		if err != nil || c.Value == "" {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		if !validateToken(c.Value, pass) {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
