package middleware

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"github.com/Krynegal/url_shortener.git/internal"
	"github.com/google/uuid"
	"net/http"
)

var (
	key = []byte{193, 175, 17, 153, 220, 178, 229, 188, 18, 205, 215, 225, 202,
		239, 181, 31, 53, 150, 51, 111, 44, 36, 103, 199, 135, 185, 180, 234, 145, 255, 53, 93}
	nonce = []byte{188, 53, 153, 211, 53, 29, 174, 45, 48, 153, 251, 227}
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := uuid.NewString()
		var isExistedUser bool
		if ic, err := r.Cookie(string(internal.UserIDSessionKey)); err == nil {
			if dc, err := decrypt(ic.Value); err == nil {
				userID = dc
				isExistedUser = true
			}
		}
		if !isExistedUser {
			ec, err := encrypt(userID)
			if err != nil {
				http.Error(w, "Internal server error", 500)
			}
			oc := &http.Cookie{
				Name:  string(internal.UserIDSessionKey),
				Value: ec,
				Path:  `/`,
			}
			http.SetCookie(w, oc)
		}

		ctx := context.WithValue(r.Context(), internal.UserIDSessionKey, internal.Session{UserID: userID})
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func decrypt(c string) (string, error) {
	b, err := hex.DecodeString(c)
	if err != nil {
		return "", err
	}
	aesGCM, err := makeGSM()
	if err != nil {
		return "", err
	}

	d, err := aesGCM.Open(nil, nonce, b, nil)
	if err != nil {
		return "", nil
	}

	return string(d), nil
}

func encrypt(c string) (string, error) {
	aesGCM, err := makeGSM()
	if err != nil {
		return "", err
	}

	e := aesGCM.Seal(nil, nonce, []byte(c), nil)

	return hex.EncodeToString(e), nil
}

func makeGSM() (cipher.AEAD, error) {
	newCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(newCipher)
	if err != nil {
		return nil, err
	}

	return aesGCM, nil
}
