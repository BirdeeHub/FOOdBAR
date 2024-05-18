package db

import (
	foodlib "FOOdBAR/FOOlib"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// TODO: Make this function work for new db scheme
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
	// TODO: get multiple dietary somehow
	// TODO: Check if row exists already, if so, do update instead?
	name := c.FormValue("itemName")
	dietary := c.FormValue("itemDietary_0")
	amount := c.FormValue("itemAmount")
	units := c.FormValue("itemUnits")
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
