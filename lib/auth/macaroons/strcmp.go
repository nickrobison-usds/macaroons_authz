package macaroons

import (
	"fmt"

	"gopkg.in/macaroon-bakery.v1/bakery/checkers"
)

// StrcmpChecker verifies that a string caveat is satisfied
type StrcmpChecker string

// CheckFirstPartyCaveat verifies that a first party, string caveat is satisfied
func (c StrcmpChecker) CheckFirstPartyCaveat(caveat string) error {
	return c.checkCaveat(caveat)
}

// CheckThirdPartyCaveat verifies that a third-party string caveat is satisfied
func (c StrcmpChecker) CheckThirdPartyCaveat(caveatId string, caveat string) ([]checkers.Caveat, error) {
	fmt.Println("Caveat ID:", caveatId)
	err := c.checkCaveat(caveat)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c StrcmpChecker) checkCaveat(caveat string) error {
	if caveat != string(c) {
		return fmt.Errorf("%s does not match %s", caveat, c)
	}
	return nil
}
