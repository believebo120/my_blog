package services

import (
	"database/sql"
	"my_blog/config"
	"my_blog/models"
	"time"
)

// CommentService 评论服务
type CommentService struct{}

// GetCommentsByArticle 获取文章的所有评论
func (s *CommentService) GetCommentsByArticle(articleID int) ([]models.Comment, error) {
	rows, err := config.DB.Query(`
		SELECT id, article_id, content, author, create_at
		FROM comments
		WHERE article_id = ?
		ORDER BY create_at DESC
	`, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.ArticleID,
			&comment.Content,
			&comment.Author,
			&comment.CreateAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

// CreateComment 创建评论
func (s *CommentService) CreateComment(comment *models.Comment) (int64, error) {
	result, err := config.DB.Exec(`
		INSERT INTO comments (article_id, content, author, create_at)
		VALUES (?, ?, ?, ?)
	`,
		comment.ArticleID,
		comment.Content,
		comment.Author,
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetCommentByID 根据ID获取评论
func (s *CommentService) GetCommentByID(id int) (*models.Comment, error) {
	var comment models.Comment
	err := config.DB.QueryRow(`
		SELECT id, article_id, content, author, create_at
		FROM comments
		WHERE id = ?
	`, id).Scan(
		&comment.ID,
		&comment.ArticleID,
		&comment.Content,
		&comment.Author,
		&comment.CreateAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// UpdateComment 更新评论
func (s *CommentService) UpdateComment(id int, content string) error {
	result, err := config.DB.Exec(`
		UPDATE comments
		SET content = ?, create_at = ?
		WHERE id = ?
	`,
		content,
		time.Now(),
		id,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteComment 删除评论
func (s *CommentService) DeleteComment(id int) error {
	result, err := config.DB.Exec("DELETE FROM comments WHERE id = ?", id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
