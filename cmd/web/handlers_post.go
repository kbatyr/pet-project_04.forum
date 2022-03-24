package main

import (
	"errors"
	"fmt"
	mod "forum/pkg/models"
	"net/http"
	"strconv"
)

func (app *application) showPost(w http.ResponseWriter, r *http.Request) {

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

	// Извлекаем значение параметра id из URL
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	post, err := app.posts.GetPost(id)
	if err != nil {
		if errors.Is(err, mod.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	post.Categories, err = app.posts.GetCategory(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	post.Comments, err = app.posts.GetComments(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	post.Rating, err = app.posts.GetRating(id, "rating_info")
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, com := range post.Comments {
		com.Rating, err = app.posts.GetRating(com.ID, "comments_rating")
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	app.render(w, r, TMPL_SHOWPOST, &templateData{
		Post:       post,
		UserID:     uID,
		Authorized: authorized,
	})
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {

	categories, err := app.posts.ShowAllCategories()
	if err != nil {
		app.serverError(w, err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		uId := getUserIDByCookie(r)
		app.render(w, r, TMPL_CREATEPOST, &templateData{
			AllCategories: categories,
			UserID:        uId,
			Authorized:    true})
	case http.MethodPost:

		if err = r.ParseForm(); err != nil {
			app.serverError(w, err)
			return
		}

		if err := IsValidData(r.FormValue("title"), r.FormValue("content"), r.PostForm["category"]); err != nil {
			app.serviceMsg(w, r, &templateData{AllCategories: categories, Authorized: true}, TMPL_CREATEPOST, err, http.StatusBadRequest)
			return
		}

		d := &templateData{Post: &mod.Post{
			Title:   r.FormValue("title"),
			Content: r.FormValue("content"),
		}, Authorized: true}

		uID := getUserIDByCookie(r)

		postID, err := app.posts.CreatePost(d.Post.Title, d.Post.Content, uID)
		if err != nil {
			app.serverError(w, err)
			return
		}

		d.Post.Categories = r.PostForm["category"]

		if err := app.posts.CreateCategory(postID, d.Post.Categories); err != nil {
			app.serverError(w, err)
			return
		}

		// Перенаправляем пользователя на соответствующую страницу заметки.
		http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postID), http.StatusSeeOther)
	default:
		// Используем метод Header().Set() для добавления заголовка 'Allow: POST' в
		// карту HTTP-заголовков. Первый параметр - название заголовка, а
		// второй параметр - значение заголовка.
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}
