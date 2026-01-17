package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserHandler_Register_Sucess(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockService := NewUserService(mockRepo)
	handler := NewUserHandler(mockService)

	
	dto := CreateUserDTO{
		Username: "test1123",
		Email:    "a@aa.aa",
		Age:      25,
	}

	jsonPayload, _ := json.Marshal(dto)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.handleRegisterUser(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	responseBytes := rr.Body.Bytes()

	var res UserResponseDTO
	err := json.Unmarshal(responseBytes, &res)
	require.NoError(t, err)

	assert.Equal(t, dto.Username, res.Username)
	assert.Equal(t, dto.Email, res.Email)
	assert.NotEmpty(t, res.ID)

}
