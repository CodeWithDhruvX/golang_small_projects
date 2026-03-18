package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Health check endpoint
func healthCheck(c *gin.Context) {
	requestID := c.GetString("request_id")
	
	logger.InfoWithRequestID("Health check accessed", map[string]interface{}{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
	}, requestID)

	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
		"version":   "1.0.0",
		"request_id": requestID,
	})
}

// Metrics endpoint
func metrics(c *gin.Context) {
	requestID := c.GetString("request_id")
	
	tasksMu.RLock()
	usersMu.RLock()
	projectsMu.RLock()
	defer tasksMu.RUnlock()
	defer usersMu.RUnlock()
	defer projectsMu.RUnlock()

	metrics := gin.H{
		"timestamp": time.Now(),
		"request_id": requestID,
		"entities": gin.H{
			"tasks":    len(tasks),
			"users":    len(users),
			"projects": len(projects),
		},
		"system": gin.H{
			"uptime": "N/A", // Could track actual uptime
			"memory": "N/A", // Could track memory usage
		},
	}

	logger.InfoWithRequestID("Metrics accessed", nil, requestID)
	c.JSON(http.StatusOK, metrics)
}

// Login endpoint
func login(c *gin.Context) {
	requestID := c.GetString("request_id")
	
	var loginReq struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		apiErr := errorHandler.NewAPIError(ErrorCodeBadRequest, "Invalid login request", requestID).
			WithDetails("validation_error", err.Error())
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	// Validate credentials (simplified)
	usersMu.RLock()
	user, exists := users[loginReq.Username]
	usersMu.RUnlock()

	if !exists || !authService.VerifyPassword(loginReq.Password, user.Password) {
		logger.LogSecurityEvent("Failed login attempt", map[string]interface{}{
			"username": loginReq.Username,
			"ip":       c.ClientIP(),
		}, requestID)

		apiErr := errorHandler.NewUnauthorizedError("Invalid credentials", requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	// Generate JWT token
	token, err := authService.GenerateJWTToken(user.ID, user.Username, user.Role)
	if err != nil {
		apiErr := errorHandler.HandleError(err, requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	// Update last login
	usersMu.Lock()
	user.LastLogin = time.Now().Format(time.RFC3339)
	users[user.Username] = user
	usersMu.Unlock()

	logger.InfoWithRequestID("User logged in successfully", map[string]interface{}{
		"username": loginReq.Username,
		"role":     user.Role,
	}, requestID)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
		"expires_in": 24 * 3600, // 24 hours in seconds
	})
}

// Register endpoint
func register(c *gin.Context) {
	requestID := c.GetString("request_id")
	
	var registerReq struct {
		Username string `json:"username" validate:"required,min=3,max=50"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&registerReq); err != nil {
		apiErr := errorHandler.NewAPIError(ErrorCodeBadRequest, "Invalid registration request", requestID).
			WithDetails("validation_error", err.Error())
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	// Check if user already exists
	usersMu.RLock()
	_, exists := users[registerReq.Username]
	usersMu.RUnlock()

	if exists {
		apiErr := errorHandler.NewAPIError(ErrorCodeConflict, "Username already exists", requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	// Hash password
	hashedPassword, err := authService.HashPassword(registerReq.Password)
	if err != nil {
		apiErr := errorHandler.HandleError(err, requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	// Create new user
	newUser := User{
		BaseModel: BaseModel{
			ID:        registerReq.Username,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username: registerReq.Username,
		Email:    registerReq.Email,
		Password: hashedPassword,
		Role:     "user",
		IsActive: true,
	}

	usersMu.Lock()
	users[registerReq.Username] = newUser
	usersMu.Unlock()

	logger.InfoWithRequestID("User registered successfully", map[string]interface{}{
		"username": registerReq.Username,
		"email":    registerReq.Email,
	}, requestID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":       newUser.ID,
			"username": newUser.Username,
			"email":    newUser.Email,
			"role":     newUser.Role,
		},
	})
}

// Task handlers
func getTasks(c *gin.Context) {
	requestID := c.GetString("request_id")
	
	tasksMu.RLock()
	defer tasksMu.RUnlock()

	taskList := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		taskList = append(taskList, task)
	}

	logger.InfoWithRequestID("Retrieved all tasks", map[string]interface{}{
		"count": len(taskList),
	}, requestID)

	c.JSON(http.StatusOK, gin.H{
		"tasks": taskList,
		"count": len(taskList),
	})
}

func getTask(c *gin.Context) {
	requestID := c.GetString("request_id")
	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		apiErr := errorHandler.NewAPIError(ErrorCodeBadRequest, "Invalid task ID", requestID).
			WithDetails("task_id", c.Param("id"))
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	tasksMu.RLock()
	task, exists := tasks[id]
	tasksMu.RUnlock()

	if !exists {
		apiErr := errorHandler.NewNotFoundError("Task", strconv.Itoa(id), requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	logger.InfoWithRequestID("Retrieved task", map[string]interface{}{
		"task_id": id,
	}, requestID)

	c.JSON(http.StatusOK, task)
}

func createTask(c *gin.Context) {
	requestID := c.GetString("request_id")
	
	// Get validated entity from context
	entity, exists := c.Get("validated_entity")
	if !exists {
		apiErr := errorHandler.NewAPIError(ErrorCodeInternal, "No validated entity found", requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	newTask, ok := entity.(*Task)
	if !ok {
		apiErr := errorHandler.NewAPIError(ErrorCodeInternal, "Invalid entity type", requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	tasksMu.Lock()
	defer tasksMu.Unlock()

	newTask.ID = nextTaskID
	newTask.CreatedAt = time.Now()
	newTask.UpdatedAt = time.Now()
	nextTaskID++

	tasks[newTask.ID] = *newTask

	logger.InfoWithRequestID("Task created", map[string]interface{}{
		"task_id": newTask.ID,
		"title":   newTask.Title,
	}, requestID)

	c.JSON(http.StatusCreated, newTask)
}

func updateTask(c *gin.Context) {
	requestID := c.GetString("request_id")
	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		apiErr := errorHandler.NewAPIError(ErrorCodeBadRequest, "Invalid task ID", requestID).
			WithDetails("task_id", c.Param("id"))
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	// Get validated entity from context
	entity, exists := c.Get("validated_entity")
	if !exists {
		apiErr := errorHandler.NewAPIError(ErrorCodeInternal, "No validated entity found", requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	updatedTask, ok := entity.(*Task)
	if !ok {
		apiErr := errorHandler.NewAPIError(ErrorCodeInternal, "Invalid entity type", requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	tasksMu.Lock()
	defer tasksMu.Unlock()

	if _, exists := tasks[id]; !exists {
		apiErr := errorHandler.NewNotFoundError("Task", strconv.Itoa(id), requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	updatedTask.ID = id
	updatedTask.UpdatedAt = time.Now()
	tasks[id] = *updatedTask

	logger.InfoWithRequestID("Task updated", map[string]interface{}{
		"task_id": id,
		"title":   updatedTask.Title,
	}, requestID)

	c.JSON(http.StatusOK, updatedTask)
}

func deleteTask(c *gin.Context) {
	requestID := c.GetString("request_id")
	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		apiErr := errorHandler.NewAPIError(ErrorCodeBadRequest, "Invalid task ID", requestID).
			WithDetails("task_id", c.Param("id"))
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	tasksMu.Lock()
	defer tasksMu.Unlock()

	if _, exists := tasks[id]; !exists {
		apiErr := errorHandler.NewNotFoundError("Task", strconv.Itoa(id), requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	delete(tasks, id)

	logger.InfoWithRequestID("Task deleted", map[string]interface{}{
		"task_id": id,
	}, requestID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
		"task_id": id,
	})
}

// User handlers (admin only)
func getUsers(c *gin.Context) {
	requestID := c.GetString("request_id")
	
	usersMu.RLock()
	defer usersMu.RUnlock()

	userList := make([]User, 0, len(users))
	for _, user := range users {
		// Don't include password in response
		user.Password = ""
		userList = append(userList, user)
	}

	logger.InfoWithRequestID("Retrieved all users", map[string]interface{}{
		"count": len(userList),
	}, requestID)

	c.JSON(http.StatusOK, gin.H{
		"users": userList,
		"count": len(userList),
	})
}

func getUser(c *gin.Context) {
	requestID := c.GetString("request_id")
	
	userID := c.Param("id")

	usersMu.RLock()
	user, exists := users[userID]
	usersMu.RUnlock()

	if !exists {
		apiErr := errorHandler.NewNotFoundError("User", userID, requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	// Don't include password in response
	user.Password = ""

	logger.InfoWithRequestID("Retrieved user", map[string]interface{}{
		"user_id": userID,
	}, requestID)

	c.JSON(http.StatusOK, user)
}

func createUser(c *gin.Context) {
	requestID := c.GetString("request_id")
	
	entity, exists := c.Get("validated_entity")
	if !exists {
		apiErr := errorHandler.NewAPIError(ErrorCodeInternal, "No validated entity found", requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	newUser, ok := entity.(*User)
	if !ok {
		apiErr := errorHandler.NewAPIError(ErrorCodeInternal, "Invalid entity type", requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	usersMu.Lock()
	defer usersMu.Unlock()

	// Check if user already exists
	if _, exists := users[newUser.Username]; exists {
		apiErr := errorHandler.NewAPIError(ErrorCodeConflict, "Username already exists", requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	// Hash password
	hashedPassword, err := authService.HashPassword(newUser.Password)
	if err != nil {
		apiErr := errorHandler.HandleError(err, requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	newUser.ID = newUser.Username
	newUser.Password = hashedPassword
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()

	users[newUser.Username] = *newUser

	// Don't include password in response
	newUser.Password = ""

	logger.InfoWithRequestID("User created", map[string]interface{}{
		"user_id": newUser.ID,
		"username": newUser.Username,
	}, requestID)

	c.JSON(http.StatusCreated, newUser)
}

func updateUser(c *gin.Context) {
	requestID := c.GetString("request_id")
	
	userID := c.Param("id")

	entity, exists := c.Get("validated_entity")
	if !exists {
		apiErr := errorHandler.NewAPIError(ErrorCodeInternal, "No validated entity found", requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	updatedUser, ok := entity.(*User)
	if !ok {
		apiErr := errorHandler.NewAPIError(ErrorCodeInternal, "Invalid entity type", requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	usersMu.Lock()
	defer usersMu.Unlock()

	if _, exists := users[userID]; !exists {
		apiErr := errorHandler.NewNotFoundError("User", userID, requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	// Hash password if provided
	if updatedUser.Password != "" {
		hashedPassword, err := authService.HashPassword(updatedUser.Password)
		if err != nil {
			apiErr := errorHandler.HandleError(err, requestID)
			c.JSON(apiErr.HTTPStatus(), apiErr)
			return
		}
		updatedUser.Password = hashedPassword
	} else {
		// Keep existing password
		updatedUser.Password = users[userID].Password
	}

	updatedUser.ID = userID
	updatedUser.Username = userID
	updatedUser.UpdatedAt = time.Now()

	users[userID] = *updatedUser

	// Don't include password in response
	updatedUser.Password = ""

	logger.InfoWithRequestID("User updated", map[string]interface{}{
		"user_id": userID,
	}, requestID)

	c.JSON(http.StatusOK, updatedUser)
}

func deleteUser(c *gin.Context) {
	requestID := c.GetString("request_id")
	
	userID := c.Param("id")

	usersMu.Lock()
	defer usersMu.Unlock()

	if _, exists := users[userID]; !exists {
		apiErr := errorHandler.NewNotFoundError("User", userID, requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	delete(users, userID)

	logger.InfoWithRequestID("User deleted", map[string]interface{}{
		"user_id": userID,
	}, requestID)

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
		"user_id": userID,
	})
}

// Project handlers
func getProjects(c *gin.Context) {
	requestID := c.GetString("request_id")
	
	projectsMu.RLock()
	defer projectsMu.RUnlock()

	projectList := make([]Project, 0, len(projects))
	for _, project := range projects {
		projectList = append(projectList, project)
	}

	logger.InfoWithRequestID("Retrieved all projects", map[string]interface{}{
		"count": len(projectList),
	}, requestID)

	c.JSON(http.StatusOK, gin.H{
		"projects": projectList,
		"count":    len(projectList),
	})
}

func getProject(c *gin.Context) {
	requestID := c.GetString("request_id")
	
	projectID := c.Param("id")

	projectsMu.RLock()
	project, exists := projects[projectID]
	projectsMu.RUnlock()

	if !exists {
		apiErr := errorHandler.NewNotFoundError("Project", projectID, requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	logger.InfoWithRequestID("Retrieved project", map[string]interface{}{
		"project_id": projectID,
	}, requestID)

	c.JSON(http.StatusOK, project)
}

func createProject(c *gin.Context) {
	requestID := c.GetString("request_id")
	
	entity, exists := c.Get("validated_entity")
	if !exists {
		apiErr := errorHandler.NewAPIError(ErrorCodeInternal, "No validated entity found", requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	newProject, ok := entity.(*Project)
	if !ok {
		apiErr := errorHandler.NewAPIError(ErrorCodeInternal, "Invalid entity type", requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	projectsMu.Lock()
	defer projectsMu.Unlock()

	// Generate unique project ID
	projectID := fmt.Sprintf("proj_%d", time.Now().UnixNano())
	
	newProject.ID = projectID
	newProject.CreatedAt = time.Now()
	newProject.UpdatedAt = time.Now()

	projects[projectID] = *newProject

	logger.InfoWithRequestID("Project created", map[string]interface{}{
		"project_id": projectID,
		"name":       newProject.Name,
	}, requestID)

	c.JSON(http.StatusCreated, newProject)
}

func updateProject(c *gin.Context) {
	requestID := c.GetString("request_id")
	
	projectID := c.Param("id")

	entity, exists := c.Get("validated_entity")
	if !exists {
		apiErr := errorHandler.NewAPIError(ErrorCodeInternal, "No validated entity found", requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	updatedProject, ok := entity.(*Project)
	if !ok {
		apiErr := errorHandler.NewAPIError(ErrorCodeInternal, "Invalid entity type", requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	projectsMu.Lock()
	defer projectsMu.Unlock()

	if _, exists := projects[projectID]; !exists {
		apiErr := errorHandler.NewNotFoundError("Project", projectID, requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	updatedProject.ID = projectID
	updatedProject.UpdatedAt = time.Now()

	projects[projectID] = *updatedProject

	logger.InfoWithRequestID("Project updated", map[string]interface{}{
		"project_id": projectID,
		"name":       updatedProject.Name,
	}, requestID)

	c.JSON(http.StatusOK, updatedProject)
}

func deleteProject(c *gin.Context) {
	requestID := c.GetString("request_id")
	
	projectID := c.Param("id")

	projectsMu.Lock()
	defer projectsMu.Unlock()

	if _, exists := projects[projectID]; !exists {
		apiErr := errorHandler.NewNotFoundError("Project", projectID, requestID)
		c.JSON(apiErr.HTTPStatus(), apiErr)
		return
	}

	delete(projects, projectID)

	logger.InfoWithRequestID("Project deleted", map[string]interface{}{
		"project_id": projectID,
	}, requestID)

	c.JSON(http.StatusOK, gin.H{
		"message":    "Project deleted successfully",
		"project_id": projectID,
	})
}

// RequireRoleMiddleware checks if user has required role
func RequireRoleMiddleware(requiredRole string, errorHandler *ErrorHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		
		userClaims, exists := c.Get("user_claims")
		if !exists {
			apiErr := errorHandler.NewUnauthorizedError("User not authenticated", requestID)
			c.JSON(apiErr.HTTPStatus(), apiErr)
			c.Abort()
			return
		}

		claims, ok := userClaims.(*JWTClaims)
		if !ok {
			apiErr := errorHandler.NewUnauthorizedError("Invalid user claims", requestID)
			c.JSON(apiErr.HTTPStatus(), apiErr)
			c.Abort()
			return
		}

		if err := authService.Authorize(claims, requiredRole, ""); err != nil {
			apiErr := errorHandler.HandleError(err, requestID)
			c.JSON(apiErr.HTTPStatus(), apiErr)
			c.Abort()
			return
		}

		c.Next()
	}
}
