package composer

import (
	"strings"
	"testing"
)

func TestSuggestions(t *testing.T) {
	cases := map[string]string{
		`find files named x in "!"`: `find './!' -type f -name 'x'`,
		"who am i":                  "whoami", "WHERE AM I": "pwd",
		`find files named "*.txt" in .`:                       `find '.' -type f -name '*.txt'`,
		`find directories named cache in -root`:               `find './-root' -type d -name 'cache'`,
		`find files named "*.log" in /tmp containing "error"`: `find '/tmp' -type f -name '*.log' -exec grep -nH -F -- 'error' {} +`,
		`grep "hello world" in "my file.txt"`:                 `grep -n -F -- 'hello world' 'my file.txt'`,
		`search -n recursively in .`:                          `grep -rn -F -- '-n' '.'`,
		`copy a to b`:                                         `cp -i -- 'a' 'b'`,
		`copy directory a to b`:                               `cp -R -i -- 'a' 'b'`,
		`chmod 755 script.sh`:                                 `chmod 755 -- 'script.sh'`,
		`chmod g-w file`:                                      `chmod g-w -- 'file'`,
		`make script.sh executable`:                           `chmod u+x -- 'script.sh'`,
		`who am i and then pwd`:                               `whoami && pwd`,
		`search "a and then b" in x`:                          `grep -n -F -- 'a and then b' 'x'`,
		`copy "$(touch /tmp/clipromt-marker);x" to "b'c"`:     `cp -i -- '$(touch /tmp/clipromt-marker);x' 'b'"'"'c'`,
	}
	for prompt, want := range cases {
		t.Run(prompt, func(t *testing.T) {
			got, err := Suggest(prompt)
			if err != nil || got != want {
				t.Fatalf("got %q, %v; want %q", got, err, want)
			}
		})
	}
}
func TestReject(t *testing.T) {
	for _, p := range []string{"", "delete everything", "copy a", `grep "oops in x`, "chmod 999 file", "chmod --reference=x file", "pwd; whoami", "pwd\nwhoami", "pwd and then delete x", strings.Repeat("pwd and then ", 4) + "pwd"} {
		if got, err := Suggest(p); err == nil || got != "" {
			t.Errorf("accepted %q: %q, %v", p, got, err)
		}
	}
}
func FuzzSuggest(f *testing.F) {
	for _, s := range []string{"pwd", `copy "a'b" to "$(whoami)"`, "pwd and then whoami"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, p string) {
		got, err := Suggest(p)
		if err != nil && got != "" {
			t.Fatal("partial suggestion")
		}
		if strings.ContainsAny(got, "\n\r\x00") {
			t.Fatal("control character output")
		}
	})
}

func TestReleaseIntents(t *testing.T) {
	for p, w := range map[string]string{
		`move "a b" to -dest`:           `mv -i -- 'a b' './-dest'`,
		`show processes`:                `ps -ef`,
		`find process "python.*"`:       `pgrep -af -- 'python.*'`,
		`stop process 234`:              `kill -TERM -- 234`,
		`follow log -log`:               `tail -n 50 -F -- './-log'`,
		`show last 30 lines of app.log`: `tail -n 30 -- 'app.log'`,
	} {
		got, e := Suggest(p)
		if e != nil || got != w {
			t.Errorf("%s: %q %v", p, got, e)
		}
	}
	for _, p := range []string{`copy "" to a`, `stop process 0`, `stop process -1`, `stop process 1`, `stop process python`, `show last 0 lines of log`, string([]byte{255}), strings.Repeat("x", 8193)} {
		if c, e := Suggest(p); e == nil || c != "" {
			t.Errorf("accepted %q", p)
		}
	}
}
