package models

import "time"

type Message struct {
	ID          uint   `gorm:"primary_key"`
	Content     string `gorm:"not null"`
	Author      User   `gorm:"ForeignKey:AuthorID"`
	AuthorID    uint
	Recipient   User `gorm:"ForeignKey:RecipientID"`
	RecipientID uint
	CreatedOn   time.Time `gorm:"not null"`
	ChatRoom    ChatRoom  `gorm:"ForeignKey:ChatRoomID"`
	ChatRoomID  uint
}

type ChatRoom struct {
	ID        uint   `gorm:"primary_key"`
	Users     []User `gorm:"many2many:chat_room_users;"`
	Messages  []Message
	CreatedOn time.Time `gorm:"not null"`
}
