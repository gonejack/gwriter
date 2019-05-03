/*
Package writer is an file library for appending file.
*/
package gwriter

import (
	"github.com/gonejack/gwriter/config"
	"github.com/gonejack/gwriter/internal/basicWriter"
)

// NewWriter creates a new writer instance of basicWriter match Writer interface
func NewWriter(name string, config config.Config) (w Writer) {
	return basicWriter.New(name, config)
}
