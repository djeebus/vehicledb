package db

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

var usersTable = `
CREATE TABLE users (
	"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	"email_address" TEXT,
	"password_hash" INTEGER NOT NULL DEFAULT 0
)`

type User struct {
	EmailAddress string `json:"email_address"`
	UserId       int64  `json:"user_id"`
}

func hashPassword(password string) string {
	passwordBytes := []byte(password) // strings are utf-8 encoded already
	hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}

func validatePassword(password string, passwordHash []byte) bool {
	return false
}

func CreateUser(emailAddress string, password string) (*User, error) {
	query := `INSERT INTO users (email_address, password_hash) VALUES (?, ?)`
	stmt, err := sqlDb.Prepare(query)
	if err != nil {
		return nil, err
	}

	result, err := stmt.Exec(emailAddress, hashPassword(password))
	if err != nil {
		return nil, err
	}

	lastInserted, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user := User{
		EmailAddress: emailAddress,
		UserId:       lastInserted,
	}
	return &user, nil
}

func GetUser(userId int64) (*User, error) {
	query := `SELECT email_address FROM users WHERE id = ?`
	stmt, err := sqlDb.Prepare(query)
	if err != nil {
		return nil, err
	}

	row, err := stmt.Query(userId)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	for row.Next() {
		var emailAddress string

		err = row.Scan(&emailAddress)
		if err != nil {
			return nil, err
		}

		user := User{
			UserId:       userId,
			EmailAddress: emailAddress,
		}
		return &user, nil
	}

	return nil, &UserNotFoundError{UserID: userId}
}

func UpdateUser(userId int64, emailAddress string) (*User, error) {
	// nothing to update
	if emailAddress == "" {
		return GetUser(userId)
	}

	query := `UPDATE users SET email_address=? WHERE id = ?`
	stmt, err := sqlDb.Prepare(query)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(emailAddress, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %s", err)
	}

	user := User{
		UserId: userId,
		EmailAddress: emailAddress,
	}

	return &user, nil
}

func DeleteUser(userId int64) error {
	query := `DELETE FROM users WHERE id = ?`
	stmt, err := sqlDb.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare delete statement: %v", err)
	}

	_, err = stmt.Exec(userId)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}
