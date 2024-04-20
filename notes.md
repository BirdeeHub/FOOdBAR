Page and tab data should be used to REQUEST from database.

They should be stored in memory and their items should be fetched, which contain just uuid, name and tab type

which should then be filled in by data queried from db based on the items.

tabtype states should be stored by user

features

TODO:
database
efficient infinite scroll for all list style tabs
all the list tabs

earnings tab
shopping price search feature?
styling/animation?


recipe list and entry
suggestion based on stock
sort recipe by category, dietary, ingredients, % ingredients in stock, (total ingredient price?)
finalize menu.

generate prep list, and shopping list (accounts for current stock), the menu, and projected and actual portions ordered.

track what was bought and then used after the day to keep track of the total stock, as actual amounts may differ based on prices for different quantities at stores.

total workflow, 3 interaction stages. Generate menu, input orders and generate the lists, input actual stock bought and used to correct suggested values.

6 tab interface? recipe database, menu tab where you have the menu and track orders, prep list, shopping list, stock list, profit review

Would be nice to have these able to be dynamically resized and put side by side or above and below.

persistent data consists of recipes, and stock, and configuration.

Will need to store recipe objects, and ingredient objects.

recipes will need category dietary, ingredients, last price sold for, fields

ingredients will need amounts in stock, time since purchased, storage method, (last queried price?)

when all data for stock and menu and orders have been finalized for the week, generate profit report, with option for adding extra incidental costs.


Interface will be written as a web page, because I want to learn how to make an interactive web page, and I will be using go + htmx to do it.

Run the program, and then visit localhost:42069 (or other configured value)

recipe entry should have completion for existing items in ingredients dietary, category, etc but not for name or instructions or amounts or prices.





# AI SPAM








To get the result of a media query for `prefers-color-scheme` in JavaScript and post it to an endpoint, you can follow these steps:

1. Use `window.matchMedia()` to check if the user's preferred color scheme matches a specific media query.
2. Extract the result of the media query.
3. Send the result to an endpoint using an HTTP request. You can use `fetch()` to perform this.

Here's a basic example:

```javascript
// Check if the user prefers dark mode
const darkModeQuery = window.matchMedia('(prefers-color-scheme: dark)');
const isDarkMode = darkModeQuery.matches;

// Define the data to be sent
const data = {
  isDarkMode: isDarkMode
};

// Define the endpoint URL
const endpoint = 'https://example.com/endpoint';

// Define the fetch options
const options = {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify(data)
};

// Send the data to the endpoint
fetch(endpoint, options)
  .then(response => {
    if (response.ok) {
      console.log('Data successfully sent to the endpoint');
    } else {
      console.error('Failed to send data to the endpoint');
    }
  })
  .catch(error => {
    console.error('Error:', error);
  });
```

In this example:
- We use `window.matchMedia('(prefers-color-scheme: dark)')` to check if the user prefers dark mode.
- The `matches` property of the media query result indicates whether the condition is met.
- We define the data to be sent as an object with a key `isDarkMode` containing the result of the media query.
- We define the endpoint URL where we want to send the data.
- We use `fetch()` to make a POST request to the endpoint with the data.
- We handle the response from the endpoint to log success or failure.

Ensure to replace `'https://example.com/endpoint'` with the actual endpoint URL where you want to send the data. Also, handle errors and responses appropriately based on your requirements.


---

# sqlite


Using SQLite in Go is relatively straightforward thanks to the `database/sql` package provided by the Go standard library and the `github.com/mattn/go-sqlite3` package, which serves as a SQLite driver for Go.

Here's a simple step-by-step guide on how to use SQLite in Go:

1. **Install SQLite3**: Ensure that SQLite3 is installed on your system. You can download it from the SQLite website if it's not already installed.

2. **Install the SQLite Driver for Go**: You can use the `go get` command to install the SQLite driver package:
   ```
   go get github.com/mattn/go-sqlite3
   ```

3. **Import Required Packages**: Import the necessary packages in your Go file.
   ```go
   import (
       "database/sql"
       _ "github.com/mattn/go-sqlite3"
   )
   ```

4. **Open a Database Connection**: Use the `sql.Open` function to establish a connection to your SQLite database.
   ```go
   db, err := sql.Open("sqlite3", "path/to/your/database.db")
   if err != nil {
       log.Fatal(err)
   }
   defer db.Close()
   ```

5. **Create Tables (if needed)**: You can execute SQL statements to create tables if they don't exist already.
   ```go
   _, err = db.Exec(`
       CREATE TABLE IF NOT EXISTS users (
           id INTEGER PRIMARY KEY AUTOINCREMENT,
           name TEXT,
           email TEXT
       )
   `)
   if err != nil {
       log.Fatal(err)
   }
   ```

6. **Perform Database Operations**: You can perform various database operations like insert, update, delete, and select using `db.Exec` and `db.Query` methods.
   ```go
   // Insert data
   _, err = db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", "John Doe", "john@example.com")
   if err != nil {
       log.Fatal(err)
   }

   // Query data
   rows, err := db.Query("SELECT id, name, email FROM users")
   if err != nil {
       log.Fatal(err)
   }
   defer rows.Close()

   for rows.Next() {
       var id int
       var name, email string
       err := rows.Scan(&id, &name, &email)
       if err != nil {
           log.Fatal(err)
       }
       fmt.Println(id, name, email)
   }
   ```

7. **Handle Errors**: Always check and handle errors properly to ensure your application behaves as expected.

8. **Close the Connection**: Don't forget to close the database connection when you're done using it.

