package db

import (
	foodlib "FOOdBAR/FOOlib"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

func SubmitPantryItem(c echo.Context, pd *foodlib.PageData, td *foodlib.TabData, item *foodlib.TabItem) error {
	if item.Ttype == foodlib.Invalid {
		return errors.New("invalid Tab Type")
	}
	userID, err := foodlib.GetUserFromClaims(foodlib.GetClaimsFromContext(c))
	if err != nil {
		return err
	}
	db, tableName, err := CreateTabTableIfNotExists(userID, td.Ttype)
	if err != nil {
		return err
	}
	defer db.Close()

	params, err := c.FormParams()
	if err != nil {
		return err
	}
	c.Logger().Print(params)
	name := c.FormValue("itemName")
	rawdietary := params["dietary[]"]
	amount := c.FormValue("itemAmount")
	units := c.FormValue("itemUnits")
	var amountFloat float64
	if amount == "" {
		amountFloat = 0.0
	} else {
		amountFloat, err = strconv.ParseFloat(amount, 64)
		if err != nil {
			return errors.New("amount is not a number")
		}
	}

	dietary, err := json.Marshal(rawdietary)
	if err != nil {
		return err
	}

	// Check if row exists already, if so, do update instead
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM "+tableName+" WHERE id = ?)", item.ItemID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		updateStmt, err := db.Prepare(fmt.Sprintf("UPDATE %s SET name = ?, dietary = ?, amount = ?, units = ? WHERE id = ?", tableName))
		if err != nil {
			return err
		}
		defer updateStmt.Close()

		_, err = updateStmt.Exec(name, dietary, amountFloat, units, item.ItemID)
		return err
	}

	insertStmt, err := db.Prepare(fmt.Sprintf(`INSERT INTO %s (id, last_author, name, dietary, amount, units) VALUES (?, ?, ?, ?, ?, ?)`, tableName))
	if err != nil {
		return err
	}
	defer insertStmt.Close()

	_, err = insertStmt.Exec(item.ItemID.String(), pd.UserID, name, dietary, amount, units)
	return err
}
