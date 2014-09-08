oplog-dump
==========
A binary that dumps the oplog from MongoDB.

This has some advantages over a direct mongodump command:
- Its interface is designed specifically for creating oplog dumps for a certain point in time.
- It can be used in conjunction with www.github.com/Clever/oplog-replay to run Mongo operations on multiple databases.


Usage
-----
Build Oplog Dump and put it in your GOPATH with:

`go get github.com/Clever/oplog-dump/cmd`

Run it as follows:
`oplog-dump --dir /tmp/out`

Params:
  -path="/dev/stdout": The path to write the dump to
  -mongoUrl="localhost": The URL to dump from
  -unixTime=0: Grab all oplog entries greater than or equal to this timestamp
