package gopg

import "os"

type Table struct {
	d    *Database
	path string
}

func (t *Table) Exists() error {
	return exists(t.path)
}

func (t *Table) Pages() (*PageIter, error) {
	// TODO: tables can have multiple files if > 1 GB
	file, err := os.Open(t.path)
	if err != nil {
		return nil, err
	}
	return &PageIter{t: t, f: file}, nil
}
