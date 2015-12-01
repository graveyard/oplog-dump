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
  -query="": (Optional) grab all oplog entries satisfying the query

The query should be a Mongo query, for example the query '{ ns : { $ne : \"database.cooldocs\" } }' would retrieve everything in database `database` outside of the `cooldocs` collection, while '{ ns : \"database.cooldocs\" }' would grab exactly the `cooldocs` collection. Take a look at the oplog collection to see the structure of ops.

## Changing Dependencies

### New Packages

When adding a new package, you can simply use `make vendor` to update your imports.
This should bring in the new dependency that was previously undeclared.
The change should be reflected in [Godeps.json](Godeps/Godeps.json) as well as [vendor/](vendor/).

### Existing Packages

First ensure that you have your desired version of the package checked out in your `$GOPATH`.

When to change the version of an existing package, you will need to use the godep tool.
You must specify the package with the `update` command, if you use multiple subpackages of a repo you will need to specify all of them.
So if you use package github.com/Clever/foo/a and github.com/Clever/foo/b, you will need to specify both a and b, not just foo.

```
# depending on github.com/Clever/foo
godep update github.com/Clever/foo

# depending on github.com/Clever/foo/a and github.com/Clever/foo/b
godep update github.com/Clever/foo/a github.com/Clever/foo/b
```

