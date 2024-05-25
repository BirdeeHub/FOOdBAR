
# AI VOMIT

Below is stuff that is probably mostly random AI garbage but I might take 1 or 2 useful things from it still to look up further

---

### SQL


> how to select items conditionally from a database?

Here are a few ways to conditionally select data from a database in SQL:

**WHERE clause**

Use a WHERE clause to filter rows based on a condition:

```sql
SELECT * FROM table WHERE column_name = some_value
```

**CASE statement** 

Use a CASE statement to selectively output columns based on conditions:

```sql
SELECT 
  CASE WHEN condition THEN column_1 ELSE column_2 END AS column
FROM table
```

**IF function**

Some databases support an IF function to return one value or another conditionally:

```sql 
SELECT 
  IF(condition, column_1, column_2) AS column
FROM table
``` 

**JOINs**

Join tables together conditionally to select related data:

```sql
SELECT * 
FROM table_1
LEFT JOIN table_2 ON table_1.id = table_2.foreign_key
WHERE table_2.id IS NULL
```

This will select rows from table_1 that don't have a match in table_2.

**UNION** 

Combine the results of multiple SELECT statements conditionally:

```sql
SELECT column FROM table WHERE condition_1
UNION 
SELECT column FROM table WHERE condition_2
```

This allows you to select rows meeting different conditions in one query.

So in summary, use WHERE, CASE, IF, JOINS, and UNION to add conditional logic to your SELECT queries.

> If, for example, I wanted to select the highest 50 values of a column, instead of all matching the criteria, how would I do this effectively?

To select only the top N values from a column, you can use the LIMIT clause along with ORDER BY.

For example, to select the 50 highest values from a "score" column:

```sql
SELECT score
FROM table
ORDER BY score DESC
LIMIT 50;
```

This will:

1. Order the rows by score in descending order (highest first)

2. Return only the first 50 rows 

So we get the top 50 highest scores.

You can also combine this with WHERE clauses to filter before selecting the top values:

```sql 
SELECT score
FROM table
WHERE condition
ORDER BY score DESC
LIMIT 50;
```

Some key points:

- ORDER BY determines the sorting order before the LIMIT is applied.

- You need to specify DESC to sort highest first, otherwise it will give the lowest values.

- LIMIT without ORDER BY will give arbitrary rows rather than top values.

So in summary, use ORDER BY and LIMIT together to get the top N rows according to a specified ordering. Add a WHERE clause to filter before selecting top values.

> what other things can I sort by?

In SQL, including SQLite, there are various ways to sort data beyond just ascending and descending order. Here are some common sorting options:

1. **Ascending Order (ASC)**: This is the default sorting order if not specified explicitly. It sorts the data from the lowest value to the highest.

    ```sql
    SELECT column_name
    FROM table_name
    ORDER BY column_name ASC;
    ```

2. **Descending Order (DESC)**: It sorts the data from the highest value to the lowest.

    ```sql
    SELECT column_name
    FROM table_name
    ORDER BY column_name DESC;
    ```

3. **Alphabetical Order (ASC/DESC)**: For text data, you can sort alphabetically. This follows the natural order of the characters in the specified collation.

    ```sql
    SELECT column_name
    FROM table_name
    ORDER BY column_name ASC;
    ```

4. **Custom Sorting**: You can also use a custom sorting order using the `CASE` statement to define specific sorting criteria.

    ```sql
    SELECT column_name
    FROM table_name
    ORDER BY
        CASE 
            WHEN condition THEN value
            WHEN condition THEN value
            ELSE value
        END;
    ```

5. **Sorting by Multiple Columns**: You can sort by multiple columns. Rows with the same value in the first column are then sorted by the second column, and so on.

    ```sql
    SELECT column_name1, column_name2
    FROM table_name
    ORDER BY column_name1 ASC, column_name2 DESC;
    ```

6. **Random Order**: To retrieve data in a random order, you can use the `RANDOM()` or `RAND()` function depending on the database system.

    ```sql
    SELECT column_name
    FROM table_name
    ORDER BY RANDOM();
    ```

