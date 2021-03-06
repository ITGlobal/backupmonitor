package model

// AuthRequest contains parameters for authentication
type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// String converts an object to string
func (p AuthRequest) String() string {
	return toJSON(&p)
}

// AuthResponse contains a generated access token
type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// String converts an object to string
func (p AuthResponse) String() string {
	return toJSON(&p)
}
