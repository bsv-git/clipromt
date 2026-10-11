# clipromt v0.1.0

Type `/c` and a supported English request in Bash. Press Enter to get a shell command in the same input line. Edit it, check the paths, and press Enter again when you want to run it.

clipromt runs locally. The binary uses only the Go standard library. It has no AI model, API key, network calls, telemetry, or background service. It generates text and never launches that text.

## Install

No repository clone, Go compiler, AI setup, or root access is needed. Download and run the installer:

```bash
curl -fsSL https://raw.githubusercontent.com/bsv-git/clipromt/main/install.sh -o /tmp/clipromt-install.sh
bash /tmp/clipromt-install.sh --shell
```

The installer selects Linux amd64 or arm64, downloads v0.1.0, verifies the selected archive's checksum, and installs to `~/.local/bin`. `--shell` adds Bash initialization to `~/.bashrc` without replacing its existing content. Open a new Bash terminal, or enable it in your current terminal:

```bash
eval "$(~/.local/bin/clipromt init bash)"
```

Without curl, use `wget -O /tmp/clipromt-install.sh https://raw.githubusercontent.com/bsv-git/clipromt/main/install.sh`. Omit `--shell` if you only want the binary. The installer downloads files during installation; the installed CLI needs no network. The downloaded installer can be inspected before running it.

