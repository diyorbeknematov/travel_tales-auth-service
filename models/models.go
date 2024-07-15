package models

type Register struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FullName  string `json:"full_name"`
	CreatedAt string `json:"created_at"`
}

type UserLogin struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResetPassword struct {
	Email string `json:"email"`
}

type UpdatePassword struct {
	ID          string `json:"id"`
	NewPassword string `json:"new_password"`
}


type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Errors struct {
	Message string `json:"message"`
}

type Success struct {
	Message string `json:"message"`
}
