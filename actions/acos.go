package actions

import "github.com/gobuffalo/buffalo"

// acoIndex default implementation.
func AcosIndex(c buffalo.Context) error {
	return c.Render(200, r.HTML("api/acos/index.html"))
}