These are some of the common sorting options available in SQL, but the specifics may vary depending on the database system you're using.


>>> In SQLite, you can create triggers for various events that occur on tables. Here are some of the most common ones:1. **INSERT**: Triggered when a new row is inserted into the table.2. **UPDATE**: Triggered whe
... n an existing row is updated in the table (like I showed you earlier).3. **DELETE**: Triggered when a row is deleted from the table. Can I get the row that these were triggered for, and put the row into the au
... dit table witht he extra info?

In SQLite, you can access the affected rows using the `NEW` table in UPDATE and DELETE triggers, and the `INSERTED` table in INSERT triggers.

Here's how you can modify your trigger to include the affected row in the audit table:
```sql
CREATE TRIGGER update_item_audits
AFTER UPDATE OF itemID ON items
BEGIN
    INSERT INTO audits (userID, itemID, updated_at)
    VALUES (NEW.userID, NEW.itemID, datetime('now'));
END;
```
In this trigger:

* `NEW` refers to the row that was updated. You can access its columns using the syntax `NEW.column_name`, such as `NEW.userID`.
* The values from the `NEW` row are used to populate the audit table with the updated data.

For INSERT and DELETE triggers, you would use the `INSERTED` or `DELETED` tables, respectively:
```sql
CREATE TRIGGER insert_item_audits
AFTER INSERT ON items
BEGIN
    INSERT INTO audits (userID, itemID, created_at)
    VALUES (NEW.userID, NEW.itemID, datetime('now'));
END;

CREATE TRIGGER delete_item_audits
AFTER DELETE ON items
BEGIN
    INSERT INTO audits (userID, itemID, deleted_at)
    VALUES (OLD.userID, OLD.itemID, datetime('now'));
END;
```
In these triggers:

* `NEW` refers to the row that was inserted. You can access its columns using the syntax `NEW.column_name`, such as `NEW.userID`.
* `OLD` refers to the row that was deleted. You can access its columns using the syntax `OLD.column_name`, such as `OLD.userID`.

Remember to replace `datetime('now')` with your preferred timestamp format.

These triggers will populate the audit table with the updated, inserted, or deleted data, along with the relevant timestamps and user IDs.

> Is there a generic trigger that will trigger on ANY event for the table, and also insert which action it was into the table?
> 
> ```markdown
> CREATE TRIGGER update_item_audits
> AFTER UPDATE OF itemID ON items
> BEGIN
>     INSERT INTO audits (userID, itemID, updated_at)
>     VALUES (NEW.userID, NEW.itemID, datetime('now'));
> END;
> ```

Yes, you can create a trigger that fires on INSERT, UPDATE, or DELETE by using `AFTER INSERT OR UPDATE OR DELETE` instead of specifying a specific operation. 

You can also check which operation triggered the trigger using the `TG_OP` variable.

Here is an example:

```sql
CREATE TRIGGER audit_changes
AFTER INSERT OR UPDATE OR DELETE ON items
FOR EACH ROW
BEGIN
  INSERT INTO audit_log (table_name, operation, rowid, changed_at) 
  VALUES ('items', TG_OP, NEW.rowid, CURRENT_TIMESTAMP);
END;
```

This will insert a row into the `audit_log` table on any INSERT, UPDATE or DELETE, and populate the `operation` column with 'INSERT', 'UPDATE' or 'DELETE' based on which operation triggered it.

The `NEW` and `OLD` references work the same as individual triggers. For INSERT, only `NEW` is populated. For DELETE, only `OLD` is populated. For UPDATE, both are populated.

So this allows you to create a single trigger that audits all changes to a table, instead of needing separate triggers for each operation.


> can I do lists in sql as a value in a table?
> 

Yes, you can store lists or arrays as a value in a SQL table column. Here are a few ways to do it in SQLite:

**Comma-separated string**

Store the list as a comma-separated string:

```sql
CREATE TABLE recipes (
  id INTEGER PRIMARY KEY,
  name TEXT,
  ingredients TEXT -- comma separated list
);

-- Example row
INSERT INTO recipes VALUES (
  1, 
  'Omelette',
  'eggs,cheese,milk,salt,pepper'
);
```

