package sqlite3

import (
	"fmt"
	mod "forum/pkg/models"
	"time"
)


func (m *PostModel) FilterByParam(param string, id string) ([]*mod.Post, error) {

	query := `
		SELECT p.post_id, p.title, p.content, p.creation_date
		,u.username 
	`

	switch param {
	case "category":
		query = fmt.Sprint(query, `
			FROM post p
			JOIN user AS u ON p.user_id = u.user_id
			JOIN category AS c ON ps.category_id = c.category_id
			JOIN post_category AS ps ON ps.post_id=p.post_id
			WHERE c.category_name= ?
		`)
	case "user":
		query = fmt.Sprint(query, `
			FROM post p
			JOIN user AS u ON p.user_id=u.user_id
			WHERE u.user_id = ?
		`)
	case "reaction":
		query = fmt.Sprint(query, `
			FROM rating_info r
			JOIN user AS u ON r.user_id = u.user_id
			JOIN post AS p ON r.post_id = p.post_id
			WHERE u.user_id = ? AND r.reaction = 'like'
		`)
	}

	rows, err := m.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*mod.Post
	for rows.Next() {
		var time time.Time
		p := &mod.Post{}

		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &time, &p.Author); err != nil {
			return nil, err
		}

		p.CreationDate = time.Format("2006-01-02 15:04")

		posts = append(posts, p)

		if err := rows.Err(); err != nil {
			return nil, err
		}
	}
	return posts, nil
}
