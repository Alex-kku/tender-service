package entity

import "time"

type Bid struct {
	Id          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Status      string    `db:"status"`
	TenderId    string    `db:"tender_id"`
	AuthorType  string    `db:"author_type"`
	AuthorId    string    `db:"author_id"`
	Version     int       `db:"version"`
	CreatedAt   time.Time `db:"created_at"`
}

type CreateBidInput struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"required,max=500"`
	TenderId    string `json:"tenderId" binding:"required,max=100"`
	AuthorType  string `json:"authorType"  binding:"required,oneof=Organization User"`
	AuthorId    string `json:"authorId" binding:"required,max=100"`
}

type EditBidInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
