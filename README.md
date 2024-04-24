# This repo is an extreme work in progress, it is NOT EVEN CLOSE TO FINISHED!!!!

## This is my first attempt at making a webpage,
## It is also my first time using go, htmx, sql, js, html, css, templ, authentication, routing, etc.

This program will be a tracker app for a friend's small catering/personal chef buisness when it is finished.

Mostly though, it is a way for me to learn web development as I have not done any before.

There will likely be some spaghetti that I will figure out how to handle better as I go.

### Install and development building instructions

For now it is packaged for nix only so install nix package manager, then clone the repo and cd into it

Then to build or run use ```nix build --show-trace``` or ```nix run --show-trace```

and to hack around with it with hot reload run the following 2 commands:

```bash
nix develop --show-trace
# then inside the shell:
air
```
Then go to localhost:42069 in a web browser

If you dont want to install nix, the following should work, and the air command might work, assuming you have go, templ, and sqlite3 installed.
```bash
go mod tidy && templ generate && go build -o dist/FOOdBAR ./FOOdBAR/main.go
```
-OR-
```bash
air
```

### Name:

It's FU'd Beyond Any Recognition

The FOOdBAR tracks the food in the FOOdb and tells you how much FOOd you have and sold and what it cost,
and reminds you of recipes that fit those criteria best. Allowing you to cook more variety with less headache

---

### Done so far:

setup go, htmx, templ, echo project

made it be buildable by nix, and included a dev shell

defined the basic api routes and templates needed for the tab skeleton

set up authentication via JWT and auth database

defined skeleton of runtime state type objects

defined part of the skeleton for the tabs themselves

---

### TODO AND IDEAS:
database for tab items

expandable elements for tab items

Would be nice to have dynamic resizing with the mouse

infinite scroll for all list style tabs

write some damn tests

make it http*S*

recipe list and entry

suggestion based on stock

sort recipe by category, dietary, ingredients, % ingredients in stock, (total ingredient price?)

generate prep list, and shopping list (accounts for current stock), the menu, and projected and actual portions ordered.

track what was bought and then used after the day to keep track of the total stock, as actual amounts may differ based on prices for different quantities at stores.

total workflow, 3 interaction stages. Generate menu, input orders and generate the lists, input actual stock bought and used to correct suggested values.

6 tab interface? recipe database, menu tab where you have the menu and track orders, prep list, shopping list, stock list, profit review

recipes will need category, dietary, ingredients, estimated cost, fields, ingredients should have an optional scaling factor for scaling recipes accurately

ingredients will need amounts in stock, time since purchased, storage method, (last queried price?)

when all data for stock and menu and orders have been finalized for the week, generate profit report, with option for adding extra incidental costs.

Run the program, and then visit localhost:42069 (or other configured value)

styling/animation?

recipe entry should have completion for existing items in ingredients dietary, category, etc but not for name or instructions or amounts or prices.

shopping price search feature?

---

# AI VOMIT

Below is stuff that is probably mostly random AI garbage but I might take 1 or 2 useful things from it still to look up further

---

### SQL


> can tables in sqlite contain pointers to other tables?
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


---

### caching

