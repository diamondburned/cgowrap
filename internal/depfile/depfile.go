package depfile

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// FileList describes a list of file paths.
type FileList []string

// ModTime returns the modification time of the file list. A zero-value time is
// returned if the list is empty or if no files can be stat'd. Otherwise, the
// latest ModTime is returned.
func (l FileList) ModTime() time.Time {
	var t time.Time

	for _, f := range l {
		s, err := os.Stat(f)
		if err != nil {
			log.Println("missing", f)
			// TODO: invalidate cache
			continue
		}

		if mod := s.ModTime(); mod.After(t) {
			t = mod
		}
	}

	return t
}

// PopFirst pops the first file off the list and returns it.
func (l *FileList) PopFirst() string {
	first := (*l)[0]
	*l = (*l)[1:]
	return first
}

// File describes a partial depfile. It is not a complete representation of a
// depfile, but it doesn't have to be.
type File struct {
	// Sources maps the input file to a list of files that are its dependencies.
	Sources map[string]FileList
}

// ParseFileOnDisk parses the given file path.
func ParseFileOnDisk(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ParseFile(f)
}

// ParseFile parses the given reader.
func ParseFile(r io.Reader) (*File, error) {
	f := File{Sources: make(map[string]FileList)}

	// states
	var currentSources FileList
	var currentFile string

	// flush flushes the old file.
	flush := func() {
		if currentFile != "" {
			f.Sources[currentFile] = currentSources
		}
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}

		if strings.HasPrefix(text, "  ") {
			currentSources = append(currentSources, parsePartialLine(text))
			continue
		}

		// Potentially a new file.
		parts := strings.SplitN(text, ": ", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("unexpected line %q", text)
		}

		flush()

		currentFile = parts[0]
		currentSources = make([]string, 1, 16)
		currentSources[0] = parsePartialLine(parts[1])
	}

	flush()

	return &f, scanner.Err()
}

func parsePartialLine(line string) string {
	line = strings.TrimPrefix(line, "  ")
	line = strings.TrimSuffix(line, " \\")
	return line
}

// ModTime returns the latest ModTime in the sources.
func (f *File) ModTime() time.Time {
	var t time.Time

	for _, src := range f.Sources {
		if mod := src.ModTime(); mod.After(t) {
			t = mod
		}
	}

	return t
}

// PopFirstSources removes the first file in all sources. This is useful for
// getting rid of the input file.
func (f *File) PopFirstSources() {
	for k, src := range f.Sources {
		src.PopFirst()
		f.Sources[k] = src
	}
}
