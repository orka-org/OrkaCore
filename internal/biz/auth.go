package biz

import (
	"context"
	"errors"
	"time"

	"github.com/orka-org/orkacore/internal/conf"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/orka-org/orkacore/pkg/tokens"
)

type TokensPair struct {
	Access  string
	Refresh string
}

type authUsecase struct {
	log     *log.Helper
	userUc  UserUsecase
	tf      tokens.TokenFactory
	jwtConf *conf.JWT
}

type AuthUsecase interface {
	Login(ctx context.Context, email, password string) (*TokensPair, error)
	Register(ctx context.Context, email, password string) (*TokensPair, error)
	ValidateToken(ctx context.Context, token string) (*User, error)
	GetUser(ctx context.Context, userID string) (*User, error)
}

func NewAuthUsecase(tf tokens.TokenFactory, userUc UserUsecase, logger log.Logger) AuthUsecase {
	return &authUsecase{
		log:    log.NewHelper(logger),
		userUc: userUc,
		tf:     tf,
	}
}

func (a *authUsecase) Login(ctx context.Context, email, password string) (*TokensPair, error) {
	user, err := a.userUc.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return a.genTokenPair(ctx, user)
}

func (a *authUsecase) Register(ctx context.Context, email, password string) (*TokensPair, error) {
	a.log.Debug("register user")
	a.log.Debug("email", email)
	a.log.Debug("password", password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user, err := a.userUc.CreateUser(ctx, &User{
		Email:    email,
		Password: string(hashedPassword),
	})
	if err != nil {
		return nil, err
	}
	return a.genTokenPair(ctx, user)
}

func (a *authUsecase) ValidateToken(ctx context.Context, token string) (*User, error) {
	res, err := a.tf.NewTokenPayload().Parse(token)
	if err != nil {
		return nil, err
	}

	user, err := a.userUc.GetUserByID(ctx, res.GetID())
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (a *authUsecase) GetUser(ctx context.Context, userID string) (*User, error) {
	user, err := a.userUc.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (a *authUsecase) genToken(ctx context.Context, user *User, expiry int32) (string, error) {
	a.log.Debug("gen tokens; user id = ", user.ID)
	token := a.tf.NewTokenPayload().
		SetID(user.ID).
		SetUsername(user.Username).
		SetEmail(user.Email).
		SetExtraClaims(map[string]interface{}{}).
		Build(time.Second * time.Duration(expiry))
	tokenBytes, err := token.Sign()
	if err != nil {
		return "", err
	}
	return string(tokenBytes), nil
}

func (a *authUsecase) genTokenPair(ctx context.Context, user *User) (*TokensPair, error) {
	accessExp := a.jwtConf.GetExpiry()
	if accessExp == 0 {
		accessExp = 60 * 60
	}
	refreshExp := a.jwtConf.GetRefreshExpiry()
	if refreshExp == 0 {
		refreshExp = 60 * 60 * 24 * 7
	}

	access, err := a.genToken(ctx, user, accessExp)
	if err != nil {
		return nil, err
	}
	refresh, err := a.genToken(ctx, user, refreshExp)
	if err != nil {
		return nil, err
	}

	return &TokensPair{
		Access:  access,
		Refresh: refresh,
	}, nil
}
