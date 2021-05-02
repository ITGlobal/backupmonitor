package model

import "golang.org/x/crypto/bcrypt"

// User contains information about application user
type User struct {
	ID           int    `json:"id"`
	UserName     string `json:"username"`
	PasswordHash string `json:"-"`
}

// String converts an object to string
func (p *User) String() string {
	return toJSON(&p)
}

// CheckPassword checks user's password
func (p *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.PasswordHash), []byte(password))
	return err == nil
}

// SetPassword sets user's password
func (p *User) SetPassword(password string) error {
	buff, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.PasswordHash = string(buff)
	return nil
}

// Users is a list of User
type Users []*User

// UserChangePasswordRequest contains parameters to change user's password
type UserChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// String converts an object to string
func (p *UserChangePasswordRequest) String() string {
	return toJSON(&p)
}
