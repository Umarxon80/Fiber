package db

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type User struct {
	Id           uint   `json:"id"`
	First_name   string `json:"first_name" validate:"required"`
	Last_name    string `json:"last_name" validate:"required"`
	Role         string `json:"role" default:"user"`
	Email        string `json:"email" validate:"required,email"`
	Phone_number string `json:"phone_number"`
	Age          uint8  `json:"age" validate:"gte=1,lte=120"`
	Password     string `json:"password" validate:"required,min=8"`
}

func createUserTable() error {
	_, err := DbConnection.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		role VARCHAR(255) DEFAULT 'user',
		email VARCHAR(255),
		phone_number VARCHAR(255),
		age SMALLINT,
		password VARCHAR(255) NOT NULL
	)
	`)
	return err
}

// // Users CRUD

func CreateUser(ctx fiber.Ctx) error {
	var user User
	var id uint8
	user = ctx.Locals("user").(User)
	err := DbConnection.QueryRow(context.Background(), `
		INSERT INTO users (
			first_name, last_name, role, email, phone_number, age, password
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		) RETURNING id
	`, user.First_name, user.Last_name, user.Role, user.Email, user.Phone_number, user.Age, user.Password).Scan(&id)
	if err != nil {
		log.Error("Error creating user ")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}
	user.Id = uint(id)
	log.Info("User created, id: ", id)
	return ctx.Status(fiber.StatusCreated).JSON(user)
}
func GetUsers(ctx fiber.Ctx) error {
	var users []User
	rows, err := DbConnection.Query(context.Background(), `
	SELECT * FROM users
	`)
	if err != nil {
		log.Error("Error getting users, ", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}
	for rows.Next() {
		var buffer User
		err := rows.Scan(&buffer.Id, &buffer.First_name, &buffer.Last_name, &buffer.Role, &buffer.Email, &buffer.Phone_number, &buffer.Age, &buffer.Password)
		if err != nil {
			log.Error("Error getting users, ", err)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err,
			})
		}
		users = append(users, buffer)
	}
	log.Info("Returning all users")
	return ctx.JSON(users)
}
func GetOneUser(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var user User
	err := DbConnection.QueryRow(context.Background(), `
	SELECT * FROM users where id=$1
	`, id).Scan(&user.Id, &user.First_name, &user.Last_name, &user.Role, &user.Email, &user.Phone_number, &user.Age, &user.Password)
	if err != nil {
		log.Error("Error getting user, ", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	log.Info("Returning one user, id: ", id)
	return ctx.JSON(user)
}
func GetUserByEmail(ctx fiber.Ctx, email string) (User, error) {
	var user User
	err := DbConnection.QueryRow(context.Background(), `
	SELECT * FROM users where email=$1
	`, email).Scan(&user.Id, &user.First_name, &user.Last_name, &user.Role, &user.Email, &user.Phone_number, &user.Age, &user.Password)
	if err != nil {
		log.Error("Error getting user, ", err)
		return User{}, err
	}
	return user, nil
}
func PatchUser(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var user User
	user = ctx.Locals("user").(User)

	ch, err := DbConnection.Exec(context.Background(), `
	UPDATE users
	SET  first_name=$1, last_name=$2, role=$3, email=$4, phone_number=$5, age=$6, password=$7
	WHERE id=$8
	`, user.First_name, user.Last_name, user.Role, user.Email, user.Phone_number, user.Age, user.Password, id)
	if err != nil {
		log.Error("Error updating user, ", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if ch.RowsAffected() < 1 {
		log.Error("User not found id: ", id)
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Such user does not exists",
		})
	}
	log.Info("User updated, id: ", id)
	return ctx.JSON(fiber.Map{
		"User updated, id: ": id,
	})
}
func DeleteUser(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	ch, err := DbConnection.Exec(context.Background(), `
	DELETE from users
	WHERE id=$1
	`, id)
	if err != nil {
		log.Error("Error updating user, ", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}
	if ch.RowsAffected() < 1 {
		log.Error("User not found id: ", id)
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Such user does not exists",
		})
	}
	log.Info("User deleted, id: ", id)
	return ctx.JSON(fiber.Map{
		"User deleted, id: ": id,
	})
}
