package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {

	mux := http.NewServeMux()

	// Инициализируем FileServer, он будет обрабатывать
	// HTTP-запросы к статическим файлам из папки "./ui/static".
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static", http.NotFoundHandler())

	// Используем функцию mux.Handle() для регистрации обработчика для
	// всех запросов, которые начинаются с "/static/". Мы убираем
	// префикс "/static" перед тем как запрос достигнет http.FileServer
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/signup", app.alreadyLoggedIn(app.signupUser))
	mux.HandleFunc("/login", app.alreadyLoggedIn(app.loginUser))
	mux.HandleFunc("/logout", app.logoutUser)

	mux.HandleFunc("/post", app.showPost)
	mux.HandleFunc("/create", app.hasPersmission(app.createPost))
	mux.HandleFunc("/comment", app.hasPersmission(app.leaveComment))

	mux.HandleFunc("/rating", app.hasPersmission(app.postRating))
	mux.HandleFunc("/rating/comment", app.hasPersmission(app.postRating))

	mux.HandleFunc("/category", app.filterByCategory)
	mux.HandleFunc("/user_posts", app.hasPersmission(app.filterByUser))
	mux.HandleFunc("/liked_posts", app.hasPersmission(app.filterByReaction))

	return mux
}
