package sqlite3

import (
	"database/sql"
	"errors"
	mod "forum/pkg/models"
	"time"
)

// Определяем тип который обертывает пул подключения sql.DB
type PostModel struct {
	DB *sql.DB
}

// Метод для возвращения данных заметки по её идентификатору ID.
func (m *PostModel) GetPost(id int) (*mod.Post, error) {
	query := `SELECT p.post_id, p.title, p.content, p.creation_date, u.username
	FROM post AS p
	INNER JOIN user AS u
	ON p.user_id = u.user_id
	WHERE post_id = ?;
	`

	row := m.DB.QueryRow(query, id)

	p := &mod.Post{}

	var time time.Time

	err := row.Scan(&p.ID, &p.Title, &p.Content, &time, &p.Author)
	if err != nil {
		// Если в БД отсутствует строка по запросу
		if errors.Is(err, sql.ErrNoRows) {
			return nil, mod.ErrNoRecord
		}
		return nil, err
	}
	p.CreationDate = time.Format("2006-01-02 15:04")

	return p, nil
}

// Метод возвращает все посты созданные за всё время.
func (m *PostModel) GetAllPosts() ([]*mod.Post, error) {

	query := `SELECT p.post_id, p.title, u.username, p.content, p.creation_date
	FROM post AS p
	INNER JOIN user AS u
	ON p.user_id = u.user_id`

	// В ответ получим sql.Rows, который содержит результат нашего запроса.
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Инициализируем пустой срез для хранения всех постов.
	var posts []*mod.Post

	// Пробегаемся по каждой строке полученного результата
	for rows.Next() {

		// Создаем указатель на новую структуру models.Post
		var time time.Time
		p := &mod.Post{}

		// Используем rows.Scan(), чтобы скопировать значения полей в структуру.
		if err := rows.Scan(&p.ID, &p.Title, &p.Author, &p.Content, &time); err != nil {
			return nil, err
		}

		p.CreationDate = time.Format("2006-01-02 15:04")
		// Добавляем структуру в срез.
		posts = append(posts, p)

		// Когда цикл rows.Next() завершается, вызываем метод rows.Err(), чтобы узнать
		// если в ходе работы у нас не возникла какая либо ошибка.
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}
	return posts, nil
}

// Метод для создания новой заметки в базе дынных.
func (m *PostModel) CreatePost(title, content string, user_id string) (int, error) {

	query := `INSERT INTO post (title,content,creation_date, user_id)
	VALUES(?,?,datetime('now', 'localtime'),?)`

	result, err := m.DB.Exec(query, title, content, user_id)
	if err != nil {
		return 0, err
	}

	// Используем метод LastInsertId(), чтобы получить последний ID
	// созданной записи из таблицу snippets.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
