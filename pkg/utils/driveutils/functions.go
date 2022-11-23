package driveutils

import "google.golang.org/api/drive/v3"

func GetRootFolder(service *drive.Service) *drive.File {
	folder := GetFolder(service, ROOT)
	return folder
}

func GetMediaFolder(service *drive.Service) *drive.File {
	folders := GetFolders(service, ROOT)

	for _, folder := range folders {
		if folder.Name == "Media" {
			return folder
		}
	}
	return nil
}

func GetPersonalFolder(service *drive.Service) *drive.File {
	folders := GetFolders(service, ROOT)

	for _, folder := range folders {
		if folder.Name == "Personal" {
			return folder
		}
	}
	return nil
}
