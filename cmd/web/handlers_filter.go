package main

import (
	mod "forum/pkg/models"
	"net/http"
)

func (app *application) filterByUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	var authorized bool
	if app.isAuthorized(w, r) {

		authorized = true

		if err := app.refreshCookie(w, r); err != nil {
			app.serverError(w, err)
			return
		}
	}

	uID := getUserIDByCookie(r)
	posts, err := app.posts.FilterByParam("user", uID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if err := app.posts.GetAllCategories(posts); err != nil {
		app.serverError(w, err)
		return
	}

	m := map[string][]*mod.Post{uID: posts}

	app.render(w, r, TMPL_FILTER_BY_USER, &templateData{
		FilteredPosts: m,
		UserID:        uID,
		Authorized:    authorized})
}

func (app *application) filterByReaction(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	var authorized bool
	if app.isAuthorized(w, r) {

		authorized = true

		if err := app.refreshCookie(w, r); err != nil {
			app.serverError(w, err)
			return
		}
	}

	uID := getUserIDByCookie(r)
	posts, err := app.posts.FilterByParam("reaction", uID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if err := app.posts.GetAllCategories(posts); err != nil {
		app.serverError(w, err)
		return
	}

	m := map[string][]*mod.Post{uID: posts}

	app.render(w, r, TMPL_FILTER_BY_REACTION, &templateData{
		FilteredPosts: m,
		UserID:        uID,
		Authorized:    authorized})
}

func (app *application) filterByCategory(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	var authorized bool
	var uID string
	if app.isAuthorized(w, r) {

		uID = getUserIDByCookie(r)

		authorized = true

		if err := app.refreshCookie(w, r); err != nil {
			app.serverError(w, err)
			return
		}
	}

	category := r.URL.Query().Get("name")
	if category == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	posts, err := app.posts.FilterByParam("category", category)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if err := app.posts.GetAllCategories(posts); err != nil {
		app.serverError(w, err)
		return
	}

	m := map[string][]*mod.Post{
		category: posts,
	}

	app.render(w, r, TMPL_FILTER_BY_CATEGORY, &templateData{
		FilteredPosts: m,
		UserID:        uID,
		Authorized:    authorized})
}