Prefer a manual download? Get the archive for your machine and `SHA256SUMS` from [v0.1.0 downloads](https://github.com/bsv-git/clipromt/releases/tag/v0.1.0). GitHub's “Download ZIP” is source code, not the ready-to-run binary. From the folder containing those release files, an offline install is:

```bash
bash /tmp/clipromt-install.sh --from-dir "$PWD" --shell
```

Or extract and install manually:

For a local amd64 archive (use `arm64` for an ARM machine):

```bash
# Run in the folder containing the downloaded release files.
# Check only the archive you downloaded.
awk '$2 == "clipromt_0.1.0_linux_amd64.tar.gz"' SHA256SUMS | sha256sum -c -
mkdir -p /tmp/clipromt-install
# Extract only after the checksum passes.
tar -xzf clipromt_0.1.0_linux_amd64.tar.gz -C /tmp/clipromt-install
mkdir -p "$HOME/.local/bin"
install -m 755 /tmp/clipromt-install/clipromt "$HOME/.local/bin/clipromt"
export PATH="$HOME/.local/bin:$PATH"
clipromt version
clipromt doctor
```

Checksums detect changed files; they do not authenticate a publisher. No signature is supplied for this initial local package.

Or build with Go 1.22 or newer:

```bash
go build -buildvcs=false -trimpath -o bin/clipromt ./cmd/clipromt
mkdir -p "$HOME/.local/bin"
install -m 755 bin/clipromt "$HOME/.local/bin/clipromt"
export PATH="$HOME/.local/bin:$PATH"
```

The installed binary needs no Go or Python at runtime. Bash 4.4+ with Readline, Emacs editing mode, and Linux command-line utilities are required. GNU find, grep and coreutils are the target; process lookup uses procps `pgrep`.

## Test it locally in one minute

The starter catalog is already built into the binary. These commands only print suggestions:

```bash
~/.local/bin/clipromt catalog
~/.local/bin/clipromt suggest 'who am i'
~/.local/bin/clipromt suggest 'find files named "*.txt" in .'
~/.local/bin/clipromt suggest 'who am i and then show current directory'
```

Expected results are `whoami`, `find '.' -type f -name '*.txt'`, and `whoami && pwd`. Nothing is executed by `suggest`. A short [starter catalog](examples/starter.txt) is also included in the repo; it is plain English reference data, not a shell script.

After Bash initialization, type `/c who am i` and press Enter once. The line should change to `whoami` without printing your username. Edit it if needed. Press Enter again to run it, or Ctrl-C to cancel. Then try `/c show current directory` or `/c show processes`.

To test find and grep on disposable sample data:

```bash
mkdir -p /tmp/clipromt-demo
printf 'hello\nerror: example\n' > /tmp/clipromt-demo/example.txt
~/.local/bin/clipromt suggest 'find files named "*.txt" in /tmp/clipromt-demo'
~/.local/bin/clipromt suggest 'grep "error" in /tmp/clipromt-demo/example.txt'
```

Use the same requests after `/c` to review and manually run the commands in Bash. These are examples of the supported grammar; arbitrary wording may need rephrasing.

## Set up Bash

Run in an interactive Bash terminal:

```bash
set -o emacs
eval "$(clipromt init bash)"
```

`init` prints the bundled shell integration and the absolute executable path. The `eval` above loads that fixed integration. Do not use `eval` on `suggest` output.

Try:

```text
/c find files named "*.txt" in .
```

Press Enter. The editable line becomes:

```bash
find '.' -type f -name '*.txt'
```

It has not run. A separate Enter runs it. Ctrl-X, Ctrl-G also generates from the current input, with or without `/c`. Errors leave the original line intact; Ctrl-C clears it.

To keep setup for new terminals, add these lines to `~/.bashrc` after other keybinding setup:

```bash
export PATH="$HOME/.local/bin:$PATH"
if [[ $- == *i* ]] && command -v clipromt >/dev/null && [[ ! -o vi ]]; then
    eval "$(clipromt init bash)"
fi
```

clipromt owns Enter, Ctrl-J and Ctrl-X sequences for its dispatch. Existing custom bindings on those keys are replaced. `clipromt_disable` restores ordinary Enter and removes its shortcuts; it does not restore previous custom bindings. Reinitialize after changing editing modes or keybindings. Multiline pasted commands and unusual terminal keymaps are outside the tested scope. Pasting multiple Enter characters can execute the replacement just as manually pressing Enter twice can.

## Requests

Run `clipromt catalog` for the complete grammar. These are supported examples:

| Request | Generated behavior |
| --- | --- |
| `who am i` | Show the current user |
| `show current directory` | Print the current directory |
| `find files named "*.txt" in .` | Find matching files |
| `find directories named cache in .` | Find matching directories |
| `find files named "*.log" in . containing "error"` | Find plus literal grep, with safe filename passing |
| `grep "hello world" in "my file.txt"` | Literal search with line numbers |
| `search "error" recursively in .` | Recursive literal search, without following directory symlinks |
| `copy "a.txt" to "b.txt"` | Copy with overwrite confirmation |
| `copy directory "folder" to "backup"` | Recursive copy with overwrite confirmation |
| `move "a.txt" to "b.txt"` | Move with overwrite confirmation |
| `chmod 755 "script.sh"` | Change permissions |
| `make "script.sh" executable` | Add owner execute permission |
| `show processes` | List processes |
| `find process "python"` | Match process command lines using a pgrep regular expression |
| `stop process 1234` | Suggest SIGTERM to one numeric PID greater than 1 |
| `follow log "/var/log/syslog"` | Show the last 50 lines and follow file replacement |
| `show last 100 lines of "app.log"` | Show a bounded number of recent lines |
| `who am i and then show current directory` | Compose steps with `&&` |

`cp`, `mv`, `pwd`, and `whoami` aliases also work. Up to four steps are allowed. Every step must parse before any output is returned. A later step runs only if the previous one succeeds when you manually execute the command. Grep and pgrep return failure when they find no matches.

Command words ignore case; arguments keep their case. Quote paths or search terms containing spaces inside the request. Use the other quote type for an apostrophe or double quote; backslash escape syntax is not supported. Empty arguments, invalid UTF-8, control characters, and requests above 8192 bytes are rejected.

Arguments are single-quoted with embedded apostrophes escaped. Shell substitutions, semicolons and backticks in arguments remain literal. Leading dash paths and find expression-like roots are made explicit relative paths. `~`, environment variables, and wildcards in copy/move paths remain literal. Find name patterns are globs; process patterns are regular expressions. Paths are not checked for existence or ownership. Copies and moves use the platform's normal destination-directory behavior. This generator does not sandbox the command you choose to execute.

Missing arguments, unsupported clauses, and ambiguous requests return a clarification error without a partial command. Rephrase using a catalog example. There is no conversational clarification state or inference about what a vague request meant.

## Controls and common problems

```bash
clipromt help
clipromt catalog
clipromt config
clipromt doctor
clipromt_disable                 # End shell integration in this terminal.
eval "$(clipromt init bash)"     # Restart it.
```

There is no job or daemon to stop. Ctrl-C cancels a CLI invocation or a log command you chose to run. Remove the `.bashrc` setup lines to stop initialization in future terminals.

If `/c` says “No such file,” the integration is not active: use interactive Bash, run `set -o emacs`, and initialize again. If `clipromt` is missing, add its install directory to PATH. If the binary was moved, run `init` again so the binding picks up its new path. `doctor` checks command availability; it cannot inspect the parent shell's live bindings or verify every GNU utility version.

If a shortcut is intercepted or your shell is unsupported, use the print-only alternative:

```bash
clipromt suggest 'copy "a.txt" to "backup.txt"'
```

Copy the result into your prompt and review it. Never pipe it into a shell. `config` reports fixed runtime settings as JSON; v0.1 has no mutable or persisted configuration and never reads a config file.

## Verification and releases

Development checks need Go and Python 3; Python is only for the terminal test:

```bash
go vet ./...
go test -race ./...
go build -buildvcs=false -o bin/clipromt ./cmd/clipromt
bash tests/bash-integration.sh
go test ./internal/composer -run '^$' -fuzz FuzzSuggest -fuzztime 10s
bash scripts/release.sh
```

CI runs these checks. The release workflow builds Linux amd64/arm64 archives and checksums when tag `v0.1.0` is pushed. See [release notes](docs/RELEASE-v0.1.0.md) and [architecture decisions](docs/ARCHITECTURE.md).

v0.1 is a bounded deterministic English grammar, not a general natural-language interpreter. It does not support arbitrary scripts, pipes supplied by the user, redirections, deletion, privilege escalation, AI backends, Zsh, Bash Vi mode, Windows or macOS. The real terminal integration was verified on Linux with Bash Emacs Readline. ARM artifacts are cross-compiled and need a runtime check on an ARM machine.
