package gopg

import "fmt"

type PageHeader struct {
	p *Page
	PageHeaderData
}

func (p *PageHeader) String() string {
	return fmt.Sprintf("%+v", p.PageHeaderData)
}

type PageHeaderData struct {
	LSN             uint64
	Checksum        uint16
	Flags           uint16
	Lower           uint16
	Upper           uint16
	Special         uint16
	PageSizeVersion uint16
	PruneXid        uint32
}

func (h *PageHeader) LenItems() int {
	return int(h.Lower-pageHeaderSize) / itemIdentifierSize
}
