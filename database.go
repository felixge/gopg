package gopg

import (
	"fmt"
	"os"
	"path/filepath"
)

type DatabaseIter struct {
	c *Cluster
	f *os.File
}

func (d *DatabaseIter) Next() (*Database, error) {
	infos, err := d.f.Readdir(1)
	if err != nil {
		return nil, err
	}
	db := &Database{
		c:    d.c,
		path: filepath.Join(d.f.Name(), infos[0].Name()),
	}
	return db, nil
}

func (d *DatabaseIter) Close() error {
	return d.f.Close()
}

type Database struct {
	c    *Cluster
	path string
}

func (d *Database) tablePath(relFileNode uint32) string {
	return filepath.Join(d.path, fmt.Sprintf("%d", relFileNode))
}

func (d *Database) Exists() error {
	return exists(d.path)
}

func (d *Database) TableByRelFileNode(relFileNode uint32) (*Table, error) {
	t := &Table{d: d, path: d.tablePath(relFileNode)}
	return t, t.Exists()
}
