package csvfile

import (
	"bytes"
	"encoding/csv"
	"os"

	"github.com/diamondburned/cgowrap/internal/logg"
)

// Write writes a single record into a file at the given location.
func Write(path string, rows ...string) {
	if err := write(path, rows...); err != nil {
		logg.DebugFatalErr("cannot write csv:", err)
	}
}

func write(path string, rows ...string) error {
	var buf bytes.Buffer

	r := csv.NewWriter(&buf)

	if err := r.Write(rows); err != nil {
		return err
	}

	r.Flush()

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_SYNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(buf.Bytes()); err != nil {
		return err
	}

	return f.Close()
}
