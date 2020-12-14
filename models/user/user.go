package user

// type DALRequest struct {
// 	Id       *int
// 	Name     *string
// 	Email    *string
// 	Password *string
// }

// type DALResponse struct {
// 	Id    int
// 	Name  string
// 	Email string
// }

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

// type UsersResponse struct {
// 	Users []UserResponse
// }
