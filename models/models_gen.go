// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

import (
	"time"
)

type Auth struct {
	AuthToken *AuthToken `json:"authToken"`
	User      *User      `json:"user"`
}

type AuthToken struct {
	AccessToken string    `json:"accessToken"`
	ExpiredAt   time.Time `json:"expiredAt"`
}

type MeetupFilter struct {
	Name *string `json:"name"`
}

type NewMeetup struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateMeetup struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
