package model

import (
	"wongnok/internal/model/dto"

	"gorm.io/gorm"
)

type Favorite struct {
	gorm.Model
	FoodRecipeID uint
	UserID       string
}

func (favorite Favorite) FromRequest(request dto.FavoriteRequest) Favorite {
	return Favorite{
		FoodRecipeID: request.FoodRecipeID,
	}
}
func (favorite Favorite) ToResponse() dto.FavoriteResponse {
	return dto.FavoriteResponse{
		ID:           favorite.ID,
		FoodRecipeID: favorite.FoodRecipeID,
		UserID:       favorite.UserID,
	}
}

// Ratings คือ "ชุดของ Rating หลาย ๆ อัน"
// คือเราจะใช้ "slice ของ Rating" เราจึงตั้งชื่อใหม่ให้จำง่ายขึ้น
// ชื่อใหม่ (alias)
type Favorites []Favorite

func (favorite Favorites) ToResponse() dto.FavoritesResponse {
	var results = make([]dto.FavoriteResponse, 0)

	for _, Favorite := range favorite {
		results = append(results, Favorite.ToResponse())
	}

	return dto.FavoritesResponse{
		Results: results,
	}
}