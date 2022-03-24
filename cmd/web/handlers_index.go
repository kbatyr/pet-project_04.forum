package main

import (
	"net/http"
)

const (
	TMPL_HOME               = "home.page.html"
	TMPL_SIGNUP             = "signup.page.html"
	TMPL_LOGIN              = "login.page.html"
	TMPL_SHOWPOST           = "show.page.html"
	TMPL_CREATEPOST         = "create.page.html"
	TMPL_FILTER_BY_CATEGORY = "filter.category.page.html"
	TMPL_FILTER_BY_USER     = "filter.user.page.html"
	TMPL_FILTER_BY_REACTION = "filter.reaction.page.html"

	URL_HOME        = "/"
	URL_SIGNUP      = "/signup"
	URL_LOGIN       = "/login"
	URL_POST        = "/post"
	URL_CREATE_POST = "/create"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

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

	if r.URL.Path != URL_HOME {
		app.notFound(w)
		return
	}

	posts, err := app.posts.GetAllPosts()
	if err != nil {
		app.serverError(w, err)
		return
	}

	user, _ := app.user.GetUserBySession(w, r)

	app.posts.GetAllCategories(posts)

	app.render(w, r, TMPL_HOME, &templateData{
		AllPosts:   posts,
		User:       user,
		UserID:     uID,
		Authorized: authorized,
	})
}
