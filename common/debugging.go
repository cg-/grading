/*
 *  Common Useful Functions
 */

package common

import "fmt"

// DebugLogger is an object that handles debugging messages
type DebugLogger struct {
	enabled bool
	prefix  string
}

// NewDebugLogger is a constructor
func NewDebugLogger() *DebugLogger {
	return &DebugLogger{
		enabled: false,
		prefix:  "debug: ",
	}
}

// Enable turns on debugging messages.
func (d *DebugLogger) Enable() {
	d.enabled = true
}

// Disable turns off debugging messages.
func (d *DebugLogger) Disable() {
	d.enabled = false
}

// Print outputs debugging messages.
func (d *DebugLogger) Print(s string) {
	if d.enabled {
		fmt.Printf("%s%s\n", d.prefix, s)
	}
}
