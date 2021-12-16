package main

import (
	"fmt"
	"net/http"
	"strings"
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

	err1 := UserHasRole(token.Value, role2)
	if err1 != nil {
		return c.JSON(http.StatusOK, returnData("Du har ikke tilgang her", "0"))
	}

	data := ReciveData(c)

	//Create a new journal if not exist
	res, err2 := db.Exec("INSERT into udips6.journal VALUES($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING", strings.ToUpper(data[0]), data[1], data[2], data[3], data[4], data[5])
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
	res, err2 := db.Exec("UPDATE udips6.journal SET journalsickness=$2, journalalergic=$3, journalbloodtype=$4, journalnotes=$5 WHERE journalname=$1", strings.ToUpper(data[0]), data[2], data[3], data[4], data[5])
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
		return c.JSON(http.StatusForbidden, returnData("Login mislyktes", "0"))
	}

	err1 := UserHasRole(token.Value, role2)
	if err1 != nil {
		return c.JSON(http.StatusOK, returnData("Du har ikke tilgang her", "0"))
	}

	// Add incident to database
	data := ReciveData(c)
	added := time.Now().Month().String() + " " + fmt.Sprint(time.Now().Day()) + " " + fmt.Sprint(time.Now().Year())
	rows, err := db.Exec("INSERT INTO udips6.incident VALUES($1, $2, $3)", strings.ToUpper(data[0]), data[5], added)
	if err != nil {
		return c.JSON(http.StatusOK, returnData("Hendelse ble ikke lagt til pga en feil som oppsto!", "0"))
	}

	if d, _ := rows.RowsAffected(); d < 1 {
		return c.JSON(http.StatusOK, returnData("Hendelse ble ikke lagt til pga en feil som oppsto!", "0"))
	}

	return c.JSON(http.StatusOK, returnData("Hendelse ble lagt til!", "1"))
}

// Get all journals or only selected from database
func GetJournals(c echo.Context) error {

	// Check if user key and role is valid
	token, err := c.Cookie("token")
	journalsearch := strings.ToUpper(c.Param("journal"))
	if err != nil {
		return c.JSON(http.StatusForbidden, returnData("Login mislyktes", "0"))
	}

	err1 := UserHasRole(token.Value, role)
	if err1 != nil {
		return c.JSON(http.StatusOK, returnData("Du har ikke tilgang her", "0"))
	}

	// Search for profiles
	row, err := db.Query("SELECT journalname FROM udips6.journal")
	if err != nil {
		return c.JSON(http.StatusOK, returnData("Uventet feil opp sto henting av journaler", "0"))
	}

	// Add all profiles to variable
	var rowData string
	if strings.EqualFold(journalsearch, "ALL") {
		for row.Next() {
			var rowData2 string
			row.Scan(&rowData2)
			rowData += rowData2 + "-|||-"
		}
	} else {
		for row.Next() {
			var rowData2 string
			row.Scan(&rowData2)
			doescontain := strings.ContainsAny(rowData2, strings.ToUpper(journalsearch))

			if doescontain {
				rowData += rowData2 + "-|||-"
			}
		}
	}

	return c.JSON(http.StatusOK, returnData("Journaler var hentet!", rowData))
}

// Delete the journal +  incidents from database
func DeleteJournal(c echo.Context) error {
	// Check if user key and role is valid
	token, err := c.Cookie("token")
	if err != nil {
		return c.JSON(http.StatusForbidden, returnData("Login mislyktes", "0"))
	}

	err1 := UserHasRole(token.Value, role3)
	if err1 != nil {
		return c.JSON(http.StatusOK, returnData("Du har ikke tilgang her", "0"))
	}

	// Delete Journal
	journalrolename := ReciveData(c)
	res, err := db.Exec("DELETE FROM udips6.journal WHERE journalname=$1", journalrolename[0])
	if err != nil {
		return c.JSON(http.StatusOK, returnData("Feil oppsto under sletting!", "0"))
	}

	if d, _ := res.RowsAffected(); d < 1 {
		return c.JSON(http.StatusOK, returnData("Feil oppsto under sletting!", "0"))
	}

	return c.JSON(http.StatusOK, returnData("Journal slettet!", "1"))
}

// Get the journals content
func GetJournalIncident(c echo.Context) error {
	// Check if user key and role is valid
	token, err := c.Cookie("token")
	user := strings.ToUpper(c.Param("user"))
	if err != nil {
		return c.JSON(http.StatusForbidden, returnData("Login mislyktes", "0"))
	}

	err1 := UserHasRole(token.Value, role)
	if err1 != nil {
		return c.JSON(http.StatusOK, returnData("Du har ikke tilgang her", "0"))
	}

	res, err := db.Query("SELECT incident FROM udips6.incident WHERE journalname=$1", strings.Trim(user, ":"))
	if err != nil {
		return c.JSON(http.StatusOK, returnData("Feil oppsto under henting av data!", "0"))
	}

	var data string
	for res.Next() {
		var data2 string
		res.Scan(&data2)
		data += data2 + "-|||-"
		fmt.Print(data2)
	}
	return c.JSON(http.StatusOK, returnData("Henter data....", data))
}
