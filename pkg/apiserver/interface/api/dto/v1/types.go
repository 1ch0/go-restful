package v1

import "time"

type UserBase struct {
	CreateTime    time.Time `json:"createTime"`
	LastLoginTime time.Time `json:"lastLoginTime"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Alias         string    `json:"alias,omitempty"`
	Disabled      bool      `json:"disabled"`
}

// SimpleResponse simple response model for temporary
type SimpleResponse struct {
	Status string `json:"status"`
}

type LoginRequest struct {
	Code     string `json:"code,omitempty" optional:"true"`
	Username string `json:"username,omitempty" optional:"true"`
	Password string `json:"password,omitempty" optional:"true"`
}

type LoginResponse struct {
	User         *UserBase `json:"user"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type ListUserResponse struct {
}
