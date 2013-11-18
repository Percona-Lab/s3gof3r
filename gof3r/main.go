// Command gof3r provides a command-line interface to Amazon AWS S3.
//
// Usage:
//   To upload a file to S3:
//      gof3r  --up --file_path=<file_path> --url=<public_url> -h<http_header1> -h<http_header2>...
//   To download a file from S3:
//      gof3r  --down --file_path=<file_path> --url=<public_url>
//
//   The file does not need to be seekable or stat-able.
//
//   Examples:
//     $ gof3r  --up --file_path=test_file --url=https://bucket1.s3.amazonaws.com/object -hx-amz-meta-custom-metadata:123 -hx-amz-meta-custom-metadata2:123abc -hx-amz-server-side-encryption:AES256 -hx-amz-storage-class:STANDARD
//     $ gof3r  --down --file_path=test_file --url=https://bucket1.s3.amazonaws.com/object
//
// Environment:
//
// AwS_ACCESS_KEY – an AWS Access Key Id (required)
//
// AWS_SECRET_KEY – an AWS Secret Access Key (required)
//
// Complete Usage:
//  gof3r [OPTIONS]
//
// Help Options:
//  -h, --help=      Show this help message
//
// Application Options:
//      --up         Upload to S3
//      --down       Download from S3
//  -f, --file_path= canonical path to file
//  -u, --url=       Url of S3 object
//  -h, --headers=   HTTP headers ({})
//  -c, --md5-checking   Verify integrity with  md5 checksum
package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

// Options common to both puts and gets
type CommonOpts struct {
	//Url         string      `short:"u" long:"url" description:"Url of S3 object"` //TODO: bring back url support
	Key          string      `long:"key" description:"key of s3 object"`
	Bucket       string      `long:"bucket" description:"s3 bucket"`
	Header       http.Header `long:"header" short:"m" description:"HTTP headers"`
	CheckDisable bool        `long:"md5Check-off" description:"Do not use md5 hash checking to ensure data integrity. By default, the md5 hash of is calculated concurrently during puts, stored at <bucket>.md5/<key>.md5, and verified on gets."`
	Concurrency  int         `long:"concurrency" short:"c" default:"20" description:"Concurrency of transfers"`
	PartSize     int64       `long:"partsize" short:"s" description:"initial size of concurrent parts, in bytes" default:"20 MB"`
	Debug        bool        `long:"debug" description:"Print debug statements and dump stacks."`
}

var parser = flags.NewParser(nil, flags.Default)

func main() {

	start := time.Now()
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
	log.Println("Duration:", time.Since(start))
}

func init() {
	// set the number of processes to the number of cpus for parallelization of transfers
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func debug() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Println("MEMORY STATS")
	log.Println(fmt.Printf("%d,%d,%d,%d\n", m.HeapSys, m.HeapAlloc, m.HeapIdle, m.HeapReleased))
	log.Println("NUM CPU:", runtime.NumCPU())

	//profiling
	f, err := os.Create("memprofileup.out")
	fg, err := os.Create("goprof.out")
	fb, err := os.Create("blockprof.out")
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(f)
	pprof.Lookup("goroutine").WriteTo(fg, 0)
	pprof.Lookup("block").WriteTo(fb, 0)
	f.Close()
	fg.Close()
	fb.Close()
	time.Sleep(1 * time.Second)
	panic("Dump the stacks:")
}
