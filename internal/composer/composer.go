// Package composer translates a bounded grammar into shell text. It never runs it.
package composer

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func quote(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'" }
func path(s string) string {
	if strings.HasPrefix(s, "-") || s == "!" || s == "(" || s == ")" {
		s = "./" + s
	}
	return quote(s)
}

// tokens supports quoted arguments, but deliberately does not interpret shell syntax.
func tokens(s string) ([]string, error) {
	var out []string
	var b strings.Builder
	var q rune
	active := false
	for _, r := range s {
		if unicode.IsControl(r) {
			return nil, fmt.Errorf("control characters are unsupported")
		}
		if q != 0 {
			if r == q {
				q = 0
			} else {
				b.WriteRune(r)
			}
			continue
		}
		if r == '\'' || r == '"' {
			q = r
			active = true
			continue
		}
		if unicode.IsSpace(r) {
			if active {
				out = append(out, b.String())
				b.Reset()
				active = false
			}
		} else {
			b.WriteRune(r)
			active = true
		}
	}
	if q != 0 {
		return nil, fmt.Errorf("unclosed quote")
	}
	if active {
		out = append(out, b.String())
	}
	return out, nil
}

// Suggest accepts exact supported phrases. Literal data is always shell quoted.
func Suggest(prompt string) (string, error) {
	if len(prompt) > 8192 || !utf8.ValidString(prompt) {
		return "", fmt.Errorf("request must be valid UTF-8 and at most 8192 bytes")
	}
	ts, err := tokens(prompt)
	if err != nil {
		return "", err
	}
	if len(ts) == 0 {
		return "", fmt.Errorf("empty prompt")
	}
	// Split only unquoted delimiters by scanning the original text.
	var parts []string
	var q rune
	start := 0
	for i := 0; i < len(prompt); i++ {
		c := rune(prompt[i])
		if q != 0 {
			if c == q {
				q = 0
			}
			continue
		}
		if c == '\'' || c == '"' {
			q = c
			continue
		}
		if strings.HasPrefix(strings.ToLower(prompt[i:]), " and then ") {
			parts = append(parts, prompt[start:i])
			i += len(" and then ") - 1
			start = i + 1
		}
	}
	parts = append(parts, prompt[start:])
	if len(parts) > 4 {
		return "", fmt.Errorf("at most four steps are supported")
	}
	var commands []string
	for _, p := range parts {
		t, e := tokens(p)
		if e != nil {
			return "", e
		}
		c, e := single(t)
		if e != nil {
			return "", e
		}
		commands = append(commands, c)
	}
	return strings.Join(commands, " && "), nil
}

var modeRE = regexp.MustCompile(`^([0-7]{3,4}|[ugoa]*[+=-][rwx]+)$`)

func single(t []string) (string, error) {
	l := make([]string, len(t))
	for i, v := range t {
		l[i] = strings.ToLower(v)
	}
	for _, v := range t {
		if v == "" {
			return "", fmt.Errorf("empty arguments are ambiguous; provide an explicit value")
		}
	}
	phrase := strings.Join(l, " ")
	switch phrase {
	case "whoami", "who am i", "which user", "show current user":
		return "whoami", nil
	case "show processes", "list processes":
		return "ps -ef", nil
	case "pwd", "where am i", "show current directory", "current directory":
		return "pwd", nil
	}
	if len(t) == 6 && l[0] == "find" && (l[1] == "files" || l[1] == "directories") && l[2] == "named" && l[4] == "in" {
		kind := "f"
		if l[1] == "directories" {
			kind = "d"
		}
		return "find " + path(t[5]) + " -type " + kind + " -name " + quote(t[3]), nil
	}
	if len(t) == 8 && l[0] == "find" && l[1] == "files" && l[2] == "named" && l[4] == "in" && l[6] == "containing" {
		return "find " + path(t[5]) + " -type f -name " + quote(t[3]) + " -exec grep -nH -F -- " + quote(t[7]) + " {} +", nil
	}
	if len(t) == 4 && (l[0] == "grep" || l[0] == "search") && l[2] == "in" {
		return "grep -n -F -- " + quote(t[1]) + " " + path(t[3]), nil
	}
	if len(t) == 5 && l[0] == "search" && l[2] == "recursively" && l[3] == "in" {
		return "grep -rn -F -- " + quote(t[1]) + " " + path(t[4]), nil
	}
	if len(t) == 4 && (l[0] == "move" || l[0] == "mv") && l[2] == "to" {
		return "mv -i -- " + path(t[1]) + " " + path(t[3]), nil
	}
	if len(t) == 3 && l[0] == "find" && l[1] == "process" {
		return "pgrep -af -- " + quote(t[2]), nil
	}
	if len(t) == 3 && l[0] == "stop" && l[1] == "process" {
		n, e := strconv.Atoi(t[2])
		if e != nil || n <= 1 || n > 2147483647 {
			return "", fmt.Errorf("provide one numeric process ID greater than 1")
		}
		return "kill -TERM -- " + strconv.Itoa(n), nil
	}
	if len(t) == 3 && l[0] == "follow" && l[1] == "log" {
		return "tail -n 50 -F -- " + path(t[2]), nil
	}
	if len(t) == 6 && l[0] == "show" && l[1] == "last" && l[3] == "lines" && l[4] == "of" {
		n, e := strconv.Atoi(t[2])
		if e != nil || n < 1 || n > 100000 {
			return "", fmt.Errorf("line count must be 1 through 100000")
		}
		return "tail -n " + strconv.Itoa(n) + " -- " + path(t[5]), nil
	}
	if len(t) == 4 && (l[0] == "copy" || l[0] == "cp") && l[2] == "to" {
		return "cp -i -- " + path(t[1]) + " " + path(t[3]), nil
	}
	if len(t) == 5 && l[0] == "copy" && l[1] == "directory" && l[3] == "to" {
		return "cp -R -i -- " + path(t[2]) + " " + path(t[4]), nil
	}
	if len(t) == 3 && l[0] == "chmod" && modeRE.MatchString(t[1]) && !strings.HasPrefix(t[1], "-") {
		return "chmod " + t[1] + " -- " + path(t[2]), nil
	}
	if len(t) == 3 && l[0] == "make" && l[2] == "executable" {
		return "chmod u+x -- " + path(t[1]), nil
	}
	return "", fmt.Errorf("unsupported or ambiguous request; run clipromt catalog for supported phrases; quote paths containing spaces")
}

// Quote returns a literal shell argument, never executable syntax.
func Quote(s string) string { return quote(s) }
