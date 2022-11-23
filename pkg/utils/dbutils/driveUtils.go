package dbutils

import (
	"github.com/Dhruv9449/mou/pkg/database"
	"github.com/Dhruv9449/mou/pkg/models"
)

func GetRootFolder() models.Folder {
	var rootFolder models.Folder
	database.DB.Where("parent_id = ?", "Mou").First(&rootFolder)
	return rootFolder
}

func GetFolder(id string) models.Folder {
	var folder models.Folder
	database.DB.Where("id = ?", id).First(&folder)
	return folder
}

func GetFile(id string) models.File {
	var file models.File
	database.DB.Where("id = ?", id).First(&file)
	return file
}

func DeleteFile(id string) {
	database.DB.Delete(&models.File{}, id)
}

func DeleteFolderAndChildren(id string) {
	var folder models.Folder
	database.DB.Where("id = ?", id).First(&folder)
	children := folder.SubFolders()
	for _, child := range children {
		DeleteFolderAndChildren(child.ID)
	}
	files := folder.Files()
	for _, file := range files {
		DeleteFile(file.ID)
	}
	database.DB.Delete(&folder)
}
