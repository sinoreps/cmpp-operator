package controller

import (
	"github.com/sinoreps/cmpp-operator/pkg/controller/httpproxy"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, httpproxy.Add)
}
