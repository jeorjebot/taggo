package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const changelogHeader = `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

`

const unreleasedSection = "## [Unreleased]\n"

// UpdateChangelog updates CHANGELOG.md when a new tag is created.
// It moves the [Unreleased] content into a new versioned section and updates comparison links.
// ScaffoldChangelog creates a new CHANGELOG.md with the standard header and empty [Unreleased] section.
func (r *GitRepoInfo) ScaffoldChangelog() error {
	changelogPath := filepath.Join(r.Path, "CHANGELOG.md")

	if _, err := os.Stat(changelogPath); err == nil {
		fmt.Println("[*] CHANGELOG.md already exists")
		return nil
	}

	content := changelogHeader + unreleasedSection
	err := os.WriteFile(changelogPath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write CHANGELOG.md: %w", err)
	}

	fmt.Println("[*] CHANGELOG.md created")
	return nil
}

func (r *GitRepoInfo) UpdateChangelog(newTag string) error {
	changelogPath := filepath.Join(r.Path, "CHANGELOG.md")

	// only update if CHANGELOG.md already exists
	if _, err := os.Stat(changelogPath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(changelogPath)
	if err != nil {
		return fmt.Errorf("failed to read CHANGELOG.md: %w", err)
	}

	content := ensureUnreleasedSection(string(data))
	content = insertVersionSection(content, newTag)
	content = updateComparisonLinks(content, newTag, r.LastTag, r.RemoteURL)

	err = os.WriteFile(changelogPath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write CHANGELOG.md: %w", err)
	}

	fmt.Println("[*] CHANGELOG.md updated: [Unreleased] -> [" + newTag + "]")
	return nil
}

// RevertChangelog moves the content of a deleted tag's section back into [Unreleased].
func (r *GitRepoInfo) RevertChangelog(deletedTag string) error {
	changelogPath := filepath.Join(r.Path, "CHANGELOG.md")

	data, err := os.ReadFile(changelogPath)
	if err != nil {
		if os.IsNotExist(err) {
			// no changelog, nothing to revert
			return nil
		}
		return fmt.Errorf("failed to read CHANGELOG.md: %w", err)
	}

	content := string(data)

	// find the deleted tag's section
	deletedHeader := fmt.Sprintf("## [%s]", deletedTag)
	deletedIdx := strings.Index(content, deletedHeader)
	if deletedIdx < 0 {
		// tag section not found in changelog, nothing to revert
		return nil
	}

	// find end of the deleted section header line
	afterDeletedHeader := deletedIdx + len(deletedHeader)
	nlIdx := strings.Index(content[afterDeletedHeader:], "\n")
	if nlIdx < 0 {
		return nil
	}
	afterDeletedLine := afterDeletedHeader + nlIdx + 1

	// find the next versioned section after the deleted one
	nextSectionIdx := strings.Index(content[afterDeletedLine:], "## [")

	var deletedContent string
	var afterDeleted string
	if nextSectionIdx >= 0 {
		deletedContent = content[afterDeletedLine : afterDeletedLine+nextSectionIdx]
		afterDeleted = content[afterDeletedLine+nextSectionIdx:]
	} else {
		// no next section — content until comparison links or end
		rest := content[afterDeletedLine:]
		linkIdx := findComparisonLinksStart(rest)
		if linkIdx >= 0 {
			deletedContent = rest[:linkIdx]
			afterDeleted = rest[linkIdx:]
		} else {
			deletedContent = rest
			afterDeleted = ""
		}
	}

	// find [Unreleased] section and append the deleted content to it
	unreleasedIdx := strings.Index(content, "## [Unreleased]")
	if unreleasedIdx < 0 {
		// no unreleased section, add one with the content
		content = content[:deletedIdx] + unreleasedSection + deletedContent + afterDeleted
	} else {
		// find end of [Unreleased] line
		afterUnreleased := unreleasedIdx + len("## [Unreleased]")
		unreleasedNl := strings.Index(content[afterUnreleased:], "\n")
		if unreleasedNl < 0 {
			return nil
		}
		afterUnreleasedLine := afterUnreleased + unreleasedNl + 1

		// get existing unreleased content (between [Unreleased] and the deleted section)
		existingUnreleased := content[afterUnreleasedLine:deletedIdx]

		// merge: existing unreleased + deleted section content + rest after deleted
		var sb strings.Builder
		sb.WriteString(content[:afterUnreleasedLine])
		// append existing unreleased content (trimmed)
		trimmedExisting := strings.TrimSpace(existingUnreleased)
		if trimmedExisting != "" {
			sb.WriteString("\n")
			sb.WriteString(trimmedExisting)
			sb.WriteString("\n")
		}
		// append deleted section content
		trimmedDeleted := strings.TrimSpace(deletedContent)
		if trimmedDeleted != "" {
			sb.WriteString("\n")
			sb.WriteString(trimmedDeleted)
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
		sb.WriteString(afterDeleted)
		content = sb.String()
	}

	// remove comparison link for the deleted tag and update [unreleased] link
	content = removeComparisonLink(content, deletedTag, r.RemoteURL)

	err = os.WriteFile(changelogPath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write CHANGELOG.md: %w", err)
	}

	fmt.Println("[*] CHANGELOG.md reverted: [" + deletedTag + "] -> [Unreleased]")
	return nil
}

// removeComparisonLink removes the link for a deleted tag and updates the [unreleased] link.
func removeComparisonLink(content, deletedTag, remoteURL string) string {
	if remoteURL == "" {
		return content
	}

	lines := strings.Split(content, "\n")
	var result []string
	var prevTag string

	// find what the deleted tag pointed to (its base) from its comparison link
	deletedLinkPrefix := "[" + deletedTag + "]: "
	for _, line := range lines {
		if strings.HasPrefix(line, deletedLinkPrefix) {
			// extract the base tag from compare/BASE...TAG
			if idx := strings.Index(line, "/compare/"); idx >= 0 {
				rest := line[idx+len("/compare/"):]
				if dotIdx := strings.Index(rest, "..."); dotIdx >= 0 {
					prevTag = rest[:dotIdx]
				}
			}
			continue // skip this link
		}
		result = append(result, line)
	}

	// update [unreleased] link to point to prevTag instead of deletedTag
	if prevTag != "" {
		for i, line := range result {
			lower := strings.ToLower(line)
			if strings.HasPrefix(lower, "[unreleased]:") {
				result[i] = fmt.Sprintf("[unreleased]: %s/compare/%s...HEAD", remoteURL, prevTag)
				break
			}
		}
	}

	return strings.Join(result, "\n")
}

func readOrScaffoldChangelog(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// scaffold a new changelog
			return changelogHeader + unreleasedSection, nil
		}
		return "", fmt.Errorf("failed to read CHANGELOG.md: %w", err)
	}
	return string(data), nil
}

