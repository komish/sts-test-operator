package controller

import (
	"github.com/komish/sts-test-operator/pkg/controller/ststester"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, ststester.Add)
}
