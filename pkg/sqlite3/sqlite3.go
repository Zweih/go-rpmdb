package sqlite3

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"sync"

	dbi "github.com/Zweih/go-rpmdb/pkg/db"
	"github.com/Zweih/go-rpmdb/pkg/worker"
)

type SQLite3 struct {
	path string
}

type decodedBlob struct {
	data  []byte
	index int
}

type hexBatch struct {
	lines      []string
	startIndex int
}

var (
	// https://www.sqlite.org/fileformat.html
	SQLite3_HeaderMagic = []byte("SQLite format 3\x00")
	ErrorInvalidSQLite3 = fmt.Errorf("invalid or unsupported SQLite3 format")
)

func Open(path string) (*SQLite3, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b := make([]byte, 16)
	if _, err := file.Read(b); err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	if !bytes.Equal(b, SQLite3_HeaderMagic) {
		return nil, ErrorInvalidSQLite3
	}

	if _, err := exec.LookPath("sqlite3"); err != nil {
		return nil, fmt.Errorf("sqlite3 command not found: %w", err)
	}

	return &SQLite3{path: path}, nil
}

func (db *SQLite3) Read() <-chan dbi.Entry {
	entries := make(chan dbi.Entry)

	go func() {
		defer close(entries)

		hexBatches := make(chan hexBatch, 20)
		errChan := make(chan error, 10)
		var errGroup sync.WaitGroup

		go func() {
			defer close(hexBatches)

			uri := fmt.Sprintf("file:%s?mode=ro&immutable=1&cache=private", db.path)

			script := `
      PRAGMA cache_size = -64000;
      PRAGMA page_size = 4096;
      PRAGMA temp_store = MEMORY;
      PRAGMA synchronous = OFF;
      SELECT hex(blob) FROM Packages;
      `

			cmd := exec.Command("sqlite3", "-readonly", "-batch", uri, script)

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				errChan <- fmt.Errorf("failed to create pipe: %w", err)
				return
			}

			if err := cmd.Start(); err != nil {
				errChan <- fmt.Errorf("failed to start sqlite3: %w", err)
				return
			}

			scanner := bufio.NewScanner(stdout)
			buf := make([]byte, 0, 64*1024)
			scanner.Buffer(buf, 8*1024*1024)
			batch := make([]string, 0, 10)
			startIndex := 0
			currentIndex := 0

			for scanner.Scan() {
				hexData := scanner.Text()
				if len(hexData) == 0 {
					continue
				}

				batch = append(batch, hexData)

				if len(batch) >= 10 {
					hexBatches <- hexBatch{
						lines:      batch,
						startIndex: startIndex,
					}

					batch = make([]string, 0, 10)
					startIndex = currentIndex + 1
				}

				currentIndex++
			}

			if len(batch) > 0 {
				hexBatches <- hexBatch{
					lines:      batch,
					startIndex: startIndex,
				}
			}
		}()

		decodedChan := worker.RunWorkers(
			hexBatches,
			errChan,
			&errGroup,
			decodeBatch,
			0,
			50,
		)

		for batch := range decodedChan {
			for _, entry := range batch {
				entries <- entry
			}
		}

		errGroup.Wait()
		close(errChan)

		for err := range errChan {
			entries <- dbi.Entry{Err: err}
		}
	}()

	return entries
}

func decodeBatch(batch hexBatch) ([]dbi.Entry, error) {
	results := make([]dbi.Entry, 0, len(batch.lines))

	for i, line := range batch.lines {
		data, err := hex.DecodeString(line)
		if err != nil {
			return nil, fmt.Errorf("failed to decode hex at index %d: %w", batch.startIndex+i, err)
		}

		results = append(results, dbi.Entry{Value: data})
	}

	return results, nil
}

func (db *SQLite3) Close() error {
	return nil
}
