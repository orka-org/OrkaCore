package data

import "go.mongodb.org/mongo-driver/bson/primitive"

type Permission struct {
	Object string `json:"object" bson:"object"`
	Action string `json:"action" bson:"action"`
}

type Role struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	OrgID       string             `json:"org_id" bson:"org_id"`
	Name        string             `json:"name" bson:"name"`
	Permisisons []Permission       `json:"permissions" bson:"permissions"`

	CreatedAt string `json:"created_at" bson:"created_at"`
	UpdatedAt string `json:"updated_at" bson:"updated_at"`
	DeletedAt string `json:"deleted_at" bson:"deleted_at"`
}
