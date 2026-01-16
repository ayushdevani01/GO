package main

import "errors"

type UserRepository interface {
	Save(user User) error
	FindByID(id string) (User, error)
	Findall() ([]User, error)
} // interface defining user repository methods

type MockUserRepository struct {
	data map[string]User
} // mock implementation of UserRepository

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		data: make(map[string]User), // Db connection
	}
}

func (r *MockUserRepository) Save(user User) error {
	if _, exists := r.data[user.ID]; exists {
		return errors.New("user already exists")
	}
	r.data[user.ID] = user
	return nil
}

func (r *MockUserRepository) FindByID(id string) (User, error) {
	user, exists := r.data[id]
	if !exists {
		return User{}, errors.New("user not found")
	}
	return user, nil
}
func (r *MockUserRepository) Findall() ([]User, error) {
	users := make([]User, 0, len(r.data))
	for _, user := range r.data {
		users = append(users, user)
	}
	return users, nil
}
