package bytetrap

import (
	"math/rand"

	"embed"
	"io"
	"log"
	"path/filepath"
	"strings"
	"sync"
)

type Tag uint8

const (
	TagZero Tag = iota
	TagText
	TagHTML
)

//go:embed copypasta
var pastaFS embed.FS

// prefix in pastaFS
const pastaPrefix = "copypasta"

func (p *Pasta) HasTag(t Tag) bool {
	if p.Tags == nil {
		return false
	}

	_, ok := p.Tags[t]
	return ok
}

type Pasta struct {
	Content string

	Name      string
	Tags      map[Tag]struct{}
	NativeTag Tag
}

// Try to make one tag the other
func (p *Pasta) To(target *Tag) string {
	if target == nil { // cant cast to nil
		return p.Content
	}

	if p.NativeTag == TagText && *target == TagHTML {
		return "<p><code style=\"color:white,background-color:black\">\n" + strings.ReplaceAll(p.Content, "\n", "<br>") + "\n</code></p>"
	}

	return p.Content
}

var (
	pastaMap   map[string]*Pasta
	pastaSlice []*Pasta

	pastaOnce sync.Once
)

// read pastaFS and read into pastaMap and pastaSlice
func readPasta() {
	pastaOnce.Do(func() {
		go startStats()

		dir, err := pastaFS.ReadDir(pastaPrefix)
		if err != nil {
			log.Fatalf("Failed to read pastaPrefix('%s'): %s", pastaPrefix, err)
		}

		// init map&slice
		pastaMap = make(map[string]*Pasta, len(dir))
		pastaSlice = make([]*Pasta, 0, len(dir))

		for i := 0; i < len(dir); i++ {
			// cant read dir
			if dir[i].IsDir() {
				continue
			}

			name := dir[i].Name()
			path := pastaPrefix + "/" + name
			file, err := pastaFS.Open(path)
			if err != nil {
				log.Fatalf("Failed to read copypasta from '%s': %s", path, err)
			}

			b, err := io.ReadAll(file)
			if err != nil {
				log.Fatalf("Failed to read copypasta from '%s': %s", path, err)
			}

			tag := TagFromFile(name)
			pasta := &Pasta{
				Content:   string(b),
				Name:      name,
				Tags:      map[Tag]struct{}{tag: struct{}{}},
				NativeTag: tag,
			}

			pastaMap[name] = pasta
			pastaSlice = append(pastaSlice, pasta)
		}
	})
}

func TagFromFile(n string) Tag {
	switch filepath.Ext(n) {
	case ".html":
		return TagHTML
	case ".txt":
		return TagText
	}

	return TagZero
}

// do not write to the slice
func PastaSlice() []*Pasta {
	readPasta()

	return pastaSlice
}

// do not write to the slice
func PastaMap() map[string]*Pasta {
	readPasta()

	return pastaMap
}

// returns a random pasta
func GetPasta() *Pasta {
	return GetPastaTag(nil)
}

func GetPastaTag(t []Tag) *Pasta {
	s := PastaSlice()

	if t == nil {
		return s[rand.Intn(len(s))]
	}

	for {
		p := s[rand.Intn(len(s))]
		for _, t := range t {
			if p.HasTag(t) {
				return p
			}
		}
	}
}

// writes random stuff ignoring tags
func Write(w io.Writer) error {
	return write(w, false, nil, nil)
}

// writes random copy pasta thats just plain text
func WritePlain(w io.Writer) error {
	return write(w, false, []Tag{TagText}, nil)
}

// writes random copy pasta with one of ts tags
func WriteTag(w io.Writer, t []Tag) error {
	return write(w, false, t, nil)
}

func write(w io.Writer, stat bool, t []Tag, target *Tag) error {
	var i int
	var err error

	for {
		p := GetPastaTag(t)

		sagthi := p.To(target)
		i, err = w.Write([]byte(sagthi))
		if err != nil {
			return err
		}

		if stat {
			statsCh <- int64(i)
		}
	}
}
