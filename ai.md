
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

The `NEW` and `OLD` references work the same as individual triggers. For INSERT, only `NEW` is populated. For DELETE, only `OLD` is populated. For UPDATE, both are populated.

So this allows you to create a single trigger that audits all changes to a table, instead of needing separate triggers for each operation.


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

To sort by values within JSON blobs, you can use SQLite's JSON functions. Here's an example of how you can order by a specific key within a JSON blob:

```sql
SELECT *
FROM your_table
ORDER BY json_extract(your_json_column, '$.key_to_sort_by');
```

Replace `your_table` with the name of your table, `your_json_column` with the name of the column containing the JSON blob, and `key_to_sort_by` with the key within the JSON blob that you want to sort by.

For example, if you have a table `recipes` with a column `ingredients` containing a JSON blob like `{"ingredient1": "apple", "ingredient2": "banana"}`, and you want to sort by the `ingredient1` key, you can use:

```sql
SELECT *
FROM recipes
ORDER BY json_extract(ingredients, '$.ingredient1');
```

If you need to sort by multiple keys within the JSON blob, you can use a combination of `json_extract` and `||` (string concatenation) like this:

```sql
SELECT *
FROM your_table
ORDER BY json_extract(your_json_column, '$.key1') || json_extract(your_json_column, '$.key2');
```

This will sort by the concatenated values of `key1` and `key2` within the JSON blob.

Note that SQLite's JSON functions are available starting from version 3.38.0 (2022-01-10). If you're using an older version of SQLite, you may need to use a different approach, such as extracting the JSON values into separate columns or using a user-defined function.

> They are lists not objects for the most part, how would I sort by the values in the list rather than the index of the list element?
>

If the JSON data in your column is a list instead of an object, you can use the `json_extract` function with an array index to access the elements of the list. However, since you want to sort by the values in the list, you'll need to use a subquery or a window function to extract the values and then order by them.

Here's an example using a subquery:

```sql
SELECT *
FROM your_table
ORDER BY (
    SELECT json_extract(your_json_column, '$[0]')
    FROM your_table
    WHERE rowid = outer.rowid
);
```

This query selects all rows from `your_table` and orders them by the first element of the JSON list in the `your_json_column` column. The subquery `SELECT json_extract(your_json_column, '$[0]')` extracts the first element of the JSON list for each row, and the outer query orders by the result of this subquery.

If you want to order by multiple elements in the list, you can concatenate the `json_extract` calls in the subquery:

```sql
SELECT *
FROM your_table
ORDER BY (
    SELECT json_extract(your_json_column, '$[0]') || json_extract(your_json_column, '$[1]')
    FROM your_table
    WHERE rowid = outer.rowid
);
```

This will order by the concatenation of the first and second elements of the JSON list.

Alternatively, you can use the `json_group_array` function with a window function to extract the list elements into separate rows, and then order by the values:

```sql
SELECT *
FROM (
    SELECT rowid, json_extract(your_json_column, '$[' || idx || ']') AS value
    FROM your_table
    CROSS JOIN json_each(json_group_array(your_json_column))
    ORDER BY rowid, idx
)
ORDER BY value;
```

This query first creates a derived table that extracts each element of the JSON list into a separate row, along with the original row ID (`rowid`). It then orders the derived table by the extracted values (`value`).

Note that these examples assume that your JSON data is a list of scalar values (e.g., strings or numbers). If your JSON data is a list of objects or nested structures, you may need to adjust the `json_extract` calls accordingly.

> what is a window function?

A window function in SQL is a function that performs a calculation across a set of rows that are somehow related to the current row. Unlike aggregate functions like `SUM` or `AVG`, which operate on groups of rows, window functions operate on a "window" of data, which is a subset of rows related to the current row.

Window functions are useful for performing calculations that require context from other rows, such as ranking, running totals, moving averages, and more. They allow you to perform complex data analysis and calculations within a single SQL query, without the need for self-joins or subqueries.

Some common window functions in SQL include:

1. **`ROW_NUMBER()`**: Assigns a sequential row number to each row within a partition.
2. **`RANK()`**: Assigns a rank to each row within a partition, with ties receiving the same rank.
3. **`DENSE_RANK()`**: Assigns a rank to each row within a partition, with ties receiving the same rank, but without gaps in the ranking sequence.
4. **`NTILE()`**: Divides the rows in a partition into a specified number of groups and assigns a group number to each row.
5. **`LEAD()`** and **`LAG()`**: Access data from the next or previous row within a partition.
6. **`SUM()`**, **`AVG()`**, **`MIN()`**, **`MAX()`**: Perform aggregate calculations over a window of rows.

Window functions are often used in combination with the `OVER` clause, which defines the window or partition over which the function operates. The `OVER` clause can include an `ORDER BY` clause to specify the order in which rows are processed, and a `PARTITION BY` clause to divide the result set into partitions.

Here's an example that uses the `ROW_NUMBER()` window function to assign a sequential row number to each row within a partition:

