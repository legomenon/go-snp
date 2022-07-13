package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Snippet struct {
	ID      int       `json:"id" `
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Created time.Time `json:"created" db:"created"`
	Expires time.Time `json:"expires"`
}

type SnippetModel struct {
	DB *pgxpool.Pool
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	query := `
		INSERT INTO snippets (title, content, expires)
		VALUES ($1, $2, NOW() + INTERVAL '1 day' * $3)
		RETURNING id;
	`
	id := 0
	err := m.DB.QueryRow(context.Background(), query, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	query := `
		SELECT title,content,date_created,expires FROM snippets
		WHERE expires > date_created AND id = $1;
	`
	data := &Snippet{}
	err := m.DB.QueryRow(context.Background(), query, id).Scan(&data.Title, &data.Content, &data.Created, &data.Expires)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return data, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	query := `
		SELECT id,title,content,date_created,expires FROM snippets
		WHERE expires > date_created ORDER BY id DESC LIMIT 10;
	`
	rows, err := m.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	snippets := []*Snippet{}
	for rows.Next() {
		s := &Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
