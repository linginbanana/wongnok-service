package dto

type UserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Nickname  string `json:"nickName"`
	ImageUrl  string `json:"imageUrl"`
}

type UserRequest struct {
	NickName string `validate:"required"`
	ImageUrl string `validate:"required"`
}
