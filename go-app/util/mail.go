package util

import (
	"fmt"
	"go-app/models"
)

// SendResetPasswordEmail sends a mail to the user with a validation request link
func SendResetPasswordEmail(user *models.User, token string) error {
	// send mail with encoded token
	fmt.Println("Send mail here")
	return nil
}
