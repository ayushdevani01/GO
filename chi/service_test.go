package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterUserSucess(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewUserService(mockRepo)

	dto := CreateUserDTO{
		Username: "test1123",
		Email:    "a@aa.aa",
		Age:      25,
	}

	response, err := service.RegisterUser(dto)

	require.NoError(t, err)
	require.NotNil(t, response)

}

func TestRegisterUser_ValidationFailure_Age(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewUserService(mockRepo)

	dto := CreateUserDTO{
		Username: "test1123",
		Email:    "a@aa.aa",
		Age:      11,
	}
	response, err := service.RegisterUser(dto)

	require.Error(t, err)
	assert.Nil(t, response)

	assert.Contains(t, err.Error(), "user must be at least 18 years old")
}

func TestRegisterUser_ValidationFailure_EmptyFields(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewUserService(mockRepo)

	dto := CreateUserDTO{
		Username: "test1123",
		Email:    "",
		Age:      25,
	}

	response, err := service.RegisterUser(dto)

	require.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "username and email cannot be empty")
}

func TestRegisterUser_DuplicateUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewUserService(mockRepo)

	dto := CreateUserDTO{
		Username: "test1123",
		Email:    "a@aaa.aaa",
		Age:      25,
	}

	_, err := service.RegisterUser(dto)
	require.NoError(t, err, "This first request should pass!")

	response, err := service.RegisterUser(dto)
	require.Error(t, err, "This should fail!!!")
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to save user: user already exists")
}
