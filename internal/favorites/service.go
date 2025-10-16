package favorite

import (
	"wongnok/internal/global"
	"wongnok/internal/model"
	"wongnok/internal/users"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type IUserService user.IService

type IService interface {
	Get(userID string) (model.Favorites, error)
	GetByUser(foodRecipeQuery model.FoodRecipeQuery, claims model.Claims) (model.FoodRecipes, int64, error)
	Create(id int, claims model.Claims) (model.Favorite, error)
	Delete(id int, claims model.Claims) error
}

type Service struct {
	Repository  IRepository
	UserService IUserService
}

func NewService(db *gorm.DB) IService {
	return &Service{
		Repository:  NewRepository(db),
		UserService: user.NewService(db),
	}
}

func (service Service) GetByUser(foodRecipeQuery model.FoodRecipeQuery, claims model.Claims) (model.FoodRecipes, int64, error) {
	total, err := service.Repository.Count(claims.ID, foodRecipeQuery.Search)
	if err != nil {
		return nil, 0, err
	}
	recipes, err := service.Repository.GetByUser(foodRecipeQuery, claims.ID)
	if err != nil {
		return nil, 0, err
	}
	recipes = recipes.CalculateAverageRatings()
	return recipes, total, nil
}

func (service Service) Get(userID string) (model.Favorites, error) {
	favorites, err := service.Repository.Get(userID)
	if err != nil {
		return nil, err
	}

	return favorites, nil
}

func (service Service) Create(id int, claims model.Claims) (model.Favorite, error) {

	userID, err := service.UserService.GetByID(claims)
	if err != nil {
		return model.Favorite{}, errors.Wrap(err, "create Favorite")
	}

	favoriteFromGet, err := service.Repository.GetDeleteByID(id, claims.ID)
	if err != nil {
		var favorite model.Favorite
		favorite.FoodRecipeID = uint(id)
		favorite.UserID = userID.ID

		if err := service.Repository.Create(&favorite); err != nil {
			return model.Favorite{}, errors.Wrap(err, "create favorite")
		}

		return favorite, nil
	}

	if err := service.Repository.Update(int(favoriteFromGet.ID)); err != nil {
		return model.Favorite{}, errors.Wrap(err, "update Favorite")
	}

	return favoriteFromGet, err

}

func (service Service) Delete(id int, claims model.Claims) error {
	favorite, err := service.Repository.GetByID(id, claims.ID)
	if err != nil {
		// กรณีไม่พบ id ที่ต้องการ update
		return errors.Wrap(err, "find recipe")

	}

	if favorite.UserID != claims.ID {
		// กรณี user ที่ login ไม่ตรงกับ user ที่สร้าง recipe
		return global.ErrForbidden
	}

	return service.Repository.Delete(int(favorite.ID))
}