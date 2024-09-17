package entity

import "time"

type Tender struct {
	Id             string    `db:"id"`
	Name           string    `db:"name"`
	Description    string    `db:"description"`
	ServiceType    string    `db:"service_type"`
	Status         string    `db:"status"`
	OrganizationId string    `db:"organization_id"`
	Version        int       `db:"version"`
	CreatedAt      time.Time `db:"created_at"`
}

type CreateTenderInput struct {
	Name            string `json:"name" binding:"required,max=100"`
	Description     string `json:"description" binding:"required,max=500"`
	ServiceType     string `json:"serviceType" binding:"required,oneof=Construction Delivery Manufacture"`
	OrganizationId  string `json:"organizationId" binding:"required,max=100"`
	CreatorUsername string `json:"creatorUsername" binding:"required"`
}

type EditTenderInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ServiceType *string `json:"serviceType"`
}
