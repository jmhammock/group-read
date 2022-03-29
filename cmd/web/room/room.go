package room

type Room struct {
	Id       string  `json:"id"`
	Book     string  `json:"book"`
	UsersIds []int64 `json:"user_ids"`
}
