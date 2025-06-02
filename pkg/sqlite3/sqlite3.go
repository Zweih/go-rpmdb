package sqlite3

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	dbi "github.com/Zweih/go-rpmdb/pkg/db"
)

type SQLite3 struct {
	path string
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

		cmd := exec.Command("sqlite3", db.path, "SELECT hex(blob) FROM Packages")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			entries <- dbi.Entry{Err: fmt.Errorf("failed to create pipe: %w", err)}
			return
		}

		if err := cmd.Start(); err != nil {
			entries <- dbi.Entry{Err: fmt.Errorf("failed to start sqlite3: %w", err)}
			return
		}

		scanner := bufio.NewScanner(stdout)
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, 10*1024*1024)

		for scanner.Scan() {
			hexData := strings.TrimSpace(scanner.Text())
			if hexData == "" {
				continue
			}

			data, err := hexDecode(hexData)
			if err != nil {
				entries <- dbi.Entry{Err: fmt.Errorf("failed to decode hex data: %w", err)}
				return
			}

			entries <- dbi.Entry{Value: data, Err: nil}
		}

		if err := scanner.Err(); err != nil {
			entries <- dbi.Entry{Err: fmt.Errorf("scanner error: %w", err)}
			return
		}

		if err := cmd.Wait(); err != nil {
			entries <- dbi.Entry{Err: fmt.Errorf("sqlite3 command failed: %w", err)}
		}
	}()

	return entries
}

func (db *SQLite3) Close() error {
	return nil
}

func hexDecode(s string) ([]byte, error) {
	if len(s)%2 != 0 {
		return nil, fmt.Errorf("odd length hex string")
	}

	result := make([]byte, len(s)/2)
	for i := 0; i < len(s); i += 2 {
		high, err := hexCharToByte(s[i])
		if err != nil {
			return nil, err
		}
		low, err := hexCharToByte(s[i+1])
		if err != nil {
			return nil, err
		}
		result[i/2] = high<<4 | low
	}
	return result, nil
}

func hexCharToByte(c byte) (byte, error) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', nil
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, nil
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, nil
	default:
		return 0, fmt.Errorf("invalid hex character: %c", c)
	}
}
