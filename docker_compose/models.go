package main

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CoinOutput struct {
	gorm.Model
	Id        string         `json: "id,omitempty"`
	Exchanges pq.StringArray `gorm:"type:text[]" json:"exchanges"`
	TaskRun   int            `json: "taskrun,omitempty"`
}

type Response struct {
	Id        string         `json: "id,omitempty"`
	Exchanges pq.StringArray `gorm:"type:text[]" json:"exchanges"`
	TaskRun   int            `json: "taskrun,omitempty"`
}
