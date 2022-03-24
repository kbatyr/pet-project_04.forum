package main

import (
	"errors"
	"forum/pkg/models"
	"net/http"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func IsValidCredentials(u *models.CreateUserDTO) error {
	// Проверка на корректность логина и почты

	if u.UserName == "" || len(u.UserName) < 5 || len(u.UserName) > 20 || len(strings.TrimSpace(u.UserName)) == 0 {
		return errors.New("the name field must be between 5-20 chars")
	}

	if u.Email == "" || len(strings.TrimSpace(u.Email)) == 0 {
		return errors.New("the email field is required")
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return errors.New("the email field should be a valid email address")
	}

	if u.Password == "" || len(u.Password) < 5 || len(strings.TrimSpace(u.Password)) == 0 {
		return errors.New("the password field must be at least 5 chars")
	}

	if u.Password != u.PasswordRPT {
		return errors.New("the entered passwords do not match")
	}

	return nil
}

func (app *application) isAuthorized(w http.ResponseWriter, r *http.Request) bool {

	_, err := r.Cookie("session_id")
	if err != nil {
		return false
	}

	_, err = app.user.GetUserBySession(w, r)

	return err == nil
}

func (app *application) checkPassword(w http.ResponseWriter, r *http.Request, u *models.User) bool {

	psw := r.FormValue("password")
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(psw))

	return err == nil
}
