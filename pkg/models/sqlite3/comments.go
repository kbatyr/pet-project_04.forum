package sqlite3

import (
	"forum/pkg/models"
	"time"
)

func (m *PostModel) CreateComment(postID int, userID, msg string) error {
	query := `
		INSERT INTO comment (msg, creation_date, user_id, post_id)
		VALUES (?, datetime('now', 'localtime'), ?, ?)
	`
	_, err := m.DB.Exec(query, msg, userID, postID)
	if err != nil {
		return err
	}
	return nil
}

func (m *PostModel) GetComments(postID int) ([]*models.Comment, error) {

	query := `
		SELECT c.comment_id, c.msg, c.creation_date, u.username
		FROM comment AS c
		INNER JOIN user AS u ON c.user_id = u.user_id
		WHERE post_id = ?
	`

	rows, err := m.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment

	for rows.Next() {

		var time time.Time
		c := &models.Comment{}

		if err := rows.Scan(&c.ID, &c.Message, &time, &c.Author); err != nil {
			return nil, err
		}

		c.Date = time.Format("2006-01-02 15:04")

		comments = append(comments, c)

		if err := rows.Err(); err != nil {
			return nil, err
		}
	}
	return comments, nil
}
