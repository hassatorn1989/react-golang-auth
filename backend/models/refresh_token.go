package models

import "time"

type RefreshToken struct {
	ID               uint64     `gorm:"primaryKey" json:"id"`
	UserID           uint64     `gorm:"index;not null" json:"user_id"`
	TokenHash        string     `gorm:"size:255;not null" json:"token_hash"`
	JTI              string     `gorm:"size:100;uniqueIndex;not null" json:"jti"`
	FamilyID         string     `gorm:"size:100;index;not null" json:"family_id"`
	ParentJTI        *string    `gorm:"size:100" json:"parent_jti"`
	ExpiresAt        time.Time  `gorm:"index;not null" json:"expires_at"`
	RevokedAt        *time.Time `json:"revoked_at"`
	ReplacedByJTI    *string    `gorm:"size:100" json:"replaced_by_jti"`
	CreatedIP        *string    `gorm:"size:100" json:"created_ip"`
	CreatedUserAgent *string    `gorm:"size:255" json:"created_user_agent"`
	LastUsedAt       *time.Time `json:"last_used_at"`
	ReuseDetectedAt  *time.Time `json:"reuse_detected_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}
