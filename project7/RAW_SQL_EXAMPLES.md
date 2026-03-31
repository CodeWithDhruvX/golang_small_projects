# Raw SQL Queries with GORM in Project7

This document demonstrates various ways to write raw SQL queries using GORM in your project.

## Available Raw SQL Methods

### 1. Basic Raw SQL with Struct Mapping

**Function**: `GetUsersByAgeRange(minAge, maxAge int)`

```go
func (p *PostgresDB) GetUsersByAgeRange(minAge, maxAge int) ([]models.User, error) {
	var users []models.User
	query := "SELECT * FROM users WHERE age BETWEEN ? AND ? AND deleted_at IS NULL"
	err := p.DB.Raw(query, minAge, maxAge).Scan(&users).Error
	return users, err
}
```

**API Endpoint**: `GET /api/v1/sql/users/age-range?min_age=20&max_age=30`

### 2. Raw SQL for Specific Fields

**Function**: `GetUserEmails()`

```go
func (p *PostgresDB) GetUserEmails() ([]string, error) {
	var emails []string
	query := "SELECT email FROM users WHERE deleted_at IS NULL AND active = true"
	err := p.DB.Raw(query).Scan(&emails).Error
	return emails, err
}
```

**API Endpoint**: `GET /api/v1/sql/users/emails`

### 3. Raw SQL with JOINs

**Function**: `GetUsersWithPostCount()`

```go
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
```

**API Endpoint**: `GET /api/v1/sql/users/post-count`

### 4. Exec for INSERT/UPDATE/DELETE

**Function**: `DeactivateInactiveUsers()`

```go
func (p *PostgresDB) DeactivateInactiveUsers() (int64, error) {
	query := "UPDATE users SET active = false WHERE last_login < NOW() - INTERVAL '30 days'"
	result := p.DB.Exec(query)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
```

### 5. Raw SQL with Parameters

**Function**: `SearchPosts(searchTerm string)`

```go
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
```

**API Endpoint**: `GET /api/v1/posts/search?q=keyword`

### 6. Complex Aggregation Query

**Function**: `GetCategoryStats()`

```go
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
```

**API Endpoint**: `GET /api/v1/sql/categories/stats`

### 7. GORM SQL Builder (Safer Alternative)

**Function**: `GetActiveUsersWithFilters(ageMin, ageMax *int, emailDomain string)`

```go
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
```

**API Endpoint**: `GET /api/v1/sql/users/filters?min_age=20&max_age=30&email_domain=example.com`

## Usage Examples

### Test with PowerShell

```powershell
# Get all user emails
Invoke-RestMethod -Uri "http://localhost:8083/api/v1/sql/users/emails" -Method GET

# Get users with post counts
Invoke-RestMethod -Uri "http://localhost:8083/api/v1/sql/users/post-count" -Method GET

# Get users by age range
Invoke-RestMethod -Uri "http://localhost:8083/api/v1/sql/users/age-range?min_age=20&max_age=30" -Method GET

# Search posts
Invoke-RestMethod -Uri "http://localhost:8083/api/v1/posts/search?q=example" -Method GET

# Get category statistics
Invoke-RestMethod -Uri "http://localhost:8083/api/v1/sql/categories/stats" -Method GET

# Get users with filters
Invoke-RestMethod -Uri "http://localhost:8083/api/v1/sql/users/filters?min_age=25&email_domain=example.com" -Method GET
```

### Test with curl

```bash
# Get all user emails
curl "http://localhost:8083/api/v1/sql/users/emails"

# Get users with post counts
curl "http://localhost:8083/api/v1/sql/users/post-count"

# Get users by age range
curl "http://localhost:8083/api/v1/sql/users/age-range?min_age=20&max_age=30"

# Search posts
curl "http://localhost:8083/api/v1/posts/search?q=example"

# Get category statistics
curl "http://localhost:8083/api/v1/sql/categories/stats"

# Get users with filters
curl "http://localhost:8083/api/v1/sql/users/filters?min_age=25&email_domain=example.com"
```

## Best Practices

1. **Use Parameterized Queries**: Always use `?` placeholders to prevent SQL injection
2. **Handle Errors**: Always check for errors when executing raw SQL
3. **Use GORM Builder When Possible**: For complex dynamic queries, consider using GORM's query builder
4. **Test Queries**: Test your raw SQL queries in a database client first
5. **Document Complex Queries**: Add comments explaining complex SQL logic
6. **Consider Performance**: Use appropriate indexes for raw SQL queries

## GORM Raw SQL Methods

- `db.Raw()`: Execute raw SQL and scan results into structs
- `db.Exec()`: Execute raw SQL without returning results (INSERT/UPDATE/DELETE)
- `db.Model().Raw()`: Combine GORM model with raw SQL
- `db.Find().Raw()`: Use raw SQL within GORM query chains

These examples show that you have complete flexibility to write any SQL query while still using GORM as your ORM!
