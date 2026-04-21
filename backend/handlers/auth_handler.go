package handlers

import (
	"auth-backend/config"
	"auth-backend/dto"
	"auth-backend/models"
	"auth-backend/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB  *gorm.DB
	Cfg config.AppConfig
}

func NewAuthHandler(db *gorm.DB, cfg config.AppConfig) *AuthHandler {
	return &AuthHandler{DB: db, Cfg: cfg}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	familyID := uuid.NewString()
	jti := uuid.NewString()

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, h.Cfg.AccessTokenSecret, h.Cfg.AccessTokenMinutes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot generate access token"})
		return
	}

	refreshToken, exp, err := utils.GenerateRefreshToken(user.ID, jti, familyID, h.Cfg.RefreshTokenSecret, h.Cfg.RefreshTokenDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot generate refresh token"})
		return
	}

	ip := c.ClientIP()
	ua := c.Request.UserAgent()

	rt := models.RefreshToken{
		UserID:           user.ID,
		TokenHash:        utils.SHA256(refreshToken),
		JTI:              jti,
		FamilyID:         familyID,
		ExpiresAt:        exp,
		CreatedIP:        &ip,
		CreatedUserAgent: &ua,
	}

	if err := h.DB.Create(&rt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot save refresh token"})
		return
	}

	c.SetCookie(
		"refresh_token",
		refreshToken,
		h.Cfg.RefreshTokenDays*24*60*60,
		"/",
		h.Cfg.CookieDomain,
		h.Cfg.CookieSecure,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "missing refresh token"})
		return
	}

	claims, err := utils.ParseRefreshToken(refreshToken, h.Cfg.RefreshTokenSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid refresh token"})
		return
	}

	var current models.RefreshToken
	if err := h.DB.Where("jti = ?", claims.JTI).First(&current).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "refresh token not found"})
		return
	}

	// ตรวจ hash ตรงกันไหม
	if current.TokenHash != utils.SHA256(refreshToken) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "refresh token mismatch"})
		return
	}

	// หมดอายุ
	if time.Now().After(current.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "refresh token expired"})
		return
	}

	// reuse detection
	if current.RevokedAt != nil {
		now := time.Now()
		h.DB.Model(&models.RefreshToken{}).
			Where("family_id = ? AND revoked_at IS NULL", current.FamilyID).
			Updates(map[string]interface{}{
				"revoked_at":        now,
				"reuse_detected_at": now,
			})

		c.JSON(http.StatusUnauthorized, gin.H{"message": "refresh token reuse detected"})
		return
	}

	var user models.User
	if err := h.DB.First(&user, current.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "user not found"})
		return
	}

	newJTI := uuid.NewString()

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, h.Cfg.AccessTokenSecret, h.Cfg.AccessTokenMinutes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot generate access token"})
		return
	}

	newRefreshToken, exp, err := utils.GenerateRefreshToken(user.ID, newJTI, current.FamilyID, h.Cfg.RefreshTokenSecret, h.Cfg.RefreshTokenDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot generate refresh token"})
		return
	}

	now := time.Now()

	tx := h.DB.Begin()

	if err := tx.Model(&current).Updates(map[string]interface{}{
		"revoked_at":      now,
		"replaced_by_jti": newJTI,
		"last_used_at":    now,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot revoke old refresh token"})
		return
	}

	parentJTI := current.JTI
	ip := c.ClientIP()
	ua := c.Request.UserAgent()

	newRT := models.RefreshToken{
		UserID:           user.ID,
		TokenHash:        utils.SHA256(newRefreshToken),
		JTI:              newJTI,
		FamilyID:         current.FamilyID,
		ParentJTI:        &parentJTI,
		ExpiresAt:        exp,
		CreatedIP:        &ip,
		CreatedUserAgent: &ua,
	}

	if err := tx.Create(&newRT).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot save new refresh token"})
		return
	}

	tx.Commit()

	c.SetCookie(
		"refresh_token",
		newRefreshToken,
		h.Cfg.RefreshTokenDays*24*60*60,
		"/",
		h.Cfg.CookieDomain,
		h.Cfg.CookieSecure,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh_token")
	if refreshToken != "" {
		hash := utils.SHA256(refreshToken)
		now := time.Now()

		h.DB.Model(&models.RefreshToken{}).
			Where("token_hash = ? AND revoked_at IS NULL", hash).
			Updates(map[string]interface{}{
				"revoked_at": now,
			})
	}

	c.SetCookie("refresh_token", "", -1, "/", h.Cfg.CookieDomain, h.Cfg.CookieSecure, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
