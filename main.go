package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/diamondburned/cgowrap/internal/cgowrap"
	"github.com/diamondburned/cgowrap/internal/logg"
)

func init() {
	logg.SetEnabled(os.Getenv("CGOWRAP_FATAL") == "1")
}

type state struct {
	args      []string
	input     []byte
	cache     cacheState
	cacheable bool
}

type cacheState struct {
	*cgowrap.Cache
	depfileKey  string
	depfilePath string
	cacheKey    string
}

func main() {
	out := run()
	out.Print()
	os.Exit(out.Status)
}

func run() cgowrap.Output {
	s := state{args: os.Args[1:]}
	s.init()
	defer s.close()

	if out, ok := s.cached(); ok {
		return out
	}

	return s.run()
}

func (s *state) init() {
	if s.args[len(s.args)-1] == "-" {
		// TODO: stdin support
		return
	}

	input, err := os.ReadFile(s.args[len(s.args)-1])
	if err != nil {
		logg.DebugFatal("cannot open input file:", err)
		return
	}

	if !isGuessKinds(input) {
		return
	}

	if !s.openCache() {
		return
	}

	s.input = input
	s.cacheable = true
	return
}

// cached returns true if the command is cached and still valid.
func (s *state) cached() (cgowrap.Output, bool) {
	if !s.cacheable {
		return cgowrap.Output{}, false
	}

	s.cache.depfileKey = fmt.Sprintf("cgo.%s.d", hashInputs(
		// Use the arguments (excluded input path) and the input file content as
		// the hash for the depfile name.
		s.args[:len(s.args)-1],
		s.input,
	))

	if !s.cache.Depfile.IsValid(s.cache.depfileKey) {
		// Depfile not found, so avoid this cache and ask for a new one.
		s.args = append([]string{"-MD", "-MF", s.cache.Depfile.Path(s.cache.depfileKey)}, s.args...)
		log.Println("invalid depfile")
		return cgowrap.Output{}, false
	}

	// We'll only check the cached output if our depfile is up to date. We don't
	// need to account for this in the input hash, though.
	out, ok := s.cache.GuessKinds.Load(s.cache.depfileKey)
	if ok {
		log.Println("cache hit")
		return out, true
	}

	log.Println("missing guesskinds")
	return cgowrap.Output{}, false
}

// openCache initializes the cache.
func (s *state) openCache() bool {
	if s.cache.Cache != nil {
		return true
	}

	cache, err := cgowrap.OpenCache()
	if err == nil {
		s.cache.Cache = cache
		return true
	}

	logg.DebugFatalErr("cannot open cache database:", err)
	return false
}

func (s *state) close() {
	if s.cacheable {
		err := s.cache.Close()
		logg.DebugFatalErr("cannot close db:", err)
	}
}

// run runs the compiler and caches it if available.
func (s *state) run() cgowrap.Output {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command(CC(), s.args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	cmd.Run()

	out := cgowrap.Output{
		Stdout: stdout.Bytes(),
		Stderr: stderr.Bytes(),
		Status: cmd.ProcessState.ExitCode(),
	}

	s.save(out)
	return out
}

func (s *state) save(out cgowrap.Output) {
	if !s.cacheable || s.cache.depfileKey == "" {
		return
	}

	var err error

	err = s.cache.Depfile.Save(s.cache.depfileKey)
	logg.DebugFatalErr("cannot save depfile:", err)

	err = s.cache.GuessKinds.Save(s.cache.depfileKey, out)
	logg.DebugFatalErr("cannot save guessKinds:", err)
}

func CC() string {
	return envOr("CGOWRAP_CC", envOr("GCC", "gcc"))
}

func envOr(env, or string) string {
	if v := os.Getenv(env); v != "" {
		return v
	}
	return or
}

func hashStr(str string) string {
	h := sha256.New()
	return base64.URLEncoding.EncodeToString(h.Sum([]byte(str)))
}

func hashInputs(v ...interface{}) string {
	h := sha256.New()
	for _, v := range v {
		switch v := v.(type) {
		case string:
			io.WriteString(h, v)
		case []string:
			for _, v := range v {
				io.WriteString(h, v)
			}
		case []byte:
			h.Write(v)
		default:
			fmt.Fprint(h, v)
		}
	}
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func isGuessKinds(stdin []byte) bool {
	return bytesContainsLines(stdin,
		`#line 1 "cgo-generated-wrapper"`,
		`#line 1 "completed"`,
		`int __cgo__1 = __cgo__2;`,
	)
}

func bytesContainsLines(b []byte, strs ...string) bool {
	for _, str := range strs {
		if !bytes.Contains(b, []byte(str)) {
			return false
		}
	}
	return true
}
