package handlers

import (
	"github.com/Dhruv9449/mou/pkg/database"
	"github.com/Dhruv9449/mou/pkg/serializers"
	"github.com/Dhruv9449/mou/pkg/utils"
	"github.com/Dhruv9449/mou/pkg/utils/dbutils"
	"github.com/Dhruv9449/mou/pkg/utils/driveutils"
	"github.com/gofiber/fiber/v2"
)

func getDriveHome(c *fiber.Ctx) error {
	user, err := utils.GetUserFromHeader(c)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	if user.Role != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You are not authorized to perform this action",
		})
	}

	return c.Status(fiber.StatusOK).JSON(serializers.FolderFullSerializer(dbutils.GetRootFolder()))
}

func getDriveFolder(c *fiber.Ctx) error {
	user, err := utils.GetUserFromHeader(c)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	if user.Role != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You are not authorized to perform this action",
		})
	}

	return c.Status(fiber.StatusOK).JSON(serializers.FolderFullSerializer(dbutils.GetFolder(c.Params("id"))))
}

func getDriveFile(c *fiber.Ctx) error {
	user, err := utils.GetUserFromHeader(c)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	if user.Role != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You are not authorized to perform this action",
		})
	}

	return c.Status(fiber.StatusOK).JSON(serializers.FileSerializer(dbutils.GetFile(c.Params("id"))))
}

func createDriveFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	parent_folder := c.FormValue("parent_folder")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file",
		})
	}

	if parent_folder == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid parent folder",
		})
	}

	user, err := utils.GetUserFromHeader(c)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	if user.Role != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You are not authorized to perform this action",
		})
	}

	fileContent, err := file.Open()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file",
		})
	}

	drivefile := driveutils.CreateFile(driveutils.SERVICE, file.Filename, parent_folder, fileContent)

	if drivefile == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error while creating file",
		})
	}

	dbfile := serializers.DriveToDBFileSerializer(drivefile)

	database.DB.Create(dbfile)

	return c.Status(fiber.StatusOK).JSON(serializers.FileSerializer(dbfile))

}

func createDriveFolder(c *fiber.Ctx) error {
	name := c.FormValue("name")
	parent_folder := c.FormValue("parent_folder")

	if name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid name",
		})
	}

	if parent_folder == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid parent folder",
		})
	}

	user, err := utils.GetUserFromHeader(c)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	if user.Role != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You are not authorized to perform this action",
		})
	}

	drivefolder := driveutils.CreateFolder(driveutils.SERVICE, name, parent_folder)

	if drivefolder == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error while creating folder",
		})
	}

	dbfolder := serializers.DriveToDBFolderSerializer(drivefolder)

	database.DB.Create(dbfolder)

	return c.Status(fiber.StatusOK).JSON(serializers.FolderBlockSerializer(dbfolder))
}

func deleteDriveFile(c *fiber.Ctx) error {
	user, err := utils.GetUserFromHeader(c)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	if user.Role != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You are not authorized to perform this action",
		})
	}

	file := dbutils.GetFile(c.Params("id"))

	if file.ID == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	driveutils.DeleteFile(driveutils.SERVICE, file.ID)

	database.DB.Delete(file)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "File deleted",
	})
}

func deleteDriveFolder(c *fiber.Ctx) error {
	user, err := utils.GetUserFromHeader(c)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	if user.Role != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You are not authorized to perform this action",
		})
	}

	folder := dbutils.GetFolder(c.Params("id"))

	if folder.ID == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Folder not found",
		})
	}

	driveutils.DeleteFolder(driveutils.SERVICE, folder.ID)

	dbutils.DeleteFolderAndChildren(folder.ID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Folder deleted",
	})
}

func DriveRouter(app fiber.Router) {
	group := app.Group("/drive")
	group.Get("/", getDriveHome)
	group.Get("/folder/:id", getDriveFolder)
	group.Get("/file/:id", getDriveFile)
	group.Post("/file", createDriveFile)
	group.Post("/folder", createDriveFolder)
	group.Delete("/file/:id", deleteDriveFile)
	group.Delete("/folder/:id", deleteDriveFolder)
}
