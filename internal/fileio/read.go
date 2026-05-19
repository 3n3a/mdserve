package fileio

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
)

var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

func Read(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", err
		}
		return "", err
	}

	trimmed := bytes.TrimPrefix(raw, utf8BOM)
	if utf8.Valid(trimmed) {
		return string(trimmed), nil
	}

	decoded, err := charmap.ISO8859_1.NewDecoder().Bytes(raw)
	if err != nil {
		return string(raw), nil
	}
	return string(decoded), nil
}