```sql
SELECT
    product_name,
    category,
    price,
    ROW_NUMBER() OVER (PARTITION BY category ORDER BY price DESC) AS row_num
FROM
    products;
```

This query assigns a row number to each product within its category, ordered by price in descending order. The `PARTITION BY` clause groups the rows by category, and the `ORDER BY` clause specifies the order in which the row numbers are assigned within each partition.

Window functions are powerful tools for data analysis and can significantly simplify complex queries that would otherwise require self-joins or subqueries.

> If I had a sorted query, how would I get values 20 - 50 from it?

To get values 20 - 50 from a sorted query, you can use the `OFFSET` and `FETCH` clauses in SQL. Here's an example:

```sql
SELECT *
FROM (
    SELECT
        column1,
        column2,
        ...
        ROW_NUMBER() OVER (ORDER BY sort_column1, sort_column2, ...) AS row_num
    FROM
        your_table
) AS subquery
WHERE row_num BETWEEN 20 AND 50;
```

Here's how it works:

1. The inner query uses the `ROW_NUMBER()` window function to assign a sequential row number to each row based on the specified sort order (`ORDER BY sort_column1, sort_column2, ...`).
2. The outer query treats the result of the inner query as a subquery and filters the rows based on the `row_num` column.
3. The `WHERE` clause in the outer query selects only the rows where `row_num` is between 20 and 50, effectively returning rows 20 through 50 from the sorted result set.

Alternatively, you can use the `OFFSET` and `FETCH` clauses directly in the query:

```sql
SELECT
    column1,
    column2,
    ...
FROM
    your_table
ORDER BY
    sort_column1,
    sort_column2,
    ...
OFFSET 19 ROWS
FETCH NEXT 31 ROWS ONLY;
```

In this approach:

1. The `ORDER BY` clause sorts the rows based on the specified columns.
2. The `OFFSET 19 ROWS` skips the first 19 rows (rows 1 through 19).
3. The `FETCH NEXT 31 ROWS ONLY` clause retrieves the next 31 rows (rows 20 through 50).

Note that the `OFFSET` and `FETCH` clauses are standard SQL syntax, but some database systems may use different syntax or functions to achieve the same result. For example, in MySQL, you would use `LIMIT` and `OFFSET` like this:

```sql
SELECT
    column1,
    column2,
    ...
FROM
    your_table
ORDER BY
    sort_column1,
    sort_column2,
    ...
LIMIT 31 OFFSET 19;
```

This query retrieves 31 rows (rows 20 through 50) after skipping the first 19 rows.

---

### notes

Non-AI notes

hyperscript update current value of input and get via javascript event
```html
templ ExtraField(val string, field string, itemid uuid.UUID) {
	<input
		id={ fmt.Sprintf("%s_%s", templ.URL(field), itemid.String()) }
		class="text-black"
		value={ val }
		type="text"
		name={ fmt.Sprintf("%s[]", templ.URL(field)) }
		_="on input trigger myExtraFieldEvent(value: (value of me), id: @id)"
	/>
}
```
```javascript
document.addEventListener('myExtraFieldEvent', (event) => {
    console.log("myExtraFieldEvent");
    console.log(event.detail.id);
    console.log(event.detail.value);
});
```

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

```go


	e.POST("/api/mediaQuery", func(c echo.Context) error {
		pageData, err := foodlib.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		if c.FormValue("query") == "(prefers-color-scheme: dark)" && c.FormValue("value") == "light" {
			pageData.Palette = foodlib.Light
			pageData.SavePageData(c)
		} else {
			pageData.Palette = foodlib.Dark
			pageData.SavePageData(c)
		}
		return c.NoContent(http.StatusOK)
	})

```

```templ

script getColorScheme(mediaQueryEndpoint string) {
	// Function to send the media query value to the endpoint
	function sendMediaQueryValue(querystring, value) {
		// Define the endpoint URL
		const endpoint = mediaQueryEndpoint;
		const data = new URLSearchParams();
		data.append('query', querystring);
		data.append('value', value);
		const options = {
			method: 'POST',
			headers: {
				'Content-Type': 'application/x-www-form-urlencoded' // Specify the content type as JSON
			},
			body: data // Convert the data object to a JSON string
		}
		// Send a GET request to the endpoint
		fetch(endpoint, options)
			.then(response => {
				if (response.ok) {
					console.log('Media query value sent successfully');
				} else {
					console.error('Failed to send media query value');
				}
			})
			.catch(error => {
				console.error('Error occurred while sending media query value:', error);
			});
	}
	function checkColorscheme() {
		// Define your media query here
		const querystring = '(prefers-color-scheme: dark)';
		const mediaQuery = window.matchMedia(querystring);
		if (mediaQuery.matches) {
			sendMediaQueryValue(querystring, 'dark');
		} else {
			sendMediaQueryValue(querystring, 'light');
		}
	}

	checkColorscheme()
}
```
