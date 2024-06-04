package downloader_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/sleepdeprecation/tfmgr/internal/downloader"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDownloader(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Downloader Suite")
}

var _ = Describe("Downloader", func() {
	Describe("Using the actual endpoints", func() {
		version := "1.4.0"

		dl := downloader.New()

		var downloadDir string

		BeforeEach(func() {
			var err error
			downloadDir, err = os.MkdirTemp("", "tfmgr-downloader")
			Expect(err).To(Succeed())

			DeferCleanup(func() {
				os.RemoveAll(downloadDir)
			})
		})

		It("Downloads a given terraform version", func() {
			release, err := dl.GetRelease(version)
			Expect(err).To(Succeed())

			Expect(dl.Download(release, downloadDir)).To(Succeed())

			cmd := exec.Command(filepath.Join(downloadDir, version, "terraform"), "-v")

			stdoutReader, err := cmd.StdoutPipe()
			Expect(err).To(Succeed())
			stderrReader, err := cmd.StderrPipe()
			Expect(err).To(Succeed())

			Expect(cmd.Start()).To(Succeed())

			stdout, err := io.ReadAll(stdoutReader)
			Expect(err).To(Succeed())
			stderr, err := io.ReadAll(stderrReader)
			Expect(err).To(Succeed())

			Expect(cmd.Wait()).To(Succeed())

			Expect(string(stdout)).To(HavePrefix(fmt.Sprintf("Terraform v%s\non %s_%s", version, runtime.GOOS, runtime.GOARCH)))
			Expect(string(stderr)).To(Equal(""))
		})
	})
})
