package tokens

import "time"

type TokenPayload interface {
	Build(duration time.Duration) TokenPayload
	Sign() ([]byte, error)
	Parse(token string) (TokenPayload, error)
	SetID(id string) TokenPayload
	SetUsername(username string) TokenPayload
	SetEmail(email string) TokenPayload
	SetExtraClaims(extraClaims map[string]interface{}) TokenPayload
	SetExtraClaim(key string, value interface{}) TokenPayload
	ToMap() map[string]interface{}
	GetID() string
	GetUsername() string
	GetEmail() string
	GetExtraClaims() map[string]interface{}
	GetExp() int64
}

func (tf *TokenFactory) NewTokenPayload() TokenPayload {
	return &tokenPayload{
		builder:     nil,
		tf:          tf,
		ExtraClaims: make(map[string]interface{}),
	}
}

func (p *tokenPayload) SetID(id string) TokenPayload {
	p.ID = id
	return p
}

func (p *tokenPayload) SetUsername(username string) TokenPayload {
	p.Username = username
	return p
}

func (p *tokenPayload) SetEmail(email string) TokenPayload {
	p.Email = email
	return p
}

func (p *tokenPayload) SetExtraClaims(extraClaims map[string]interface{}) TokenPayload {
	p.ExtraClaims = extraClaims
	return p
}

func (p *tokenPayload) SetExtraClaim(key string, value interface{}) TokenPayload {
	p.ExtraClaims[key] = value
	return p
}

func (p *tokenPayload) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":           p.ID,
		"username":     p.Username,
		"email":        p.Email,
		"extra_claims": p.ExtraClaims,
	}
}

func (p *tokenPayload) GetID() string {
	if p == nil {
		return ""
	}
	return p.ID
}

func (p *tokenPayload) GetUsername() string {
	if p == nil {
		return ""
	}
	return p.Username
}

func (p *tokenPayload) GetEmail() string {
	if p == nil {
		return ""
	}
	return p.Email
}

func (p *tokenPayload) GetExtraClaims() map[string]interface{} {
	if p == nil {
		return make(map[string]interface{})
	}
	return p.ExtraClaims
}

func (p *tokenPayload) GetExp() int64 {
	if p == nil {
		return 0
	}
	return p.Exp
}
