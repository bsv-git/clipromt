#!/usr/bin/env bash
# Download and install the released binary. No Go, root, or AI service required.
set -euo pipefail
version=0.1.0
prefix="${HOME}/.local/bin"
rc_file="${HOME}/.bashrc"
setup_shell=false
from_dir=
while (($#)); do
    case "$1" in
        --shell) setup_shell=true; shift ;;
        --prefix|--rc-file|--from-dir)
            option=$1
            (($# >= 2)) || { echo "Missing value for $option" >&2; exit 2; }
            case "$option" in
                --prefix) prefix=$2 ;;
                --rc-file) rc_file=$2 ;;
                --from-dir) from_dir=$2 ;;
            esac
            shift 2 ;;
        --help|-h)
            echo 'Usage: bash install.sh [--shell] [--prefix DIRECTORY] [--rc-file FILE] [--from-dir RELEASE_DIRECTORY]'
            echo 'Installs v0.1.0; --shell adds an idempotent Bash startup block. --from-dir installs offline.'
            exit 0 ;;
        *) echo "Unknown option: $1" >&2; exit 2 ;;
    esac
done
[[ $(uname -s) == Linux ]] || { echo 'This release supports Linux only.' >&2; exit 1; }
case "$(uname -m)" in
    x86_64|amd64) arch=amd64 ;;
    aarch64|arm64) arch=arm64 ;;
    *) echo 'Supported architectures: x86_64 and aarch64.' >&2; exit 1 ;;
esac
for utility in tar sha256sum install mktemp cp grep; do
    command -v "$utility" >/dev/null || { echo "Missing utility: $utility" >&2; exit 1; }
done
stage=$(mktemp -d)
trap 'rm -rf "$stage"' EXIT
archive="clipromt_${version}_linux_${arch}.tar.gz"
fetch() {
    if [[ -n $from_dir ]]; then
        cp -- "$from_dir/$1" "$stage/$1"
    elif command -v curl >/dev/null; then
        curl --fail --location --silent --show-error "https://github.com/bsv-git/clipromt/releases/download/v${version}/$1" -o "$stage/$1"
    elif command -v wget >/dev/null; then
        wget -q "https://github.com/bsv-git/clipromt/releases/download/v${version}/$1" -O "$stage/$1"
    else
        echo 'Install curl or wget, or use --from-dir with downloaded release files.' >&2
        exit 1
    fi
}
fetch "$archive"
fetch SHA256SUMS
# Verify only the archive selected for this machine; other assets need not be downloaded.
checksum=
while read -r digest name; do
    if [[ $name == "$archive" && $digest =~ ^[0-9a-f]{64}$ ]]; then
        [[ -z $checksum ]] || { echo 'Duplicate checksum entry.' >&2; exit 1; }
        checksum=$digest
    fi
done < "$stage/SHA256SUMS"
[[ -n $checksum ]] || { echo 'Missing or invalid archive checksum.' >&2; exit 1; }
printf '%s  %s\n' "$checksum" "$archive" > "$stage/selected-checksum"
(cd "$stage" && sha256sum -c selected-checksum)
# Extract only the expected binary from the checked archive.
tar -xzf "$stage/$archive" -C "$stage" clipromt
mkdir -p -- "$prefix"
prefix=$(cd "$prefix" && pwd -P)
binary="$prefix/clipromt"
install -m 755 -- "$stage/clipromt" "$binary"
"$binary" version
if $setup_shell; then
    [[ ! -L $rc_file ]] || { echo 'Refusing a symlink rc file; configure Bash manually.' >&2; exit 1; }
    mkdir -p -- "$(dirname "$rc_file")"
    touch -- "$rc_file"
    if ! grep -Fqx '# clipromt setup v0.1' "$rc_file"; then
        {
            printf '\n# clipromt setup v0.1\n'
            printf 'if [[ $- == *i* && ! -o vi ]] && [[ -x %q ]]; then\n' "$binary"
            printf '    eval "$(%q init bash)"\n' "$binary"
            printf 'fi\n# end clipromt setup\n'
        } >> "$rc_file"
    fi
    printf 'Bash setup saved in %s. Open a new Bash terminal to use /c.\n' "$rc_file"
fi
printf '\nTry now:\n  %q catalog\n  %q suggest '\''who am i'\''\n' "$binary" "$binary"
printf '\nEnable /c in this terminal:\n  eval "$(%q init bash)"\n' "$binary"
