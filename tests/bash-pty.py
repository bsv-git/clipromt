#!/usr/bin/env python3
"""Exercise actual Readline under a terminal, including two-Enter execution."""
import os, pty, select, time, tempfile
from pathlib import Path
root = Path(__file__).resolve().parents[1]

def run(mode, trigger):
    with tempfile.TemporaryDirectory() as d:
        marker = Path(d) / 'destination'
        src = Path(d) / 'source'
        src.write_text('safe')
        pid, fd = pty.fork()
        if pid == 0:
            os.environ['PS1'] = 'TEST> '
            os.environ['TERM'] = 'xterm'
            os.execvp('bash', ['bash', '--noprofile', '--norc', '-i'])
        def read_until(needle, timeout=5):
            data = b''
            until = time.monotonic() + timeout
            while time.monotonic() < until:
                if select.select([fd], [], [], .05)[0]:
                    data += os.read(fd, 65536)
                    if needle in data:
                        return data
            raise AssertionError((needle, data))
        def send(s): os.write(fd, s.encode())
        try:
            read_until(b'TEST> ')
            send(f'export CLIPROMT_BIN="{root}/bin/clipromt"; set -o {mode}; eval "$("{root}/bin/clipromt" init bash)"\n')
            read_until(b'TEST> ')
            prefix = '' if trigger == '\x18\x07' else '/c '
            send(f'{prefix}copy "{src}" to "{marker}"' + trigger)
            read_until(b'cp -i --')
            time.sleep(.15)
            assert not marker.exists(), 'trigger executed generated command'
            # Edit the generated command by appending a harmless shell comment.
            send(' # reviewed\r')
            read_until(b'TEST> ')
            time.sleep(.2)
            assert marker.read_text() == 'safe', 'second Enter failed'
            send('/c unsupported request\r')
            read_until(b'unsupported or ambiguous')
            # Cancel preserved input; ordinary command must still execute.
            send('\x03')
            read_until(b'TEST> ')
            send('printf "NORMAL_OK\\n"\r')
            read_until(b'NORMAL_OK\r\n')
            send('clipromt_disable\r')
            read_until(b'TEST> ')
            send('echo DISABLED_OK\r')
            read_until(b'DISABLED_OK\r\n')
            send('exit\r')
        finally:
            os.close(fd)
            try: os.kill(pid, 15)
            except ProcessLookupError: pass
            os.waitpid(pid, 0)
for trigger in ('\r', '\n', '\x18\x07'):
    run('emacs', trigger)
    print(f'PASS: real Bash, trigger {trigger!r}, editable line, no auto-execution, error, normal Enter, disable')
