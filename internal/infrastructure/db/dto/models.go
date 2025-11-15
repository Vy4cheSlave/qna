package dto

import (
	"time"
)

type Question struct {
	Id        int `gorm:"primaryKey;autoIncrement"`
	Text      string
	CreatedAt time.Time
}

type Answer struct {
	Id         int `gorm:"primaryKey;autoIncrement"`
	QuestionId int
	UserId     string
	Text       string
	CreatedAt  time.Time

	Question Question `gorm:"foreignKey:QuestionId;references:Id;constraint:OnDelete:CASCADE"`
	User     User     `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE"`
}

type User struct {
	Id        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string `gorm:"type:varchar(100);not null"`
	CreatedAt time.Time
}
