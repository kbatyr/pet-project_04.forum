package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/mattn/go-sqlite3"
)

// Помощник serverError записывает сообщение об ошибке в errorLog и
// затем отправляет пользователю ответ 500 "Внутренняя ошибка сервера".
func (app *application) serverError(w http.ResponseWriter, err error) {

	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	// Дает увидеть название проблематичного файла и номер строки из трассировки стека,
	// вместо указания строки логера который данную ошибку выводит!
	app.errLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Помощник clientError отправляет определенный код состояния и соответствующее описание
// пользователю.
func (app *application) serviceMsg(w http.ResponseWriter, r *http.Request, tmplData *templateData, tmpl string, err error, status int) {

	w.WriteHeader(status)
	tmplData.ServiceMsg = err.Error()
	app.render(w, r, tmpl, tmplData)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Мы также реализуем помощник notFound. Это просто
// удобная оболочка вокруг clientError, которая отправляет пользователю ответ "404 Страница не найдена".
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {

	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("template %s doesn't exist", name))
		return
	}
	var b bytes.Buffer
	if err := ts.Execute(&b, td); err != nil {
		app.serverError(w, err)
		return
	}
	b.WriteTo(w)
}

func (app *application) isConstraintErr(err error) bool {
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		if errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) {
			return true
		}
	}
	return false
}
