package commands

import (
	"fmt"

	"github.com/hashicorp/go-version"
)

func parseVersion(tag string) (major, minor, patch int, err error) {
	v, err := version.NewSemver(tag)
	if err != nil {
		return 0, 0, 0, err
	}
	vSegments := v.Segments()
	major = vSegments[0]
	minor = vSegments[1]
	patch = vSegments[2]
	return major, minor, patch, nil
}

func (r *GitRepoInfo) IncMajor() (newTag string, err error) {
	major, _, _, err := parseVersion(r.LastTag)
	if err != nil {
		return "", err
	}
	newTag = fmt.Sprintf("v%d.%d.%d", major+1, 0, 0)
	return newTag, nil
}

func (r *GitRepoInfo) IncMinor() (newTag string, err error) {
	major, minor, _, err := parseVersion(r.LastTag)
	if err != nil {
		return "", err
	}
	newTag = fmt.Sprintf("v%d.%d.%d", major, minor+1, 0)
	return newTag, nil
}

func (r *GitRepoInfo) IncPatch() (newTag string, err error) {
	major, minor, patch, err := parseVersion(r.LastTag)
	if err != nil {
		return "", err
	}
	newTag = fmt.Sprintf("v%d.%d.%d", major, minor, patch+1)
	return newTag, nil
}
