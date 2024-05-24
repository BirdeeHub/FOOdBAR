package db

import (
	foodlib "FOOdBAR/FOOlib"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
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

func GetTabItemData(userID uuid.UUID, item *foodlib.TabItem) (map[string]interface{}, error) {
	if item == nil {
		return nil, errors.New("nil tab target")
	}
	if item.Ttype == foodlib.Invalid {
		return nil, errors.New("invalid Tab Type")
	}
	db, tableName, err := CreateTabTableIfNotExists(userID, item.Ttype)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var data map[string]interface{}
	err = db.QueryRow("SELECT * FROM "+tableName+" WHERE id = ?", item.ItemID).Scan(&data)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return data, nil
}

func GetTabItemDataValue[T any](raw map[string]interface{}, key string, out *T) error {
	if raw == nil {
		return errors.New("no data to search")
	}
	if rawval, ok := raw[key]; ok && rawval != nil {
		switch rawval.(type) {
		case T:
			outval := rawval.(T)
			out = &outval
		case []byte:
			return json.Unmarshal(rawval.([]byte), out)
		default:
			return errors.New("value's type does not match out, nor is it unmarshalable to match the type of out")
		}
	}
	return errors.New("lookup failed")
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
