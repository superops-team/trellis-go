package upgrade

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
)

// Checker checks for and downloads new versions from GitHub Releases.
type Checker struct {
	CurrentVersion string
	RepoOwner      string
	RepoName       string
	HTTPClient     *http.Client
	BaseURL        string // Default: https://api.github.com
}

// Release represents a GitHub release.
type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
	Prerelease bool `json:"prerelease"`
}

// Asset represents a release asset.
type Asset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
	Size        int64  `json:"size"`
}

// CheckLatest fetches the latest release matching the given tag filter.
// tag: "latest" (stable), "beta" (pre-release), or a specific version tag.
func (c *Checker) CheckLatest(tag string) (*Release, error) {
	client := c.httpClient()

	baseURL := c.BaseURL
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}
	var url string
	switch tag {
	case "latest":
		url = fmt.Sprintf("%s/repos/%s/%s/releases/latest", baseURL, c.RepoOwner, c.RepoName)
	case "beta":
		url = fmt.Sprintf("%s/repos/%s/%s/releases", baseURL, c.RepoOwner, c.RepoName)
	default:
		url = fmt.Sprintf("%s/repos/%s/%s/releases/tags/%s", baseURL, c.RepoOwner, c.RepoName, tag)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned %d: %s", resp.StatusCode, string(body))
	}

	if tag == "beta" {
		// Get the latest pre-release
		var releases []Release
		if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
			return nil, fmt.Errorf("decode releases: %w", err)
		}
		for _, r := range releases {
			if r.Prerelease {
				return &r, nil
			}
		}
		if len(releases) > 0 {
			return &releases[0], nil // Fallback to latest stable
		}
		return nil, fmt.Errorf("no releases found")
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("decode release: %w", err)
	}
	return &release, nil
}

// NeedsUpdate checks if the current version is older than the latest release.
func (c *Checker) NeedsUpdate(latest *Release) bool {
	current := strings.TrimPrefix(c.CurrentVersion, "v")
	latestVer := strings.TrimPrefix(latest.TagName, "v")
	return current != latestVer
}

// Download downloads the matching asset for the current OS/arch.
func (c *Checker) Download(release *Release, dest string) error {
	asset := c.findAsset(release)
	if asset == nil {
		return fmt.Errorf("no asset found for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	client := c.httpClient()
	req, err := http.NewRequest("GET", asset.DownloadURL, nil)
	if err != nil {
		return fmt.Errorf("create download request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("download asset: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned %d", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer out.Close()

	written, err := io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("write download: %w", err)
	}
	if written != asset.Size {
		return fmt.Errorf("downloaded %d bytes, expected %d", written, asset.Size)
	}

	if err := os.Chmod(dest, 0755); err != nil {
		return fmt.Errorf("chmod download: %w", err)
	}

	return nil
}

// Replace replaces the binary at target with the downloaded one.
func (c *Checker) Replace(downloaded, target string) error {
	src, err := os.Open(downloaded)
	if err != nil {
		return fmt.Errorf("open downloaded: %w", err)
	}
	defer src.Close()

	dest, err := os.OpenFile(target, os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("open target binary: %w", err)
	}
	defer dest.Close()

	if _, err := io.Copy(dest, src); err != nil {
		return fmt.Errorf("replace binary: %w", err)
	}
	return nil
}

func (c *Checker) findAsset(release *Release) *Asset {
	osArch := fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)
	for i := range release.Assets {
		if strings.Contains(release.Assets[i].Name, osArch) {
			return &release.Assets[i]
		}
	}
	// Fallback: first asset
	if len(release.Assets) > 0 {
		return &release.Assets[0]
	}
	return nil
}

func (c *Checker) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return http.DefaultClient
}
