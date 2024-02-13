package internal

import (
	"golang.org/x/tools/imports"
)

func ImportsCode(code string) ([]byte, error) {
	opts := &imports.Options{
		TabIndent: true,
		TabWidth:  2,
		Fragment:  true,
		Comments:  true,
	}

	formatcode, err := imports.Process("", []byte(code), opts)
	if err != nil {
		return nil, err
	}

	if string(formatcode) == code {
		return formatcode, err
	}

	return ImportsCode(string(formatcode))
}
