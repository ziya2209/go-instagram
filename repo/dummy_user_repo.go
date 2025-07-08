package repo

import (
	"errors"
	"instagram/models"
	"sync"
)

type DummyUserRepo struct {
	users  []models.User
	mutex  sync.RWMutex
	nextID int
}

func NewDummyUserRepo() UserRepo {
	return &DummyUserRepo{
		users:  make([]models.User, 0),
		nextID: 1,
	}
}

func (r *DummyUserRepo) Create(user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if email already exists
	for _, existingUser := range r.users {
		if existingUser.Email == user.Email {
			return errors.New("email already exists")
		}
	}

	// Check if username already exists
	for _, existingUser := range r.users {
		if existingUser.Username == user.Username {
			return errors.New("username already exists")
		}
	}

	// Assign ID and add to slice
	user.Id = r.nextID
	r.nextID++
	r.users = append(r.users, *user)

	return nil
}

func (r *DummyUserRepo) GetById(id int) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Id == id {
			return &user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (r *DummyUserRepo) GetByUsername(username string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			return &user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (r *DummyUserRepo) Update(user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, existingUser := range r.users {
		if existingUser.Id == user.Id {
			// Check if email is being changed and if new email already exists
			if existingUser.Email != user.Email {
				for j, otherUser := range r.users {
					if j != i && otherUser.Email == user.Email {
						return errors.New("email already exists")
					}
				}
			}

			// Check if username is being changed and if new username already exists
			if existingUser.Username != user.Username {
				for j, otherUser := range r.users {
					if j != i && otherUser.Username == user.Username {
						return errors.New("username already exists")
					}
				}
			}

			r.users[i] = *user
			return nil
		}
	}

	return errors.New("user not found")
}

func (r *DummyUserRepo) Delete(id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, user := range r.users {
		if user.Id == id {
			// Remove user from slice
			r.users = append(r.users[:i], r.users[i+1:]...)
			return nil
		}
	}

	return errors.New("user not found")
}

func (r *DummyUserRepo) GetAll() ([]*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make([]*models.User, len(r.users))
	for i := range r.users {
		result[i] = &r.users[i]
	}

	return result, nil
}
