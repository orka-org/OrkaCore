package biz

import (
	"context"
	"time"

	kerr "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type User struct {
	ID        string
	Username  string
	FirstName *string
	LastName  *string
	Phone     *string
	Email     string
	Password  string

	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	SuspendedAt *time.Time
}

type UserRepo interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetAllUsers(ctx context.Context) ([]*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	SuspendUser(ctx context.Context, user *User) (*User, error)
	DeleteUser(ctx context.Context, user *User) (*User, error)
}

type UserUsecase interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetAllUsers(ctx context.Context) ([]*User, error)

	UpdateUserPassword(ctx context.Context, password string, email string) (*User, error)
	UpdateUserEmail(ctx context.Context, id string, email string) (*User, error)
	UpdateUserProfile(ctx context.Context, user *User) (*User, error)
	SuspendUser(ctx context.Context, user *User) (*User, error)
}

type userUsecase struct {
	log      *log.Helper
	userRepo UserRepo
}

func NewUserUsecase(userRepo UserRepo, logger log.Logger) UserUsecase {
	return &userUsecase{
		log:      log.NewHelper(logger),
		userRepo: userRepo,
	}
}

func (u *userUsecase) CreateUser(ctx context.Context, user *User) (*User, error) {
	return u.userRepo.CreateUser(ctx, user)
}

func (u *userUsecase) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return u.userRepo.GetUserByEmail(ctx, email)
}

func (u *userUsecase) GetUserByID(ctx context.Context, id string) (*User, error) {
	return u.userRepo.GetUserByID(ctx, id)
}

func (u *userUsecase) GetAllUsers(ctx context.Context) ([]*User, error) {
	return u.userRepo.GetAllUsers(ctx)
}

func (u *userUsecase) UpdateUserPassword(ctx context.Context, id string, password string) (*User, error) {
	payload := &User{
		ID:       id,
		Password: password,
	}
	return u.userRepo.UpdateUser(ctx, payload)
}

func (u *userUsecase) UpdateUserEmail(ctx context.Context, id string, email string) (*User, error) {
	payload := &User{
		ID:    id,
		Email: email,
	}
	return u.userRepo.UpdateUser(ctx, payload)
}

func (u *userUsecase) UpdateUserProfile(ctx context.Context, user *User) (*User, error) {
	if user.Password != "" {
		return nil, kerr.BadRequest("Password cannot be upudated", "If you wish to update the password, please follow the right procedure")
	}
	if user.Email != "" {
		return nil, kerr.BadRequest("Email cannot be upudated", "If you wish to update the email, please follow the right procedure")
	}
	payload := &User{}
	if user.Username != "" {
		payload.Username = user.Username
	}
	if user.FirstName != nil {
		payload.FirstName = user.FirstName
	}
	if user.LastName != nil {
		payload.LastName = user.LastName
	}
	if user.Phone != nil {
		payload.Phone = user.Phone
	}
	return u.userRepo.UpdateUser(ctx, payload)
}

func (u *userUsecase) SuspendUser(ctx context.Context, user *User) (*User, error) {
	payload := &User{
		SuspendedAt: user.SuspendedAt,
	}
	return u.userRepo.UpdateUser(ctx, payload)
}

func (u *userUsecase) DeleteUser(ctx context.Context, user *User) (*User, error) {
	payload := &User{
		DeletedAt: user.DeletedAt,
	}
	return u.userRepo.UpdateUser(ctx, payload)
}
