package macaroons

import (
	"context"
	"fmt"
)

type ContextCheck struct {
	Key string
}

func (c ContextCheck) StrCheck(ctx context.Context, cond, args string) error {
	expect, _ := ctx.Value(c.Key).(string)
	if args != expect {
		return fmt.Errorf("%s doesn't match %s", cond, expect)
	}
	return nil
}
