package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	unixTime := flag.Int("time", 0, "Only get the entries greater than or equal to this unix timestamp")
	mongoUrl := flag.String("host", "localhost", "The mongo url")
	dumpDir := flag.String("dir", "", "The directory where the data is dumped")
	flag.Parse()

	if len(*dumpDir) == 0 {
		log.Printf("Error: 'dir' not set")
		flag.PrintDefaults()
		os.Exit(2)
	}

	if err := runDump(*dumpDir, *mongoUrl, *unixTime); err != nil {
		// Try to return the same exit code as mongodump. This doesn't work on all platforms,
		// so if we can't figure out the exit code then we just use the exit code 2.
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		}
		os.Exit(2)
	}
}

// runDump runs the mongodump command. It's factored out so that it can be unit tested easily.
func runDump(dumpDir, host string, unixTime int) error {
	cmd := exec.Command("mongodump",
		"--db", "local",
		"--collection", "oplog.rs",
		"--out", dumpDir,
		"--host", host,
		"--query", fmt.Sprintf("{ts : { $gte : Timestamp(%d, 0)}}", unixTime))
	// Forward stderr for logging / debugging. Note that we forward stdout to stderr so that it doesn't
	// polluate the command's return value (ie stdout)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stderr
	return cmd.Run()
}
