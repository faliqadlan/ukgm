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
	DeletedBy AuthUser   `json:"-" bson:"deleted_by"`
}

type TimeRange struct {
	From *time.Time `json:"from" bson:"from"`
	To   *time.Time `json:"to" bson:"to"`
}

func (timeRange *TimeRange) Validate() error {
	if timeRange.From.Before(*timeRange.To) {
		return nil
	}
	return errors.New("from cannot be after To")
}

func (timestamp *Timestamp) CreateNew(auth *AuthUser) {
	location, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(location).Local()
	timestamp.CreatedAt = &now
	timestamp.UpdatedAt = &now
	timestamp.CreatedBy = *auth
	timestamp.UpdatedBy = *auth
}

func (timestamp *Timestamp) Update(auth *AuthUser) {
	location, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(location).Local()
	timestamp.UpdatedAt = &now
	timestamp.UpdatedBy = *auth
}

func (timestamp *Timestamp) Delete(auth *AuthUser) {
	location, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(location).Local()
	timestamp.DeletedAt = &now
	timestamp.DeletedBy = *auth
}
