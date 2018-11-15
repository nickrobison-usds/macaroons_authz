package macaroons

import (
	"context"
	"fmt"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
)

type ContextCheck struct {
	Key string
}

func (c ContextCheck) StrCheck(ctx context.Context, cond, args string) error {
	expect, _ := ctx.Value(c.Key).(string)
	return strCmp(cond, expect, args)
}

type CMSAssociationCheck struct {
	ContextKey       string
	AssociationTable string
	DB               *pop.Connection
}

func (c CMSAssociationCheck) Check(ctx context.Context, cond, args string) error {

	// Get the context value
	associationID, _ := ctx.Value(c.ContextKey).(string)

	var result uuid.UUID

	queryString := fmt.Sprintf("SELECT user_id from %s WHERE %s = ? and user_id = ?", c.AssociationTable, c.ContextKey)
	err := c.DB.RawQuery(queryString, associationID, args).First(&result)
	if err != nil {
		return err
	}

	return strCmp(cond, result.String(), args)

}

func strCmp(cond, expect, args string) error {
	if args != expect {
		return fmt.Errorf("%s doesn't match %s", cond, expect)
	}
	return nil
}
