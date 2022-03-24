package main

import (
	"net/http"
)

func (app *application) alreadyLoggedIn(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Code for middleware
		if app.isAuthorized(w, r) {

			if err := app.refreshCookie(w, r); err != nil {
				app.serverError(w, err)
				return
			}

			http.Redirect(w, r, URL_HOME, http.StatusSeeOther)
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) hasPersmission(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !app.isAuthorized(w, r) {
			http.Redirect(w, r, URL_LOGIN, http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
