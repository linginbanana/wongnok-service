package profile

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type Handler struct {
    db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
    return &Handler{db: db}
}

func (h *Handler) Update(c *gin.Context) {
    userID := c.Param("id")
    var req struct {
        Name     string `json:"name"`
        Nickname string `json:"nickname"`
        Avatar   string `json:"avatar"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    if err := h.db.Model(&User{}).Where("id = ?", userID).Updates(User{
        Name:     req.Name,
        Nickname: req.Nickname,
        Avatar:   req.Avatar,
    }).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
}

func (h *Handler) Get(c *gin.Context) {
    userID := c.Param("id")
    var user User
    if err := h.db.First(&user, "id = ?", userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.JSON(http.StatusOK, user)
}
