package main_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/swce/metadata-resource/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"path/filepath"
	"github.com/magiconair/properties"
)

var _ = Describe("In", func() {
	var tmpdir string
	var destination string

	var inCmd *exec.Cmd

	BeforeEach(func() {
		var err error

		tmpdir, err = ioutil.TempDir("", "in-destination")
		Expect(err).NotTo(HaveOccurred())

		destination = path.Join(tmpdir, "in-dir")

		inCmd = exec.Command(inPath, destination)

		inCmd.Env = append(
			"BUILD_ID=1",
			"BUILD_NAME=2",
			"BUILD_JOB_NAME=3",
			"BUILD_PIPELINE_NAME=4",
			"ATC_EXTERNAL_URL=5"
		)
	})

	AfterEach(func() {
		os.RemoveAll(tmpdir)
	})

	Context("when executed", func() {
		var request models.InRequest
		var response models.InResponse

		BeforeEach(func() {

			request = models.InRequest{
				Version: models.TimestampVersion{
					"version": "1"
				},
				Source: models.Source{},
			}

			response = models.InResponse{}
		})

		JustBeforeEach(func() {
			stdin, err := inCmd.StdinPipe()
			Expect(err).NotTo(HaveOccurred())

			session, err := gexec.Start(inCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			err = json.NewEncoder(stdin).Encode(request)
			Expect(err).NotTo(HaveOccurred())

			<-session.Exited
			Expect(session.ExitCode()).To(Equal(0))

			err = json.Unmarshal(session.Out.Contents(), &response)
			Expect(err).NotTo(HaveOccurred())
		})

		It("reports the version to be the input version", func() {
			Expect(len(response.Version)).To(Equal(1))
			Expect(response.Version["version"]).To(Equal("1"))
		})

		It("writes the requested data the destination", func() {

			checkProp(destination, "build-id", "BUILD_ID", "1", response.Metadata)
			checkProp(destination, "build-name", "BUILD_NAME", "2", response.Metadata)
			checkProp(destination, "build-job-name", "BUILD_JOB_NAME", "3", response.Metadata)
			checkProp(destination, "build-pipeline-name", "BUILD_PIPELINE_NAME", "4", response.Metadata)
			checkProp(destination, "atc-external-url", "ATC_EXTERNAL_URL", "5", response.Metadata)
		})

	})
})

func checkProp(destination string, filename string, prop string, valueToCheck string, meta models.Metadata) {
	output := filepath.Join(destination, filename)
	file, err := ioutil.ReadFile(output)
	if err != nil {
		fatal("reading output file "+output, err)
	}
	defer file.Close()
	val := string(file)
	Expect(val).To(Equal(valueToCheck))
	Expect(meta[prop]).To(Equal(valueToCheck))
}
