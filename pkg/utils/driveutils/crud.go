package driveutils

import (
	"mime/multipart"

	"google.golang.org/api/drive/v3"
)

func CreateFolder(service *drive.Service, name string, parent string) *drive.File {
	folder := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parent},
	}

	folder, err := service.Files.Create(folder).SupportsAllDrives(true).Do()

	if err != nil {
		panic(err)
	}

	return folder
}

func CreateFile(service *drive.Service, name string, parent string, r multipart.File) *drive.File {
	file := &drive.File{
		MimeType: "application/octet-stream",
		Name:     name,
		Parents:  []string{parent},
	}

	file, err := service.Files.Create(file).Media(r).Do()

	if err != nil {
		panic(err)
	}

	return file
}

func GetFile(service *drive.Service, id string) *drive.File {
	file, err := service.Files.Get(id).SupportsAllDrives(true).Do()

	if err != nil {
		panic(err)
	}

	return file
}

func GetFolder(service *drive.Service, id string) *drive.File {
	folder := GetFile(service, id)

	if folder.MimeType != "application/vnd.google-apps.folder" {
		panic("Not a folder")
	}

	return folder
}

func GetFiles(service *drive.Service, parent string) []*drive.File {
	q := "'" + parent + "' in parents"
	files, err := service.Files.List().Q(q).Do()

	if err != nil {
		panic(err)
	}

	return files.Files
}

func GetFolders(service *drive.Service, parent string) []*drive.File {
	q := "'" + parent + "' in parents and mimeType = 'application/vnd.google-apps.folder'"
	folders, err := service.Files.List().Q(q).Do()

	if err != nil {
		panic(err)
	}

	return folders.Files
}

func DeleteFile(service *drive.Service, id string) {
	err := service.Files.Delete(id).Do()

	if err != nil {
		panic(err)
	}
}

func DeleteFolder(service *drive.Service, id string) {
	folder := GetFolder(service, id)
	files := GetFiles(service, folder.Id)

	for _, file := range files {
		DeleteFile(service, file.Id)
	}

	folders := GetFolders(service, folder.Id)

	for _, folder := range folders {
		DeleteFolder(service, folder.Id)
	}

	DeleteFile(service, folder.Id)
}

func RenameFile(service *drive.Service, id string, name string) *drive.File {
	file := GetFile(service, id)
	file.Name = name
	file, err := service.Files.Update(id, file).Do()

	if err != nil {
		panic(err)
	}

	return file
}

func RenameFolder(service *drive.Service, id string, name string) *drive.File {
	folder := GetFolder(service, id)
	folder.Name = name
	folder, err := service.Files.Update(id, folder).Do()

	if err != nil {
		panic(err)
	}

	return folder
}
