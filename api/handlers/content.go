package handlers

import (
	"net/http"
	"strconv"

	"aubergine/internal/database"
	"aubergine/internal/models"

	"github.com/gin-gonic/gin"
)

// ListContent retrieves all content available, potentially filtering out based on tier client-side or server-side.
func ListContent(c *gin.Context) {
	var contents []models.Content

	// Optional: Get user plan if authenticated, and filter. Currently returning all for display.
	if err := database.DB.Find(&contents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch content"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"contents": contents})
}

// StreamContent acts as a mock signed URL generator or direct manifest server
func StreamContent(c *gin.Context) {
	contentIDStr := c.Param("id")
	contentID, err := strconv.Atoi(contentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content id"})
		return
	}

	var content models.Content
	if err := database.DB.First(&content, contentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "content not found"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Check access level
	var subscription models.UserSubscription
	if err := database.DB.Where("user_id = ? AND is_active = true", userID).Preload("Plan").First(&subscription).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "no active subscription"})
		return
	}

	tierRanks := map[string]int{"free": 0, "basic": 1, "premium": 2}
	userLevel := tierRanks[subscription.Plan.AccessLevel]
	contentLevel := tierRanks[content.AccessLevel]

	if userLevel < contentLevel {
		c.JSON(http.StatusForbidden, gin.H{"error": "upgrade required to stream this content"})
		return
	}

	// In production, this would generate a signed AWS CloudFront Cookie or Token
	signedURL := content.VideoURL + "?token=mock_signed_token_allowing_streaming"

	c.JSON(http.StatusOK, gin.H{
		"manifest_url": signedURL,
		"message":      "Access granted",
	})
}

func GetContentRecommendations(c *gin.Context) {
	// Simple implementation: return all content as recommendations
	var contents []models.Content
	if err := database.DB.Find(&contents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch recommendations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recommendations": contents})
}

func CreateContent(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil || user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	var content models.Content
	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&content).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create content"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"content": content})
}

func UpdateContent(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil || user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	contentIDStr := c.Param("id")
	contentID, err := strconv.Atoi(contentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content id"})
		return
	}

	var content models.Content
	if err := database.DB.First(&content, contentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "content not found"})
		return
	}

	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&content).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update content"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"content": content})
}

func DeleteContent(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil || user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	contentIDStr := c.Param("id")
	contentID, err := strconv.Atoi(contentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content id"})
		return
	}

	if err := database.DB.Delete(&models.Content{}, contentID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete content"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "content deleted successfully"})
}
