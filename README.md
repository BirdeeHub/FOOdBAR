[![Frontend Masters](https://static.frontendmasters.com/assets/brand/logos/full.png)](https://frontendmasters.com)

This was a companion repo for the [HTMX & Go with ThePrimeagen](https://frontendmasters.com/courses/htmx) course on [Frontend Masters](https://frontendmasters.com).


> in this repo, I have a go and htmx simple app, and a template html site separately. How do I serve the template html site with go and htmx?

Here is one way to serve the template HTML site with Go and htmx:

1. Create a handler in Go that will serve the static HTML files. For example:

```go
func serveStatic(c echo.Context) error {
  path := c.Param("filepath")
  return c.File(path)
}
```

This will serve files from the filesystem.

2. Register the route and path to match your HTML files:

```go 
e.GET("/static/*", serveStatic)
```

3. Put your HTML templates in a folder like `static/templates/`. 

4. Load the index page to start:

```
http://localhost:8080/static/templates/index.html
```

5. Any links to other pages should be prefixed with `/static` to get routed correctly.

6. Use htmx attributes on elements to call backend APIs as needed.

So in summary:

- Serve static files from Go
- Organize HTML templates in a folder
- Load index page to start 
- Prefix links with /static
- Use htmx for AJAX calls to API

This will allow you to serve a static HTML site while still using htmx for dynamic functionality.
