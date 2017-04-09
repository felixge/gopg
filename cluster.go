package gopg

import (
	"fmt"
	"os"
	"path/filepath"
)

type ClusterConfig struct {
	// Path is the path to the postgres data directory (aka PGDATA)
	Path string
	// PageSize is the page size postgres was compiled to use (defaults
	// to 8KB).
	PageSize int
}

func (c ClusterConfig) WithDefaults() ClusterConfig {
	if c.PageSize == 0 {
		c.PageSize = 8 * 1024
	}
	return c
}

func (c ClusterConfig) Open() (*Cluster, error) {
	c = c.WithDefaults()
	cluster := &Cluster{c: c}
	return cluster, cluster.Exists()
}

type Cluster struct {
	c ClusterConfig
}

func (c *Cluster) Exists() error {
	return exists(c.c.Path)
}

func (c *Cluster) Databases() (*DatabaseIter, error) {
	file, err := os.Open(c.databasesPath())
	if err != nil {
		return nil, err
	}
	return &DatabaseIter{c: c, f: file}, nil
}

func (c *Cluster) databasesPath() string {
	return filepath.Join(c.c.Path, "base")
}

func (c *Cluster) databasePath(oid uint32) string {
	return filepath.Join(c.databasesPath(), fmt.Sprintf("%d", oid))
}

func (c *Cluster) DatabaseByOID(oid uint32) (*Database, error) {
	db := &Database{c: c, path: c.databasePath(oid)}
	return db, db.Exists()
}
