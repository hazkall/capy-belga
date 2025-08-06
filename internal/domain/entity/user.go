package entity

type User struct {
	ID    int64  `json:"id,omitempty"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
