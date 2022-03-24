package sqlite3

import (
	"database/sql"
	"errors"
	mod "forum/pkg/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Определяем тип который обертывает пул подключения sql.DB
type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) GetUserBySession(w http.ResponseWriter, r *http.Request) (*mod.User, error) {

	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, err
	}

	query := `
	SELECT u.user_id, u.username, u.email
	FROM user AS u
	INNER JOIN user_session AS us
	ON u.user_id = us.user_id
	WHERE uuid = ?
	`
	row := u.DB.QueryRow(query, cookie.Value)
	usr := &mod.User{}
	err = row.Scan(&usr.ID, &usr.UserName, &usr.Email)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, mod.ErrNoRecord
		}
		return nil, err
	}
	return usr, nil
}

func (u *UserModel) GetUserByName(w http.ResponseWriter, r *http.Request) (*mod.User, error) {

	userName := r.FormValue("username")

	query := `
	SELECT u.user_id, u.username, u.password
	FROM user AS u
	WHERE username = ?
	`

	row := u.DB.QueryRow(query, userName)
	usr := &mod.User{}
	err := row.Scan(&usr.ID, &usr.UserName, &usr.Password)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, mod.ErrNoRecord
		}
		return nil, err
	}
	return usr, nil
}

func (u *UserModel) CreateUserDB(usr *mod.User) error {

	query := `INSERT INTO user (username, email, password, registration_date) VALUES (?,?,?,datetime('now', 'localtime'))`

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usr.Password), 10)
	if err != nil {
		return err
	}

	_, err = u.DB.Exec(query, usr.UserName, usr.Email, hashedPassword)
	if err != nil {
		return err
	}

	// id, err := result.LastInsertId()
	// if err != nil {
	// 	return err
	// }

	return nil
}
