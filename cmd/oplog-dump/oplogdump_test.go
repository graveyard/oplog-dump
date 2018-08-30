package main

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/stretchr/testify/assert"

	bsonScanner "github.com/Clever/oplog-replay/bson"
)

// simpleDocStruct is for inserting simple Mongo documents
type simpleDocStruct struct {
	key string
}

func getPaddedTime() int {
	time.Sleep(time.Duration(1) * time.Second)
	n := int(time.Now().Unix())
	time.Sleep(time.Duration(2) * time.Second)
	return n
}

func TestComposingWithOplogReplay(t *testing.T) {
	// This test assumes that we're connecting to Mongo replica set. Otherwise, the oplog won't
	// be generated.
	unixTime := getPaddedTime()
	session, err := mgo.Dial("localhost")
	if err != nil {
		t.Log(err.Error())
	}

	assert.Nil(t, err)
	db := session.DB("myTestDb")
	c := db.C("myCollection")
	assert.Nil(t, c.Insert(&simpleDocStruct{key: "key"}))

	time.Sleep(time.Duration(5) * time.Second)

	assert.Nil(t, c.Insert(&simpleDocStruct{key: "key2"}))
	// Dump at the time we started operations. Should get both operations
	dumpAtTime(t, unixTime, 2, "{ns: {$ne: \"myTestDb.$cmd\"}}")
	// 5 seconds later we should get only one.
	dumpAtTime(t, unixTime+5, 1, "")
}

func TestCollectionFiltering(t *testing.T) {
	// This test assumes that we're connecting to Mongo replica set. Otherwise, the oplog won't
	// be generated.

	unixTime := getPaddedTime()
	session, err := mgo.Dial("localhost")

	if err != nil {
		t.Log(err.Error())
	}
	assert.Nil(t, err)
	db := session.DB("myTestDb")
	c1, c2, c3 := db.C("myCollection"), db.C("yourCollection"), db.C("ourCollection")

	assert.Nil(t, c1.Insert(&simpleDocStruct{key: "key"}))
	assert.Nil(t, c2.Insert(&simpleDocStruct{key: "key2"}))
	assert.Nil(t, c3.Insert(&simpleDocStruct{key: "key3"}))

	dumpAtTime(t, unixTime, 3, "{ns: {$ne: \"myTestDb.$cmd\"}}")
	dumpAtTime(t, unixTime, 1, "{ns: \"myTestDb.myCollection\"}")
	dumpAtTime(t, unixTime, 2, "{ns: {$nin : [\"myTestDb.myCollection\", \"myTestDb.$cmd\"]}}")
}

func dumpAtTime(t *testing.T, unixTime, expectedResults int, query string) {
	tempDir, err := ioutil.TempDir("/tmp", "oplogDumpTest")

	if err != nil {
		t.Log(err.Error())
	}
	assert.Nil(t, err)
	defer os.RemoveAll(tempDir)
	err = runDump(tempDir, "localhost", query, unixTime)
	if err != nil {
		t.Log(err.Error())
	}
	assert.Nil(t, err)

	file, err := os.Open(tempDir + "/local/oplog.rs.bson")

	if err != nil {
		t.Log(err.Error())
	}
	assert.Nil(t, err)
	scanner := bsonScanner.New(file)
	count := 0
	for scanner.Scan() {
		count = count + 1
	}
	assert.Equal(t, expectedResults, count)
}

func TestCopyBsonFile(t *testing.T) {
	tempDir, err := ioutil.TempDir("/tmp", "oplogDumpTest")

	if err != nil {
		t.Log(err.Error())
	}
	assert.Nil(t, err)
	defer os.RemoveAll(tempDir)

	// Create a directory structure that mirrors the oplog one the code expects
	err = os.Mkdir(tempDir+"/local", 0744)

	if err != nil {
		t.Log(err.Error())
	}
	assert.Nil(t, err)
	file, err := os.Create(tempDir + "/local/oplog.rs.bson")

	if err != nil {
		t.Log(err.Error())
	}
	assert.Nil(t, err)
	assert.Nil(t, ioutil.WriteFile(file.Name(), []byte("test-bson-file"), 0644))

	// Create a file to copy to
	toFile, err := ioutil.TempFile(tempDir, "bsonCopyTest")

	if err != nil {
		t.Log(err.Error())
	}
	assert.Nil(t, err)
	assert.Nil(t, copyBsonFile(tempDir, toFile.Name()))

	// Check that the data matches
	fileData, err := ioutil.ReadFile(toFile.Name())

	if err != nil {
		t.Log(err.Error())
	}
	assert.Nil(t, err)
	assert.Equal(t, "test-bson-file", string(fileData))
}
