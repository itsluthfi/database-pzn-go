package repository

import (
	"context"
	"database-go-pzn/entity"
	"database/sql"
	"errors"
	"strconv"
)

type commentRepositoryImpl struct {
	DB *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepositoryImpl{DB: db}
}

func (r *commentRepositoryImpl) Insert(ctx context.Context, comment entity.Comment) (entity.Comment, error) {
	script := "INSERT INTO comments(email, comment) VALUES (?, ?)"

	result, err := r.DB.ExecContext(ctx, script, comment.Email, comment.Comment)
	if err != nil {
		return comment, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return comment, err
	}

	comment.Id = int32(id)
	return comment, nil
}

func (r *commentRepositoryImpl) FindById(ctx context.Context, id int32) (entity.Comment, error) {
	script := "SELECT id, email, comment FROM comments WHERE id = ? LIMIT 1"

	rows, err := r.DB.QueryContext(ctx, script, id)
	comment := entity.Comment{}
	if err != nil {
		return comment, err
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&comment.Id, &comment.Email, &comment.Comment)
		return comment, nil
	} else {
		return comment, errors.New("Id " + strconv.Itoa(int(id)) + "not found")
	}
}

func (r *commentRepositoryImpl) FindAll(ctx context.Context) ([]entity.Comment, error) {
	script := "SELECT id, email, comment FROM comments"

	rows, err := r.DB.QueryContext(ctx, script)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []entity.Comment
	for rows.Next() {
		comment := entity.Comment{}
		rows.Scan(&comment.Id, &comment.Email, &comment.Comment)
		comments = append(comments, comment)
	}

	return comments, nil
}
