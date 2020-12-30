package dto

import "go_free/model"

type UserDto struct {
	Username  string `json:"username"`
	Telephone string `json:"telephone"`
}

func ToUserDto(user model.UserInfos) UserDto {
	return UserDto{
		Username:  user.Username,
		Telephone: user.Telephone,
	}
}
