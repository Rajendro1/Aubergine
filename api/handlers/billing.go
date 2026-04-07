package handlers

import (
	"net/http"
	"os"

	"aubergine/internal/database"
	"aubergine/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

// Initialize Stripe Key
func init() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	if stripe.Key == "" {
		// Mock key for local dev
		stripe.Key = "sk_test_mockkey"
	}
}

// CreateCheckoutSession creates a Stripe checkout session for upgrading to premium.
func CreateCheckoutSession(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	// In a real app, you would fetch Price IDs dynamically or from env vars
	priceID := "price_mock_premium_tier"

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL:    stripe.String("http://localhost:8080/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:     stripe.String("http://localhost:8080/cancel"),
		CustomerEmail: stripe.String(user.Email),
	}

	s, err := session.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create checkout session", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": s.URL})
}
