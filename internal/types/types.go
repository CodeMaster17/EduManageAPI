package types

type Student struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
