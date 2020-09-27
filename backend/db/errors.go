package db

import "fmt"

type UserNotFoundError struct {
	UserID RowID
}

func (u *UserNotFoundError) Error() string {
	return fmt.Sprintf("User #%d not found", u.UserID)
}

type EmailAddressNotFoundError struct {
	EmailAddress string
}

func (err *EmailAddressNotFoundError) Error() string {
	return fmt.Sprintf("User with email address %s not found", err.EmailAddress)
}
