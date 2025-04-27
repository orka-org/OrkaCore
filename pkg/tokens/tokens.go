package tokens

import (
	"errors"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/orka-org/orkacore/internal/conf"
)

type TokenFactory struct {
	name   string
	secret string
	alg    jwa.KeyAlgorithm
	key    []byte
}

func NewTokenProvider(conf *conf.Bootstrap) TokenFactory {
	name := conf.GetService().GetName()
	if name == "" {
		name = "orka"
	}
	secret := conf.GetJwt().GetSecret()
	if secret == "" {
		secret = "orka"
	}
	return TokenFactory{
		name:   name,
		secret: secret,
		alg:    jwa.HS256(),
		key:    []byte(secret),
	}
}

func NewTokenFactory(name, secret string) *TokenFactory {
	return &TokenFactory{
		name:   name,
		secret: secret,
		alg:    jwa.HS256(),
		key:    []byte(secret),
	}
}

type tokenPayload struct {
	ID          string                 `json:"id"`
	Username    string                 `json:"username"`
	Email       string                 `json:"email"`
	ExtraClaims map[string]interface{} `json:"extra_claims"`
	Exp         int64                  `json:"exp"`

	builder *jwt.Builder
	tf      *TokenFactory
}

func (p tokenPayload) Build(duration time.Duration) TokenPayload {
	builder := jwt.NewBuilder().IssuedAt(time.Now()).
		NotBefore(time.Now()).
		Issuer(p.tf.name).
		Expiration(time.Now().Add(duration)).
		Subject(p.ID).
		Claim("username", p.Username).
		Claim("email", p.Email).
		Claim("extraClaims", p.ExtraClaims)

	p.builder = builder
	return &p
}

func (p tokenPayload) Sign() ([]byte, error) {
	if p.builder == nil {
		return nil, errors.New("builder is nil")
	}
	if p.tf.secret == "" {
		return nil, errors.New("secret is empty")
	}
	token, err := p.builder.Build()
	if err != nil {
		return nil, err
	}
	opts := []jwt.SignOption{jwt.WithKey(p.tf.alg, p.tf.key)}
	tok, err := jwt.Sign(token, opts...)
	if err != nil {
		return nil, err
	}
	return tok, nil
}

func (p *tokenPayload) Parse(token string) (TokenPayload, error) {
	if p == nil {
		return nil, errors.New("token payload is nil")
	}
	if token == "" {
		return nil, errors.New("token is empty")
	}
	tok := []byte(token)
	opts := []jwt.ParseOption{jwt.WithValidate(true), jwt.WithKey(p.tf.alg, p.tf.key)}

	payload, err := jwt.Parse(tok, opts...)
	if err != nil {
		return nil, err
	}
	id, ok := payload.Subject()
	if !ok {
		return nil, errors.New("subject is not string")
	}

	p.ID = id
	var username string
	err = payload.Get("username", &username)
	if err != nil {
		return nil, err
	}
	p.Username = username

	var email string
	err = payload.Get("email", &email)
	if err != nil {
		return nil, err
	}
	p.Email = email

	var extraClaims map[string]interface{}
	err = payload.Get("extraClaims", &extraClaims)
	if err != nil {
		return nil, err
	}
	p.ExtraClaims = extraClaims

	exp, _ := payload.Expiration()
	return &tokenPayload{
		ID:          id,
		Username:    username,
		Email:       email,
		ExtraClaims: extraClaims,
		Exp:         exp.Unix(),
		builder:     nil,
		tf:          p.tf,
	}, nil
}
