package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"aubergine/internal/database"
	"aubergine/internal/models"
)

// ListVideos retrieves all videos available, potentially filtering out based on tier client-side or server-side.
func ListVideos(c *gin.Context) {
	var videos []models.Video
	
	// Optional: Get user plan if authenticated, and filter. Currently returning all for display.
	if err := database.DB.Find(&videos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch videos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"videos": videos})
}

// StreamVideo acts as a mock signed URL generator or direct manifest server
func StreamVideo(c *gin.Context) {
	videoID := c.Param("id")
	
	var video models.Video
	if err := database.DB.First(&video, videoID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
		return
	}

	// This route is protected by MinimumTier dynamically based on video.RequiresTier
	// However, since video info is dynamic, we perform an inline tier check:
	userPlan, exists := c.Get("plan")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	
	tierRanks := map[string]int{"free": 0, "basic": 1, "premium": 2}
	planStr := userPlan.(string)
	
	if tierRanks[planStr] < tierRanks[video.RequiresTier] {
		c.JSON(http.StatusForbidden, gin.H{"error": "upgrade required to stream this content"})
		return
	}

	// In production, this would generate a signed AWS CloudFront Cookie or Token
	signedURL := video.URL + "?token=mock_signed_token_allowing_streaming"

	c.JSON(http.StatusOK, gin.H{
		"manifest_url": signedURL,
		"message":      "Access granted",
	})
}
