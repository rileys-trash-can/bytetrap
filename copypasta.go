package bytetrap

import (
	"math/rand"
	"time"

	"embed"
	"io"
	"log"
	"sync"
)

//go:embed copypasta
var pastaFS embed.FS

// prefix in pastaFS
const pastaPrefix = "copypasta"

var (
	pastaMap   map[string]string
	pastaSlice []string

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
		pastaMap = make(map[string]string, len(dir))
		pastaSlice = make([]string, 0, len(dir))

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

			str := string(b)

			pastaMap[name] = str
			pastaSlice = append(pastaSlice, str)
		}
	})
}

// do not write to the slice
func PastaSlice() []string {
	readPasta()

	return pastaSlice
}

// do not write to the slice
func PastaMap() map[string]string {
	readPasta()

	return pastaMap
}

var Rand = rand.New(rand.NewSource(time.Now().UnixMilli() + time.Now().UnixMicro()))

// returns a random pasta
func GetPasta() string {
	s := PastaSlice()

	return s[Rand.Intn(len(s))]
}

func Write(w io.Writer) error {
	return write(w, false)
}

func write(w io.Writer, stat bool) error {
	for {
		i, err := w.Write([]byte(GetPasta()))
		if err != nil {
			return err
		}

		if stat {
			statsCh <- int64(i)
		}
	}
}
