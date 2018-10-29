package actions

import "github.com/gobuffalo/buffalo"

func RouteHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("routes.html"))
}
