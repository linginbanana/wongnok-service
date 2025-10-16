package favorite

type Favorite struct {
    ID       uint   `gorm:"primaryKey" json:"id"`
    UserID   string `json:"userId"`
    RecipeID string `json:"recipeId"`
}
