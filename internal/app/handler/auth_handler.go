package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/olkonon/shortener/internal/app/common"
	"net/http"
	"time"
)

func (h *Handler) WithAuth(handle http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		//Извлечение юзера из куки если есть
		user := h.readSessionToken(r)

		if mux.Vars(r) == nil {
			r = mux.SetURLVars(r, map[string]string{common.MuxUserVarName: user})
		} else {
			mux.Vars(r)[common.MuxUserVarName] = user
		}
		handle.ServeHTTP(w, r)
	}
	return http.HandlerFunc(logFn)
}

func (h *Handler) AnonymousAuthHandler(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		//Извлечение юзера из куки если есть
		user := mux.Vars(r)[common.MuxUserVarName]
		if user == common.AnonymousUser {
			user = uuid.New().String()
			mux.Vars(r)[common.MuxUserVarName] = user
		}
		h.writeSessionToken(w, user)
		f(w, r)
	}
	return http.HandlerFunc(logFn)
}

func (h *Handler) RequireAuthHandler(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		//Извлечение юзера из куки если есть
		user := mux.Vars(r)[common.MuxUserVarName]
		if user == common.AnonymousUser {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		h.writeSessionToken(w, user)
		f(w, r)
	}
	return http.HandlerFunc(logFn)
}

func (h *Handler) readSessionToken(r *http.Request) string {
	cookie, err := r.Cookie(common.SessionCookieName)
	if err != nil {
		return common.AnonymousUser
	}
	decodedValue, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return common.AnonymousUser
	}
	signedValue := string(decodedValue)
	//Check size contain signature
	if len(signedValue) < sha256.Size {
		return common.AnonymousUser
	}

	signature := signedValue[:sha256.Size]
	value := signedValue[sha256.Size:]

	mac := hmac.New(sha256.New, h.secretKey)
	mac.Write([]byte(value))
	expectedSignature := mac.Sum(nil)
	if hmac.Equal([]byte(signature), expectedSignature) {
		return value
	}
	return common.AnonymousUser
}

func (h *Handler) writeSessionToken(w http.ResponseWriter, user string) {
	mac := hmac.New(sha256.New, h.secretKey)
	mac.Write([]byte(user))
	signature := mac.Sum(nil)

	cookie := new(http.Cookie)
	cookie.Name = common.SessionCookieName
	cookie.Value = base64.StdEncoding.EncodeToString([]byte(string(signature) + user))

	cookie.Expires = time.Now().Add(time.Hour)
	http.SetCookie(w, cookie)
}

func (h *Handler) MockTestUserCookie() *http.Cookie {
	mac := hmac.New(sha256.New, h.secretKey)
	mac.Write([]byte(common.TestUser))
	signature := mac.Sum(nil)

	cookie := new(http.Cookie)
	cookie.Name = common.SessionCookieName
	cookie.Value = base64.StdEncoding.EncodeToString([]byte(string(signature) + common.TestUser))

	cookie.Expires = time.Now().Add(time.Hour)
	return cookie
}
