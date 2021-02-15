package main

//UserInfo represents information about a user
type UserInfo struct {
	Fullname string `json:"fullname,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

//AuthRequest represents a request packet received from the HTTP Server
type AuthRequest struct {
	Method   string   `json:"method,omitempty"`
	UserInfo UserInfo `json:"user_info,omitempty"`
}

//AuthResponse represents a request packet sent to the HTTP Server
type AuthResponse struct {
	HTTPCode string
	Message  string
}
