package data

type Settings struct {
	AllowMembersToInvite   bool   `json:"allow_members_to_invite" bson:"allow_members_to_invite"`
	AdminRoleConfirmation  bool   `json:"admin_role_confirmation" bson:"admin_role_confirmation"`
	InviteConfirmation     bool   `json:"invite_confirmation" bson:"invite_confirmation"`
	InviteConfirmationRole string `json:"invite_confirmation_role" bson:"invite_confirmation_role"`
	DefualtInviteRole      string `json:"default_invite_role" bson:"default_invite_role"`
}

type Org struct {
	ID          string   `json:"id" bson:"_id,omitempty"`
	Name        string   `json:"name" bson:"name"`
	Description string   `json:"description" bson:"description"`
	OwnerID     string   `json:"owner_id" bson:"owner_id"`
	Settings    Settings `json:"settings" bson:"settings"`

	CreatedAt int64 `json:"created_at" bson:"created_at"`
	UpdatedAt int64 `json:"updated_at" bson:"updated_at"`
	DeletedAt int64 `json:"deleted_at" bson:"deleted_at"`
}
