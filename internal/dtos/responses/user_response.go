package responses

type UserResponse struct {
	Token string `json:"token"`
}

func (res *UserResponse) FormatterLoginUser(token string) UserResponse {
	return UserResponse{Token: token}
}
