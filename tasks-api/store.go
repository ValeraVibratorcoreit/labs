package main

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("task not found")

type TaskStore struct {
	db *pgxpool.Pool
}

func (s *TaskStore) Create(ctx context.Context, title string) (Task, error) {
	var t Task
	err := s.db.QueryRow(ctx, `
		INSERT INTO tasks (title)
		VALUES ($1)
		RETURNING id, title, done, created_at`, title).
		Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt)
	return t, err
}

func (s *TaskStore) List(ctx context.Context) ([]Task, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, title, done, created_at FROM tasks ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (s *TaskStore) Get(ctx context.Context, id int64) (Task, error) {
	var t Task
	err := s.db.QueryRow(ctx,
		`SELECT id, title, done, created_at FROM tasks WHERE id = $1`, id).
		Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return Task{}, ErrNotFound
	}
	if err != nil {
		return Task{}, err
	}
	return t, nil
}

func (s *TaskStore) Update(ctx context.Context, id int64, title *string, done *bool) (Task, error) {
	var t Task
	err := s.db.QueryRow(ctx, `
		UPDATE tasks
		SET title = COALESCE($2, title),
		    done = COALESCE($3, done)
		WHERE id = $1
		RETURNING id, title, done, created_at`, id, title, done).
		Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return Task{}, ErrNotFound
	}
	if err != nil {
		return Task{}, err
	}
	return t, nil
}

func (s *TaskStore) Delete(ctx context.Context, id int64) error {
	tag, err := s.db.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
