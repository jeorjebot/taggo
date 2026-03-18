#!/bin/bash
set -e

# Release script for Taggo
# Builds binaries for all platforms, extracts changelog for the tag,
# and creates a GitHub release with the binaries attached.
#
# Prerequisites: gh CLI (https://cli.github.com/)
#
# Usage:
#   ./scripts/release.sh          # release the latest tag
#   ./scripts/release.sh v1.2.0   # release a specific tag

# Determine tag
TAG="${1:-$(git describe --tags --abbrev=0)}"
echo "[*] Releasing ${TAG}"

# Verify tag exists
if ! git rev-parse "${TAG}" >/dev/null 2>&1; then
    echo "==> Error: tag ${TAG} does not exist"
    exit 1
fi

# Check gh is installed
if ! command -v gh >/dev/null 2>&1; then
    echo "==> Error: gh CLI is required. Install it from https://cli.github.com/"
    exit 1
fi

# Extract changelog for this tag from CHANGELOG.md
# Grabs everything between ## [TAG] and the next ## [ section
CHANGELOG=$(awk "/^## \\[${TAG}\\]/{found=1; next} /^## \\[/{if(found) exit} found{print}" CHANGELOG.md)

if [ -z "${CHANGELOG}" ]; then
    echo "[*] Warning: no changelog entry found for ${TAG}, using empty description"
    CHANGELOG="Release ${TAG}"
fi

echo "[*] Changelog:"
echo "${CHANGELOG}"
echo ""

# Build binaries from the tagged commit
BUILD_DIR="$(pwd)/dist"
rm -rf "${BUILD_DIR}"
mkdir -p "${BUILD_DIR}"

WORK_DIR=$(mktemp -d)
echo "[*] Checking out ${TAG} in ${WORK_DIR}..."
git worktree add --detach "${WORK_DIR}" "${TAG}" --quiet

PLATFORMS=(
    "darwin/arm64"
    "darwin/amd64"
    "linux/arm64"
    "linux/amd64"
    "windows/arm64"
    "windows/amd64"
)

echo "[*] Building binaries from ${TAG} commit..."
for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS="${PLATFORM%/*}"
    GOARCH="${PLATFORM#*/}"
    OSNAME="${GOOS}"
    if [ "${GOOS}" = "darwin" ]; then
        OSNAME="macos"
    fi
    OUTPUT="${BUILD_DIR}/taggo-${TAG}-${OSNAME}-${GOARCH}"
    if [ "${GOOS}" = "windows" ]; then
        OUTPUT="${OUTPUT}.exe"
    fi

    echo "    ${GOOS}/${GOARCH} -> ${OUTPUT}"
    (cd "${WORK_DIR}" && GOOS="${GOOS}" GOARCH="${GOARCH}" go build \
        -ldflags "-X main.version=${TAG}" \
        -o "${OUTPUT}" .)
done

# Cleanup worktree
git worktree remove "${WORK_DIR}" --force

echo "[*] Binaries built:"
ls -lh "${BUILD_DIR}/"
echo ""

# Create GitHub release
echo "[*] Creating GitHub release for ${TAG}..."
gh release create "${TAG}" \
    "${BUILD_DIR}"/taggo-"${TAG}"-* \
    --title "${TAG}" \
    --notes "${CHANGELOG}"

echo "[*] Release created: https://github.com/jeorjebot/taggo/releases/tag/${TAG}"

# Print sha256 checksums (useful for updating Homebrew Cask)
echo ""
echo "[*] SHA256 checksums (for Homebrew Cask):"
for f in "${BUILD_DIR}"/taggo-"${TAG}"-macos-*; do
    echo "    $(basename "$f"): $(shasum -a 256 "$f" | awk '{print $1}')"
done

# Cleanup build artifacts
rm -rf "${BUILD_DIR}"
echo ""
echo "[*] Done!"