> So, I have a PageData struct type, and i would like to create a new one for a user upon login, and send it to them to save in session storage. I would then like to get it from context on every request from the client. Is this possible if so how
> 
> ```go
> 	e.GET("", func(c echo.Context) error {
> 		userID, err := GetUserFromToken(c)
> 		if err != nil {
> 			return echo.NewHTTPError(http.StatusUnauthorized, err)
> 		}
> 		c.Logger().Print(c)
> 		return HTML(c, http.StatusOK, views.Homepage(viewutils.GetPageData(userID)))
> 	})
> 
> 	e.POST("", func(c echo.Context) error {
> 		userID, err := GetUserFromToken(c)
> 		if err != nil {
> 			return echo.NewHTTPError(http.StatusUnauthorized, err)
> 		}
> 		c.Logger().Print(c)
> 		return HTML(c, http.StatusOK, views.Homepage(viewutils.GetPageData(userID)))
> 	})
> 
> 	e.POST("/api/mediaQuery", func(c echo.Context) error {
> 		userID, err := GetUserFromToken(c)
> 		if err != nil {
> 			return c.NoContent(http.StatusOK)
> 		}
> 		pageData := viewutils.GetPageData(userID)
> 		if c.FormValue("query") == "(prefers-color-scheme: dark)" && c.FormValue("value") == "dark" {
> 			pageData.Palette = viewutils.Dark
> 		} else {
> 			pageData.Palette = viewutils.Light
> 		}
> 		return c.NoContent(http.StatusOK)
> 	})
> 
> 	e.DELETE("/api/tabButton/deactivate/:type", func(c echo.Context) error {
> 		userID, err := GetUserFromToken(c)
> 		if err != nil {
> 			return echo.NewHTTPError(http.StatusUnauthorized, err)
> 		}
> 		c.Logger().Print(c)
> 		tt, err := viewutils.String2TabType(c.Param("type"))
> 		if err != nil {
> 			return echo.NewHTTPError(
> 				http.StatusInternalServerError,
> 				errors.New("not a valid tab type"),
> 			)
> 		}
> 		pageData := viewutils.GetPageData(userID)
> 		tabdata, err := pageData.GetTabDataByType(*tt)
> 		return RenderTab(TabDeactivateRenderer, c, pageData, tabdata)
> 	})
> 
> 	e.GET("/api/tabButton/activate/:type", func(c echo.Context) error {
> 		userID, err := GetUserFromToken(c)
> 		if err != nil {
> 			return echo.NewHTTPError(http.StatusUnauthorized, err)
> 		}
> 		c.Logger().Print(c)
> 		tt, err := viewutils.String2TabType(c.Param("type"))
> 		if err != nil {
> 			return echo.NewHTTPError(
> 				http.StatusInternalServerError,
> 				errors.New("not a valid tab type"),
> 			)
> 		}
> 
> 		// TODO: fetch appropriate TabData.Items from database
> 		// based on sort. Implement infinite scroll for them.
> 		pageData := viewutils.GetPageData(userID)
> 		tabdata, err := pageData.GetTabDataByType(*tt)
> 		return RenderTab(TabActivateRenderer, c, pageData, tabdata)
> 	})
> 
> 	e.POST("/api/tabButton/maximize/:type", func(c echo.Context) error {
> 		userID, err := GetUserFromToken(c)
> 		if err != nil {
> 			return echo.NewHTTPError(http.StatusUnauthorized, err)
> 		}
> 		c.Logger().Print(c)
> 		tt, err := viewutils.String2TabType(c.Param("type"))
> 		if err != nil {
> 			return echo.NewHTTPError(
> 				http.StatusInternalServerError,
> 				errors.New("not a valid tab type"),
> 			)
> 		}
> 
> 		// TODO: fetch appropriate TabData.Items from database
> 		// based on sort. Implement infinite scroll for them.
> 		pageData := viewutils.GetPageData(userID)
> 		tabdata, err := pageData.GetTabDataByType(*tt)
> 		return RenderTab(TabMaximizeRenderer, c, pageData, tabdata)
> 	})
> 	return nil
> }
> ```

Yes, it is possible to create a new PageData struct for the user on login and send it to the client to save in sessionStorage. Here is one way to do it:

1. When handling the login route, create the PageData:

```go
func HandleLogin(c echo.Context) error {

  // login logic

  userID := // get user ID

  pageData := viewutils.InitPageData(userID)

  // save pageData in session
  session := sessions.Default(c)
  session.Set("pageData", pageData)
  session.Save()

  return c.JSON(http.StatusOK, pageData)
}
```

2. The client saves the PageData in sessionStorage:

```js
// in login handler 
fetch('/login', {
  // login request 
}).then(res => res.json())
  .then(pageData => {
    sessionStorage.setItem('pageData', JSON.stringify(pageData))
  }) 
```

3. On subsequent requests, get PageData from session and attach to context:

```go
func AttachPageData(next echo.HandlerFunc) echo.HandlerFunc {
  return func(c echo.Context) error {

    pageDataJSON := session.Get("pageData")
    if pageDataJSON != nil {
      var pageData PageData 
      // decode pageDataJSON to pageData struct

      c.Set("pageData", pageData)
    }

    return next(c)
  }
}

// Usage
e.Use(AttachPageData) 

e.GET("/someRoute", func(c echo.Context) error {
  pageData := c.Get("pageData").(*PageData)
  
  // use pageData
})
```

So the PageData can be accessed from context on subsequent requests after login.

> func HandleLogin(c echo.Context) error {
> 
>   // login logic
> 
>   userID := // get user ID
> 
>   pageData := viewutils.InitPageData(userID)
> 
>   // save pageData in session
>   session := sessions.Default(c)
>   session.Set("pageData", pageData)
>   session.Save()
> 
>   return c.JSON(http.StatusOK, pageData)
> }
> ```
> 
> 2. The client saves the PageData in sessionStorage:
> 
> ```js
> // in login handler 
> fetch('/login', {
>   // login request 
> }).then(res => res.json())
>   .then(pageData => {
>     sessionStorage.setItem('pageData', JSON.stringify(pageData))
>   }) 
> 
> 
> 
>   why do you need to do the save on the client if I instruct it to save via echo?

Good catch! You are right, saving the PageData to the session on the server side with `session.Save()` is sufficient. 

