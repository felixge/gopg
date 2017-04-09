package gopg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type ItemIdentifierIter struct {
	p      *Page
	offset uint16
}

func (ii *ItemIdentifierIter) Next() (*ItemIdentifier, error) {
	header, err := ii.p.Header()
	if err != nil {
		return nil, err
	}
	if ii.offset >= header.Lower {
		return nil, io.EOF
	}
	// Item identifiers are a bit of a PITA to parse in Go, because
	// they use non-8 bit aligned fields :(
	a := binary.LittleEndian.Uint16(ii.p.data[ii.offset:])
	ii.offset += 2
	b := binary.LittleEndian.Uint16(ii.p.data[ii.offset:])
	ii.offset += 2
	itemIdentifier := &ItemIdentifier{
		p: ii.p,
		ItemIdentifierData: ItemIdentifierData{
			Offset: a & 0x7fff,
			Flags:  uint8((a >> 15) | ((b & 0x0001) << 1)),
			Len:    b >> 1,
		},
	}
	return itemIdentifier, nil
}

type ItemIdentifier struct {
	p *Page
	ItemIdentifierData
}

func (ii *ItemIdentifier) String() string {
	return fmt.Sprintf("%+v", ii.ItemIdentifierData)
}

func (ii *ItemIdentifier) Tuple() (*Tuple, error) {
	t := &Tuple{ii: ii}
	r := bytes.NewReader(ii.p.data[ii.Offset:])
	if err := binary.Read(r, binary.LittleEndian, &t.TupleHeader); err != nil {
		return nil, err
	}
	return t, nil
}

type ItemIdentifierData struct {
	Offset uint16 // 15 bit
	Flags  uint8  // 2 bit
	Len    uint16 // 15 bit
}
