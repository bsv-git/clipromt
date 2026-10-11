# Architecture and scope, v0.1.0

The referenced `/mnt/data/clipromt_architecture_product_spec_v0_1.md` and `.docx` were unavailable in this workspace. The cached source conversation did not expose their content. This implementation follows the current user's explicit requirements and extends the existing Go grammar prototype. Exact parity with that missing document cannot be claimed.

The CLI dispatches help, catalog, suggest, init, doctor, config, and version. The composer validates bounded UTF-8 input, tokenizes literal quoted arguments, splits unquoted `and then`, matches allowlisted intent shapes, and renders shell text. It rejects the entire request if any step fails. It uses no shell evaluator or subprocess to generate a command. Doctor checks executable paths without launching utilities.

Bash integration is embedded in the binary. Init emits that integration with the absolute binary path. Readline Enter is a macro that first calls a shell function and then a dispatch key. For `/c` the callback replaces the input buffer and removes the dispatch action; for other input it binds the dispatch to accept-line. This avoids accepting generated text on the trigger keypress. The shortcut writes READLINE_LINE and its byte cursor position. No DEBUG trap, command-not-found interception, background daemon, or history replay is used.

Quoted literal data remains separate from fixed command syntax. GNU find uses -exec argument passing instead of an unsafe filename pipeline. Copy/move include interactive overwrite flags. Process termination requires one PID greater than 1 and proposes TERM. These choices reduce mistakes but cannot guarantee the safety of an action the user manually executes.

The default Linux Bash/Emacs setup is the release target. Vi mode is rejected during initialization because the Readline macro did not work reliably in the terminal test. Zsh needs a separate ZLE widget and terminal validation; it is deferred. Configuration is a read-only report of fixed behavior, avoiding a mutable execute or network setting.

Future versions can add more deterministic intents, structured clarification, a customizable trigger and binding restoration, and separately tested Zsh integration. Optional model providers remain future work and must not become a runtime requirement for offline tasks. Any future backend must keep the same editable-buffer and separate-execution contract.
