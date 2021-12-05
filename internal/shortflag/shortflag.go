package shortflag

import (
	"fmt"
	"strings"
)

// Flags describes the parsed flags and arguments from the given strings.
type Flags struct {
	// Args contains non-flags.
	Args []string
	// Flags is a list of flags mapping to its values.
	Flags []Flag
}

// Flag returns the flag inside Flags, or nil if none.
func (f *Flags) Flag(name string) *Flag {
	for i, flag := range f.Flags {
		if flag.Name == name {
			return &f.Flags[i]
		}
	}
	return nil
}

func (f *Flags) addFlag(name string) *Flag {
	if flag := f.Flag(name); flag != nil {
		return flag
	}

	f.Flags = append(f.Flags, Flag{Name: name})
	return &f.Flags[len(f.Flags)-1]
}

// Flag describes a flag and all its values.
type Flag struct {
	// Name is the flag name including the prefixing dash.
	Name string
	// Values contains the flag values, or nil if the flag is a blank flag.
	Values []string
}

func (f *Flag) addValue(value string) {
	f.Values = append(f.Values, value)
}

// Opts describes the options for parsing.
type Opts struct {
	ValueFlags []string
	BlankFlags []string
}

type flagType uint8

const (
	notFlag flagType = iota
	valueFlag
	blankFlag
)

func (o Opts) lookupFlag(arg string) flagType {
	if findStr(o.ValueFlags, arg) > -1 {
		return valueFlag
	}
	if findStr(o.BlankFlags, arg) > -1 {
		return blankFlag
	}
	return notFlag
}

// Omit omits flags and flag values listed inside Opts.
func Omit(args []string, opts Opts) []string {
	f, err := Parse(args, opts)
	if err != nil {
		return args
	}
	return f.Args
}

// NonFlags omits all arguments that start with a dash.
func NonFlags(args []string) []string {
	return filter(args, func(v string) bool { return !strings.HasPrefix(v, "-") })
}

// OmitNonFlags omits all arguments that don't start with a dash.
func OmitNonFlags(args []string) []string {
	return filter(args, func(v string) bool { return strings.HasPrefix(v, "-") })
}

func filter(args []string, f func(string) bool) []string {
	strv := make([]string, 0, len(args))
	for _, arg := range args {
		if f(arg) {
			strv = append(strv, arg)
		}
	}
	return strv
}

// Parse parses the given list of arguments into Flags. The function will always
// return a valid Flags instance.
func Parse(args []string, opts Opts) (*Flags, error) {
	f := Flags{
		Args:  make([]string, 0, len(args)),
		Flags: make([]Flag, 0, len(opts.ValueFlags)+len(opts.BlankFlags)),
	}

	var currentFlag string

	for _, arg := range args {
		if currentFlag != "" {
			f.addFlag(currentFlag).addValue(arg)
			currentFlag = ""
			continue
		}

		if !strings.HasPrefix(arg, "-") {
			f.Args = append(f.Args, arg)
			continue
		}

		numDashes := numDashes(arg)
		flag := arg

		if numDashes == 1 && len(flag) > len("-X") {
			// Ensure that we're looking up the exact shortflag's name.
			flag = flag[:2]
		}

		switch opts.lookupFlag(flag) {
		case notFlag:
			f.Args = append(f.Args, arg)
		case blankFlag:
			f.addFlag(flag)
		case valueFlag:
			switch numDashes {
			case 0:
				panic("unreachable")
			case 1:
				// Single-dash, so the value can be right after. We can just do
				// a length check.
				if len(arg) > len("-X") {
					// Value is right after.
					f.addFlag(flag).addValue(arg[2:])
				} else {
					currentFlag = arg
				}
			default:
				// Definitely a long flag, meaning that the value might be on
				// the next value.
				currentFlag = arg
			}
		}
	}

	if currentFlag != "" {
		return &f, fmt.Errorf("flag %q missing value", currentFlag)
	}

	return &f, nil
}

func numDashes(str string) int {
	var i int
	for i < len(str) && str[i] == '-' {
		i++
	}
	return i
}

func findStr(strv []string, find string) int {
	for i, str := range strv {
		if strings.HasPrefix(str, find) {
			return i
		}
	}
	return -1
}
