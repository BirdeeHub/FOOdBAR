package db

import (
	foodlib "FOOdBAR/FOOlib"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
)

func FillXTabItems(userID uuid.UUID, tbd *foodlib.TabData, number int) error {
	db, tableName, err := CreateTabTableIfNotExists(userID, tbd.Ttype)
	defer db.Close()
	if err != nil {
		return err
	}
	// TODO: fill tbd.Items with X number of items based on SortMethod numbers stored in slice in TabData
	// get the actual SortMethod via GetSortMethodByNumber(k int)
	rows, err := db.Query("SELECT id FROM " + tableName)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return err
		}
		tbd.AddTabItem(&foodlib.TabItem{ItemID: id})
	}

	return nil
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

	data := make(map[string]interface{})
	err = db.QueryRowx("SELECT * FROM "+tableName+" WHERE id = ?", item.ItemID).MapScan(data)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	cols, err := db.Queryx("SELECT name, type FROM pragma_table_info(?)", tableName)
	if err != nil {
		return nil, err 
	}

	colTypes := make(map[string]string)
	for cols.Next() {
		var name, typ string
		if err := cols.Scan(&name, &typ); err != nil {
			return nil, err
		}
		colTypes[name] = typ
	}

	for col, typ := range colTypes {
		var ok bool
		switch typ {
		case "INTEGER":
			data[col], ok = data[col].(int)
		case "TEXT":
			data[col], ok = data[col].(string)
		case "FLOAT":
			data[col], ok = data[col].(float64)
		case "BLOB":
			data[col], ok = data[col].([]byte)
		case "DATETIME":
			data[col], ok = data[col].(time.Time)
		default:
			data[col], ok = data[col].([]byte)
		}
		if !ok {
			return nil, errors.New("column type does not match")
		}
	}

	return data, nil
}

func GetTabItemDataValue[T any](raw map[string]interface{}, key string, out *T) error {
	defer func() error {
		if r := recover(); r != nil {
			// recover from panic if type convert fails and return error
			err := fmt.Errorf("panic: %s", r) 
			return err
		} else {
			return nil
		}
	}()
	if raw == nil {
		return errors.New("no data to search")
	}
	if rawval, ok := raw[key]; ok && rawval != nil {
		switch rawval.(type) {
		case bool:
			reflect.ValueOf(out).Elem().Set(reflect.ValueOf(rawval))
			return nil
		case float64:
			reflect.ValueOf(out).Elem().Set(reflect.ValueOf(rawval))
			return nil
		case int:
			reflect.ValueOf(out).Elem().Set(reflect.ValueOf(rawval))
			return nil
		case int64:
			reflect.ValueOf(out).Elem().Set(reflect.ValueOf(rawval))
			return nil
		case time.Time:
			reflect.ValueOf(out).Elem().Set(reflect.ValueOf(rawval))
			return nil
		case string:
			reflect.ValueOf(out).Elem().Set(reflect.ValueOf(rawval))
			return nil
		case []byte:
			return json.Unmarshal(rawval.([]byte), out)
		default:
			return errors.New("value's type does not match out, nor is it unmarshalable to match the type of out")
		}
	}
	return errors.New("lookup failed")
}
