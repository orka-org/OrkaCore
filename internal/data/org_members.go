package data

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrgMember struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	OrgID       string             `json:"org_id" bson:"org_id"`
	UserID      string             `json:"user_id" bson:"user_id"`
	RoleID      string             `json:"role_id" bson:"role_id"`
	InvitedByID string             `json:"invited_by_id" bson:"invited_by_id"`

	CreatedAt string `json:"created_at" bson:"created_at"`
	UpdatedAt string `json:"updated_at" bson:"updated_at"`
	DeletedAt string `json:"deleted_at" bson:"deleted_at"`
}