Then you can split the string on commas when retrieving the data.

**JSON array** 

Store the list as a JSON array string:

```sql
CREATE TABLE recipes (
  id INTEGER PRIMARY KEY, 
  name TEXT,
  ingredients TEXT -- json array
);

-- Example row 
INSERT INTO recipes VALUES (
  1,
  'Pasta',
  '["noodles", "tomato sauce", "parmesan"]' 
);
```

SQLite has JSON functions like JSON_EXTRACT to query JSON values.

**Separate table**

Use a separate table to store the list elements, with a foreign key to the parent row:

```sql
CREATE TABLE recipe_ingredients (
  recipe_id INTEGER,
  ingredient TEXT,
  FOREIGN KEY (recipe_id) REFERENCES recipes(id)
);

INSERT INTO recipe_ingredients VALUES
  (1, 'eggs'),
  (1, 'cheese'),
  (1, 'milk'); 
```

Then join the tables to get the ingredients for a recipe.

So in summary, you have several options to store lists in SQLite columns, with trade-offs.

> I want to sort based on values in json lists that are values in the db

To sort based on values in a JSON list column in SQLite, you can use the JSON functions to extract the values and sort on them.

For example, if you have a table like:

```sql
CREATE TABLE recipes (
  id INTEGER PRIMARY KEY,
  name TEXT,
  ingredients TEXT -- json array 
);

INSERT INTO recipes VALUES
  (1, 'Omelette', '["eggs", "cheese", "milk"]'),
  (2, 'Pasta', '["noodles", "tomato sauce", "parmesan"]'); 
```

You can sort by the length of the ingredients array like this:

```sql
SELECT name, ingredients
FROM recipes
ORDER BY JSON_ARRAY_LENGTH(ingredients) DESC;
```

This will extract the length of the JSON array and sort by it.

To access specific element values, you can use JSON_EXTRACT:

```sql
SELECT name, ingredients
FROM recipes
ORDER BY JSON_EXTRACT(ingredients, '$[0]') DESC; 
```

This will sort by the first element in the array. The $[index] syntax lets you access array elements.

You can also filter based on array values:

```sql
SELECT * 
FROM recipes
WHERE JSON_EXTRACT(ingredients, '$[0]') = 'eggs'
```

So in summary, use the SQLite JSON functions like JSON_ARRAY_LENGTH, JSON_EXTRACT and JSON_TYPE to query and sort by values within JSON columns.



---


> db is a sqlx db and I wish to read the rows directly into a map of rowname, value, where some of them may be json arrays
> 
> ```go
> func GetTabItemData(pd *foodlib.PageData, item *foodlib.TabItem) (map[string]interface{}, error) {
> 	if item.Ttype == foodlib.Invalid {
> 		return nil, errors.New("Invalid Tab Type")
> 	}
> 	db, tableName, err := CreateTabTableIfNotExists(pd.UserID, item.Ttype)
> 	if err != nil {
> 		return nil, err
> 	}
> 	defer db.Close()
> 	return nil, errors.New("not yet implemented"+tableName)
> }
> ```

Here is one way to read the rows into a map with sqlx:

```go
func GetTabItemData(pd *foodlib.PageData, item *foodlib.TabItem) (map[string]interface{}, error) {
	if item.Ttype == foodlib.Invalid {
		return nil, errors.New("Invalid Tab Type")
	}
	db, tableName, err := CreateTabTableIfNotExists(pd.UserID, item.Ttype)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	
	rows, err := db.Queryx("SELECT * FROM " + tableName)
	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	results := make(map[string]interface{})
	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}

		for i, col := range columns {
			if b, ok := values[i].([]byte); ok {
				err = json.Unmarshal(b, &results[col])
				if err != nil {
					results[col] = string(b)
				}
			} else {
				results[col] = values[i]
			}
		}
	}
	
	return results, nil
}
```

This scans each row into a slice of interface{}, then iterates over the columns and inserts into the map, handling JSON parsing for []byte columns.

> how do I make sure that a value in a table gets read as []byte type? do I make the type in the database BYTEARRAY? Is that a thing in sqlite?

