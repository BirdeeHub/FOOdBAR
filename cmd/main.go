package main

import (
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"
	"errors"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func HTML(c echo.Context, code int, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	c.Response().Status = code
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}
var id int = 0
type Contact struct {
	Name  string
	Email string
	Id int
}

type Contacts = []Contact

func (d *Data) hasEmail(email string) bool {
	for _, contact := range d.Contacts {
		if contact.Email == email {
			return true
		}
	}
	return false
}

func newContact(name string, email string) Contact {
	id++
	return Contact{
		Name: name,
		Email: email,
		Id: id,
	}
}

type Data struct {
	Contacts Contacts
}

func (d *Data) indexOf(id int) (int, error) {
	for i, contact := range d.Contacts {
		if contact.Id == id {
			return i, nil
		}
	}
	return 0, errors.New("contact not found")
}

func newData() Data {
	return Data{
		Contacts: []Contact{
			newContact("Foo", "foo@bar"),
			newContact("bar", "beyopnd@bar"),
			newContact("FUBAR", "foo@test"),
		},
	}
}

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

func newFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

type Page struct {
	Data Data
	Form FormData
}

func newPage() Page {
	return Page{
		Data: newData(),
		Form: newFormData(),
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	thePage := newPage()
	e.Renderer = newTemplate()

	e.Static("/images", "images")
	e.Static("/css", "css")

	e.GET("/templtest", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", thePage)
	})

	e.GET("/contactList", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", thePage)
	})

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		if thePage.Data.hasEmail(email) {
			formData := newFormData()
			formData.Values["name"] = name
			formData.Values["email"] = email
			formData.Errors["email"] = "Email already exists"

			return c.Render(http.StatusUnprocessableEntity, "form", formData)
		}

		contact := newContact(name, email)
		thePage.Data.Contacts = append(thePage.Data.Contacts, contact)
		c.Render(http.StatusOK, "form", newFormData())
		return c.Render(http.StatusOK, "oob-contact", contact)
	})

	e.DELETE("/contacts/:id", func(c echo.Context) error {
		time.Sleep(1 * time.Second)
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid id")
		}

		index, err := thePage.Data.indexOf(id)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}
		thePage.Data.Contacts = append(thePage.Data.Contacts[:index], thePage.Data.Contacts[index+1:]...)

		return c.NoContent(http.StatusOK)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
