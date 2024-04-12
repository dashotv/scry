package app

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/dashotv/scry/search"
)

func (a *Application) EsIndex(c echo.Context) error {
	return c.JSON(http.StatusOK, H{})
}

func (a *Application) EsMedia(c echo.Context) error {
	m := &search.Media{}
	if err := c.Bind(&m); err != nil {
		return c.JSON(http.StatusBadRequest, &Response{Error: true, Message: err.Error()})
	}

	resp, err := a.Client.IndexMedia(m)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: resp})
}

func (a *Application) EsRelease(c echo.Context) error {
	r := &search.Release{}
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, &Response{Error: true, Message: err.Error()})
	}

	resp, err := a.Client.IndexRelease(r)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: resp})
}
