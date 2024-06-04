package manager_test

import (
	"os"
	"path/filepath"

	"github.com/sleepdeprecation/tfmgr/internal/manager"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test version selector", func() {
	var workingDir string
	var detector *manager.VersionDetector

	BeforeEach(func() {
		tmpDir, err := os.MkdirTemp("", "tfmgr-version")
		Expect(err).To(Succeed())

		cwd, err := os.Getwd()
		Expect(err).To(Succeed())

		workingDir = tmpDir

		Expect(os.Chdir(workingDir)).To(Succeed())

		detector = manager.NewDetector(workingDir)

		DeferCleanup(func() {
			Expect(os.Chdir(cwd)).To(Succeed())
			Expect(os.RemoveAll(workingDir)).To(Succeed())
		})
	})

	It("Detects the version from the .terraform-version file", func() {
		Expect(os.WriteFile(
			filepath.Join(workingDir, ".terraform-version"),
			[]byte("1.2.3"),
			0644)).To(Succeed())

		version, err := detector.CheckVersionFile()
		Expect(err).To(Succeed())
		Expect(version).To(Equal("1.2.3"))
	})
})
