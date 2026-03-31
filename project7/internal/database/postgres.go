package database

import (
	"fmt"
	"log"
	"project7/internal/config"
	"project7/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresDB struct {
	DB *gorm.DB
}

func NewPostgresDB(cfg *config.AppConfig) (*PostgresDB, error) {
	dsn := cfg.GetPostgresDSN()
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")

	// Auto migrate the schema
	err = db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Tag{},
		&models.Category{},
		&models.Product{},
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to migrate PostgreSQL schema: %w", err)
	}

	log.Println("PostgreSQL schema migration completed")

	return &PostgresDB{DB: db}, nil
}

func (p *PostgresDB) Close() error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Helper methods for common operations
func (p *PostgresDB) CreateUser(user *models.User) error {
	return p.DB.Create(user).Error
}

func (p *PostgresDB) GetUser(id uint) (*models.User, error) {
	var user models.User
	err := p.DB.Preload("Posts").Preload("Posts.Tags").First(&user, id).Error
	return &user, err
}

func (p *PostgresDB) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := p.DB.Find(&users).Error
	return users, err
}

func (p *PostgresDB) UpdateUser(user *models.User) error {
	return p.DB.Save(user).Error
}

func (p *PostgresDB) DeleteUser(id uint) error {
	return p.DB.Delete(&models.User{}, id).Error
}

func (p *PostgresDB) CreatePost(post *models.Post) error {
	return p.DB.Create(post).Error
}

func (p *PostgresDB) GetPost(id uint) (*models.Post, error) {
	var post models.Post
	err := p.DB.Preload("User").Preload("Tags").First(&post, id).Error
	return &post, err
}

func (p *PostgresDB) GetAllPosts() ([]models.Post, error) {
	var posts []models.Post
	err := p.DB.Preload("User").Preload("Tags").Find(&posts).Error
	return posts, err
}

// Raw SQL Query Examples

// Method 1: Raw SQL with struct mapping
func (p *PostgresDB) GetUsersByAgeRange(minAge, maxAge int) ([]models.User, error) {
	var users []models.User
	query := "SELECT * FROM users WHERE age BETWEEN ? AND ? AND deleted_at IS NULL"
	err := p.DB.Raw(query, minAge, maxAge).Scan(&users).Error
	return users, err
}

// Method 2: Raw SQL for specific fields
func (p *PostgresDB) GetUserEmails() ([]string, error) {
	var emails []string
	query := "SELECT email FROM users WHERE deleted_at IS NULL AND active = true"
	err := p.DB.Raw(query).Scan(&emails).Error
	return emails, err
}

// Method 3: Raw SQL with JOIN
func (p *PostgresDB) GetUsersWithPostCount() ([]struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	PostCount int    `json:"post_count"`
}, error) {
	var results []struct {
		ID        uint   `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		PostCount int    `json:"post_count"`
	}
	
	query := `
		SELECT u.id, u.username, u.email, COUNT(p.id) as post_count 
		FROM users u 
		LEFT JOIN posts p ON u.id = p.user_id 
		WHERE u.deleted_at IS NULL 
		GROUP BY u.id, u.username, u.email 
		ORDER BY post_count DESC`
	
	err := p.DB.Raw(query).Scan(&results).Error
	return results, err
}

// Method 4: Exec for INSERT/UPDATE/DELETE
func (p *PostgresDB) DeactivateInactiveUsers() (int64, error) {
	query := "UPDATE users SET active = false WHERE last_login < NOW() - INTERVAL '30 days'"
	result := p.DB.Exec(query)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// Method 5: Raw SQL with parameters
func (p *PostgresDB) SearchPosts(searchTerm string) ([]models.Post, error) {
	var posts []models.Post
	query := `
		SELECT p.* FROM posts p 
		INNER JOIN users u ON p.user_id = u.id 
		WHERE (p.title ILIKE ? OR p.content ILIKE ?) 
		AND p.deleted_at IS NULL 
		AND u.deleted_at IS NULL`
	
	searchPattern := "%" + searchTerm + "%"
	err := p.DB.Raw(query, searchPattern, searchPattern).Scan(&posts).Error
	return posts, err
}

// Method 6: Complex aggregation query
func (p *PostgresDB) GetCategoryStats() ([]struct {
	CategoryName string  `json:"category_name"`
	ProductCount int     `json:"product_count"`
	AvgPrice     float64 `json:"avg_price"`
	TotalValue   float64 `json:"total_value"`
}, error) {
	var results []struct {
		CategoryName string  `json:"category_name"`
		ProductCount int     `json:"product_count"`
		AvgPrice     float64 `json:"avg_price"`
		TotalValue   float64 `json:"total_value"`
	}
	
	query := `
		SELECT 
			c.name as category_name,
			COUNT(p.id) as product_count,
			COALESCE(AVG(p.price), 0) as avg_price,
			COALESCE(SUM(p.price * p.stock), 0) as total_value
		FROM categories c 
		LEFT JOIN products p ON c.id = p.category_id AND p.deleted_at IS NULL
		WHERE c.deleted_at IS NULL 
		GROUP BY c.id, c.name 
		ORDER BY total_value DESC`
	
	err := p.DB.Raw(query).Scan(&results).Error
	return results, err
}

// Method 7: Using GORM's SQL Builder (safer alternative)
func (p *PostgresDB) GetActiveUsersWithFilters(ageMin, ageMax *int, emailDomain string) ([]models.User, error) {
	var users []models.User
	query := p.DB.Model(&models.User{}).Where("active = ? AND deleted_at IS NULL", true)
	
	if ageMin != nil {
		query = query.Where("age >= ?", *ageMin)
	}
	if ageMax != nil {
		query = query.Where("age <= ?", *ageMax)
	}
	if emailDomain != "" {
		query = query.Where("email LIKE ?", "%@"+emailDomain)
	}
	
	err := query.Find(&users).Error
	return users, err
}
