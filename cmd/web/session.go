package main

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
)

const sessionLength = 600

func getUserIDByCookie(r *http.Request) string {
	cookie, _ := r.Cookie("user_id")
	return cookie.Value
}

func setUserIDCookie(id string, w http.ResponseWriter) *http.Cookie {

	cookie := &http.Cookie{
		Name:  "user_id",
		Value: id,
	}

	http.SetCookie(w, cookie)
	return cookie
}

func setCookie(w http.ResponseWriter) *http.Cookie {

	sessionID := uuid.NewV4()
	cookie := &http.Cookie{
		Name:   "session_id",
		Value:  sessionID.String(),
		MaxAge: sessionLength,
	}

	http.SetCookie(w, cookie)
	return cookie
}

func (app *application) deleteCookie(w http.ResponseWriter, r *http.Request, name string) {

	cookie, _ := r.Cookie(name)

	cookie = &http.Cookie{
		Name:   name,
		Value:  "",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
}

func (app *application) refreshCookie(w http.ResponseWriter, r *http.Request) error {

	c, _ := r.Cookie("session_id")
	c.MaxAge = sessionLength
	http.SetCookie(w, c)

	if err := app.user.UpdateSession(c); err != nil {
		return err
	}
	return nil
}
