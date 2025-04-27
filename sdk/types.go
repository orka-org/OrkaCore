package sdk

import (
	pb "github.com/orka-org/orkacore/api/auth/v1"
)

type User struct {
	ID          string
	Email       string
	Username    string
	FirstName   *string
	LastName    *string
	Phone       *string
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   *string
	SuspendedAt *string
}

func protoToUser(u *pb.User) *User {
	if u == nil {
		return nil
	}
	return &User{
		ID:          u.Id,
		Email:       u.Email,
		Username:    u.Username,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Phone:       u.Phone,
		CreatedAt:   *u.CreatedAt,
		UpdatedAt:   *u.UpdatedAt,
		DeletedAt:   u.DeletedAt,
		SuspendedAt: u.SuspendedAt,
	}
}
