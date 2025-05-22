package services

import (
	"database/sql"
	"my_blog/config"
	"my_blog/models"
	"time"
)

// ArticleService 文章服务
type ArticleService struct{}

// GetAllArticles 获取所有文章
func (s *ArticleService) GetAllArticles() ([]models.Article, error) {
	rows, err := config.DB.Query(`
		SELECT a.id, a.author, a.title, a.content, a.create_at, a.image_path, a.views,
			   c.id, c.name, c.description
		FROM articles a
		LEFT JOIN categories c ON a.category_id = c.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		err := rows.Scan(
			&article.ID,
			&article.Author,
			&article.Title,
			&article.Content,
			&article.CreateAt,
			&article.ImagePath,
			&article.Views,
			&article.Category.ID,
			&article.Category.Name,
			&article.Category.Description,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}

// GetArticleByID 根据ID获取文章
func (s *ArticleService) GetArticleByID(id int) (*models.Article, error) {
	var article models.Article
	err := config.DB.QueryRow(`
		SELECT a.id, a.author, a.title, a.content, a.create_at, a.image_path, a.views,
			   c.id, c.name, c.description
		FROM articles a
		LEFT JOIN categories c ON a.category_id = c.id
		WHERE a.id = ?
	`, id).Scan(
		&article.ID,
		&article.Author,
		&article.Title,
		&article.Content,
		&article.CreateAt,
		&article.ImagePath,
		&article.Views,
		&article.Category.ID,
		&article.Category.Name,
		&article.Category.Description,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// CreateArticle 创建文章
func (s *ArticleService) CreateArticle(article *models.Article, categoryName string) (int64, error) {
	var categoryID int
	err := config.DB.QueryRow("SELECT id FROM categories WHERE name = ?", categoryName).Scan(&categoryID)
	if err != nil {
		return 0, err
	}

	result, err := config.DB.Exec(`
		INSERT INTO articles (title, content, author, create_at, image_path, category_id, views)
		VALUES (?, ?, ?, ?, ?, ?, 0)
	`,
		article.Title,
		article.Content,
		article.Author,
		time.Now(),
		article.ImagePath,
		categoryID,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// UpdateArticle 更新文章
func (s *ArticleService) UpdateArticle(id int, article *models.Article) error {
	_, err := config.DB.Exec(`
		UPDATE articles 
		SET title = ?, content = ?, image_path = ?
		WHERE id = ?
	`,
		article.Title,
		article.Content,
		article.ImagePath,
		id,
	)
	return err
}

// DeleteArticle 删除文章
func (s *ArticleService) DeleteArticle(id int) error {
	result, err := config.DB.Exec("DELETE FROM articles WHERE id = ?", id)
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

// GetArticlesByCategory 获取分类下的所有文章
func (s *ArticleService) GetArticlesByCategory(categoryID int) ([]models.Article, error) {
	rows, err := config.DB.Query(`
		SELECT a.id, a.author, a.title, a.content, a.create_at, a.image_path, a.views,
			   c.id, c.name, c.description
		FROM articles a
		LEFT JOIN categories c ON a.category_id = c.id
		WHERE a.category_id = ?
	`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		err := rows.Scan(
			&article.ID,
			&article.Author,
			&article.Title,
			&article.Content,
			&article.CreateAt,
			&article.ImagePath,
			&article.Views,
			&article.Category.ID,
			&article.Category.Name,
			&article.Category.Description,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}
