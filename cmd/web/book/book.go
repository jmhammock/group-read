package book

type Book struct {
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	Location string `json:"location"`
}
