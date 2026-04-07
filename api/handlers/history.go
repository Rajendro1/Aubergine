package handlers

import (
	"net/http"

	"aubergine/internal/database"
	"aubergine/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateProgressRequest struct {
	ContentID       uint `json:"content_id" binding:"required"`
	ProgressSeconds int  `json:"progress_seconds"`
	IsCompleted     bool `json:"is_completed"`
}

func UpdateProgress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find or create watch history
	var history models.WatchHistory
	if err := database.DB.Where("user_id = ? AND content_id = ?", userID, req.ContentID).First(&history).Error; err != nil {
		// Create new
		history = models.WatchHistory{
			UserID:          uint(userID.(uint)),
			ContentID:       req.ContentID,
			ProgressSeconds: req.ProgressSeconds,
			IsCompleted:     req.IsCompleted,
			LastWatchedAt:   time.Now(),
		}
		if err := database.DB.Create(&history).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update progress"})
			return
		}
	} else {
		// Update existing
		history.ProgressSeconds = req.ProgressSeconds
		history.IsCompleted = req.IsCompleted
		history.LastWatchedAt = time.Now()
		if err := database.DB.Save(&history).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update progress"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "progress updated", "history": history})
}

func GetContinueWatching(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var histories []models.WatchHistory
	if err := database.DB.Where("user_id = ? AND is_completed = false", userID).Preload("Content").Order("last_watched_at DESC").Find(&histories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch continue watching"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"continue_watching": histories})
}
