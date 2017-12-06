package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"

	"github.com/swce/metadata-resource/models"
	"fmt"
	"bufio"
	"sort"
)

func main() {
	if len(os.Args) < 2 {
		fatalNoErr("usage: " + os.Args[0] + " <destination>")
	}

	destination := os.Args[1]

	log("creating destination dir " + destination)
	err := os.MkdirAll(destination, 0755)
	if err != nil {
		fatal("creating destination", err)
	}

	meta := make(models.Metadata)

	var request models.InRequest

	err = json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fatal("reading request", err)
	}

	var inVersion = request.TimestampVersion

	handleProp(destination, "build-id", "BUILD_ID", meta)
	handleProp(destination, "build-name", "BUILD_NAME", meta)
	handleProp(destination, "build-job-name", "BUILD_JOB_NAME", meta)
	handleProp(destination, "build-pipeline-name", "BUILD_PIPELINE_NAME", meta)
	handleProp(destination, "atc-external-url", "ATC_EXTERNAL_URL", meta)

	json.NewEncoder(os.Stdout).Encode(models.InResponse{
		Version:  inVersion,
		Metadata: meta
	})

	log("Done")
}

func fatal(doing string, err error) {
	fmt.Fprintln(os.Stderr, "error "+doing+": "+err.Error())
	os.Exit(1)
}

func log(doing string) {
	fmt.Fprintln(os.Stderr, doing)
}

func fatalNoErr(doing string) {
	log(doing)
	os.Exit(1)
}

func handleProp(destination string, filename string, prop string, meta models.Metadata) {
	output := filepath.Join(destination, filename)
	log("creating output file " + output)
	file, err := os.Create(output)
	if err != nil {
		fatal("creating output file "+output, err)
	}
	defer file.Close()

	val := os.Getenv(prop)
	meta[prop] = val
	w := bufio.NewWriter(file)
	fmt.Fprintf(w, "%s", val)

	err = w.Flush()

	if err != nil {
		fatal("writing output file"+output, err)
	}
}
