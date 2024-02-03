package models

import "time"

type User struct {
	ID        uint      `json:"ID,omitempty"`
	CreatedAt time.Time `json:"CreatedAt"`
	Username  string    `json:"username" json:"Username,omitempty"`
	Password  string    `json:"password" json:"Password,omitempty"`
	UserType  string    `json:"userType"`
}
