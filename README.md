# gopg

This is a toy project for verifying my understanding of how [postgres stores
data on
disk](https://www.postgresql.org/docs/9.6/static/storage-page-layout.html).

## What's Implemented

Given a postgres data directory, a database `oid` and a table `relfilenode`
this project will read all heap pages, and follow their item identifiers, as
well as heap tuple headers and data.

## Run it yourself

Clone the project, make sure you have postgres installed (`initdb` and `pg_ctl`
should be in your PATH), then simply type make:

```
$ make
# lots of output related to initializing a standalone postgres instance
# in the local `data` directory (this will not impact your system install) ...
=== RUN   TestExample
page 0:
  header &gopg.PageHeader{p:(*gopg.Page)(nil), PageHeaderData:gopg.PageHeaderData{LSN:0x0, Checksum:0x0, Flags:0x0, Lower:0x24, Upper:0x1fa0, Special:0x2000, PageSizeVersion:0x2004, PruneXid:0x0}}:
  tuple 1
    item identifier: &gopg.ItemIdentifier{p:(*gopg.Page)(0xc42010e1e0), ItemIdentifierData:gopg.ItemIdentifierData{Offset:0x1fe0, Flags:0x1, Len:0x1c}}
    tuple header: &gopg.Tuple{ii:(*gopg.ItemIdentifier)(0xc420108200), TupleHeader:gopg.TupleHeader{XMin:0x361, XMax:0x0, Field3:0x1, CTID:[6]uint8{0x0, 0x0, 0x0, 0x0, 0x1, 0x0}, XInfomask2:0x1, XInfomask:0x800, Offset:0x18}}
    data: 01000000
  tuple 2
    item identifier: &gopg.ItemIdentifier{p:(*gopg.Page)(0xc42010e1e0), ItemIdentifierData:gopg.ItemIdentifierData{Offset:0x1fc0, Flags:0x1, Len:0x1c}}
    tuple header: &gopg.Tuple{ii:(*gopg.ItemIdentifier)(0xc420108340), TupleHeader:gopg.TupleHeader{XMin:0x361, XMax:0x0, Field3:0x1, CTID:[6]uint8{0x0, 0x0, 0x0, 0x0, 0x2, 0x0}, XInfomask2:0x1, XInfomask:0x800, Offset:0x18}}
    data: 02000000
  tuple 3
    item identifier: &gopg.ItemIdentifier{p:(*gopg.Page)(0xc42010e1e0), ItemIdentifierData:gopg.ItemIdentifierData{Offset:0x1fa0, Flags:0x1, Len:0x1c}}
    tuple header: &gopg.Tuple{ii:(*gopg.ItemIdentifier)(0xc420108480), TupleHeader:gopg.TupleHeader{XMin:0x361, XMax:0x0, Field3:0x1, CTID:[6]uint8{0x0, 0x0, 0x0, 0x0, 0x3, 0x0}, XInfomask2:0x1, XInfomask:0x800, Offset:0x18}}
    data: 03000000
--- PASS: TestExample (1.02s)
```
