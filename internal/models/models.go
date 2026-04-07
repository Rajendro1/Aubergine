package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Email          string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash   string         `gorm:"not null" json:"-"`
	Name           string         `json:"name"`
	Phone          string         `json:"phone"`
	Bio            string         `json:"bio"`
	ProfilePicture string         `json:"profile_picture"`
	Role           string         `gorm:"default:'user'" json:"role"` // user, admin
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type Plan struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Name              string    `gorm:"not null" json:"name"`
	Price             float64   `json:"price"`
	ValidityDays      int       `json:"validity_days"`
	AccessLevel       string    `json:"access_level"` // free, basic, premium
	MaxDevicesAllowed int       `json:"max_devices_allowed"`
	Resolution        string    `json:"resolution"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type UserSubscription struct {
	ID                   uint           `gorm:"primaryKey" json:"id"`
	UserID               uint           `gorm:"index" json:"user_id"`
	PlanID               uint           `gorm:"index" json:"plan_id"`
	StartDate            time.Time      `json:"start_date"`
	EndDate              time.Time      `json:"end_date"`
	IsActive             bool           `gorm:"default:true" json:"is_active"`
	StripeSubscriptionID string         `gorm:"uniqueIndex" json:"stripe_subscription_id"`
	Status               string         `json:"status"` // active, past_due, canceled
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
	User                 User           `gorm:"foreignKey:UserID" json:"-"`
	Plan                 Plan           `gorm:"foreignKey:PlanID" json:"-"`
}

type Content struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	TrailerURL  string         `json:"trailer_url"`
	VideoURL    string         `json:"video_url"`
	AccessLevel string         `json:"access_level"` // free, basic, premium
	Genre       string         `json:"genre"`
	ReleaseDate time.Time      `json:"release_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type WatchHistory struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"index" json:"user_id"`
	ContentID       uint      `gorm:"index" json:"content_id"`
	ProgressSeconds int       `json:"progress_seconds"`
	IsCompleted     bool      `json:"is_completed"`
	LastWatchedAt   time.Time `json:"last_watched_at"`
	User            User      `gorm:"foreignKey:UserID" json:"-"`
	Content         Content   `gorm:"foreignKey:ContentID" json:"-"`
}

type ActiveStreamSession struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"index" json:"user_id"`
	DeviceID      string    `json:"device_id"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
	User          User      `gorm:"foreignKey:UserID" json:"-"`
}
