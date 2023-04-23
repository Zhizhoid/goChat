package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

type header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type payload struct {
	Exp      time.Time `json:"exp"`
	Username string    `json:"username"`
}

func GenerateJWT(username string, lifetime time.Duration, key string) (string, error) {
	h := header{
		Alg: "HS256",
		Typ: "JWT",
	}

	p := payload{
		Exp:      time.Now().Add(lifetime),
		Username: username,
	}

	hJson, err := json.Marshal(h)
	if err != nil {
		return "", err
	}

	pJson, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	hmacMessage := base64.RawURLEncoding.EncodeToString(hJson) + "." + base64.RawURLEncoding.EncodeToString(pJson)

	s := hmac.New(sha256.New, []byte(key))
	_, err = s.Write([]byte(hmacMessage))
	if err != nil {
		return "", err
	}

	return hmacMessage + "." + base64.RawURLEncoding.EncodeToString(s.Sum(nil)), nil
}

func IsValid(token, key string) (bool, error) {
	lastDotPos := strings.LastIndexByte(token, '.')
	if lastDotPos == -1 {
		return false, nil
	}

	hmacMessage := token[:lastDotPos]
	s := hmac.New(sha256.New, []byte(key))
	_, err := s.Write([]byte(hmacMessage))
	if err != nil {
		return false, err
	}

	validSignature64 := base64.RawURLEncoding.EncodeToString(s.Sum(nil))
	inputSignature64 := token[lastDotPos+1:]

	return inputSignature64 == validSignature64, nil
}

func IsExpired(token string) (bool, error) {
	firstDotPos := strings.IndexByte(token, '.')
	lastDotPos := strings.LastIndexByte(token, '.')

	p64 := token[firstDotPos+1 : lastDotPos]
	pJson, err := base64.RawURLEncoding.DecodeString(p64)
	if err != nil {
		return false, err
	}

	var p payload
	err = json.Unmarshal(pJson, &p)
	if err != nil {
		return false, err
	}

	return time.Now().After(p.Exp), nil
}

func GenerateKey() {

}
