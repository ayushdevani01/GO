package main

import (
	"errors"
	"fmt"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(dto CreateUserDTO) (*UserResponseDTO, error) {
	if dto.Username == "" || dto.Email == "" {
		return nil, errors.New("username and email cannot be empty")
	}

	if dto.Age < 18 {
		return nil, errors.New("user must be at least 18 years old")
	}

	user := User{
		ID:       "tempid",
		Username: dto.Username,
		Email:    dto.Email,
		Age:      dto.Age,
	}

	if err := s.repo.Save(user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	response := &UserResponseDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	return response, nil
}

func (s *UserService) GetAllUsers() ([]UserResponseDTO, error) {
	users, err := s.repo.Findall()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}
	var responses []UserResponseDTO
	for _, user := range users {
		responses = append(responses, UserResponseDTO{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		})
	}
	return responses, nil
}
