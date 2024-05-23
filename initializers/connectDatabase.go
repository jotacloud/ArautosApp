package initializers

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {

	err := godotenv.Load(".env")
	if err != nil {
		panic("Erro ao carregar arquivo .env: " + err.Error())
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	portStr := os.Getenv("DB_PORT")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic("Erro ao converter a porta para inteiro: " + err.Error())
	}

	// Construindo a string DSN
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + strconv.Itoa(port) + " sslmode=disable TimeZone=America/Fortaleza"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Erro ao conectar com o bando de dados" + err.Error())
	} else {
		fmt.Println("Conexão com o banco de dados estabelecida!")
	}

	DB = db
}
