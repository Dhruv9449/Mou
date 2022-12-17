package models

import (
	"time"
)

type BlogPost struct {
	ID        uint   `gorm:"primary_key"`
	Title     string `gorm:"not null,unique"`
	Content   string `gorm:"not null"`
	Thumbnail string
	CreatedOn time.Time `gorm:"not null"`
	UpdatedOn time.Time `gorm:"not null"`
	Author    User      `gorm:"ForeignKey:AuthorID"`
	AuthorID  uint      `gorm:"not null"`
	Published bool      `gorm:"default:true"`
	Comments  []Comment `gorm:"ForeignKey:BlogPostID"`
}

type Comment struct {
	ID         uint      `gorm:"primary_key"`
	Content    string    `gorm:"not null"`
	Author     User      `gorm:"ForeignKey:AuthorID"`
	AuthorID   uint      `gorm:"not null"`
	BlogPost   BlogPost  `gorm:"ForeignKey:BlogPostID"`
	BlogPostID uint      `gorm:"not null"`
	CreatedOn  time.Time `gorm:"not null"`
	Replies    []Comment `gorm:"ForeignKey:ParentID"`
	ParentID   uint      `gorm:"default:null"`
}
