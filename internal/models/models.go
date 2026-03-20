package models

type Note struct {
	ID        int    `json:"id"`
	Note      string `json:"note"`
	CreatedAt string `json:"created_at"`
}
