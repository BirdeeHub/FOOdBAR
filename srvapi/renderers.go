package srvapi

import (
	"foodbar/views"
	"foodbar/views/viewutils"
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
	func(*viewutils.TabType, echo.Context, *viewutils.PageData, *viewutils.TabData) error
}

func RenderTab[TR TabRenderer](tr TR, tt *viewutils.TabType, c echo.Context, data *viewutils.PageData, td *viewutils.TabData) error {
	return tr(tt, c, data, td)
}

func TabDeactivateRenderer(tt *viewutils.TabType, c echo.Context, data *viewutils.PageData, td *viewutils.TabData) error {
	if tt.IsActive() {
		tt.SetActive(false)
		for i, v := range data.TabDatas {
			if v.Ttype == *tt {
				data.TabDatas = append(data.TabDatas[:i], data.TabDatas[i+1:]...)
				break
			}
		}
		return HTML(c, http.StatusOK, views.OOBtabButtonToggle(*tt))
	}
	return nil
}

func TabActivateRenderer(tt *viewutils.TabType, c echo.Context, data *viewutils.PageData, td *viewutils.TabData) error {
	if !tt.IsActive() {
		tt.SetActive(true)
		data.TabDatas = append(data.TabDatas, *td)
		HTML(c, http.StatusOK, views.OOBtabViewContainer(*td))
		return HTML(c, http.StatusOK, views.TabButton(*tt))
	} else {
		return HTML(c, http.StatusOK, views.TabButton(*tt))
	}
}

func TabMaximizeRenderer(tt *viewutils.TabType, c echo.Context, data *viewutils.PageData, td *viewutils.TabData) error {
	for _, v := range data.TabDatas {
		if (v.Ttype.IsActive() && v.Ttype != *tt) {
			v.Ttype.SetActive(false)
			HTML(c, http.StatusOK, views.OOBtabButtonToggle(v.Ttype))
		} else if (v.Ttype.IsActive() && v.Ttype == *tt) {
			data.TabDatas = []viewutils.TabData{v}
		}
	}
	if !tt.IsActive() {
		data.TabDatas = []viewutils.TabData{*td}
		tt.SetActive(true)
		HTML(c, http.StatusOK, views.OOBtabButtonToggle(*tt))
	}
	return HTML(c, http.StatusOK, views.TabContainer(*td))
}

