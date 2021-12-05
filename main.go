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
	"strconv"
	"strings"
	"time"

	"github.com/diamondburned/cgowrap/internal/cgowrap"
	"github.com/diamondburned/cgowrap/internal/csvfile"
	"github.com/diamondburned/cgowrap/internal/logg"
	"github.com/diamondburned/cgowrap/internal/shortflag"
)

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

var pwd, _ = os.Getwd()

var (
	profileOut = os.Getenv("CGOWRAP_PROFILE")
	runFatal   = os.Getenv("CGOWRAP_FATAL") == "1"
	mustCache  = os.Getenv("CGOWRAP_MUST_CACHE") == "1"
)

func main() {
	logg.SetEnabled(runFatal || mustCache)

	out := run()
	out.Print()
	os.Exit(out.Status)
}

func run() cgowrap.Output {
	s := state{args: os.Args[1:]}
	s.init()
	defer s.close()

	var out cgowrap.Output

	f := func() bool {
		o, ok := s.cached()
		if ok {
			out = o
			return false
		}

		out = s.run()
		return true
	}

	if profileOut != "" {
		s.record(f)
	} else {
		f()
	}

	return out
}

func (s *state) record(f func() bool) {
	start := time.Now()
	status := "cached"
	if f() {
		status = "uncached"
	}
	end := time.Now()

	csvfile.Write(profileOut,
		strings.Join(s.args, " "),
		strconv.FormatFloat(end.Sub(start).Seconds(), 'f', -1, 64),
		status,
	)
}

func (s *state) init() {
	// isGuessKinds relies on the assumption that when the check passes, only
	// one input will ever be given, which is cgo's crafted input file. This
	// assumption simplifies the code a lot.

	input, err := os.ReadFile(s.args[len(s.args)-1])
	if err != nil {
		// Probably an incorrect assumption about arguments.
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

	// Parse never returns nil.
	args, err := shortflag.Parse(s.args, shortflag.Opts{
		ValueFlags: []string{"-o"},
	})
	if err != nil {
		logg.DebugFatalErr("error parsing flags:", err)
		return cgowrap.Output{}, false
	}

	// Ues the arguments without the -o flag and all input files. The input file
	// is assumed to only be 1, and the -o flag is not deterministic.
	neededArgs := shortflag.OmitNonFlags(args.Args)
	// Use neededArgs as the hash input for depfileKey along with the input
	// file's content and the current working directory.
	hash := hashAll([]interface{}{neededArgs, pwd, s.input})
	s.cache.depfileKey = fmt.Sprintf("cgo.%s.d", hash)

	if err := s.cache.Depfile.Validate(s.cache.depfileKey); err != nil {
		// Depfile not found, so avoid this cache and ask for a new one.
		s.args = append([]string{"-MD", "-MF", s.cache.Depfile.Path(s.cache.depfileKey)}, s.args...)
		cacheMissed(neededArgs, hash, "invalid depfile:", err)
		return cgowrap.Output{}, false
	}

	// We'll only check the cached output if our depfile is up to date. We don't
	// need to account for this in the input hash, though.
	out, ok := s.cache.GuessKinds.Load(s.cache.depfileKey)
	if ok {
		return out, true
	}
	cacheMissed(neededArgs, hash, "missing guessKinds")
	return cgowrap.Output{}, false
}

func cacheMissed(args []string, hash string, v ...interface{}) {
	if mustCache {
		log.Printf("args: %q", args)
		log.Printf("hash: %q", hash)
		v = append([]interface{}{"cache missed:"}, v...)
		log.Fatalln(v...)
	}
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

func hashAll(v ...interface{}) string {
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
