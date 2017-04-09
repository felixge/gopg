package gopg

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
)

const (
	pageHeaderSize     = 24
	itemIdentifierSize = 4
	tupleHeaderSize    = 23
)

type PageIter struct {
	t *Table
	f *os.File
}

func (p *PageIter) Next() (*Page, error) {
	pageSize := p.t.d.c.c.PageSize
	buf := make([]byte, pageSize)
	if _, err := io.ReadAtLeast(p.f, buf, pageSize); err != nil {
		return nil, err
	}
	return &Page{data: buf}, nil
}

func (p *PageIter) Close() error {
	return p.f.Close()
}

type Page struct {
	data   []byte
	header *PageHeader
}

func (p *Page) Header() (*PageHeader, error) {
	if p.header != nil {
		return p.header, nil
	}
	p.header = &PageHeader{p: p}
	if err := binary.Read(
		bytes.NewBuffer(p.data),
		binary.LittleEndian,
		&p.header.PageHeaderData,
	); err != nil {
		return nil, err
	}
	return p.header, nil
}

func (p *Page) ItemIdentifiers() *ItemIdentifierIter {
	return &ItemIdentifierIter{
		p:      p,
		offset: pageHeaderSize,
	}
}
