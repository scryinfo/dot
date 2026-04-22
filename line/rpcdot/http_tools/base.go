package httptools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/lib/kits"
)

const (
	TokenName         = "Token" // the name "Token" be used many places, so use the name "TokenName"
	TokenGame         = "TokenGame"
	Authorization     = "Authorization"
	AuthorizationGame = "AuthorizationGame"
)

type ReqBase struct {
	Ts int64 `json:"ts"`
}

type ResBase struct {
	ErrorIdd int    `json:"errorId"`
	Message  string `json:"message"`
	Ts       int64  `json:"ts"`
}

func ApiWithBaseUrl(preUrl, url string) string {
	if strings.HasPrefix(url, "/") {
		return preUrl + url
	} else {
		return preUrl + "/" + url
	}
}

func HttpErrorResponse(w http.ResponseWriter, message string, code int) {
	type ResUpload struct {
		Base ResBase `json:"base"`
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(&ResUpload{
		Base: ResBase{
			Message:  message,
			ErrorIdd: code,
			Ts:       kits.Times.SocondTs(),
		},
	})
	if err != nil {
		dot.Logger.Error().AnErr("cant encode error response: ", err).Send()
	}
}

func HttpOkResponse[T any](w http.ResponseWriter, res T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(&res)
	if err != nil {
		dot.Logger.Error().AnErr("cant encode error response: ", err).Send()
	}
}

func GetToken(r *http.Request) string {
	token := ""
	tokenCookie, err := r.Cookie(TokenName)
	if err != nil || len(tokenCookie.Value) < 1 {
		token = r.Header.Get(TokenName)
		if len(token) < 1 {
			token = r.Header.Get(Authorization)
		}
	} else {
		token = tokenCookie.Value
	}
	return token
}

func GetTokenByHeader(header http.Header) string {
	token := ""
	cookie := header.Get("Cookie")
	if len(cookie) > 1 {
		r := http.Request{
			Header: http.Header{
				"Cookie": []string{cookie},
			},
		}
		tokenCookie, err := r.Cookie(TokenName)
		if err == nil {
			token = tokenCookie.Value
		}
	}

	if len(token) < 1 {
		token = header.Get(TokenName)
		if len(token) < 1 {
			token = header.Get(Authorization)
		}
	}
	return token
}

func GetTokenGame(r *http.Request) string {
	token := ""
	tokenCookie, err := r.Cookie(TokenGame)
	if err != nil || len(tokenCookie.Value) < 1 {
		token = r.Header.Get(TokenGame)
		if len(token) < 1 {
			token = r.Header.Get(AuthorizationGame)
		}
	} else {
		token = tokenCookie.Value
	}
	return token
}

func GetTokenGameByHeader(header http.Header) string {
	token := ""
	cookie := header.Get("Cookie")
	if len(cookie) > 1 {
		r := http.Request{
			Header: http.Header{
				"Cookie": []string{cookie},
			},
		}
		tokenCookie, err := r.Cookie(TokenGame)
		if err == nil {
			token = tokenCookie.Value
		}
	}

	if len(token) < 1 {
		token = header.Get(TokenGame)
		if len(token) < 1 {
			token = header.Get(AuthorizationGame)
		}
	}
	return token
}

type userClaims struct {
	UserName string `json:"userName"`
	jwt.RegisteredClaims
}

var _kk_j = []byte("_dot_super_$#@082347")

func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &userClaims{
		UserName: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(_kk_j)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func ValidToken(tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(token *jwt.Token) (any, error) {
		return _kk_j, nil
	})
	if err != nil {
		dot.Logger.Info().Err(err).Send()
		return err
	}
	if _, ok := token.Claims.(*userClaims); !ok || !token.Valid {
		return fmt.Errorf("Unauthorized")
	}
	return nil
}
