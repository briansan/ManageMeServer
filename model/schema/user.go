package schema

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/briansan/ManageMeServer/errors"
)

type UserSecure struct {
	ID             bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Username       string        `bson:"username" json:"username"`
	Email          string        `bson:"email" json:"email"`
	Role           int           `bson:"role" json:"role"`
	PreferredHours *TimeRange    `bson:"preferredHours" json:"preferredHours"`
}

type User struct {
	ID             bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Username       *string       `bson:"username,omitempty" json:"username,omitempty"`
	OldPassword    *string       `bson:"-" json:"oldPassword,omitempty"`
	Password       *string       `bson:"password,omitempty" json:"password,omitempty"`
	Email          *string       `bson:"email,omitempty" json:"email,omitempty"`
	Role           *int          `bson:"role,omitempty" json:"role"`
	PreferredHours *TimeRange    `bson:"preferredHours" json:"preferredHours"`
}

func (u *User) Validate() error {
	if u.Email == nil || len(*u.Email) == 0 {
		return errors.NewValidationError("email", "string")
	}
	if u.Username == nil || len(*u.Username) == 0 {
		return errors.NewValidationError("username", "string")
	}
	if u.Password == nil || len(*u.Password) == 0 {
		return errors.NewValidationError("password", "string")
	}
	if u.PreferredHours != nil {
		return u.PreferredHours.Validate()
	}
	return nil
}
