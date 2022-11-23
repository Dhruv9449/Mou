package serializers

import (
	"github.com/Dhruv9449/mou/pkg/models"
	"google.golang.org/api/drive/v3"
)

func DriveToDBFolderSerializer(folder *drive.File) models.Folder {
	return models.Folder{
		ID:       folder.Id,
		Name:     folder.Name,
		ParentID: folder.Parents[0],
	}
}

func DriveToDBFileSerializer(file *drive.File) models.File {
	return models.File{
		ID:       file.Id,
		Name:     file.Name,
		ParentID: file.Parents[0],
	}
}

func FolderBlockSerializer(folder models.Folder) map[string]interface{} {
	return map[string]interface{}{
		"id":          folder.ID,
		"name":        folder.Name,
		"description": folder.Description,
		"url":         folder.URL(),
	}
}

func FolderListSerializer(folders []models.Folder) []map[string]interface{} {
	var folderList []map[string]interface{}
	for _, folder := range folders {
		folderList = append(folderList, FolderBlockSerializer(folder))
	}
	return folderList
}

func FileSerializer(file models.File) map[string]interface{} {
	return map[string]interface{}{
		"id":          file.ID,
		"name":        file.Name,
		"description": file.Description,
		"url":         file.URL(),
	}
}

func FileListSerializer(files []models.File) []map[string]interface{} {
	var fileList []map[string]interface{}
	for _, file := range files {
		fileList = append(fileList, FileSerializer(file))
	}
	return fileList
}

func FolderFullSerializer(folder models.Folder) map[string]interface{} {
	return map[string]interface{}{
		"id":          folder.ID,
		"name":        folder.Name,
		"description": folder.Description,
		"url":         folder.URL(),
		"subFolders":  FolderListSerializer(folder.SubFolders()),
		"files":       FileListSerializer(folder.Files()),
	}
}
