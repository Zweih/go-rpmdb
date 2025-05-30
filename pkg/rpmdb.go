package pkg

import (
	"errors"
	"fmt"

	"github.com/Zweih/go-rpmdb/pkg/bdb"
	dbi "github.com/Zweih/go-rpmdb/pkg/db"
	"github.com/Zweih/go-rpmdb/pkg/ndb"
	"github.com/Zweih/go-rpmdb/pkg/sqlite3"
)

type RpmDB struct {
	db dbi.RpmDBInterface
}

func Open(path string) (*RpmDB, error) {
	// SQLite3 Open() returns nil, nil in case of DB format other than SQLite3
	sqldb, err := sqlite3.Open(path)
	if err != nil && !errors.Is(err, sqlite3.ErrorInvalidSQLite3) {
		return nil, err
	}
	if sqldb != nil {
		return &RpmDB{db: sqldb}, nil
	}

	// NDB Open() returns nil, nil in case of DB format other than NDB
	ndbh, err := ndb.Open(path)
	if err != nil && !errors.Is(err, ndb.ErrorInvalidNDB) {
		return nil, err
	}
	if ndbh != nil {
		return &RpmDB{db: ndbh}, nil
	}

	odb, err := bdb.Open(path)
	if err != nil {
		return nil, err
	}

	return &RpmDB{
		db: odb,
	}, nil
}

func (d *RpmDB) Close() error {
	return d.db.Close()
}

func (d *RpmDB) Package(name string) (*PackageInfo, error) {
	pkgs, err := d.ListPackages()
	if err != nil {
		return nil, fmt.Errorf("unable to list packages: %w", err)
	}

	for _, pkg := range pkgs {
		if pkg.Name == name {
			return pkg, nil
		}
	}
	return nil, fmt.Errorf("%s is not installed", name)
}

func (d *RpmDB) ListPackages() ([]*PackageInfo, error) {
	var pkgList []*PackageInfo

	for entry := range d.db.Read() {
		if entry.Err != nil {
			return nil, entry.Err
		}

		indexEntries, err := HeaderImport(entry.Value)
		if err != nil {
			return nil, fmt.Errorf("error during importing header: %w", err)
		}
		pkg, err := GetNEVRA(indexEntries)
		if err != nil {
			return nil, fmt.Errorf("invalid package info: %w", err)
		}
		pkgList = append(pkgList, pkg)
	}

	return pkgList, nil
}
