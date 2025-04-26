package data

import (
	"context"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/orka-org/orkacore/internal/biz"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id" `
	Username  string             `bson:"username"`
	Password  string             `bson:"password"`
	FirstName *string            `bson:"first_name"`
	LastName  *string            `bson:"last_name"`
	Email     string             `bson:"email"`
	Phone     *string            `bson:"phone"`

	CreatedAt   int64  `bson:"created_at"`
	UpdatedAt   int64  `bson:"updated_at"`
	DeletedAt   *int64 `bson:"deleted_at"`
	SuspendedAt *int64 `bson:"suspended_at"`
}

type userRepo struct {
	coll *mongo.Collection
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) (biz.UserRepo, error) {
	if data.db == nil {
		log.NewHelper(logger).Error("mongodb not initialized")
		return nil, errors.New("mongodb not initialized")
	}
	coll := data.db.Collection("users")
	return &userRepo{
		coll: coll,
		log:  log.NewHelper(logger),
	}, nil
}

func (r *userRepo) CreateUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	data, err := BiztoDB(user)
	if err != nil {
		return nil, err
	}

	data.ID = primitive.NewObjectID()
	res, err := r.coll.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	id := res.InsertedID.(primitive.ObjectID).Hex()
	r.log.Info("ID: ", id)
	user.ID = id
	return user, nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*biz.User, error) {
	res := r.coll.FindOne(ctx, bson.M{"email": email})
	if res.Err() != nil {
		return nil, res.Err()
	}

	var user biz.User
	err := res.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetUserByID(ctx context.Context, id string) (*biz.User, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	res := r.coll.FindOne(ctx, bson.M{"_id": objId})
	if res.Err() != nil {
		return nil, res.Err()
	}

	var user biz.User
	err = res.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetAllUsers(ctx context.Context) ([]*biz.User, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"suspended_at": nil, "deleted_at": nil})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*biz.User

	for cursor.Next(ctx) {
		var user biz.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	objId, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return nil, err
	}

	updateFields := bson.M{"updated_at": time.Now().Unix()}

	if user.Username != "" {
		updateFields["username"] = user.Username
	}
	if user.FirstName != nil {
		updateFields["first_name"] = user.FirstName
	}
	if user.LastName != nil {
		updateFields["last_name"] = user.LastName
	}
	if user.Email != "" {
		updateFields["email"] = user.Email
	}
	if user.Phone != nil {
		updateFields["phone"] = user.Phone
	}
	if user.Password != "" {
		updateFields["password"] = user.Password
	}
	if user.Email != "" {
		updateFields["email"] = user.Email
	}
	user.UpdatedAt = time.Now()

	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": updateFields})
	if err != nil {
		return nil, err
	}
	if res.ModifiedCount == 0 {
		return nil, errors.New("user not found")
	}

	// Fetch the updated user
	result := r.coll.FindOne(ctx, bson.M{"_id": objId})
	if result.Err() != nil {
		return nil, result.Err()
	}

	var updatedUser biz.User
	if err := result.Decode(&updatedUser); err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (r *userRepo) SuspendUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	res := r.coll.FindOneAndUpdate(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"suspended_at": time.Now().Unix()}})
	if res.Err() != nil {
		return nil, res.Err()
	}

	return user, nil
}

func (r *userRepo) DeleteUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	res := r.coll.FindOneAndDelete(ctx, bson.M{"_id": user.ID})
	if res.Err() != nil {
		return nil, res.Err()
	}

	var old biz.User
	err := res.Decode(&old)
	if err != nil {
		return nil, err
	}

	return &biz.User{
		ID:          old.ID,
		Username:    old.Username,
		FirstName:   old.FirstName,
		LastName:    old.LastName,
		Email:       old.Email,
		Phone:       old.Phone,
		CreatedAt:   old.CreatedAt,
		UpdatedAt:   old.UpdatedAt,
		SuspendedAt: old.SuspendedAt,
		DeletedAt:   old.DeletedAt,
	}, nil
}