In SQLite, you can store binary data in a BLOB column type. This will be returned from sqlite as []byte when reading the data in Go.

So in your table definition, you would specify the column that needs to store binary data as:

```sql
CREATE TABLE mytable (
  id INTEGER PRIMARY KEY, 
  name TEXT,
  data BLOB
)
```

And then when reading that back in Go with sqlx, the data column would be returned as []byte which you can then json.Unmarshal.

Some key points:

- BLOB is the column type for binary data in SQLite
- sqlx will return BLOB columns as []byte 
- You can then use json.Unmarshal to parse the []byte if it contains JSON

So using BLOB columns and scanning into []byte is the way to handle JSON encoding/decoding when storing in SQLite from Go.

> I want to get the sql types of the columns and cast the values of the map to the correct type before returning
> 
> ```go
> func GetTabItemData(userID uuid.UUID, item *foodlib.TabItem) (map[string]interface{}, error) {
> 	if item == nil {
> 		return nil, errors.New("nil tab target")
> 	}
> 	if item.Ttype == foodlib.Invalid {
> 		return nil, errors.New("invalid Tab Type")
> 	}
> 	db, tableName, err := CreateTabTableIfNotExists(userID, item.Ttype)
> 	if err != nil {
> 		return nil, err
> 	}
> 	defer db.Close()
> 
> 	data := make(map[string]interface{})
> 	err = db.QueryRowx("SELECT * FROM "+tableName+" WHERE id = ?", item.ItemID).MapScan(data)
> 	if err == sql.ErrNoRows {
> 		return nil, nil
> 	} else if err != nil {
> 		return nil, err
> 	}
> 
> 	return data, nil
> }
> ```

Here is how you can get the column types and cast the values to the correct types:

```go
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

	// Get column names and types
	cols, err := db.Queryx("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = ?", tableName)
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

	// Query data and cast values to correct types
	data := make(map[string]interface{})
	row := db.QueryRowx("SELECT * FROM "+tableName+" WHERE id = ?", item.ItemID)
	err = row.MapScan(data)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	for col, typ := range colTypes {
		switch typ {
		case "int":
			val, _ := data[col].(int64)
			data[col] = int(val)
		case "float":
			val, _ := data[col].(float64)
			data[col] = float32(val)
		// add cases for other types
		}
	}

	return data, nil
}
```

This queries the column information to get a mapping of column name to data type, uses that when scanning the result row to cast values to the correct types.

> why did you convert the float64 into a float32?
> 
> 	for col, typ := range colTypes {
> switch typ {
> case "int":
>     val, _ := data[col].(int64)
>     data[col] = int(val)
> case "float":
>     val, _ := data[col].(float64)
>     data[col] = float32(val)
> // add cases for other types
> }
> }

You're right, converting the float64 to float32 is not necessary in this case. I did that by default assuming the destination type would be float32, but since we are just storing into an interface{} map there is no need to convert the type.

A better way to handle the floats would be:

```go
for col, typ := range colTypes {
  switch typ {
  case "int":
    val, _ := data[col].(int64) 
    data[col] = int(val)

  case "float":
    val := data[col].(float64)
    data[col] = val

  // add cases for other types
  }
}
```

So just leave floats as float64 rather than converting. Thanks for catching that!

> im using sqlite. It says
> 
> no such table: information_schema.columns

You're right, the `information_schema` tables are not present in SQLite by default. 

Instead of querying `information_schema`, we can use the `sqlite_master` table to get column info:

```go
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
```

The `pragma_table_info` pragma query returns the column name and type for a given table.

So in SQLite we can use `pragma_table_info` instead of `information_schema` to get the column metadata.

---

### notes


```go
// render all items in a list with a component
// lmao I coulda just done a for loop. Saving this here in case I want to reference it.
 func RenderListWithComponent[T interface{}](list []T, component func(T) templ.Component) templ.Component {
 	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
 		for _, item := range list {
 			err := component(item).Render(ctx, w)
 			if err != nil {
 				return err
 			}
 		}
 		return nil
 	})
 }
```
