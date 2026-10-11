package main

import (
	"clipromt/internal/composer"
	"clipromt/internal/shell"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

var version = "0.1.0"

const catalog = `who am i
show current directory
find files named "*.txt" in .
find directories named cache in .
find files named "*.log" in . containing "error"
grep "hello world" in "my file.txt"
search "error" recursively in .
copy "a.txt" to "b.txt"
copy directory "folder" to "backup"
move "a.txt" to "b.txt"
chmod 755 "script.sh"
make "script.sh" executable
show processes
find process "python"
stop process 1234
follow log "/var/log/syslog"
show last 100 lines of "app.log"
Join up to four steps with: and then`
const help = `clipromt — offline English to editable Bash commands
Usage: clipromt suggest "request" | catalog | init bash | doctor | config | version | help
Setup: eval "$(clipromt init bash)"
Type /c followed by a supported request, then Enter. Review the replacement.
A second Enter runs it. Ctrl-X Ctrl-G also generates from the current line.
Generation never executes the suggested command. Unsupported requests return errors.
Use catalog for examples. No AI, telemetry, daemon, or network access.`

func run(args []string, out, errOut io.Writer) int {
	fail := func(e interface{}) int { fmt.Fprintln(errOut, "clipromt:", e); return 1 }
	if len(args) == 0 {
		fmt.Fprintln(errOut, help)
		return 2
	}
	switch args[0] {
	case "help", "--help", "-h":
		if len(args) != 1 {
			return fail("help takes no arguments")
		}
		fmt.Fprintln(out, help)
		return 0
	case "version", "--version":
		if len(args) != 1 {
			return fail("version takes no arguments")
		}
		fmt.Fprintln(out, "clipromt", version)
		return 0
	case "catalog":
		if len(args) != 1 {
			return fail("catalog takes no arguments")
		}
		fmt.Fprintln(out, catalog)
		return 0
	case "suggest":
		if len(args) != 2 {
			fmt.Fprintln(errOut, `Usage: clipromt suggest "request"`)
			return 2
		}
		command, e := composer.Suggest(args[1])
		if e != nil {
			return fail(e)
		}
		fmt.Fprintln(out, command)
		return 0
	case "init":
		if len(args) != 2 || args[1] != "bash" {
			return fail("v0.1 supports init bash only")
		}
		exe, e := os.Executable()
		if e != nil {
			return fail(e)
		}
		fmt.Fprintf(out, "export CLIPROMT_BIN=%s\n", composer.Quote(exe))
		fmt.Fprint(out, shell.Bash)
		return 0
	case "config":
		if len(args) != 1 {
			return fail("v0.1 config is read-only; settings are fixed for offline safety")
		}
		root, e := os.UserConfigDir()
		if e != nil {
			return fail(e)
		}
		data := map[string]interface{}{"version": version, "backend": "deterministic", "shell": "bash", "auto_execute": false, "network": false, "max_steps": 4, "max_prompt_bytes": 8192, "config_path": filepath.Join(root, "clipromt", "config.json"), "persistent_config": false}
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		if e := enc.Encode(data); e != nil {
			return fail(e)
		}
		return 0
	case "doctor":
		if len(args) != 1 {
			return fail("doctor takes no arguments")
		}
		missing := false
		for _, name := range []string{"bash", "find", "grep", "cp", "mv", "chmod", "ps", "pgrep", "kill", "tail"} {
			p, e := exec.LookPath(name)
			if e != nil {
				fmt.Fprintln(out, "MISSING", name)
				missing = true
			} else {
				fmt.Fprintln(out, "OK", name, p)
			}
		}
		fmt.Fprintln(out, "Shell binding: verify in interactive Bash with /c who am i. Offline backend; no service required.")
		if missing {
			return 1
		}
		return 0
	}
	fmt.Fprintln(errOut, help)
	return 2
}
func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }
