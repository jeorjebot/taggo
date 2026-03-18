# Taggo justfile

# Default: show available commands
default:
    @just --list

# Build taggo binary
build:
    go build -o taggo .

# Build with version from latest tag
build-versioned:
    go build -ldflags "-X main.version=$(git describe --tags --abbrev=0)" -o taggo .

# Install taggo via go install with version
install:
    go install -ldflags "-X main.version=$(git describe --tags --abbrev=0)" .

# Run tests
test:
    go test ./...

# Check release prerequisites (gh CLI, git clean state, tag exists)
release-check tag:
    #!/usr/bin/env bash
    set -e
    echo "[*] Checking release prerequisites for {{tag}}..."

    # Check gh CLI
    if command -v gh >/dev/null 2>&1; then
        echo "    gh CLI: $(gh --version | head -1)"
    else
        echo "    gh CLI: NOT FOUND - install with 'brew install gh'"
        exit 1
    fi

    # Check gh auth
    if gh auth status >/dev/null 2>&1; then
        echo "    gh auth: OK"
    else
        echo "    gh auth: NOT AUTHENTICATED - run 'gh auth login'"
        exit 1
    fi

    # Check go
    if command -v go >/dev/null 2>&1; then
        echo "    go: $(go version | awk '{print $3}')"
    else
        echo "    go: NOT FOUND"
        exit 1
    fi

    # Check tag exists
    if git rev-parse "{{tag}}" >/dev/null 2>&1; then
        echo "    tag {{tag}}: OK"
    else
        echo "    tag {{tag}}: NOT FOUND"
        exit 1
    fi

    # Check tag not already released
    if gh release view "{{tag}}" >/dev/null 2>&1; then
        echo "    release {{tag}}: ALREADY EXISTS"
        exit 1
    else
        echo "    release {{tag}}: not yet created (OK)"
    fi

    # Check working tree is clean
    if [ -z "$(git status --porcelain)" ]; then
        echo "    working tree: clean"
    else
        echo "    working tree: DIRTY - commit or stash changes first"
        exit 1
    fi

    echo "[*] All checks passed!"

# Create a GitHub release with binaries and changelog
release tag: (release-check tag)
    #!/usr/bin/env bash
    set -e
    echo ""
    echo "[*] Releasing {{tag}}..."

    # Extract changelog for this tag
    CHANGELOG=$(awk '/^## \[{{tag}}\]/{found=1; next} /^## \[/{if(found) exit} found{print}' CHANGELOG.md)
    if [ -z "${CHANGELOG}" ]; then
        echo "[*] Warning: no changelog entry found for {{tag}}"
        CHANGELOG="Release {{tag}}"
    else
        echo "[*] Changelog:"
        echo "${CHANGELOG}"
    fi

    # Build binaries from the tagged commit
    BUILD_DIR="$(pwd)/dist"
    rm -rf "${BUILD_DIR}"
    mkdir -p "${BUILD_DIR}"

    WORK_DIR=$(mktemp -d)
    echo "[*] Checking out {{tag}} in ${WORK_DIR}..."
    git worktree add --detach "${WORK_DIR}" "{{tag}}" --quiet

    PLATFORMS=("darwin/arm64" "darwin/amd64" "linux/arm64" "linux/amd64" "windows/arm64" "windows/amd64")

    echo ""
    echo "[*] Building binaries from {{tag}} commit..."
    for PLATFORM in "${PLATFORMS[@]}"; do
        GOOS="${PLATFORM%/*}"
        GOARCH="${PLATFORM#*/}"
        OSNAME="${GOOS}"
        if [ "${GOOS}" = "darwin" ]; then
            OSNAME="macos"
        fi
        OUTPUT="${BUILD_DIR}/taggo-{{tag}}-${OSNAME}-${GOARCH}"
        if [ "${GOOS}" = "windows" ]; then
            OUTPUT="${OUTPUT}.exe"
        fi
        echo "    ${GOOS}/${GOARCH}"
        cd "${WORK_DIR}" && GOOS="${GOOS}" GOARCH="${GOARCH}" go build \
            -ldflags "-X main.version={{tag}}" \
            -o "${OUTPUT}" .
    done

    # Cleanup worktree
    cd - >/dev/null
    git worktree remove "${WORK_DIR}" --force

    # Create release
    echo ""
    echo "[*] Creating GitHub release..."
    gh release create "{{tag}}" \
        "${BUILD_DIR}"/taggo-"{{tag}}"-* \
        --title "{{tag}}" \
        --notes "${CHANGELOG}"

    echo "[*] Release created: https://github.com/jeorjebot/taggo/releases/tag/{{tag}}"

    # Print checksums for Homebrew Cask
    echo ""
    echo "[*] SHA256 checksums (for Homebrew Cask):"
    for f in "${BUILD_DIR}"/taggo-"{{tag}}"-darwin-*; do
        echo "    $(basename "$f"): $(shasum -a 256 "$f" | awk '{print $1}')"
    done
