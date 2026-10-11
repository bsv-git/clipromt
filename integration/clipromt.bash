# Bash 4.4+ with Readline. Source only in an interactive shell.
[[ $- == *i* ]] || return 0
if [[ -o vi ]]; then
    printf "%s\n" "clipromt: switch to set -o emacs, then initialize again; v0.1 requires Emacs Readline mode" >&2
    return 1
fi
_clipromt_edit_line() {
    local request=$READLINE_LINE suggestion
    if [[ $request == /c || $request == '/c '* ]]; then
        request=${request:2}
    fi
    suggestion=$("${CLIPROMT_BIN:-clipromt}" suggest "$request") || return
    READLINE_LINE=$suggestion
    local LC_ALL=C
    READLINE_POINT=${#READLINE_LINE}
}
# The macro routes Enter through a callback and an otherwise unused dispatch key.
# Trigger requests clear that key's macro; ordinary input restores accept-line.
_clipromt_enter() {
    if [[ $READLINE_LINE == /c || $READLINE_LINE == '/c '* ]]; then
        bind '"\C-x\C-]":""'
        _clipromt_edit_line
    else
        bind '"\C-x\C-]":accept-line'
    fi
}
clipromt_disable() {
    local map
    for map in emacs-standard; do
        bind -m "$map" '"\C-m":accept-line'
        bind -m "$map" '"\C-j":accept-line'
        bind -m "$map" -r '\C-x\C-g'
        bind -m "$map" -r '\C-x\C-]'
        bind -m "$map" -r '\C-x\C-^'
    done
    unset -f _clipromt_enter _clipromt_edit_line
}
for _clipromt_map in emacs-standard; do
    bind -m "$_clipromt_map" -x '"\C-x\C-^":_clipromt_enter'
    bind -m "$_clipromt_map" -x '"\C-x\C-g":_clipromt_edit_line'
    bind -m "$_clipromt_map" '"\C-m":"\C-x\C-^\C-x\C-]"'
    bind -m "$_clipromt_map" '"\C-j":"\C-x\C-^\C-x\C-]"'
    bind -m "$_clipromt_map" '"\C-x\C-]":accept-line'
done
unset _clipromt_map
