package data

import (
	"time"

	"github.com/orka-org/orkacore/internal/biz"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func BiztoDB(user *biz.User) (*User, error) {
	var data User
	if user.ID != "" {
		oid, err := primitive.ObjectIDFromHex(user.ID)
		if err != nil {
			return nil, err
		}
		data.ID = oid
	}

	if user.Username != "" {
		data.Username = user.Username
	}
	if user.FirstName != nil {
		data.FirstName = user.FirstName
	}
	if user.LastName != nil {
		data.LastName = user.LastName
	}
	if user.Email != "" {
		data.Email = user.Email
	}
	if user.Phone != nil {
		data.Phone = user.Phone
	}
	if user.Password != "" {
		data.Password = user.Password
	}

	data.CreatedAt = user.CreatedAt.Unix()
	data.UpdatedAt = user.UpdatedAt.Unix()
	if user.DeletedAt != nil {
		unix := user.DeletedAt.Unix()
		data.DeletedAt = &unix
	}
	if user.SuspendedAt != nil {
		unix := user.SuspendedAt.Unix()
		data.SuspendedAt = &unix
	}

	return &data, nil
}

func DBtoBiz(user *User) (*biz.User, error) {
	var bizUser biz.User

	str := user.ID.Hex()
	bizUser.ID = str

	bizUser.Username = user.Username
	bizUser.FirstName = user.FirstName
	bizUser.LastName = user.LastName
	bizUser.Email = user.Email
	bizUser.Phone = user.Phone
	bizUser.Password = user.Password

	if user.CreatedAt != 0 {
		bizUser.CreatedAt = time.Unix(user.CreatedAt, 0)
	}
	if user.UpdatedAt != 0 {
		bizUser.UpdatedAt = time.Unix(user.UpdatedAt, 0)
	}
	if user.DeletedAt != nil {
		unix := *user.DeletedAt
		date := time.Unix(unix, 0)
		bizUser.DeletedAt = &date
	}
	if user.SuspendedAt != nil {
		unix := *user.SuspendedAt
		date := time.Unix(unix, 0)
		bizUser.SuspendedAt = &date
	}

	return &bizUser, nil
}
