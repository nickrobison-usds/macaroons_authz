package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/nickrobison/cms_authz/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
