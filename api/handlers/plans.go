package handlers

import (
	"net/http"
	"strconv"

	"aubergine/internal/database"
	"aubergine/internal/models"

	"github.com/gin-gonic/gin"
)

func GetPlans(c *gin.Context) {
	var plans []models.Plan
	if err := database.DB.Find(&plans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch plans"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plans": plans})
}

func CreatePlan(c *gin.Context) {
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

	var plan models.Plan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create plan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"plan": plan})
}

func UpdatePlan(c *gin.Context) {
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

	planIDStr := c.Param("id")
	planID, err := strconv.Atoi(planIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan id"})
		return
	}

	var plan models.Plan
	if err := database.DB.First(&plan, planID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
		return
	}

	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plan": plan})
}

func DeletePlan(c *gin.Context) {
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

	planIDStr := c.Param("id")
	planID, err := strconv.Atoi(planIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan id"})
		return
	}

	if err := database.DB.Delete(&models.Plan{}, planID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "plan deleted successfully"})
}
