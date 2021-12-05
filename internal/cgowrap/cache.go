package cgowrap

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/diamondburned/cgowrap/internal/depfile"
	"github.com/diamondburned/cgowrap/internal/logg"
	"github.com/peterbourgon/diskv/v3"
	"go.etcd.io/bbolt"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrMismatchModTime = errors.New("modTime mismatch")
)

var (
	depfileBucket    = "depfile"
	guessKindsBucket = "guessKinds"
)

func joinKeys(parts ...string) string {
	return strings.Join(parts, "$")
}

func getKV(db *diskv.Diskv, keys []string) (io.ReadCloser, error) {
	return db.ReadStream(joinKeys(keys...), true)
}

func getKVJSON(db *diskv.Diskv, keys []string, v interface{}) error {
	r, err := db.ReadStream(joinKeys(keys...), true)
	if err != nil {
		return err
	}
	defer r.Close()
	return json.NewDecoder(r).Decode(v)
}

func getKVBytes(db *diskv.Diskv, keys []string) ([]byte, error) {
	return db.Read(joinKeys(keys...))
}

func getKVCompressed(db *diskv.Diskv, keys []string) ([]byte, error) {
	r, err := db.ReadStream(joinKeys(keys...), true)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	z, err := zlib.NewReader(r)
	if err != nil {
		logg.DebugFatalErr("zlib: NewReader:", err)
		return nil, err
	}

	b, err := io.ReadAll(z)
	if err != nil {
		logg.DebugFatalErr("zlib: read:", err)
		return nil, err
	}

	return b, nil
}

func setKV(db *diskv.Diskv, keys []string, v []byte) error {
	return db.Write(joinKeys(keys...), v)
}

type Cache struct {
	db         *diskv.Diskv
	Depfile    *DepfileCache
	GuessKinds *GuessKindsCache
}

func OpenCache() (*Cache, error) {
	opt := *bbolt.DefaultOptions
	opt.Timeout = time.Minute
	opt.FreelistType = bbolt.FreelistMapType

	kv := diskv.New(diskv.Options{
		BasePath:     WorkDir("cache"),
		TempDir:      WorkDir(".cache.tmp"),
		Transform:    func(s string) []string { return nil },
		CacheSizeMax: 0,
	})

	c := &Cache{db: kv}
	c.Depfile = (*DepfileCache)(c)
	c.GuessKinds = (*GuessKindsCache)(c)

	return c, nil
}

type DepfileCache Cache

type depfileValue struct {
	File   depfile.File
	Latest time.Time
}

// IsValid returns true if the depfile cache is still valid.
func (c *DepfileCache) Validate(id string) error {
	var value depfileValue

	if err := getKVJSON(c.db, []string{depfileBucket, id}, &value); err != nil {
		return err
	}

	t := value.File.ModTime()
	// Verify the file list's modification time.
	if t.IsZero() || !t.Equal(value.Latest) {
		return ErrMismatchModTime
	}

	return nil
}

func (c *DepfileCache) Path(id string) string {
	return WorkFile("depfiles", id)
}

func (c *DepfileCache) Save(id string) error {
	f, err := depfile.ParseFileOnDisk(c.Path(id))
	if err != nil {
		return err
	}

	// Get rid of the first input file.
	f.PopFirstSources()

	v, err := json.Marshal(depfileValue{
		File:   *f,
		Latest: f.ModTime(),
	})
	if err != nil {
		return err
	}

	return setKV(c.db, []string{depfileBucket, id}, v)
}

type GuessKindsCache Cache

type Output struct {
	Stdout []byte `json:"-"`
	Stderr []byte `json:"-"`
	Status int    `json:"status"`
}

func (o Output) IsEmpty() bool {
	return len(o.Stdout) == 0 && len(o.Stderr) == 0
}

// Print prints the output.
func (o Output) Print() {
	os.Stderr.Write(o.Stderr)
	os.Stdout.Write(o.Stdout)
}

func (c *GuessKindsCache) Load(k string) (Output, bool) {
	var out Output
	keys := []string{guessKindsBucket, k, "json"}

	if err := getKVJSON(c.db, keys, &out); err != nil {
		return out, false
	}

	keys[2] = "out"
	out.Stdout, _ = getKVCompressed(c.db, keys)

	keys[2] = "err"
	out.Stderr, _ = getKVCompressed(c.db, keys)

	return out, !out.IsEmpty()
}

func (c *GuessKindsCache) Save(k string, out Output) error {
	j, err := json.Marshal(out)
	if err != nil {
		return err
	}

	errs := []error{
		setKV(c.db, []string{guessKindsBucket, k, "json"}, j),
		setKV(c.db, []string{guessKindsBucket, k, "out"}, compressBytes(out.Stdout)),
		setKV(c.db, []string{guessKindsBucket, k, "err"}, compressBytes(out.Stderr)),
	}

	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func compressBytes(b []byte) []byte {
	var out bytes.Buffer
	w := zlib.NewWriter(&out)

	_, err := w.Write(b)
	if err != nil {
		log.Panicln("zlib error:", err)
	}

	if err := w.Close(); err != nil {
		log.Panicln("zlib error: close:", err)
	}

	return out.Bytes()
}

func cpyBytes(b []byte) []byte {
	if b == nil {
		return nil
	}
	return append([]byte(nil), b...)
}

type bucketter interface {
	CreateBucketIfNotExists([]byte) (*bbolt.Bucket, error)
	Bucket([]byte) *bbolt.Bucket
}

func bucketTx(tx *bbolt.Tx, paths ...string) (*bbolt.Bucket, error) {
	return bucket(tx.Writable(), tx, paths...)
}

func bucket(create bool, b bucketter, paths ...string) (*bbolt.Bucket, error) {
	if len(paths) == 0 {
		log.Panicln("no paths given")
	}

	if create {
		var err error
		for _, path := range paths {
			b, err = b.CreateBucketIfNotExists([]byte(path))
			if err != nil {
				return nil, err
			}
		}
	} else {
		for _, path := range paths {
			bucket := b.Bucket([]byte(path))
			if bucket == nil {
				return nil, ErrNotFound
			}
			// beware of typed nil interfaces
			b = bucket
		}
	}

	return b.(*bbolt.Bucket), nil
}
