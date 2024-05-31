package srvapi

import (
	"FOOdBAR/views"
	foodlib "FOOdBAR/FOOlib"
	"FOOdBAR/db"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func HTML(c echo.Context, code int, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	c.Response().Status = code
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func GZscript(c echo.Context, code int, filebytes []byte) error {
	c.Response().Header().Set(echo.HeaderContentEncoding, "gzip")
	c.Response().Header().Set("Content-Encoding", "gzip")
	c.Response().Status = code
	return c.Blob(code, "application/javascript", filebytes)
}

type TabRenderer interface {
	func(echo.Context, *foodlib.PageData, *foodlib.TabData) error
}

func RenderTab[TR TabRenderer](tr TR, c echo.Context, data *foodlib.PageData, td *foodlib.TabData) error {
	return tr(c, data, td)
}

func TabDeactivateRenderer(c echo.Context, data *foodlib.PageData, td *foodlib.TabData) error {
	data.SetActive(td, false)
	err := db.SavePageData(c, data)
	if err != nil {
		echo.NewHTTPError(http.StatusTeapot, "Cannot unmarshal page data")
	}
	var tt foodlib.TabType
	if td == nil {
		tt = foodlib.Invalid
	} else {
		tt = td.Ttype
	}
	return HTML(c, http.StatusOK, views.OOBtabButtonToggle(foodlib.TabButtonData{Ttype: tt, Active: false}))
}

func TabActivateRenderer(c echo.Context, data *foodlib.PageData, td *foodlib.TabData) error {
	data.SetActive(td, true)
	err := db.SavePageData(c, data)
	if err != nil {
		echo.NewHTTPError(http.StatusTeapot, "Cannot unmarshal page data")
	}
	HTML(c, http.StatusOK, views.OOBtabViewContainer(data, td))
	var tt foodlib.TabType
	if td == nil {
		tt = foodlib.Invalid
	} else {
		tt = td.Ttype
	}
	return HTML(c, http.StatusOK, views.TabButton(foodlib.TabButtonData{Ttype: tt, Active: true}))
}

func TabMaximizeRenderer(c echo.Context, data *foodlib.PageData, td *foodlib.TabData) error {
	var toMin []*foodlib.TabData
	data.SetActive(td, true)
	if td != nil {
		for _, v := range data.TabDatas {
			if v.Ttype != td.Ttype {
				toMin = append(toMin, v)
			}
		}
	}
	for _, v := range toMin {
		data.SetActive(v, false)
	}
	err := db.SavePageData(c, data)
	if err != nil {
		echo.NewHTTPError(http.StatusTeapot, "Cannot unmarshal page data")
	}
	for _, v := range toMin {
		HTML(c, http.StatusOK, views.OOBtabButtonToggle(foodlib.TabButtonData{Ttype: v.Ttype, Active: false}))
	}
	var tt foodlib.TabType
	if td == nil {
		tt = foodlib.Invalid
	} else {
		tt = td.Ttype
	}
	HTML(c, http.StatusOK, views.OOBtabButtonToggle(foodlib.TabButtonData{Ttype: tt, Active: true}))
	return HTML(c, http.StatusOK, views.TabContainer(data, td))
}
