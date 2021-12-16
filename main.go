package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

var db *sql.DB
var role = "udips6"
var role2 = "udips6_update"
var role3 = "udips6_admin"

// var role3 = "udips6_admin"

// Initialize database
func InitDB() (*sql.DB, error) {

	// Connection to database
	const (
		host     = "80.208.226.78"
		port     = "44320"
		user     = "zerorootzero"
		password = "zerorootzerozeroroorootzeroninehundredFiftyFive905327895642905327895642905327895642905327895642905327895642"
		dbname   = "zerorootzero"
	)

	con := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", con)
	return db, err
}

func main() {
	DB, _ := InitDB()
	db = DB
	err2 := db.Ping()
	if err2 != nil {
		fmt.Print("\nDatabase failed ping test...")
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	// User management
	e.POST("/CreateJournal", CreateJournal)               // Create a journal in database for selected user // role: udips6_update
	e.POST("/UpdateJournal", UpdateJournal)               // Update the desired users journal with correct information // role: udips6_update
	e.POST("/CreateIndicentJournal", AddIncidentJournal)  // Create a incident in database with the desired user // role: udips6_update
	e.GET("/GetJournals/:journal", GetJournals)           // Get all journals // role: udips6
	e.DELETE("/DeleteJournal", DeleteJournal)             // Delete Journal from database // role: udips6_admin
	e.GET("/GetJournalContent/:user", GetJournalIncident) // Delete Journal from database // role: udips6_admin

	e.Start(":5002")
}

// Easy function to recive data in return-
func ReciveData(c echo.Context) [10]string {
	type MyData struct {
		JournalName      string
		JournalData      string
		JournalSickness  string
		JournalAlergic   string
		JournalBloodtype string
		JournalNotes     string
	}

	var myData MyData
	data, _ := (io.ReadAll(c.Request().Body))
	json.Unmarshal(data, &myData)

	if myData.JournalName == "" {
		myData.JournalName = "demo"
	}

	if myData.JournalData == "" {
		myData.JournalData = time.Now().String()
	}

	return [10]string{myData.JournalName, myData.JournalData, myData.JournalSickness, myData.JournalAlergic, myData.JournalBloodtype, myData.JournalNotes}
}

// Easy function to generate data in return-
func returnData(message string, data string) string {
	type MyData struct {
		Message string
		Data    string
	}
	var mydata MyData
	mydata.Message = message
	mydata.Data = data
	d, _ := json.Marshal(mydata)
	return string(d)
}

func UserHasRole(key string, role string) error {
	// check if token is valid
	message, data := CheckIfExist(key)
	if data == "" {
		fmt.Print("\n")
		fmt.Print(message)
		return fmt.Errorf("key invalid")
	}

	fmt.Print("\nChecking if user has role")
	var UserRole string
	db.QueryRow("SELECT rolename FROM useraccounts.invites WHERE email=$1 AND rolename=$2", data, role).Scan(&UserRole)

	var UserRole2 string
	db.QueryRow("SELECT rolename FROM useraccounts.roles WHERE email=$1 AND rolename=$2", data, role).Scan(&UserRole)

	// Is the role correct?
	if strings.EqualFold(UserRole, role) {
		return nil
	}

	// Is the role correct?
	if strings.EqualFold(UserRole2, role) {
		return nil
	}
	return fmt.Errorf("you are not a part of this role")
}

// Check if token has not expired
func CheckIfExist(key string) (string, string) {

	// Get date
	year := time.Now().Year()
	month := time.Now().Month()
	day := time.Now().Day()

	// Find data in DB
	var email string
	var timer string
	db.QueryRow("SELECT idec, added FROM useraccounts.sessions WHERE sessiontoken=$1", key).Scan(&email, &timer)

	// Check if email is found if not the key dosnt exists
	if email == "" {
		return "Invalid key", ""
	}

	if timer == fmt.Sprint(year)+" "+fmt.Sprint(month)+" "+fmt.Sprint(day) {
		return "Login OK\r", email
	}

	return "Invalid key", ""
}
