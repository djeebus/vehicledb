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
	"password_hash" TEXT
)`

type User struct {
	EmailAddress string `json:"email_address"`
	PasswordHash []byte `json:"-"`
	UserId       RowID  `json:"user_id"`
}

func (u *User) DoesPasswordMatch(password string) bool {
	passwordHash := hashPassword(password)
	if len(passwordHash) != len(u.PasswordHash) {
		return false
	}

	for idx, a := range passwordHash {
		if a != u.PasswordHash[idx] {
			return false
		}
	}
	return true
}

func hashPassword(password string) []byte {
	passwordBytes := []byte(password) // strings are utf-8 encoded already
	hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	return hash
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
		UserId:       RowID(lastInserted),
	}
	return &user, nil
}

func FindUserByEmailAddress(emailAddress string) (*User, error) {
	query := `SELECT id, password_hash FROM users WHERE email_address = ?`
	stmt, err := sqlDb.Prepare(query)
	if err != nil {
		return nil, err
	}

	row, err := stmt.Query(emailAddress)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	var (
		userId       RowID
		passwordHash []byte
	)

	for row.Next() {
		err = row.Scan(&userId)
		if err != nil {
			return nil, err
		}

		user := User{
			UserId:       userId,
			EmailAddress: emailAddress,
			PasswordHash: passwordHash,
		}
		return &user, nil
	}

	return nil, &EmailAddressNotFoundError{EmailAddress: emailAddress}
}

func GetUser(userId RowID) (*User, error) {
	query := `SELECT email_address, password_hash FROM users WHERE id = ?`
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
		var (
			passwordHash []byte
			emailAddress string
		)

		err = row.Scan(&emailAddress, &passwordHash)
		if err != nil {
			return nil, err
		}

		user := User{
			UserId:       userId,
			EmailAddress: emailAddress,
			PasswordHash: passwordHash,
		}
		return &user, nil
	}

	return nil, &UserNotFoundError{UserID: userId}
}

func UpdateUser(userId RowID, emailAddress string) (*User, error) {
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
		UserId:       userId,
		EmailAddress: emailAddress,
	}

	return &user, nil
}

func DeleteUser(userId RowID) error {
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
