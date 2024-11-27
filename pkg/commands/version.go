package commands

import (
	"fmt"

	"github.com/hashicorp/go-version"
)

func (r *GitRepoInfo) parseVersion(tag string) (major, minor, patch int, err error) {
	v, err := version.NewSemver(tag)

	if err != nil {
		return 0, 0, 0, err
	}

	// check if tag has not a prefix
	if tag[0] != 'v' {
		r.NoPrefix = true
	}

	vSegments := v.Segments()
	major = vSegments[0]
	minor = vSegments[1]
	patch = vSegments[2]
	return major, minor, patch, nil
}

func (r *GitRepoInfo) FormatTag(tag string) string {
	if r.NoPrefix {
		return tag
	}
	return "v" + tag
}

func (r *GitRepoInfo) CheckTagFormat(tag string) (err error) {
	_, err = version.NewSemver(tag)
	if err != nil {
		return err
	}
	return nil
}

func (r *GitRepoInfo) IncMajor() (newTag string, err error) {
	major, _, _, err := r.parseVersion(r.LastTag)
	if err != nil {
		return "", err
	}
	newTag = r.FormatTag(fmt.Sprintf("%d.%d.%d", major+1, 0, 0))
	return newTag, nil
}

func (r *GitRepoInfo) IncMinor() (newTag string, err error) {
	major, minor, _, err := r.parseVersion(r.LastTag)
	if err != nil {
		return "", err
	}
	newTag = r.FormatTag(fmt.Sprintf("%d.%d.%d", major, minor+1, 0))
	return newTag, nil
}

func (r *GitRepoInfo) IncPatch() (newTag string, err error) {
	major, minor, patch, err := r.parseVersion(r.LastTag)
	if err != nil {
		return "", err
	}
	newTag = r.FormatTag(fmt.Sprintf("%d.%d.%d", major, minor, patch+1))
	return newTag, nil
}

func (r *GitRepoInfo) CreatePreRelease(PreReleaseName string) (newTag string, err error) {
	major, minor, patch, err := r.parseVersion(r.LastTag)
	if err != nil {
		return "", err
	}
	newTag = r.FormatTag(fmt.Sprintf("%d.%d.%d-%s", major, minor, patch, PreReleaseName))
	return newTag, nil
}
