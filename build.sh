#!/usr/bin/env bash
set -euo pipefail

PROJECT_NAME="go-rename"
RELEASE_DIR="releases"
mkdir -p "${RELEASE_DIR}"

VERSION=$(awk -F'"' '/const Version = "/{gsub(/^v/,"",$2); print $2}' common/common.go)
VER_SUFFIX=".${VERSION}"

RAW=(
    darwin:amd64
    darwin:arm64
    linux:amd64
    linux:386
    linux:arm64
    linux:arm
    windows:amd64
    windows:386
    windows:arm64
)

NAME=(
    macos-amd64
    macos-arm64
    linux-x86-64
    linux-x86
    linux-arm64
    linux-arm
    windows-x86-64
    windows-x86
    windows-arm64
)

echo "Building ${PROJECT_NAME}${VER_SUFFIX} ..."

for i in "${!RAW[@]}"; do
    IFS=: read -r GOOS GOARCH <<< "${RAW[$i]}"
    EXT=""
    [[ $GOOS == windows ]] && EXT=.exe

    OUT_FILE="${PROJECT_NAME}${VER_SUFFIX}.${NAME[$i]}${EXT}"
    OUT_PATH="${RELEASE_DIR}/${OUT_FILE}"

    echo "  -> ${GOOS}/${GOARCH} => ${OUT_FILE}"
    CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" \
        go build -ldflags="-s -w" -trimpath -o "$OUT_PATH" .
done

echo "All done! Binaries are in '${RELEASE_DIR}/'"