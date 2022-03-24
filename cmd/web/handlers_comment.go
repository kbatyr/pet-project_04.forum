package main

import (
	"errors"
	"fmt"
	mod "forum/pkg/models"
	"net/http"
	"strconv"
)

func (app *application) leaveComment(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			app.notFound(w)
			return
		}

		uID := getUserIDByCookie(r)

		p, err := app.posts.GetPost(id)
		if err != nil {
			if errors.Is(err, mod.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
			}
			return
		}

		p.Rating, err = app.posts.GetRating(id, "rating_info")
		if err != nil {
			app.serverError(w, err)
			return
		}

		p.Categories, err = app.posts.GetCategory(id)
		if err != nil {
			app.serverError(w, err)
			return
		}

		p.Comments, err = app.posts.GetComments(id)
		if err != nil {
			app.serverError(w, err)
			return
		}

		for _, com := range p.Comments {
			com.Rating, err = app.posts.GetRating(id, "comments_rating")
			if err != nil {
				app.serverError(w, err)
				return
			}
		}

		msg := r.FormValue("comment")
		if err := IsValidComment(msg); err != nil {
			app.serviceMsg(w, r, &templateData{Post: p, Authorized: true}, TMPL_SHOWPOST, err, http.StatusBadRequest)
			return
		}

		if err := app.posts.CreateComment(p.ID, uID, msg); err != nil {
			app.serverError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/post?id=%d", id), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, URL_HOME, http.StatusSeeOther)
	}
}
