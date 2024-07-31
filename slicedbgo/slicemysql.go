package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User represents the user model
type User struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:255"`
	Age  int
}

func main() {
	// DSN (Data Source Name) format for MySQL
	dsn := "root:Suresh@123@tcp(127.0.0.1:3306)/learndb?charset=utf8mb4&parseTime=True&loc=Local"

	// Open a connection to the MySQL database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Failed to migrate the database: ", err)
	}

	fmt.Println("Database connection successful and schema migrated")

	// Create a new user
	//newUser := User{Name: "vaibhav viswanathan", Age: 14}
	newusers := []User{
		{Name: "Manoj", Age: 42},
		{Name: "Gayathiri", Age: 34},
		{Name: "Varsha", Age: 22},
	}

	// Insert the slice of users into the database
	result := db.Create(&newusers)
	//result := db.Create(&newUser)
	if result.Error != nil {
		log.Fatal("Failed to create a new user: ", result.Error)
	}

	fmt.Printf("Inserted %d users\n", len(newusers))

	// Query users
	var users []User
	result = db.Find(&users)
	if result.Error != nil {
		log.Fatal("Failed to query users: ", result.Error)
	}

	for _, user := range users {
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}
}
