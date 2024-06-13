package main

import (
	_ "fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/rahulgubili3003/postgres-hrms-go/model"
	"github.com/rahulgubili3003/postgres-hrms-go/storage"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

type Employee struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Department string `json:"department"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateEmployee(context *fiber.Ctx) error {
	employee := Employee{}

	err := context.BodyParser(&employee)

	if err != nil {
		err := context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		if err != nil {
			return err
		}
		return err
	}

	err = r.DB.Create(&employee).Error
	if err != nil {
		err := context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create book"})
		if err != nil {
			return err
		}
		return err
	}

	err = context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book has been added"})
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetEmployees(context *fiber.Ctx) error {
	bookModels := &[]model.Employee{}

	err := r.DB.Find(bookModels).Error
	if err != nil {
		err := context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get books"})
		if err != nil {
			return err
		}
		return err
	}

	err = context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "books fetched successfully",
		"data":    bookModels,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateEmployee)
	api.Get("/books", r.GetEmployees)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("could not load the database")
	}
	err = model.MigrateEmployee(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	err = app.Listen(":8081")
	if err != nil {
		return
	}
}
