package goproxmox

import (
	"fmt"
	"regexp"
)

var (
	vmDoesNotExistRegexp   = regexp.MustCompile(`500 Configuration file \S+\/(\d+).conf' does not exist$`)
	nodeDoesNotExistRegexp = regexp.MustCompile(`500 hostname lookup '(\S+)' failed - failed to get address info for: \S+: Name or service not known$`)
)

// ArgError is an error that represents an error with an input to goproxmox. It
// identifies the argument and the cause (if possible).
type ArgError struct {
	arg    string
	reason string
}

var _ error = &ArgError{}

// NewArgError creates an InputError.
func NewArgError(arg, reason string) *ArgError {
	return &ArgError{
		arg:    arg,
		reason: reason,
	}
}

func (e *ArgError) Error() string {
	return fmt.Sprintf("%s is invalid because %s", e.arg, e.reason)
}

type NodeDoesNotExistError struct {
	Node string
}

func (e *NodeDoesNotExistError) Error() string {
	return fmt.Sprintf("Node %s doesn't exist", e.Node)
}

type VMDoesNotExistError struct {
	VMID string
}

func (e *VMDoesNotExistError) Error() string {
	return fmt.Sprintf("VM with id %s doesn't exist", e.VMID)
}
