package sqlite3

import (
	"fmt"
	mod "forum/pkg/models"
	"log"
	"net/http"
	"time"
)

func (u *UserModel) CreateSession(usr *mod.User, c *http.Cookie) error {
	query := `
		INSERT INTO user_session (uuid, expires, user_id)
		VALUES (?, ?, ?)
	`
	fmt.Println(time.Now().Add(time.Second * time.Duration(c.MaxAge)))
	_, err := u.DB.Exec(query, c.Value, time.Now().Add(time.Second*time.Duration(c.MaxAge)), usr.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserModel) UpdateSession(c *http.Cookie) error {

	query := `
		UPDATE user_session
		SET expires= ?
		WHERE uuid= ?;
	`
	_, err := u.DB.Exec(query, time.Now().Add(time.Second*time.Duration(c.MaxAge)), c.Value)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserModel) DeleteSession(id int) error {
	query := `DELETE FROM user_session WHERE user_id= ?`

	_, err := u.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserModel) CleanExpiredSessions() {

	for {
		time.Sleep(time.Minute * 10)

		query := `DELETE FROM user_session WHERE expires < ?`

		_, err := u.DB.Exec(query, time.Now())
		if err != nil {
			log.Fatal(err)
		}
	}
}
