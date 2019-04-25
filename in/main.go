package main

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"

	"github.com/swce/metadata-resource/models"
	"fmt"
	"bufio"
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

	meta := make(models.Metadata,8)

	var request models.InRequest

	err = json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fatal("reading request", err)
	}

	var inVersion = request.Version

	// pretend that Concourse set these metadata envvars (maybe it would be
	// a good idea to implement them on the Concourse side?)
	os.Setenv("BUILD_URL", fmt.Sprintf("%s/teams/%s/pipelines/%s/jobs/%s/builds/%s",
		url.PathEscape(os.Getenv("ATC_EXTERNAL_URL")),
		url.PathEscape(os.Getenv("BUILD_TEAM_NAME")),
		url.PathEscape(os.Getenv("BUILD_PIPELINE_NAME")),
		url.PathEscape(os.Getenv("BUILD_JOB_NAME")),
		url.PathEscape(os.Getenv("BUILD_NAME"))))
	os.Setenv("BUILD_URL_SHORT", fmt.Sprintf("%s/builds/%s",
		url.PathEscape(os.Getenv("ATC_EXTERNAL_URL")),
		url.PathEscape(os.Getenv("BUILD_ID"))))

	handleProp(destination, "build-id", "BUILD_ID", meta, 0)
	handleProp(destination, "build-name", "BUILD_NAME", meta, 1)
	handleProp(destination, "build-job-name", "BUILD_JOB_NAME", meta, 2)
	handleProp(destination, "build-pipeline-name", "BUILD_PIPELINE_NAME", meta, 3)
	handleProp(destination, "build-team-name", "BUILD_TEAM_NAME", meta, 4)
	handleProp(destination, "atc-external-url", "ATC_EXTERNAL_URL", meta, 5)
	handleProp(destination, "build-url", "BUILD_URL", meta, 6)
	handleProp(destination, "build-url-short", "BUILD_URL_SHORT", meta, 7)

	json.NewEncoder(os.Stdout).Encode(models.InResponse{
		Version:  inVersion,
		Metadata: meta,
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

func handleProp(destination string, filename string, prop string, meta models.Metadata, index int) {
	output := filepath.Join(destination, filename)
	log("creating output file " + output)
	file, err := os.Create(output)
	if err != nil {
		fatal("creating output file "+output, err)
	}
	defer file.Close()

	val := os.Getenv(prop)
	meta[index] = models.MetadataField{
		Name: prop,
		Value: val,
	}
	w := bufio.NewWriter(file)
	fmt.Fprintf(w, "%s", val)

	err = w.Flush()

	if err != nil {
		fatal("writing output file"+output, err)
	}
}
