package main

import (
	"fmt"
	mod "forum/pkg/models"
	"net/http"
	"strconv"
)

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		user, err := app.user.GetUserByName(w, r)
		if err != nil {
			app.serviceMsg(w, r, &templateData{}, TMPL_LOGIN, mod.ErrNoSuchUserOrPSW, http.StatusForbidden)
			return
		}

		if !app.checkPassword(w, r, user) {
			app.serviceMsg(w, r, &templateData{}, TMPL_LOGIN, mod.ErrNoSuchUserOrPSW, http.StatusForbidden)
			return
		}

		// Deletes client's sessions from another browsers
		if err := app.user.DeleteSession(user.ID); err != nil {
			app.serverError(w, err)
			return
		}

		app.deleteCookie(w, r, "session_id")
		app.deleteCookie(w, r, "user_id")

		cookie := setCookie(w)
		setUserIDCookie(strconv.Itoa(user.ID), w)

		if err := app.user.CreateSession(user, cookie); err != nil {
			app.serverError(w, err)
			return
		}

		http.Redirect(w, r, URL_HOME, http.StatusSeeOther)
	}

	app.render(w, r, TMPL_LOGIN, &templateData{})
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {

	if !app.isAuthorized(w, r) {
		http.Redirect(w, r, URL_HOME, http.StatusSeeOther)
		return
	}

	user, err := app.user.GetUserBySession(w, r)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.deleteCookie(w, r, "session_id")
	app.deleteCookie(w, r, "user_id")

	if err := app.user.DeleteSession(user.ID); err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, URL_LOGIN, http.StatusSeeOther)
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		userDTO := &mod.CreateUserDTO{
			UserName:    r.FormValue("username"),
			Email:       r.FormValue("email"),
			Password:    r.FormValue("psw"),
			PasswordRPT: r.FormValue("psw-repeat"),
		}

		if err := IsValidCredentials(userDTO); err != nil {
			app.serviceMsg(w, r, &templateData{}, TMPL_SIGNUP, err, http.StatusForbidden)
			return
		}

		user := &mod.User{
			UserName: r.FormValue("username"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("psw"),
		}

		// Добавляем нового юзера в БД, иначе уточнем какие пар-ы введены неверно
		err := app.user.CreateUserDB(user)
		if err != nil {
			// TODO
			fmt.Println(err)
			app.serviceMsg(w, r, &templateData{}, TMPL_SIGNUP, mod.ErrUserAlreadyExists, http.StatusConflict)
			return
		}
		http.Redirect(w, r, URL_LOGIN, http.StatusSeeOther)
	}

	app.render(w, r, TMPL_SIGNUP, &templateData{})
}