The client side sessionStorage saving is redundant in this case. The server-side session middleware will handle loading the PageData from session and attaching it to context on subsequent requests.

So you can remove the client-side sessionStorage code. The server-side code should be enough:

```go
func HandleLogin(c echo.Context) error {

  // login logic
  
  userID := // get user ID

  pageData := viewutils.InitPageData(userID)

  // save in session
  session := sessions.Default(c) 
  session.Set("pageData", pageData)
  session.Save()

  return c.NoContent(http.StatusOK) 
}

// AttachPageData middleware loads from session 

e.Use(AttachPageData)
```

The AttachPageData middleware will load the PageData from session on each request and attach it to context.

So the server-side session handling is sufficient, no need for the client to also save to sessionStorage. Good catch!





# Absolute bewilderment




> when I open the tabs in the order in the first line of these, and then maximize the one named in the last line, as you can see in the second line, an extra one does not get removed from the list.
> 
> ```markdown
> 
> ```json
> {"time":"2024-04-24T02:42:35.296028019-07:00","level":"-","prefix":"echo","file":"types.go","line":"201","message":"{\"user_id\":\"8a6d4aab-aff2-484d-967d-f56040c25a1d\",\"tab_datas\":[{\"tab_type\":\"Menu\",\"items\":{},\"order_by\":{}},{\"tab_type\":\"Pantry\",\"items\":{},\"order_by\":{}},{\"tab_type\":\"Recipe\",\"items\":{},\"order_by\":{}},{\"tab_type\":\"Preplist\",\"items\":{},\"order_by\":{}},{\"tab_type\":\"Shopping\",\"items\":{},\"order_by\":{}}],\"palette\":\"dark\"}"}
> {"time":"2024-04-24T02:42:35.296150272-07:00","level":"-","prefix":"echo","file":"types.go","line":"217","message":"{\"user_id\":\"8a6d4aab-aff2-484d-967d-f56040c25a1d\",\"tab_datas\":[{\"tab_type\":\"Pantry\",\"items\":{},\"order_by\":{}},{\"tab_type\":\"Preplist\",\"items\":{},\"order_by\":{}}],\"palette\":\"dark\"}"}
> {"time":"2024-04-24T02:42:35.296355367-07:00","id":"","remote_ip":"127.0.0.1","host":"localhost:42069","method":"POST","uri":"/FOOdBAR/api/tabButton/maximize/Pantry","user_agent":"Mozilla/5.0 (X11; Linux x86_64; rv:125.0) Gecko/20100101 Firefox/125.0","status":200,"error":"","latency":363616,"latency_human":"363.616µs","bytes_in":0,"bytes_out":2838}
> ```
> 
> 
> ```json
> {"time":"2024-04-24T02:56:24.905488188-07:00","level":"-","prefix":"echo","file":"types.go","line":"201","message":"{\"user_id\":\"8a6d4aab-aff2-484d-967d-f56040c25a1d\",\"tab_datas\":[{\"tab_type\":\"Preplist\",\"items\":{},\"order_by\":{}},{\"tab_type\":\"Earnings\",\"items\":{},\"order_by\":{}},{\"tab_type\":\"Shopping\",\"items\":{},\"order_by\":{}},{\"tab_type\":\"Pantry\",\"items\":{},\"order_by\":{}},{\"tab_type\":\"Recipe\",\"items\":{},\"order_by\":{}}],\"palette\":\"dark\"}"}
> {"time":"2024-04-24T02:56:24.905591738-07:00","level":"-","prefix":"echo","file":"types.go","line":"217","message":"{\"user_id\":\"8a6d4aab-aff2-484d-967d-f56040c25a1d\",\"tab_datas\":[{\"tab_type\":\"Preplist\",\"items\":{},\"order_by\":{}},{\"tab_type\":\"Shopping\",\"items\":{},\"order_by\":{}}],\"palette\":\"dark\"}"}
> {"time":"2024-04-24T02:56:24.90569925-07:00","id":"","remote_ip":"127.0.0.1","host":"localhost:42069","method":"POST","uri":"/FOOdBAR/api/tabButton/maximize/Preplist","user_agent":"Mozilla/5.0 (X11; Linux x86_64; rv:125.0) Gecko/20100101 Firefox/125.0","status":200,"error":"","latency":232139,"latency_human":"232.139µs","bytes_in":0,"bytes_out":2653}
> ```
> ```

Based on the logs you provided, it looks like there may be an issue with the tab maximization logic not properly removing tabs when maximizing. A few things to check:

- In the TabMaximizeRenderer function, it is looping through all tabs and setting them to inactive except the one being maximized. Verify this logic is working as expected.

- When a tab is set to inactive, check that it is properly being removed from the tab data list. It may still be in the list but just set as inactive.

