package db

import (
	foodlib "FOOdBAR/FOOlib"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)
func SubmitPantryItem(c echo.Context, pd *foodlib.PageData, td *foodlib.TabData, item *foodlib.TabItem) error {
	if item.Ttype == foodlib.Invalid {
		return errors.New("Invalid Tab Type")
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

	// TODO: implement getting more fields to fill in for dietary and then
	// make this be able to recieve them all and make into json array for saving
	name := c.FormValue("itemName")
	dietary := c.FormValue("itemDietary_0")
	amount := c.FormValue("itemAmount")
	units := c.FormValue("itemUnits")
	
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

		_, err = updateStmt.Exec(name, dietary, amount, units, item.ItemID)
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

// TODO: This function
func GetTabItemData(c echo.Context, item foodlib.TabItem) (interface{}, error) {
	if item.Ttype == foodlib.Invalid {
		return nil, errors.New("Invalid Tab Type")
	}
	userID, err := foodlib.GetUserFromClaims(foodlib.GetClaimsFromContext(c))
	if err != nil {
		return nil, err
	}
	db, tableName, err := CreateTabTableIfNotExists(userID, item.Ttype)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	return nil, errors.New("not yet implemented"+tableName)
}

// TODO: should not fetch data, but instead, which tabItems to fetch data from
func FillXTabItems(userID uuid.UUID, tbd *foodlib.TabData, number int) error {
	db, _, err := CreateTabTableIfNotExists(userID, tbd.Ttype)
	defer db.Close()
	if err != nil {
		return err
	}
	// TODO: fill tbd.Items with X number of items based on tbd.OrderBy: SortMethod
	// where the key string is a column name

	return nil
}
