package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Plan      string         `gorm:"default:'free'" json:"plan"` // free, basic, premium
	StripeID  string         `gorm:"index" json:"stripe_id"`     // Stripe Customer ID
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Subscription struct {
	ID                   uint           `gorm:"primaryKey" json:"id"`
	UserID               uint           `gorm:"index" json:"user_id"`
	StripeSubscriptionID string         `gorm:"uniqueIndex" json:"stripe_subscription_id"`
	Status               string         `json:"status"` // active, past_due, canceled
	CurrentPeriodEnd     time.Time      `json:"current_period_end"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
}

type Video struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	URL         string         `json:"url"`          // CDN URL for the manifest (e.g., .m3u8)
	Thumbnail   string         `json:"thumbnail"`    // CDN URL for thumbnail
	RequiresTier string        `json:"require_tier"` // basic, premium
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
