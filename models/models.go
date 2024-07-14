package models

type UserLogin struct {
	ID       string
	Username string
	Email    string
	Password string
}

type UpdatePassword struct {
	ID          string
	NewPassword string
}

type Errors struct {
	Message string
}

type Token struct {
	AccessToken  string
	RefreshToken string
}