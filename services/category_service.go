package services

import (
	"database/sql"
	"my_blog/config"
	"my_blog/models"
)

// CategoryService 分类服务
type CategoryService struct{}

// GetAllCategories 获取所有分类
func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	rows, err := config.DB.Query("SELECT id, name, description FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// GetCategoryByID 根据ID获取分类
func (s *CategoryService) GetCategoryByID(id int) (*models.Category, error) {
	var category models.Category
	err := config.DB.QueryRow(`
		SELECT id, name, description
		FROM categories
		WHERE id = ?
	`, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &category, nil
}

// CreateCategory 创建分类
func (s *CategoryService) CreateCategory(category *models.Category) (int64, error) {
	result, err := config.DB.Exec(`
		INSERT INTO categories (name, description)
		VALUES (?, ?)
	`,
		category.Name,
		category.Description,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// UpdateCategory 更新分类
func (s *CategoryService) UpdateCategory(id int, category *models.Category) error {
	result, err := config.DB.Exec(`
		UPDATE categories
		SET name = ?, description = ?
		WHERE id = ?
	`,
		category.Name,
		category.Description,
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

// DeleteCategory 删除分类
func (s *CategoryService) DeleteCategory(id int) error {
	result, err := config.DB.Exec("DELETE FROM categories WHERE id = ?", id)
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
