#!/bin/bash
set -e

cmd="$1"
app="rjv"

if test ! "$cmd"; then
    echo "command required."
    echo
    echo "available commands:"
    # alphabetical order
    echo "  build               build project"
    echo "  build.all           build project, ignore cache"
    echo "  build.release       build project for distribution"
    echo "  clean               deletes all generated files"
    exit 1
fi

shift
rest=$*

if test "$cmd" = "build"; then
    ./manage.sh clean
    # CGO_ENABLED=0 skips CGO and linking against glibc to build static binaries.
    # -v 'verbose'
    CGO_ENABLED=0 go build \
        -v
    echo "wrote $app"
    exit 0

elif test "$cmd" = "build.all"; then
    # CGO_ENABLED=0 skips CGO and linking against glibc to build static binaries.
    # -a 'build all'
    # -v 'verbose'
    ./manage.sh clean
    CGO_ENABLED=0 go build \
        -a \
        -v
    exit 0

elif test "$cmd" = "build.release"; then
    # GOOS is 'Go OS' and is being explicit in which OS to build for.
    # CGO_ENABLED=0 skips CGO and linking against glibc to build static binaries.
    # ld -s is 'disable symbol table'
    # ld -w is 'disable DWARF generation'
    # -trimpath removes leading paths to source files
    # -v 'verbose'
    # -o 'output'
    set -u
    version="$1" # 1.0.0
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
        -ldflags "-s -w -X main.APP_VERSION=$version" \
        -trimpath \
        -v \
        -o linux-amd64
    upx linux-amd64
    sha256sum linux-amd64 > linux-amd64.sha256

    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build \
        -ldflags "-s -w -X main.APP_VERSION=$version" \
        -trimpath \
        -v \
        -o linux-arm64
    upx linux-arm64
    sha256sum linux-arm64 > linux-arm64.sha256
    echo ---
    go version
    echo ---
    upx --version
    echo ---
    du -sh linux-a*
    echo ---
    echo "done"
    exit 0

elif test "$cmd" = "clean"; then
    # -f 'force' don't fail if file doesn't exist.
    # -v 'verbose' print the name of the file that was deleted.
    tbd=(
        "main" "$app" # generated by 'build'
        "linux-amd64" "linux-amd64.sha256 linux-arm64" "linux-arm64.sha256" # generated by 'release'
    )
    rm -fv ${tbd[@]}
    exit 0

fi

echo "unknown command: $cmd"
exit 1
