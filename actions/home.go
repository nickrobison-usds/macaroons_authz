package actions

import (
	"fmt"

	"github.com/gobuffalo/buffalo"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {

	fmt.Println(c.Session())

	c.Set("user_id", c.Session().Get("user_id"))

	fmt.Println(c.Session().Get("user_id"))
	return c.Render(200, r.HTML("index.html"))
}
