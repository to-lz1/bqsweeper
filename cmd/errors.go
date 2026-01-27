package cmd

import "errors"

var (
	ErrInvalidArgs = errors.New("invalid arguments")
	ErrNoTables    = errors.New("no tables found")
)
