package main

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"labix.org/v2/mgo"

	"github.com/stretchr/testify/assert"

	bsonScanner "github.com/Clever/oplog-replay/bson"
)

// simpleDocStruct is for inserting simple Mongo documents
type simpleDocStruct struct {
	key string
}

func TestComposingWithOplogReplay(t *testing.T) {
	// This test assumes that we're connecting to Mongo replica set. Otherwise, the oplog won't
	// be generated.

	unixTime := int(time.Now().Unix())

	time.Sleep(time.Duration(1) * time.Second)

	session, err := mgo.Dial("localhost")
	assert.Nil(t, err)
	db := session.DB("myTestDb")
	c := db.C("myCollection")
	assert.Nil(t, c.Insert(&simpleDocStruct{key: "key"}))

	time.Sleep(time.Duration(3) * time.Second)

	assert.Nil(t, c.Insert(&simpleDocStruct{key: "key2"}))
	// Dump at the time we started operations. Should get both operations
	dumpAtTime(t, unixTime, 2)
	// Three seconds later we should get only one.
	dumpAtTime(t, unixTime+3, 1)
}

func dumpAtTime(t *testing.T, unixTime, expectedResults int) {
	tempDir, err := ioutil.TempDir("/tmp", "oplogDumpTest")
	assert.Nil(t, err)
	defer os.RemoveAll(tempDir)
	assert.Nil(t, runDump(tempDir, "localhost", unixTime))

	file, err := os.Open(tempDir + "/local/oplog.rs.bson")
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
	assert.Nil(t, err)
	defer os.RemoveAll(tempDir)

	// Create a directory structure that mirrors the oplog one the code expects
	err = os.Mkdir(tempDir+"/local", 0744)
	assert.Nil(t, err)
	file, err := os.Create(tempDir + "/local/oplog.rs.bson")
	assert.Nil(t, err)
	assert.Nil(t, ioutil.WriteFile(file.Name(), []byte("test-bson-file"), 0644))

	// Create a file to copy to
	toFile, err := ioutil.TempFile(tempDir, "bsonCopyTest")
	assert.Nil(t, err)
	assert.Nil(t, copyBsonFile(tempDir, toFile.Name()))

	// Check that the data matches
	fileData, err := ioutil.ReadFile(toFile.Name())
	assert.Nil(t, err)
	assert.Equal(t, "test-bson-file", string(fileData))
}
