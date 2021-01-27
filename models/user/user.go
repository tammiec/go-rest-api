package user

type UserRequest struct {
	Id       *int
	Name     *string
	Email    *string
	Password *string
}

type UserResponse struct {
	Id    int
	Name  string
	Email string
}
