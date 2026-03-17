package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var (
	tasks   = make(map[int]Task)
	nextID  = 1
	tasksMu sync.RWMutex
)

func main() {
	r := gin.Default()

	// Initialize with some sample data
	initializeTasks()

	// Routes
	r.GET("/tasks", getTasks)
	r.GET("/tasks/:id", getTask)
	r.POST("/tasks", createTask)
	r.PUT("/tasks/:id", updateTask)
	r.DELETE("/tasks/:id", deleteTask)

	r.Run(":8080")
}

func initializeTasks() {
	tasksMu.Lock()
	defer tasksMu.Unlock()
	
	tasks[1] = Task{ID: 1, Title: "Learn Go", Completed: false}
	tasks[2] = Task{ID: 2, Title: "Build REST API", Completed: true}
	nextID = 3
}

func getTasks(c *gin.Context) {
	tasksMu.RLock()
	defer tasksMu.RUnlock()

	taskList := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		taskList = append(taskList, task)
	}

	c.JSON(http.StatusOK, taskList)
}

func getTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	tasksMu.RLock()
	defer tasksMu.RUnlock()

	task, exists := tasks[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func createTask(c *gin.Context) {
	var newTask Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasksMu.Lock()
	defer tasksMu.Unlock()

	newTask.ID = nextID
	nextID++
	tasks[newTask.ID] = newTask

	c.JSON(http.StatusCreated, newTask)
}

func updateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var updatedTask Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasksMu.Lock()
	defer tasksMu.Unlock()

	if _, exists := tasks[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	updatedTask.ID = id
	tasks[id] = updatedTask

	c.JSON(http.StatusOK, updatedTask)
}

func deleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	tasksMu.Lock()
	defer tasksMu.Unlock()

	if _, exists := tasks[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	delete(tasks, id)
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
