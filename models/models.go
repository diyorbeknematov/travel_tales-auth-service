package models

type UserLogin struct {
	ID       string
	Username string
	Email    string
}

type UpdatePassword struct {
	ID          string
	NewPassword string
}
