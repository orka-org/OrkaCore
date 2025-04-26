package tokens

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenPayload_Build(t *testing.T) {
	tf := NewTokenFactory("test", "secret")
	p := tf.NewTokenPayload().
		SetID("123").
		SetUsername("test").
		SetEmail("test@test.com").
		SetExtraClaim("extra", "test").
		Build(time.Minute * 5)

	assert.Equal(t, "123", p.GetID())
	assert.Equal(t, "test", p.GetUsername())
	assert.Equal(t, "test@test.com", p.GetEmail())
	assert.Equal(t, "test", p.GetExtraClaims()["extra"])
}

func TestTokenPayload_Sign(t *testing.T) {
	tf := NewTokenFactory("test", "secret")
	p := tf.NewTokenPayload().
		SetID("123").
		SetUsername("test").
		SetEmail("test@test.com").
		SetExtraClaim("extra", "test").
		Build(time.Minute * 5)

	token, err := p.Sign()
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestTokenPayload_Parse(t *testing.T) {
	tf := NewTokenFactory("test", "secret")
	p := tf.NewTokenPayload().
		SetID("123").
		SetUsername("test").
		SetEmail("test@test.com").
		SetExtraClaim("extra", "test").
		Build(time.Minute * 5)

	token, err := p.Sign()
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	tp := tf.NewTokenPayload()
	p2, err := tp.Parse(string(token))

	assert.NoError(t, err)
	assert.Equal(t, "123", p2.GetID())
	assert.Equal(t, "test", p2.GetUsername())
	assert.Equal(t, "test@test.com", p2.GetEmail())
	assert.Equal(t, "test", p2.GetExtraClaims()["extra"])
}
