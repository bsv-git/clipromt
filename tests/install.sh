#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
stage=$(mktemp -d)
trap 'rm -rf "$stage"' EXIT
printf '# existing setup\n' > "$stage/bashrc"
for attempt in 1 2; do
    bash install.sh --from-dir "$PWD/dist" --prefix "$stage/bin with spaces" --shell --rc-file "$stage/bashrc" >/dev/null
done
[[ $(grep -Fc '# clipromt setup v0.1' "$stage/bashrc") == 1 ]]
grep -Fqx '# existing setup' "$stage/bashrc"
bash -n "$stage/bashrc"
[[ $("$stage/bin with spaces/clipromt" suggest 'who am i') == whoami ]]
mkdir "$stage/bad"
cp dist/SHA256SUMS "$stage/bad/"
for archive in dist/clipromt_0.1.0_linux_*.tar.gz; do
    cp "$archive" "$stage/bad/"
    printf corrupt >> "$stage/bad/$(basename "$archive")"
done
if bash install.sh --from-dir "$stage/bad" --prefix "$stage/rejected" >"$stage/error" 2>&1; then
    echo 'Corrupted archive accepted' >&2
    exit 1
fi
[[ ! -e $stage/rejected/clipromt ]]
while IFS= read -r request; do
    "$stage/bin with spaces/clipromt" suggest "$request" >/dev/null
done < examples/starter.txt
printf '%s\n' 'PASS: offline install, paths with spaces, idempotent setup, checksum rejection, starter catalog'
