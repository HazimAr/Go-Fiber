package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
	"gorm.io/gorm"
)

// User Auth Model
type User struct {
	ID        guuid.UUID `gorm:"primaryKey" json:"-"`
	Username  string     `json:"username"`
	Password  string     `json:"-"`
	Sessions  []Session  `gorm:"foreignKey:UserRefer; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	CreatedAt int64      `gorm:"autoCreateTime" json:"-" `
	UpdatedAt int64      `gorm:"autoUpdateTime:milli" json:"-"`
}

// Session Model for the user
type Session struct {
	Sessionid guuid.UUID `gorm:"primaryKey" json:"sessionid"`
	UserRefer guuid.UUID `json:"-"`
	CreatedAt int64      `gorm:"autoCreateTime" json:"-" `
	UpdatedAt int64      `gorm:"autoUpdateTime:milli" json:"-"`
}

// Initalize and set the authentication and authorization routes
func AuthRoutes(router fiber.Router, db *gorm.DB) {
	auth := router.Group("/auth", securityMiddleware)
	auth.Post("/login", func(c *fiber.Ctx) error {
		json := new(User)
		if err := c.BodyParser(json); err != nil {
			return c.SendStatus(500)
		}
		empty := User{}
		if json.Username == empty.Username || empty.Password == json.Password {
			return c.Status(401).SendString("Invalid Data Sent")
		}

		foundUser := User{}
		queryUser := User{Username: json.Username}
		err := db.First(&foundUser, &queryUser).Error
		if err == gorm.ErrRecordNotFound {
			return c.Status(401).SendString("User not Found")
		}
		if foundUser.Password != json.Password {
			return c.Status(401).SendString("Incorrect Password")
		}
		newSession := Session{UserRefer: foundUser.ID, Sessionid: guuid.New()}
		CreateErr := db.Create(&newSession).Error
		if CreateErr != nil {
			return c.Status(500).SendString("Creation Error")
		}
		return c.Status(200).JSON(newSession)
	})

	auth.Post("/logout", func(c *fiber.Ctx) error {
		json := new(Session)
		if err := c.BodyParser(json); err != nil {
			return c.SendStatus(500)
		}
		if json.Sessionid == new(Session).Sessionid {
			return c.Status(401).SendString("Invalid Data Sent")
		}
		session := Session{}
		query := Session{Sessionid: json.Sessionid}
		err := db.First(&session, query).Error
		if err == gorm.ErrRecordNotFound {
			return c.Status(401).SendString("Session Not Found")
		}
		db.Delete(&session)
		return c.SendStatus(200)
	})
	auth.Post("/create", func(c *fiber.Ctx) error {
		json := new(User)
		if err := c.BodyParser(json); err != nil {
			return c.SendStatus(500)
		}
		empty := User{}
		if json.Username == empty.Username || empty.Password == json.Password {
			return c.Status(401).SendString("Invalid Data Sent")
		}
		newUser := User{
			Username: json.Username,
			Password: json.Password,
			ID:       guuid.New(),
		}
		foundUser := User{}
		query := User{Username: json.Username}
		err := db.First(&foundUser, &query).Error
		if err != gorm.ErrRecordNotFound {
			return c.Status(401).SendString("User Already Exists")
		}
		db.Create(&newUser)
		return c.SendStatus(200)
	})
	auth.Post("/user", func(c *fiber.Ctx) error {
		user := User{}
		myUser := User{Username: "NikSchaefer"}
		Sessions := []Session{}
		db.First(&user, &myUser)
		db.Model(&user).Association("Sessions").Find(&Sessions)
		user.Sessions = Sessions
		return c.JSON(user)
	})
	auth.Post("/delete", func(c *fiber.Ctx) error {
		json := new(User)
		if err := c.BodyParser(json); err != nil {
			return c.SendStatus(500)
		}
		empty := User{}
		if json.Username == empty.Username || empty.Password == json.Password {
			return c.Status(401).SendString("Invalid Data Sent")
		}
		foundUser := User{}
		query := User{Username: json.Username}
		err := db.First(&foundUser, &query).Error
		if err == gorm.ErrRecordNotFound {
			return c.Status(401).SendString("User Not Found")
		}
		if json.Password != foundUser.Password {
			return c.Status(401).SendString("Invalid Credentials")
		}
		db.Model(&foundUser).Association("Sessions").Clear()
		createErr := db.Delete(&foundUser).Error
		if createErr != nil {
			fmt.Println(createErr)
		}
		return c.SendStatus(200)
	})
	auth.Post("/update", func(c *fiber.Ctx) error {
		json := new(User)
		if err := c.BodyParser(json); err != nil {
			return c.SendStatus(500)
		}
		empty := User{}
		if json.Username == empty.Username || empty.Password == json.Password {
			return c.Status(401).SendString("Invalid Data Sent")
		}
		foundUser := User{}
		query := User{Username: json.Username}
		err := db.First(&foundUser, &query).Error
		if err == gorm.ErrRecordNotFound {
			return c.Status(401).SendString("User Not Found")
		}
		return c.SendStatus(200)
	})

}

func securityMiddleware(c *fiber.Ctx) error {
	c.Set("X-XSS-Protection", "1; mode=block")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("X-Download-Options", "noopen")
	c.Set("Strict-Transport-Security", "max-age=5184000")
	c.Set("X-Frame-Options", "DENY")
	c.Set("X-DNS-Prefetch-Control", "off")
	c.Accepts("application/json")
	return c.Next()
}