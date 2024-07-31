package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID        int
	Username  string
	Email     string
	CreatedAt time.Time
}

func main() {
	// Connect to MySQL
	db, err := connectDB()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Query all users
	users, err := queryAllUsers(db)
	if err != nil {
		panic(err.Error())
	}

	// Print users
	for _, user := range users {
		fmt.Printf("ID: %d, Username: %s, Email: %s, CreatedAt: %s\n", user.ID, user.Username, user.Email, user.CreatedAt.Format("2006-01-02 15:04:05"))
	}
}

func connectDB() (*sql.DB, error) {
	// MySQL database connection parameters
	dbHost := "localhost"
	dbPort := "3307"
	dbUser := "suresh"
	dbPassword := "suresh1982"
	dbName := "testdb"

	// Create a DSN (Data Source Name) string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Open database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MySQL database!")
	return db, nil
}

func queryAllUsers(db *sql.DB) ([]User, error) {
	// Prepare query
	query := "SELECT id, username, email, created_at FROM users"

	// Execute query
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize slice to hold users
	var users []User

	// Iterate over the result set
	for rows.Next() {
		var user User

		// Scan each row into variables
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &myTime{})
		if err != nil {
			return nil, err
		}

		// Append user to users slice
		users = append(users, user)
	}

	// Check for errors during rows iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

type myTime struct {
	time.Time
}

func (m *myTime) Scan(value interface{}) error {
	if value == nil {
		m.Time = time.Time{}
		return nil
	}
	t, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte but got %T", value)
	}
	parsedTime, err := time.Parse("2006-01-02 15:04:05", string(t))
	if err != nil {
		return err
	}
	m.Time = parsedTime
	return nil
}