func ensureUnreleasedSection(content string) string {
	if strings.Contains(content, "## [Unreleased]") {
		return content
	}

	// insert after the header block (find the first empty line after the preamble)
	// look for the first ## [ versioned section
	idx := strings.Index(content, "## [")
	if idx > 0 {
		return content[:idx] + unreleasedSection + "\n" + content[idx:]
	}

	// no versioned sections, append at the end
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return content + "\n" + unreleasedSection
}

func insertVersionSection(content, newTag string) string {
	date := time.Now().Format("2006-01-02")
	versionHeader := fmt.Sprintf("## [%s] - %s\n", newTag, date)

	unreleasedIdx := strings.Index(content, "## [Unreleased]")
	if unreleasedIdx < 0 {
		return content
	}

	// find the end of the [Unreleased] line
	afterUnreleased := unreleasedIdx + len("## [Unreleased]")
	// skip the rest of the line (e.g. trailing whitespace or newline)
	nlIdx := strings.Index(content[afterUnreleased:], "\n")
	if nlIdx < 0 {
		// [Unreleased] is the last line
		return content + "\n\n" + versionHeader
	}
	afterUnreleasedLine := afterUnreleased + nlIdx + 1

	// find the next versioned section ## [
	nextSectionIdx := strings.Index(content[afterUnreleasedLine:], "## [")

	var unreleasedContent string
	if nextSectionIdx >= 0 {
		unreleasedContent = content[afterUnreleasedLine : afterUnreleasedLine+nextSectionIdx]
	} else {
		// no next section — everything until comparison links or end
		rest := content[afterUnreleasedLine:]
		linkIdx := findComparisonLinksStart(rest)
		if linkIdx >= 0 {
			unreleasedContent = rest[:linkIdx]
		} else {
			unreleasedContent = rest
		}
	}

	// build the new content:
	// ## [Unreleased]\n\n## [newTag] - date\n<old unreleased content>...<rest>
	var sb strings.Builder
	sb.WriteString(content[:afterUnreleasedLine])
	sb.WriteString("\n")
	sb.WriteString(versionHeader)

	// write the unreleased content under the new version section
	trimmed := strings.TrimSpace(unreleasedContent)
	if trimmed != "" {
		sb.WriteString(unreleasedContent)
	} else {
		sb.WriteString("\n")
	}

	// write the rest (next sections + links)
	if nextSectionIdx >= 0 {
		sb.WriteString(content[afterUnreleasedLine+nextSectionIdx:])
	} else {
		rest := content[afterUnreleasedLine:]
		linkIdx := findComparisonLinksStart(rest)
		if linkIdx >= 0 {
			sb.WriteString(rest[linkIdx:])
		}
	}

	return sb.String()
}

