package database

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (db *DB) CreateUser(email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 16)
	id := len(dbStructure.Users) + 1
	user := User{
		ID:       id,
		Email:    email,
		Password: string(encryptedPassword),
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUser(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}

	return user, nil
}

func (db *DB) LoginUser(email, password string) (User, string, error) {
	dbStructure, err := db.loadDB()
	var foundUser *User
	if err != nil {
		return User{}, "401 Unauthorized", fmt.Errorf("Incorrect email or password")
	}
	pwByte := []byte(password)
	for _, user := range dbStructure.Users {
		if user.Email == email {
			foundUser = &user
			err := bcrypt.CompareHashAndPassword(pwByte, []byte(foundUser.Password))
			if err != nil {
				return User{}, "401 Unauthorized", fmt.Errorf("Incorrect password")
			}
		}
	}
	return User{
		ID:    foundUser.ID,
		Email: foundUser.Email,
	}, "200 OK", nil
}
