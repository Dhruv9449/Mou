package driveutils

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/Dhruv9449/mou/pkg/database"
	"github.com/Dhruv9449/mou/pkg/models"
	"github.com/Dhruv9449/mou/pkg/serializers"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var ROOT string
var SERVICE *drive.Service
var MEDIA_FOLDER string
var PERSONAL_FOLDER string

func GetClient(filename string) *http.Client {
	bfile, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var creds map[string]interface{}

	err = json.Unmarshal(bfile, &creds)

	if err != nil {
		panic(err)
	}

	ROOT = creds["root_dir"].(string)

	config := &jwt.Config{
		Email:      creds["client_email"].(string),
		PrivateKey: []byte(creds["private_key"].(string)),
		Scopes:     []string{drive.DriveScope},
		TokenURL:   google.JWTTokenURL,
	}

	client := config.Client(context.Background())
	return client
}

func GetService(filename string) *drive.Service {
	client := GetClient(filename)
	service, err := drive.NewService(context.Background(), option.WithHTTPClient(client))

	if err != nil {
		panic(err)
	}

	return service
}

func InitializeDrive() {
	service := GetService("./.env/drive-credentials.json")
	SERVICE = service

	if database.DB.Find(&models.Folder{}).Where("id = ?", ROOT).RowsAffected == 0 {
		root_folder := GetRootFolder(service)
		folder := serializers.DriveToDBFolderSerializer(root_folder)
		database.DB.Create(&folder)
	}

	MEDIA_FOLDER = GetMediaFolder(service).Id

	if MEDIA_FOLDER == "" {
		media_folder := CreateFolder(service, "Media", ROOT)
		folder := serializers.DriveToDBFolderSerializer(media_folder)
		database.DB.Create(&folder)
		MEDIA_FOLDER = media_folder.Id
	}

	PERSONAL_FOLDER = GetPersonalFolder(service).Id

	if PERSONAL_FOLDER == "" {
		personal_folder := CreateFolder(service, "Personal", ROOT)

		folder := serializers.DriveToDBFolderSerializer(personal_folder)
		database.DB.Create(&folder)
		PERSONAL_FOLDER = personal_folder.Id
	}
}
