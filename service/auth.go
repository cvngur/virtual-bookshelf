package service

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
	"time"
	"virtual-bookshelf/database"
	"virtual-bookshelf/model"
)

func Login(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	query := fmt.Sprintf("SELECT users.* FROM users WHERE email = '%s'", email)
	res, err := database.GetCluster().Query(query, nil)

	if err != nil {
		fmt.Print("Error")
	}

	var user map[string]interface{}
	err = res.One(&user)
	if err != nil {
		fmt.Print(err)
	}
	err = res.Err()
	if err != nil {
		fmt.Print(err)
	}

	hashedPassword, _ := json.Marshal(user["password"])
	s, _ := strconv.Unquote(string(hashedPassword))

	err = bcrypt.CompareHashAndPassword([]byte(s), []byte(password))
	if err != nil {
		return err
	}
	c.Cookie(&fiber.Cookie{
		Name:     "username",
		Value:    email,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		SameSite: "lax",
	})
	return nil
}

func Register(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	document := model.User{Email: email, Password: string(hashedPassword)}
	// Generating uuid
	b := make([]byte, 16)
	_, err = rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	_, err = database.GetCollection().Insert(uuid, &document, nil)
	if err != nil {
		panic(err)
	}
	return nil
}