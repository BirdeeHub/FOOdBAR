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
	if td.IsActive() {
		td.SetActive(false)
		for i, v := range data.TabDatas {
			if v.Ttype == td.Ttype {
				data.TabDatas[i].SetActive(false)
				break
			}
		}
		data.SavePageData(c)
		return HTML(c, http.StatusOK, views.OOBtabButtonToggle(td))
	}
	return nil
}

func TabActivateRenderer(c echo.Context, data *viewutils.PageData, td *viewutils.TabData) error {
	if !td.IsActive() {
		td.SetActive(true)
		data.SavePageData(c)
		HTML(c, http.StatusOK, views.OOBtabViewContainer(td))
		return HTML(c, http.StatusOK, views.TabButton(td))
	} else {
		return HTML(c, http.StatusOK, views.TabButton(td))
	}
}

func TabMaximizeRenderer(c echo.Context, data *viewutils.PageData, td *viewutils.TabData) error {
	var toMin []int
	for i, v := range data.TabDatas {
		if (v.IsActive() && v.Ttype != td.Ttype) {
			v.SetActive(false)
			toMin = append(toMin, i)
		}
	}
	if !td.IsActive() {
		td.SetActive(true)
	}
	data.SavePageData(c)
	for _, i := range toMin {
		HTML(c, http.StatusOK, views.OOBtabButtonToggle(data.TabDatas[i]))
	}
	HTML(c, http.StatusOK, views.OOBtabButtonToggle(td))
	return HTML(c, http.StatusOK, views.TabContainer(td))
}

