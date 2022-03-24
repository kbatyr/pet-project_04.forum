package main

import (
	"errors"
	"flag"
	sq "forum/pkg/models/sqlite3"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/mattn/go-sqlite3"
)

type application struct {
	errLog        *log.Logger
	infoLog       *log.Logger
	posts         *sq.PostModel
	user          *sq.UserModel
	templateCache map[string]*template.Template
}

func main() {

	// Создаем новый флаг командной строки, значение по умолчанию: ":8080".
	// Добавляем небольшую справку, объясняющая, что содержит данный флаг.
	// Значение флага будет сохранено в переменной addr.
	addr := flag.String("addr", ":8080", "HTTP network address")

	// Мы вызываем функцию flag.Parse() для извлечения флага из командной строки
	// и присовения полученного значения переменной addr
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB()
	if err != nil {
		errLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errLog.Fatal(err)
	}

	// Инициализируем новую структуру с зависимостями приложения.
	// Чтобы логи были доступны за пределами функции main
	app := &application{
		errLog:        errLog,
		infoLog:       infoLog,
		posts:         &sq.PostModel{DB: db},
		user:          &sq.UserModel{DB: db},
		templateCache: templateCache,
	}

	// Создаем таблицы для БД со связями между собой
	if err := sq.LoadSqlFiles("./pkg/models/sqlite3/schema/up/", db); err != nil {
		// Check whether err is Constraint unique type
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if !errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) {
				errLog.Println(err)
				return
			}
		}
	}

	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errLog,
	}

	go app.user.CleanExpiredSessions()

	infoLog.Printf("listening on the port %s", *addr)
	err = srv.ListenAndServe()
	errLog.Fatal(err)
}