That's it! You now have a basic understanding of how to use SQLite in Go. You can build upon this foundation to create more complex database-driven applications.



> I want my database to contain user tables, containing a table for each tab type. right now it just has example stuff and what I think is how you initialize a database?
> 
> ```go
> func readTab(tt viewutils.TabType, userID uuid.UUID) viewutils.TabData {
> 	db, err := sql.Open("sqlite3", "path/to/your/database.db")
> 	if err != nil {
> 		log.Fatal(err)
> 	}
> 	defer db.Close()
> }
> ```

Here is how you could structure the database to have separate tables for each tab type per user:

1. Create a users table to store user info:

```go
CREATE TABLE users (
  id INTEGER PRIMARY KEY, 
  name TEXT,
  email TEXT
);
```

2. Create separate tables for each tab type, with a user id column to associate with a user:

```go 
CREATE TABLE user_recipes (
  id INTEGER PRIMARY KEY,
  user_id INTEGER, 
  name TEXT,
  // other recipe columns
);

CREATE TABLE user_pantry (
  id INTEGER PRIMARY KEY, 
  user_id INTEGER,
  name TEXT,
  // other pantry columns
);

// Similarly for other tab types - shopping, menu, etc
```

3. When reading/writing tab data, lookup the user id first, then query the appropriate table:

```go
func readRecipes(userID uuid.UUID) {

  // Lookup user id
  row := db.QueryRow("SELECT id FROM users WHERE id = ?", userID)
  var userID int 
  row.Scan(&userID)

  // Get recipes for user
  rows, _ := db.Query("SELECT * FROM user_recipes WHERE user_id = ?", userID)
  
  // Process rows
}

// Similarly for other tab types
```

So in summary:
- Separate tables per tab type associated with user id
- Lookup user id first 
- Query appropriate table for that user

Let me know if this helps or if you need any clarification!

> can tables in sqlite contain pointers to other tables?
> 
> Can I create a separate set of tables for each user instead of having the userdata just be rows in a user_pantry table?

Yes, SQLite supports foreign keys which allow you to reference rows in another table. This can be used to create separate tables per user.

Here is an example schema:

```sql
CREATE TABLE users (
  id INTEGER PRIMARY KEY,
  name TEXT
);

CREATE TABLE user_1_recipes (
  id INTEGER PRIMARY KEY, 
  recipe_name TEXT,
  user_id INTEGER REFERENCES users(id) 
);

CREATE TABLE user_2_recipes (
  id INTEGER PRIMARY KEY,
  recipe_name TEXT,
  user_id INTEGER REFERENCES users(id)
);

CREATE TABLE user_1_pantry (
  id INTEGER PRIMARY KEY,
  ingredient TEXT,
  user_id INTEGER REFERENCES users(id)
); 

CREATE TABLE user_2_pantry (
  id INTEGER PRIMARY KEY,
  ingredient TEXT, 
  user_id INTEGER REFERENCES users(id)
);
```

This creates a separate recipes and pantry table for each user, with a `user_id` column that references the `id` column in the `users` table.

When inserting data, you would lookup the user id, then insert into the appropriate per-user table.

The main downside is the table structure is not dynamic - you have to know all the user ids ahead of time to create the tables. An alternative is a single table per tab type with a `user_id` column.

So in summary, yes SQLite supports separate per-user tables through foreign keys, but it requires knowing user ids upfront. Let me know if you have any other questions!

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


> how do I sanitize a string before I use it in a query in go


Sanitizing strings before using them in a database query in Go is crucial to prevent SQL injection attacks. You should use parameterized queries or prepared statements instead of directly concatenating strings into your SQL queries. Most SQL database libraries for Go support this feature.

Here's an example using the `database/sql` package along with `github.com/lib/pq` for PostgreSQL:

```go
package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

func main() {
    // Connect to the database
    db, err := sql.Open("postgres", "user=username dbname=mydb sslmode=disable")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    // Example input string
    userInput := "some input"

    // Use a prepared statement to execute the query safely
    stmt, err := db.Prepare("SELECT * FROM my_table WHERE my_column = $1")
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    // Execute the prepared statement with the sanitized input
    rows, err := stmt.Query(userInput)
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    // Process the query results
    for rows.Next() {
        var id int
        var data string
        if err := rows.Scan(&id, &data); err != nil {
            panic(err)
        }
        fmt.Printf("ID: %d, Data: %s\n", id, data)
    }
    if err := rows.Err(); err != nil {
        panic(err)
    }
}
```

In this example, we're using a prepared statement with the `$1` placeholder to safely inject the user input into the SQL query. This way, the database driver handles the sanitization of the input for you, preventing SQL injection attacks.




Yes, SQLite also supports parameterized queries, which can help prevent SQL injection attacks. Here's how you can use parameterized queries with SQLite in Go:

```go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Connect to the SQLite database
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Example input string
	userInput := "some input"

	// Use a prepared statement to execute the query safely
	stmt, err := db.Prepare("SELECT * FROM my_table WHERE my_column = ?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	// Execute the prepared statement with the sanitized input
	rows, err := stmt.Query(userInput)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Process the query results
	for rows.Next() {
		var id int
		var data string
		if err := rows.Scan(&id, &data); err != nil {
			panic(err)
		}
		fmt.Printf("ID: %d, Data: %s\n", id, data)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
}
```

In this example, we're using the `?` placeholder as a parameter in the SQL query. When you call `Query`, you pass the user input as an argument, and SQLite automatically handles the proper escaping and sanitization of the input, ensuring it's treated as data rather than part of the SQL command.


