package bid

import "time"

type Bid struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	TenderId    string    `json:"tenderId"`
	AuthorType  string    `json:"authorType"`
	AuthorId    string    `json:"authorId"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"createdAt"`
}
