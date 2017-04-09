package gopg

import "fmt"

type Tuple struct {
	ii *ItemIdentifier
	TupleHeader
}

func (t *Tuple) String() string {
	return fmt.Sprintf("%+v", t.TupleHeader)
}

// TODO: we need to get access to the catalogs to get the column types and
// properly read the data.
func (t *Tuple) Data() []byte {
	start := int(t.ii.Offset) + int(t.Offset)
	end := int(t.ii.Offset) + int(t.ii.Len)
	return t.ii.p.data[start:end]
}

type TupleHeader struct {
	XMin      uint32
	XMax      uint32
	Field3    uint32 // CID or XVac (C Union)
	CTID      [6]byte
	Infomask2 uint16
	Infomask  uint16
	Offset    uint8
}
