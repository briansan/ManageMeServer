package schema

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/briansan/ManageMeServer/errors"
)

// TimeRange defines a start and finish time
type TimeRange struct {
	// Start is a Unix timestamp
	Start *int `bson:"start,omitempty" json:"start"`
	// Finish is a Unix timestamp
	Finish *int `bson:"finish,omitempty" json:"finish"`
}

func (t *TimeRange) Validate() error {
	if t.Start == nil || *t.Start == 0 {
		return errors.NewValidationError("start", "unix timestamp (int)")
	}
	if t.Finish == nil || *t.Finish == 0 {
		return errors.NewValidationError("finish", "unix timestamp (int)")
	}
	if *t.Start > *t.Finish {
		return errors.NewValidationError("start", "less than finish")
	}
	return nil
}

func NewTimeRange(start, finish int) *TimeRange {
	return &TimeRange{Start: &start, Finish: &finish}
}

type Task struct {
	TimeRange `bson:",inline" json:",inline"`

	ID     bson.ObjectId  `bson:"_id" json:"id"`
	UserID *bson.ObjectId `bson:"userID" json:"userID"`

	User        *string `json:"user,omitempty"`
	Title       string  `bson:"title" json:"title"`
	Description string  `bson:"description" json:"description"`
}

type TaskPatch struct {
	TimeRange `bson:",inline" json:",inline"`

	Title       *string `bson:"title,omitempty" json:"title,omitempty"`
	Description *string `bson:"description,omitempty" json:"description,omitempty"`
}

func (t *Task) Validate() error {
	if t.UserID == nil {
		return errors.NewValidationError("userID", "string")
	}
	if len(t.Title) == 0 {
		return errors.NewValidationError("title", "string")
	}
	return t.TimeRange.Validate()
}

func (t *TaskPatch) NoChange() bool {
	return t.Title == nil && t.Description == nil && t.TimeRange.Start == nil && t.TimeRange.Finish == nil
}
