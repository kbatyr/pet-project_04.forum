package sqlite3

import (
	"database/sql"
	"fmt"
	"forum/pkg/models"
)

func (m *PostModel) GetRating(id int, table string) (*models.Rating, error) {

	r := &models.Rating{}

	tableID := "post_id"
	if table == "comments_rating" {
		tableID = "comment_id"
	}

	likesQuery := `SELECT COUNT(*) FROM %s WHERE %s=? AND reaction='like'`

	if err := m.DB.QueryRow(fmt.Sprintf(likesQuery, table, tableID), id).Scan(&r.Likes); err != nil {
		return nil, err
	}

	dislikesQuery := `
		SELECT COUNT(*) FROM %s WHERE %s=? AND reaction='dislike'
	`

	if err := m.DB.QueryRow(fmt.Sprintf(dislikesQuery, table, tableID), id).Scan(&r.Dislikes); err != nil {
		return nil, err
	}

	return r, nil
}

func (m *PostModel) AlreadyReacted(u *models.User, id int, reaction, table string) bool {

	tableID := "post_id"
	if table == "comments_rating" {
		tableID = "comment_id"
	}
	query := `
		SELECT * FROM %s WHERE user_id=? AND %s=? AND reaction=?
	`
	row := m.DB.QueryRow(fmt.Sprintf(query, table, tableID), u.ID, id, reaction)
	err := row.Scan()

	return err != sql.ErrNoRows
}

func (m *PostModel) UserReaction(u *models.User, id int, reaction, table string) error {
	query := ` 
		INSERT INTO %s
		VALUES (?,?,?)
		ON CONFLICT DO UPDATE
		SET reaction=?
	`

	_, err := m.DB.Exec(fmt.Sprintf(query, table), u.ID, id, reaction, reaction)
	if err != nil {
		return err
	}
	return nil
}

func (m *PostModel) RevokeReaction(u *models.User, id int, reaction, table string) error {

	tableID := "post_id"
	if table == "comments_rating" {
		tableID = "comment_id"
	}

	query := `
		DELETE FROM %s WHERE user_id=? AND %s=?
	`

	_, err := m.DB.Exec(fmt.Sprintf(query, table, tableID), u.ID, id)
	if err != nil {
		return err
	}
	return nil
}