- Double check that when rendering the tab buttons, it is only rendering the active tabs in the list. An inactive tab could still be rendered if the logic checks are not correct.

- Add some additional debugging logs or breakpoints to trace the flow and see where a tab that should be removed is still persisting.

- Review the tab data structure and maximize logic end-to-end to see if anything is missed.

The key things to verify are 1) inactive tabs are properly removed from the main list, and 2) only active tabs are rendered. The issue likely lies in one of those areas. Let me know if any of those suggestions help uncover the problem!

> the inactive tab is not being properly removed from the list, but the extra ones are not rendered still
> 
> It only happens in particular orders of the tabs being opened and only the maximise function triggers it, but in those particular orders it is consistent in its misbehavior

Thanks for the additional details. Since the inactive tabs are not being properly removed from the list, that is likely the core issue. 

A few things to check on that:

- When setting a tab inactive, verify it is removed from the main tab data list in addition to just setting the active flag to false.

- Make sure the logic that removes inactive tabs from the list is called unconditionally whenever a tab is maximized.

- Double check the logic for how the main tab data list is filtered/copied when getting the list of active tabs to render. It may be preserving inactive ones in some cases.

- Add logging to trace the full lifecycle of the tab data structure during maximize and see where inactive tabs are still persisting.

- Consider changing the tab data structure to store active and inactive tabs separately to avoid mixups.

Overall it sounds like the root cause is the inactive tabs not being fully removed from the main list when they should be. Focus debugging efforts on where that removal logic is failing or not being called properly during maximize. Let me know if any other info would help narrow down the issue!

> I do not save the inactive ones at all. This is the maximise function. I cant find the problem.
> 
> func TabMaximizeRenderer(c echo.Context, data *viewutils.PageData, td *viewutils.TabData) error {
>     var toMin []viewutils.TabType
>     data.SetActive(td, true)
>     for _, v := range data.TabDatas {
>         if (v.Ttype != td.Ttype) {
>             data.SetActive(v, false)
>             toMin = append(toMin, v.Ttype)
>         }
>     }
>     err := data.SavePageData(c)
>     if err != nil {
>         echo.NewHTTPError(http.StatusTeapot, "Cannot unmarshal page data")
>     }
>     for _, v := range toMin {
>         HTML(c, http.StatusOK, views.OOBtabButtonToggle(viewutils.TabButtonData{Ttype: v, Active: false}))
>     }
>     HTML(c, http.StatusOK, views.OOBtabButtonToggle(viewutils.TabButtonData{Ttype: td.Ttype, Active: true}))
>     return HTML(c, http.StatusOK, views.TabContainer(td))
> }
> 
>

Looking at the TabMaximizeRenderer code, a few things stand out:

- It is correctly setting inactive tabs and collecting them into the toMin list.

- It is saving the updated page data after setting tabs inactive.

- It is rendering the toggle for the newly inactive tabs.

So the issue seems to be happening somewhere else in the lifecycle. A few things I would try:

- Add logging in PageData.SetActive() to trace when a tab is set inactive. Verify it is being removed from the main list there.

- Log the full PageData.TabDatas list before and after maximize to see if inactive tabs persist.

- Double check PageData.SavePageData() - verify it is saving the updated list properly.

- Check where the initial PageData is loaded from - could inactive tabs be reconstituted from a cache or prior data?

- Trace the active tab list generation logic - could that be preserving inactive tabs?

My suspicion is the root cause is either in SetActive() not fully removing tabs, or cached data reintroducing inactive tabs. Tracing the full lifecycle around those areas may reveal where the problem is introduced.

Let me know if any other ideas come to mind on how to track this down! The key is tracing the inactive tab handling thoroughly.

> It seems like it is some pointer bs because I dont save anything anywhere, I get the pageData from the cookie and I send it back and then I get it again at every request. At first I thought it was a race condition but it seems to depend on the order of the tabs and not the speed at which I press buttons.

You're right, since the page data is loaded from the cookie on each request, it does sound like some kind of pointer/reference issue where inactive tabs are persisting in memory incorrectly.

A few suggestions on tracking this down:

- Add logging/breakpoints in PageData.SetActive() and trace what is happening to the tab data when a tab is set inactive. Verify the slice itself is modified.

- You could try cloning the tab data when loading from the cookie rather than using the slice directly. This may break any problematic references.

- Iterate through the tab data and print memory addresses before and after SetActive() - see if you notice any strange pointer behaviors.

- Consider changing the data structure to use a map rather than a slice, where you delete entries on deactivate. May avoid pointer issues.

- Add stress testing with rapid tab switches - may help uncover race conditions.

This does seem like a tricky one to pin down. The key will be tracing the memory addresses and references to try to catch where inactive tabs are persisting. A data structure change could help as well. Let me know if any other ideas come to mind!
