package handlers

import (
	"net/http"
	"strconv"
	"time"

	"aubergine/internal/database"
	"aubergine/internal/models"

	"github.com/gin-gonic/gin"
)

type SubscribeRequest struct {
	PlanID uint `json:"plan_id" binding:"required"`
}

func Subscribe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var plan models.Plan
	if err := database.DB.First(&plan, req.PlanID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
		return
	}

	// Create subscription
	subscription := models.UserSubscription{
		UserID:    uint(userID.(uint)),
		PlanID:    req.PlanID,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, plan.ValidityDays),
		IsActive:  true,
		Status:    "active",
	}

	if err := database.DB.Create(&subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create subscription"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"subscription": subscription})
}

func GetSubscriptionHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var subscriptions []models.UserSubscription
	if err := database.DB.Where("user_id = ?", userID).Preload("Plan").Find(&subscriptions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch subscriptions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subscriptions": subscriptions})
}

func CancelSubscription(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	subscriptionIDStr := c.Param("id")
	subscriptionID, err := strconv.Atoi(subscriptionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id"})
		return
	}

	var subscription models.UserSubscription
	if err := database.DB.Where("id = ? AND user_id = ?", subscriptionID, userID).First(&subscription).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	subscription.IsActive = false
	subscription.Status = "canceled"

	if err := database.DB.Save(&subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cancel subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "subscription canceled successfully"})
}
