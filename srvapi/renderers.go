package srvapi

import (
	"FOOdBAR/views"
	"FOOdBAR/views/viewutils"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func HTML(c echo.Context, code int, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	c.Response().Status = code
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

type TabRenderer interface {
	func(echo.Context, *viewutils.PageData, *viewutils.TabData) error
}

func RenderTab[TR TabRenderer](tr TR, c echo.Context, data *viewutils.PageData, td *viewutils.TabData) error {
	return tr(c, data, td)
}

func TabDeactivateRenderer(c echo.Context, data *viewutils.PageData, td *viewutils.TabData) error {
	data.SetActive(td, false)
	err := data.SavePageData(c)
	if err != nil {
		echo.NewHTTPError(http.StatusTeapot, "Cannot unmarshal page data")
	}
	return HTML(c, http.StatusOK, views.OOBtabButtonToggle(viewutils.TabButtonData{Ttype: td.Ttype, Active: false}))
}

func TabActivateRenderer(c echo.Context, data *viewutils.PageData, td *viewutils.TabData) error {
	data.SetActive(td, true)
	err := data.SavePageData(c)
	if err != nil {
		echo.NewHTTPError(http.StatusTeapot, "Cannot unmarshal page data")
	}
	HTML(c, http.StatusOK, views.OOBtabViewContainer(td))
	return HTML(c, http.StatusOK, views.TabButton(viewutils.TabButtonData{Ttype: td.Ttype, Active: true}))
}

func TabMaximizeRenderer(c echo.Context, data *viewutils.PageData, td *viewutils.TabData) error {
	var toMin []*viewutils.TabData
	data.SetActive(td, true)
	for _, v := range data.TabDatas {
		if (v.Ttype != td.Ttype) {
			toMin = append(toMin, v)
		}
	}
	for _, v := range toMin {
		data.SetActive(v, false)
	}
	err := data.SavePageData(c)
	if err != nil {
		echo.NewHTTPError(http.StatusTeapot, "Cannot unmarshal page data")
	}
	for _, v := range toMin {
		HTML(c, http.StatusOK, views.OOBtabButtonToggle(viewutils.TabButtonData{Ttype: v.Ttype, Active: false}))
	}
	HTML(c, http.StatusOK, views.OOBtabButtonToggle(viewutils.TabButtonData{Ttype: td.Ttype, Active: true}))
	return HTML(c, http.StatusOK, views.TabContainer(td))
}

