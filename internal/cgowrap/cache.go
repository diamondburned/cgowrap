package cgowrap

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/diamondburned/cgowrap/internal/depfile"
	"go.etcd.io/bbolt"
)

var ErrNotFound = errors.New("not found")

var (
	depfileBucket    = "depfile"
	guessKindsBucket = "guessKinds"
)

type Cache struct {
	db         *bbolt.DB
	Depfile    *DepfileCache
	GuessKinds *GuessKindsCache
}

func OpenCache() (*Cache, error) {
	opt := *bbolt.DefaultOptions
	opt.Timeout = time.Minute
	opt.FreelistType = bbolt.FreelistMapType

	b, err := bbolt.Open(WorkFile("cache"), os.ModePerm, &opt)
	if err != nil {
		return nil, err
	}

	c := &Cache{db: b}
	c.Depfile = (*DepfileCache)(c)
	c.GuessKinds = (*GuessKindsCache)(c)

	return c, nil
}

func (c *Cache) Close() error {
	return c.db.Close()
}

type DepfileCache Cache

type depfileValue struct {
	File   depfile.File
	Latest time.Time
}

// IsValid returns true if the depfile cache is still valid.
func (c *DepfileCache) IsValid(id string) bool {
	var value depfileValue

	err := c.db.View(func(tx *bbolt.Tx) error {
		b, err := bucketTx(tx, depfileBucket)
		if err != nil {
			return ErrNotFound
		}

		v := b.Get([]byte(id))
		if v == nil {
			return ErrNotFound
		}

		return json.Unmarshal(v, &value)
	})

	if err != nil {
		return false
	}

	t := value.File.ModTime()
	log.Println(" got modTime =", t)
	log.Println("prev modTime =", value.Latest)

	// Verify the file list's modification time.
	return !t.IsZero() && t.Equal(value.Latest)
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

	return c.db.Update(func(tx *bbolt.Tx) error {
		b, err := bucketTx(tx, depfileBucket)
		if err != nil {
			return err
		}

		return b.Put([]byte(id), v)
	})
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

	c.db.View(func(tx *bbolt.Tx) error {
		b, _ := bucketTx(tx, guessKindsBucket, k)
		if b == nil {
			return ErrNotFound
		}

		if err := json.Unmarshal(b.Get([]byte("json")), &out); err != nil {
			return err
		}

		out.Stdout, _ = decompressBytes(b.Get([]byte("out")))
		out.Stderr, _ = decompressBytes(b.Get([]byte("err")))

		return nil
	})

	return out, !out.IsEmpty()
}

func (c *GuessKindsCache) Save(k string, out Output) error {
	j, err := json.Marshal(out)
	if err != nil {
		return err
	}

	return c.db.Update(func(tx *bbolt.Tx) error {
		b, err := bucketTx(tx, guessKindsBucket, k)
		if err != nil {
			return err
		}

		errs := []error{
			b.Put([]byte("json"), j),
			b.Put([]byte("out"), compressBytes(out.Stdout)),
			b.Put([]byte("err"), compressBytes(out.Stderr)),
		}

		for _, err := range errs {
			if err != nil {
				return err
			}
		}
		return nil
	})
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

func decompressBytes(b []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(b))
	if err != nil {
		log.Println("zlib: NewReader:", err)
		return nil, err
	}
	b, err = io.ReadAll(r)
	if err != nil {
		log.Println("zlib: read:", err)
		return nil, err
	}
	return b, nil
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
