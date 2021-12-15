package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// Create new journal in database
func CreateJournal(c echo.Context) error {

	// Check if user key and role is valid
	token, err := c.Cookie("token")
	if err != nil {
		return c.JSON(http.StatusForbidden, returnData("Login mislyktes", "0"))
	}

	err1 := UserHasRole(token.Value, role)
	if err1 != nil {
		return c.JSON(http.StatusOK, returnData("Du har ikke tilgang her", "0"))
	}

	data := ReciveData(c)

	//Create a new journal if not exist
	res, err2 := db.Exec("INSERT into udips6.journal VALUES($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING", data[0], data[1], data[2], data[3], data[4], data[5])
	if err2 != nil {
		return c.JSON(http.StatusOK, returnData("Journal oppretting mislykktes", "0"))
	}

	if d, _ := res.RowsAffected(); d < 1 {
		return c.JSON(http.StatusOK, returnData("Journal oppretting mislykktes", "0"))
	}

	return c.JSON(http.StatusOK, returnData("Journal var opprettet", "1"))
}

// UpdateJournal in database
func UpdateJournal(c echo.Context) error {

	// Check if user key and role is valid
	token, err := c.Cookie("token")
	if err != nil {
		return c.JSON(http.StatusForbidden, returnData("Login mislyktes", "0"))
	}

	err1 := UserHasRole(token.Value, role2)
	if err1 != nil {
		return c.JSON(http.StatusOK, returnData("Du har ikke tilgang her", "0"))
	}

	//Create a new journal if not exist
	data := ReciveData(c)
	res, err2 := db.Exec("UPDATE udips6.journal SET journalsickness=$2, journalalergic=$3, journalbloodtype=$4, journalnotes=$5 WHERE journalname=$1", data[0], data[2], data[3], data[4], data[5])
	if err2 != nil {
		return c.JSON(http.StatusOK, returnData("Journal oppdatering mislykktes", "0"))
	}

	if d, _ := res.RowsAffected(); d < 1 {
		return c.JSON(http.StatusOK, returnData("Journal oppdatering mislykktes finns brukern?", "0"))
	}

	return c.JSON(http.StatusOK, returnData("Journal var oppdatert", "1"))
}

// Add incident to journal's profile
func AddIncidentJournal(c echo.Context) error {

	// Check if user key and role is valid
	token, err := c.Cookie("token")
	if err != nil {
		fmt.Print(err)
		return c.JSON(http.StatusForbidden, returnData("Login mislyktes", "0"))
	}

	err1 := UserHasRole(token.Value, role2)
	if err1 != nil {
		return c.JSON(http.StatusOK, returnData("Du har ikke tilgang her", "0"))
	}

	// Add incident to database
	data := ReciveData(c)
	added := time.Now().Month().String() + " " + fmt.Sprint(time.Now().Day()) + " " + fmt.Sprint(time.Now().Year())
	rows, err := db.Exec("INSERT INTO udips6.incident VALUES($1, $2, $3)", data[0], data[5], added)
	if err != nil {
		return c.JSON(http.StatusOK, returnData("Hendelse ble ikke lagt til pga en feil som oppsto!", "0"))
	}

	if d, _ := rows.RowsAffected(); d < 1 {
		return c.JSON(http.StatusOK, returnData("Hendelse ble ikke lagt til pga en feil som oppsto!", "0"))
	}

	return c.JSON(http.StatusOK, returnData("Hendelse ble lagt til!", "1"))
}