// findComparisonLinksStart finds the start of the comparison links block
// (lines starting with [something]: http)
func findComparisonLinksStart(content string) int {
	lines := strings.Split(content, "\n")
	offset := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "[") && strings.Contains(line, "]: http") {
			return offset
		}
		offset += len(line) + 1
	}
	return -1
}

func updateComparisonLinks(content, newTag, prevTag, remoteURL string) string {
	if remoteURL == "" {
		return content
	}

	// remove existing comparison links
	lines := strings.Split(content, "\n")
	var contentLines []string
	var existingLinks []string
	for _, line := range lines {
		if strings.HasPrefix(line, "[") && strings.Contains(line, "]: "+remoteURL) {
			existingLinks = append(existingLinks, line)
		} else {
			contentLines = append(contentLines, line)
		}
	}

	// trim trailing empty lines from content
	for len(contentLines) > 0 && strings.TrimSpace(contentLines[len(contentLines)-1]) == "" {
		contentLines = contentLines[:len(contentLines)-1]
	}

	// build new links
	var newLinks []string

	// [unreleased] link
	newLinks = append(newLinks, fmt.Sprintf("[unreleased]: %s/compare/%s...HEAD", remoteURL, newTag))

	// [newTag] link
	if prevTag != "" {
		newLinks = append(newLinks, fmt.Sprintf("[%s]: %s/compare/%s...%s", newTag, remoteURL, prevTag, newTag))
	} else {
		newLinks = append(newLinks, fmt.Sprintf("[%s]: %s/releases/tag/%s", newTag, remoteURL, newTag))
	}

	// keep existing version links (skip the old unreleased and any link for newTag)
	for _, link := range existingLinks {
		lower := strings.ToLower(link)
		if strings.HasPrefix(lower, "[unreleased]") {
			continue
		}
		if strings.HasPrefix(link, "["+newTag+"]") {
			continue
		}
		newLinks = append(newLinks, link)
	}

	result := strings.Join(contentLines, "\n") + "\n\n" + strings.Join(newLinks, "\n") + "\n"
	return result
}
