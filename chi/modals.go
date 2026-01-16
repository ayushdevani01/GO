package main

type CreateUserDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

type UserResponseDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type User struct {
	ID       string
	Username string
	Email    string
	Age      int
}
