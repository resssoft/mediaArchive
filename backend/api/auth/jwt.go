package auth

import (
	"bytes"
	"fmt"
	"github.com/gbrlsnchs/jwt/v3"
	config "github.com/resssoft/mediaArchive/configuration"
	"github.com/valyala/fasthttp"
	"time"
)

var (
	hs           = jwt.NewHS256(config.JwtSecretAccess())
	tokenTTL     = config.JwtAtExpires()
	bearerPrefix = []byte("bearer ")
)

type JwtPayload struct {
	jwt.Payload
	Perms    string `json:"perms,omitempty"`
	UserId   string `json:"userId,omitempty"`
	UserLang string `json:"userLang,omitempty"`
}

func NewAccessToken(perms string, userId, userLang, session string) (time.Time, []byte, error) {
	expired := time.Now().Add(tokenTTL)
	payload := jwt.Payload{
		JWTID:          session,
		ExpirationTime: jwt.NumericDate(expired),
	}
	customPayload := JwtPayload{
		Payload:  payload,
		Perms:    perms,
		UserId:   userId,
		UserLang: userLang,
	}
	token, err := jwt.Sign(customPayload, hs)
	return expired, token, err
}

func VerifyToken(token []byte, validators ...jwt.Validator) (*JwtPayload, error) {
	if len(token) == 0 {
		return nil, fmt.Errorf("token is empty")
	}
	payload := new(JwtPayload)
	validatePayload := jwt.ValidatePayload(&payload.Payload, validators...)
	if _, err := jwt.Verify(token, hs, payload, validatePayload); err != nil {
		return nil, fmt.Errorf("jwt token invalid: %w", err)
	}
	return payload, nil
}

func ExtractBearerToken(ctx *fasthttp.RequestCtx) []byte {
	token := ctx.Request.Header.Peek("Authorization")
	if len(token) <= len(bearerPrefix) {
		return nil
	}
	tokenType, token := token[:len(bearerPrefix)], token[len(bearerPrefix):]
	if bytes.EqualFold(tokenType, bearerPrefix) {
		return token
	}
	return nil
}
