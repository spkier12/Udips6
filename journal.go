package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Create new journal in database
func CreateJournal(c echo.Context) error {
	token, err := c.Cookie("token")
	if err != nil {
		fmt.Print(err)
		return c.JSON(http.StatusForbidden, returnData("Login mislyktes", "0"))
	}

	err1 := UserHasRole(token.Value, role)
	if err1 != nil {
		return c.JSON(http.StatusOK, returnData("Du har ikke tilgang her", "0"))
	}

	data := ReciveData(c)
	fmt.Print(len(data))

	//Create a new journal if not exist
	res, err2 := db.Exec("INSERT into useraccounts.journal VALUES($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING", data[0], data[1], data[2], data[3], data[4], data[5])
	if err2 != nil {
		fmt.Print(err2)
		return c.JSON(http.StatusOK, returnData("Journal oppretting mislykktes", "0"))
	}

	if d, _ := res.RowsAffected(); d < 1 {
		return c.JSON(http.StatusOK, returnData("Journal oppretting mislykktes", "0"))
	}

	return c.JSON(http.StatusOK, returnData("Journal var opprettet", "1"))
}

// UpdateJournal in database
func UpdateJournal(c echo.Context) error {
	token, err := c.Cookie("token")
	if err != nil {
		fmt.Print(err)
		return c.JSON(http.StatusForbidden, returnData("Login mislyktes", "0"))
	}

	err1 := UserHasRole(token.Value, role2)
	if err1 != nil {
		return c.JSON(http.StatusOK, returnData("Du har ikke tilgang her", "0"))
	}

	//Create a new journal if not exist
	data := ReciveData(c)
	res, err2 := db.Exec("UPDATE useraccounts.journals SET journalsickness=$2, journalalergic=$3, journalbloodtype=$4, journalnotes=$5 WHERE journalname=$1", data[1], data[2], data[3], data[4], data[5])
	if err2 != nil {
		fmt.Print(err2)
		return c.JSON(http.StatusOK, returnData("Journal oppretting mislykktes", "0"))
	}

	if d, _ := res.RowsAffected(); d < 1 {
		return c.JSON(http.StatusOK, returnData("Journal oppretting mislykktes", "0"))
	}

	return c.JSON(http.StatusOK, returnData("Journal var oppdatert", "1"))
}
