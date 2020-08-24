package db

import "fmt"

type UserNotFoundError struct {
	UserID RowID
}

func (u *UserNotFoundError) Error() string {
	return fmt.Sprintf("User #%d not found", u.UserID)
}