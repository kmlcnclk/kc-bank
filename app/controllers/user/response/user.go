package response

import (
	"kc-bank/domain"
	"time"
)

type UserResponse struct {
	Id        string    `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Age       int32     `json:"age"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func ToUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Age:       user.Age,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserResponseList(users []*domain.User) []UserResponse {
	var response = make([]UserResponse, 0)

	for _, user := range users {
		response = append(response, ToUserResponse(user))
	}

	return response
}
