package initializers

import "go-crud-example/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.Post{})
	DB.AutoMigrate(&models.User{})
}
