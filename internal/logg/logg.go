package logg

import (
	"io"
	"log"
	"os"
)

var isEnabled bool

func SetEnabled(enabled bool) {
	isEnabled = enabled
	update()
}

func Debug(v ...interface{}) {
	log.Println(v...)
}

func DebugFatalErr(msg string, err error) {
	if err == nil {
		return
	}

	DebugFatal(msg, err)
}

func DebugFatal(v ...interface{}) {
	if len(v) > 0 {
		if err, ok := v[len(v)-1].(error); ok && err == nil {
			return
		}
	}

	if isEnabled {
		log.Fatalln(v...)
	}
}

func update() {
	if isEnabled {
		log.SetFlags(0)
		log.SetPrefix("cgowrap: ")
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(io.Discard)
	}
}
