package entities

import (
	"errors"
	"time"
)

type Timestamp struct {
	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	CreatedBy AuthUser   `json:"created_by" bson:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" bson:"updated_at"`
	UpdatedBy AuthUser   `json:"updated_by" bson:"updated_by"`
	DeletedAt *time.Time `json:"-" bson:"deleted_at"`
}

type TimeRange struct {
	From *time.Time `json:"from" bson:"from"`
	To   *time.Time `json:"to" bson:"to"`
}

func (timeRange *TimeRange) Validate() error {
	if timeRange.From.Before(*timeRange.To) {
		return nil
	}
	return errors.New("From cannot be after To")
}

func (timestamp *Timestamp) CreateNew(auth *Auth) {
	location, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(location)
	timestamp.CreatedAt = &now
	timestamp.UpdatedAt = &now
	timestamp.CreatedBy = auth.User
	timestamp.UpdatedBy = auth.User
	return
}

func (timestamp *Timestamp) Update(auth *Auth) {
	location, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(location)
	timestamp.UpdatedAt = &now
	timestamp.UpdatedBy = auth.User
	return
}

func (timestamp *Timestamp) Delete(auth *Auth) {
	location, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(location)
	timestamp.DeletedAt = &now
	return
}
