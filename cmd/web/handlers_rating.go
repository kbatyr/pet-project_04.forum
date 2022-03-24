package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (app *application) postRating(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		table := "rating_info"
		if strings.Contains(r.URL.Path, "comment") {
			table = "comments_rating"
		}

		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			app.notFound(w)
			return
		}

		u, err := app.user.GetUserBySession(w, r)
		if err != nil {
			app.serverError(w, err)
			return
		}

		reaction := r.FormValue("reaction")

		// Check whether user already liked/disliked the post
		if reaction == "like" {
			if app.posts.AlreadyReacted(u, id, reaction, table) {
				reaction = "unlike"
			}
		} else if reaction == "dislike" {
			if app.posts.AlreadyReacted(u, id, reaction, table) {
				reaction = "undislike"
			}
		} else {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		switch reaction {
		case "like":
			if err := app.posts.UserReaction(u, id, reaction, table); err != nil {
				if app.isConstraintErr(err) {
					app.clientError(w, http.StatusBadRequest)
				} else {
					app.serverError(w, err)
				}
				return
			}
		case "dislike":
			if err := app.posts.UserReaction(u, id, reaction, table); err != nil {
				if app.isConstraintErr(err) {
					app.clientError(w, http.StatusBadRequest)
				} else {
					app.serverError(w, err)
				}
				return
			}
		case "unlike":
			if err := app.posts.RevokeReaction(u, id, reaction, table); err != nil {
				if app.isConstraintErr(err) {
					app.clientError(w, http.StatusBadRequest)
				} else {
					app.serverError(w, err)
				}
				return
			}

		case "undislike":
			if err := app.posts.RevokeReaction(u, id, reaction, table); err != nil {
				if app.isConstraintErr(err) {
					app.clientError(w, http.StatusBadRequest)
				} else {
					app.serverError(w, err)
				}
				return
			}
		default:
			app.clientError(w, http.StatusBadRequest)
			return
		}

		if table == "comments_rating" {
			id, err = strconv.Atoi(r.URL.Query().Get("post_id"))
			if err != nil || id < 1 {
				app.notFound(w)
				return
			}
		}

		http.Redirect(w, r, fmt.Sprintf("/post?id=%d", id), http.StatusSeeOther)

	} else {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}
