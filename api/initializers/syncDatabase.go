package initializers

import (
	"ArautosApp/models"
)

func SyncDatabase() {

	DB.AutoMigrate(&models.User{})
}
