package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	result := Add(2, 3)
	assert.Equal(t, 5, result, "they should be equal")
}

func TestDevide_Sucess(t *testing.T) {
	result, err := Devide(10, 2)
	require.NoError(t, err)
	assert.Equal(t, 5, result, "they should be equal")
}

func TestDevide_ByZero(t *testing.T) {
	_, err := Devide(10, 0)
	require.Error(t, err)
	assert.Equal(t, "Cannot devide by zero", err.Error())
}
