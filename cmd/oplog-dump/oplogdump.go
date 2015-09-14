package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/cenkalti/backoff"
	"gopkg.in/Clever/pathio.v1"
)

func main() {
	unixTime := flag.Int("time", 0, "Only get the entries greater than or equal to this unix timestamp")
	mongoUrl := flag.String("host", "localhost", "The mongo url")
	path := flag.String("path", "/dev/stdout", "The path to write the dump to")
	query := flag.String("query", "", "Query selector, e.g. '{ns: \"database_name.collection_name\"}'")
	flag.Parse()

	tempDir, err := ioutil.TempDir("/tmp", "systemCopier")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDir)

	if err := runDump(tempDir, *mongoUrl, *query, *unixTime); err != nil {
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
func runDump(dumpDir, host, userQuery string, unixTime int) error {
	individualQueries := []string{fmt.Sprintf("ts : { $gte : Timestamp(%d, 0) }", unixTime)}
	userQuery = strings.TrimSpace(userQuery)
	if strings.HasPrefix(userQuery, "{") && strings.HasSuffix(userQuery, "}") {
		userQuery = userQuery[1 : len(userQuery)-1] // trim outer curly braces
		individualQueries = append(individualQueries, userQuery)
	} else if userQuery != "" {
		return errors.New("Query must be deliniated by outer curly braces")
	}
	query := fmt.Sprintf("{ %s }", strings.Join(individualQueries, ", "))
	cmd := exec.Command("mongodump",
		"--db", "local",
		"--collection", "oplog.rs",
		"--out", dumpDir,
		"--host", host,
		"--query", query)
	// Forward stderr for logging / debugging. Note that we forward stdout to stderr so that it doesn't
	// polluate the command's return value (ie stdout)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stderr
	return cmd.Run()
}
