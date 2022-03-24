package sqlite3

import (
	mod "forum/pkg/models"
	"strconv"
)

func (m *PostModel) GetCategory(id int) ([]string, error) {

	query := `
		SELECT category_name FROM category AS c
		LEFT JOIN post_category AS ps ON ps.category_id = c.category_id
		WHERE ps.post_id = ?
	`

	rows, err := m.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string

	for rows.Next() {

		var cat string

		if err := rows.Scan(&cat); err != nil {
			return nil, err
		}

		categories = append(categories, cat)

		if err := rows.Err(); err != nil {
			return nil, err
		}
	}
	return categories, nil
}

func (m *PostModel) GetAllCategories(posts []*mod.Post) error {
	for _, post := range posts {

		var err error
		post.Categories, err = m.GetCategory(post.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *PostModel) CreateCategory(postID int, categories []string) error {

	for _, cat := range categories {
		query := `
			INSERT INTO post_category (post_id, category_id)
			VALUES (?,?)
		`
		cat, err := strconv.Atoi(cat)
		if err != nil {
			return err
		}

		if _, err := m.DB.Exec(query, postID, cat); err != nil {
			return err
		}
	}
	return nil
}

func (m *PostModel) ShowAllCategories() ([]*mod.Category, error) {
	query := `SELECT * FROM category`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*mod.Category

	for rows.Next() {
		cat := &mod.Category{}

		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, err
		}

		categories = append(categories, cat)
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}
	return categories, nil
}
