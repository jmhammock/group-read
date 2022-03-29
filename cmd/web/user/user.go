package user

type User struct {
	Id          int64  `json:"id"`
	Email       string `json:"user_name"`
	Password    string `json:"-"`
	DisplayName string `json:"display_name"`
}
