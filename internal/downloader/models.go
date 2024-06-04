package downloader

import (
	"fmt"
	"runtime"
)

type TerraformRelease struct {
	Builds            []*TerraformBuild       `json:"builds"`
	IsPrerelease      bool                    `json:"is_prerelease"`
	Name              string                  `json:"name"`
	Status            *TerraformReleaseStatus `json:"status"`
	CreatedAt         string                  `json:"timestamp_created"`
	UpdatedAt         string                  `json:"timestamp_updated"`
	ChangelogUrl      string                  `json:"url_changelog"`
	ShasumsUrl        string                  `json:"url_shasums"`
	ShasumsSignatures []string                `json:"url_shasums_signatures"`
	Version           string                  `json:"version"`
}

// MyBuild returns the build information of a release that matches your current system
func (r *TerraformRelease) MyBuild() (*TerraformBuild, error) {
	for _, build := range r.Builds {
		if build.Os == runtime.GOOS && build.Arch == runtime.GOARCH {
			return build, nil
		}
	}

	return nil, fmt.Errorf("Couldn't find build for %s-%s", runtime.GOOS, runtime.GOARCH)
}

type TerraformBuild struct {
	Arch string `json:"arch"`
	Os   string `json:"os"`
	Url  string `json:"url"`
}

type TerraformReleaseStatus struct {
	State     string `json:"state"`
	UpdatedAt string `json:"timestamp_updated"`
}

type ReleaseErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r *ReleaseErrorResponse) Error() string {
	return fmt.Sprintf("%d: %s", r.Code, r.Message)
}
