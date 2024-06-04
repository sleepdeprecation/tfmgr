package downloader

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

const (
	DefaultReleaseApi = "https://api.releases.hashicorp.com"

	// status codes from the hashicorp api
	StatusReleaseNotFound = 30006
)

type Downloader struct {
	ReleaseApi string
	Client     *http.Client
}

func New() *Downloader {
	return &Downloader{
		ReleaseApi: DefaultReleaseApi,
		Client:     http.DefaultClient,
	}
}

// Download downloads the correct build for your current system into $cacheDir/$version/terraform.
func (d *Downloader) Download(release *TerraformRelease, cacheDir string) error {
	build, err := release.MyBuild()
	if err != nil {
		return err
	}

	downloadDir := filepath.Join(cacheDir, release.Version)
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return err
	}

	resp, err := d.Client.Get(build.Url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error downloading file: %d: %s", resp.StatusCode, string(body))
	}

	arch, err := zip.NewReader(bytes.NewReader(body), resp.ContentLength)
	if err != nil {
		return err
	}

	exec, err := arch.Open("terraform")
	if err != nil {
		return err
	}

	stat, err := exec.Stat()
	if err != nil {
		return err
	}

	execBody, err := io.ReadAll(exec)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(downloadDir, "terraform"), execBody, stat.Mode())
}

func (d *Downloader) GetReleases() ([]*TerraformRelease, error) {
	params := url.Values{}
	params.Set("limit", "20")

	releases := []*TerraformRelease{}

	for true {
		resp, err := d.Client.Get(fmt.Sprintf("%s/v1/releases/terraform?%s", d.ReleaseApi, params.Encode()))
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		dec := json.NewDecoder(resp.Body)

		if resp.StatusCode != http.StatusOK {
			errMessage := &ReleaseErrorResponse{}
			err := dec.Decode(errMessage)
			if err != nil {
				return nil, err
			}

			return nil, errMessage
		}

		currentReleases := []*TerraformRelease{}
		err = dec.Decode(&currentReleases)
		if err != nil {
			return nil, err
		}

		releases = append(releases, currentReleases...)

		if len(currentReleases) != 20 {
			break
		}

		params.Set("after", currentReleases[19].CreatedAt)
	}

	return releases, nil
}

func (d *Downloader) GetRelease(version string) (*TerraformRelease, error) {
	uri := fmt.Sprintf("%s/v1/releases/terraform/%s", d.ReleaseApi, version)
	resp, err := d.Client.Get(uri)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		errMessage := &ReleaseErrorResponse{}
		err := dec.Decode(errMessage)
		if err != nil {
			return nil, err
		}

		return nil, errMessage
	}

	release := &TerraformRelease{}
	err = dec.Decode(release)
	return release, err
}
