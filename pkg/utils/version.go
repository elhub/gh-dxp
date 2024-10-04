package utils

import (
	"encoding/json"
	"errors"
)

// GetLatestReleaseVersion fetches the tag of the latest release.
func GetLatestReleaseVersion(exe Executor) (string, error) {

	response, err := exe.GH("api", "-H", "Accept: application/vnd.github+json", "-H", "X-GitHub-Api-Version: 2022-11-28", "/repos/elhub/gh-dxp/releases/latest")
	if err != nil {
		return "", err
	}

	var deserializedResponse gitHubRelease

	err = json.NewDecoder(&response).Decode(&deserializedResponse)
	if err != nil {
		return "", err
	}

	if deserializedResponse.TagName == "" {
		return "", errors.New("response object contained no tag_name field")
	}

	return deserializedResponse.TagName, nil
}

// IsLatestVersion checks whether the provided version is the same as the latest release in GitHub.
func IsLatestVersion(exe Executor, localVersion string) (bool, error) {

	latestVersion, err := GetLatestReleaseVersion(exe)
	if err != nil {
		return false, err
	}
	return latestVersion == localVersion, nil
}

// IsLatestVersionOrSnapshot checks whether the provided versin is the same as the latest release in GitHub, OR "SNAPSHOT".
func IsLatestVersionOrSnapshot(exe Executor, version string) (bool, error) {
	if version == "SNAPSHOT" {
		return true, nil
	}
	return IsLatestVersion(exe, version)
}

type gitHubRelease struct {
	TagName string `json:"tag_name"`
}
