package model

import "time"

func init() {
	RegisterModel(&User{})
}

// DefaultAdminUserName default admin user name
const DefaultAdminUserName = "admin"

// DefaultAdminUserAlias default admin user alias
const DefaultAdminUserAlias = "Administrator"

// User is the model of user
type User struct {
	BaseModel
	Name          string    `json:"name"`
	Password      string    `json:"password,omitempty"`
	Disabled      bool      `json:"disabled"`
	LastLoginTime time.Time `json:"lastLoginTime,omitempty"`
}

func (u *User) TableName() string {
	return tabelNamePrefix + "user"
}

func (u *User) ShortTableName() string {
	return "usr"
}

// PrimaryKey return custom primary key
func (u *User) PrimaryKey() string {
	return u.Name
}

// Index return custom index
func (u *User) Index() map[string]string {
	index := make(map[string]string)
	if u.Name != "" {
		index["name"] = u.Name
	}
	return index
}
