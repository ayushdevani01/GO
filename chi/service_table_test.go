package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterUserTable(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockService := NewUserService(mockRepo)

	tests := []struct {
		name         string
		input        CreateUserDTO
		expectError  bool
		errorMsg     string
		expectedUser *UserResponseDTO
	}{
		{
			name: "Success --- Valid Input",
			input: CreateUserDTO{
				Username: "validuser",
				Email:    "123@1123.123",
				Age:      30,
			},
			expectError: false,
			expectedUser: &UserResponseDTO{
				ID:       "tempid",
				Username: "validuser",
				Email:    "123@1123.123",
			},
		},
		{
			name: "Failure --- Empty Username",
			input: CreateUserDTO{
				Username: "",
				Email:    "123@1123.123",
				Age:      30,
			},
			expectError: true,
			errorMsg:    "username and email cannot be empty",
		},
		{
			name: "Failure --- Underage User",
			input: CreateUserDTO{
				Username: "younguser",
				Email:    "123@1123.123",
				Age:      16,
			},
			expectError: true,
			errorMsg:    "user must be at least 18 years old",
		},
		{
			name: "Failure --- Empty Email",
			input: CreateUserDTO{
				Username: "validuser",
				Email:    "",
				Age:      30,
			},
			expectError: true,
			errorMsg:    "username and email cannot be empty",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := mockService.RegisterUser(tc.input)
			if tc.expectError {
				require.Error(t, err)
				assert.Equal(t, tc.errorMsg, err.Error())
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedUser, result)
			}
		})
	}
}
