package models

type Admin struct {
    Base
    Email    string `json:"email" gorm:"unique"`
    Username string `json:"username" gorm:"unique"`
    Password string `json:"password"`
}

type AdminErrors struct {
	Err      bool   `json:"error"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}