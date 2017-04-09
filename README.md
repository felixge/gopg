# gopg

This is a toy project for verifying my understanding of how [postgres stores
data on
disk](https://www.postgresql.org/docs/9.6/static/storage-page-layout.html) by
re-implementing parts of it as a Go library.

## What's Implemented

Given a postgres data directory, a database `oid` and a table `relfilenode`
it's possible to read all heap pages, and follow their item identifiers, and
read their tuple headers and raw tuple data.

You can see this in action in the `example_test.go` file.

## What's not Implemented

Pretty much everything you'd need in order to do something useful :). E.g.:

* Open database by name (rather than `oid`)
* Open table by name (rather than `relfilenode`)
* Read tables > 1GB
* Convert raw page data into named columns and proper types
* Respect MVCC visibility rules
* Updating of pages

I might implement some of these features in the future, but don't expect to be
able to do anything useful with this library anytime soon.

## Run it yourself

Clone the project, make sure you have postgres installed (`initdb` and `pg_ctl`
should be in your PATH), then simply type `make`:

```
$ make
# Lots of output related to initializing a standalone postgres instance. This
# will not impact your system install ...
=== RUN   TestExample
page 0:
  header: {LSN:0 Checksum:0 Flags:0 Lower:36 Upper:8072 Special:8192 PageSizeVersion:8196 PruneXid:0}:
  tuple 1
    item identifier: {Offset:8152 Flags:1 Len:35}
    tuple header: {XMin:894 XMax:0 Field3:11 CTID:[0 0 0 0 1 0] Infomask2:2 Infomask:2050 Offset:24}
    data: 010000000f6974656d2d31
  tuple 2
    item identifier: {Offset:8112 Flags:1 Len:35}
    tuple header: {XMin:894 XMax:0 Field3:11 CTID:[0 0 0 0 2 0] Infomask2:2 Infomask:2050 Offset:24}
    data: 020000000f6974656d2d32
  tuple 3
    item identifier: {Offset:8072 Flags:1 Len:35}
    tuple header: {XMin:894 XMax:0 Field3:11 CTID:[0 0 0 0 3 0] Infomask2:2 Infomask:2050 Offset:24}
    data: 030000000f6974656d2d33
--- PASS: TestExample (0.01s)
```
