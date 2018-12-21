package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/nickrobison-usds/macaroons_authz/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
