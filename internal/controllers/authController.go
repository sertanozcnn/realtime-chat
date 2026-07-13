package controllers

import (
	"context"
	"git/internal/database"
	"git/internal/models"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email, password, first name and last name.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.CreateUser true "User register details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 502 {object} map[string]interface{}
// @Router /user/signup [post]
func Register(c *fiber.Ctx) error {

	var UserSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var body models.CreateUser

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	CheckUser := UserSchema.FindOne(ctx, bson.D{{Key: "email", Value: body.Email}}).Decode(&body)

	if CheckUser == nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message": "User with  email " + body.Email + "Already exist!",
		})
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	newUser := models.UserModel{
		Name:      body.FirstName + " " + body.LastName,
		Email:     body.Email,
		Password:  string(hashPassword),
		Followers: make([]string, 0),
		Following: make([]string, 0),
	}

	result, err := UserSchema.InsertOne(ctx, &newUser)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(err)
	}

	var createdUser *models.UserModel
	query := bson.M{"_id": result.InsertedID}

	UserSchema.FindOne(ctx, query).Decode(&createdUser)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    createdUser.ID.Hex(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	JwtSecret := os.Getenv("JWT_SECRET")

	token, _ := claims.SignedString([]byte(JwtSecret))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"result": createdUser,
		"token":  token,
	})
}
