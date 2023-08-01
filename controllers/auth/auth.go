package auth

import (
	"errors"
	"time"
	"todo-api/database"
	"todo-api/models"

	"todo-api/config"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *fiber.Ctx) error {

	type RegisterInput struct {
		Name            string `json:"name" form:"name"`
		Username        string `json:"username" form:"username"`
		Password        string `json:"password" form:"password"`
		ConfirmPassword string `json:"confirm_password" form: "confirm_password"`
	}

	type NewUser struct {
		Name     string `json:"name"`
		Username string `json:"username"`
	}

	db := database.DB.Db

	var input RegisterInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
		})
	}

	if input.Password != input.ConfirmPassword {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": "Password does not match",
		})
	}

	hash, err := hashPassword(input.Password)

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid Password",
		})
	}

	user := models.User{
		Name:     input.Name,
		Username: input.Username,
		Password: hash,
	}

	//query user table for existing username, if no error, that means a user exists and username
	//is already taken
	if err := db.Where("username = ?", user.Username).First(&user).Error; err == nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": "Username has already been taken",
		})
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": "An error occured creating user",
		})
	}

	newUser := NewUser{
		Name:     user.Name,
		Username: user.Username,
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Created user", "data": newUser})
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func Login(c *fiber.Ctx) error {
	type loginInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type UserData struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var input loginInput
	var userData UserData

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid login request",
			"data":    err,
		})
	}

	username := input.Username
	password := input.Password

	user, err := getUserByUsername(username)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid credentials",
			"data":    err,
		})
	}

	if user != nil {
		userData = UserData{
			ID:       user.ID,
			Username: user.Username,
			Password: user.Password,
		}
	}

	//inputted password and stored password
	if !checkPasswordHash(password, userData.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid credentials",
			"data":    nil,
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["username"] = userData.Username
	claims["user_id"] = userData.ID
	claims["exp"] = time.Now().Add(time.Hour * 3).Unix()

	t, err := token.SignedString([]byte(config.Config("secret")))

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": t})
}

func getUserByUsername(username string) (*models.User, error) {
	db := database.DB.Db

	var user models.User
	if err := db.Where(&models.User{Username: username}).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ExtractUserId(c *fiber.Ctx) int {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	id := int(claims["user_id"].(float64))

	return id
}

func GetUserById(id int) (*models.User, error) {
	db := database.DB.Db

	var user models.User

	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func AuthenticatedUser(c *fiber.Ctx) (user *models.User, err error) {
	id := ExtractUserId(c)

	user, err = GetUserById(id)

	if err != nil {
		return nil, err
	}

	return user, nil
}
