package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/Clever/pathio"
	"github.com/cenkalti/backoff"
)

func main() {
	unixTime := flag.Int("time", 0, "Only get the entries greater than or equal to this unix timestamp")
	mongoUrl := flag.String("host", "localhost", "The mongo url")
	path := flag.String("path", "/dev/stdout", "The path to write the dump to")
	collection := flag.String("collection", "", "Collection selector")
	flag.Parse()
	fmt.Println(collection)

	tempDir, err := ioutil.TempDir("/tmp", "systemCopier")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDir)

	if err := runDump(tempDir, *mongoUrl, *collection, *unixTime); err != nil {
		// Try to return the same exit code as mongodump. This doesn't work on all platforms,
		// so if we can't figure out the exit code then we just use the exit code 2.
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		}
		os.Exit(2)
	}
	if err := writeWithRetry(tempDir, *path); err != nil {
		panic(err)
	}
}

// writeWithRetry writes to the specified path, retrying a few times on error.
func writeWithRetry(tempDir, destination string) error {
	backoffObj := backoff.ExponentialBackOff{
		InitialInterval:     5 * time.Second,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          2,
		MaxInterval:         30 * time.Second,
		MaxElapsedTime:      2 * time.Minute,
		Clock:               backoff.SystemClock,
	}
	operation := func() error {
		return copyBsonFile(tempDir, destination)
	}
	return backoff.Retry(operation, &backoffObj)
}

// copyBsonFile copies the file from the dump directory to the specified location.
func copyBsonFile(tempDir, destination string) error {
	file, err := os.Open(tempDir + "/local/oplog.rs.bson")
	if err != nil {
		return err
	}
	stats, err := file.Stat()
	if err != nil {
		return err
	}
	return pathio.WriteReader(destination, file, stats.Size())

}

// runDump runs the mongodump command. It's factored out so that it can be unit tested easily.
func runDump(dumpDir, host, collection string, unixTime int) error {
	queries := []string{fmt.Sprintf("ts : { $gte : Timestamp(%d, 0) }", unixTime)}
	if collection != "" {
		if !strings.Contains(collection, "{") {
			collection = fmt.Sprintf("\"%s\"", collection)
		}
		queries = append(queries, fmt.Sprintf("ns : %s", collection))
	}
	cmd := exec.Command("mongodump",
		"--db", "local",
		"--collection", "oplog.rs",
		"--out", dumpDir,
		"--host", host,
		"--query", fmt.Sprintf("{ %s }", strings.Join(queries, ", ")))
	// Forward stderr for logging / debugging. Note that we forward stdout to stderr so that it doesn't
	// polluate the command's return value (ie stdout)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stderr
	return cmd.Run()
}
