package models

import "github.com/Dhruv9449/mou/pkg/database"

type File struct {
	ID          string `gorm:"primaryKey"`
	Name        string
	Parent      Folder `gorm:"foreignKey:ParentID"`
	ParentID    string `gorm:"not null"`
	Description string
}

type Folder struct {
	ID          string `gorm:"primaryKey"`
	Name        string
	Description string
	Parent      *Folder `gorm:"foreignKey:ParentID"`
	ParentID    string
}

func (file *File) URL() string {
	return "https://drive.google.com/file/d/" + file.ID
}

func (folder *Folder) URL() string {
	return "https://drive.google.com/drive/folders/" + folder.ID
}

func (folder *Folder) SubFolders() []Folder {
	var subFolders []Folder
	database.DB.Where("parent_id = ?", folder.ID).Find(&subFolders)
	return subFolders
}

func (folder *Folder) Files() []File {
	var files []File
	database.DB.Where("parent_id = ?", folder.ID).Find(&files)
	return files
}
