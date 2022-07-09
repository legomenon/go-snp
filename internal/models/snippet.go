package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Snippet struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
	Expires time.Time `json:"expires"`
}

type SnippetModel struct {
	DB *pgxpool.Pool
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	query := `
		INSERT INTO snippets (title, content, expires)
		VALUES ($1, $2, $3)
		Returning id;
	`
	id := 0
	err := m.DB.QueryRow(context.Background(), query, title, content, time.Now().Add(time.Hour*time.Duration(expires))).Scan(&id)
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
	return nil, nil
}
