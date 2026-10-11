# clipromt v0.1.0

Standalone offline Go CLI for common Linux shell tasks. `/c` in interactive Bash replaces the live input with an editable generated command. The trigger does not execute it.

Includes a deterministic catalog and composer for find, grep, copy, move, chmod, processes, logs, and up to four steps joined with `&&`. Arguments are quoted as literals and unsupported requests produce errors. Init emits embedded Bash setup; doctor checks utilities; config reports fixed settings.

Validation: Go unit tests and race checks, vet, fuzz testing, and an actual Bash terminal test for no automatic execution, editing, a second Enter, error preservation, normal Enter and disabling integration. Linux amd64 is runnable locally. Linux arm64 is cross-compiled, not runtime-verified.

Limitations: exact supported phrases rather than unrestricted English; GNU/Linux utilities; Bash Emacs Readline only; no Zsh, Vi mode, mutable config, model, daemon, or network. Custom bindings are replaced. Multiple pasted Enter characters may execute a replacement. Paths and utility versions are not fully validated. The referenced external architecture document was unavailable.

