# Debugging using Remote Debug

- Run application using debug mode `bima run debug`

- Copy PID

- Run debugger `bima debug <pid>`

- Bima Framework reserved `16517` port number as debug port, you can use it for any IDE that support remote debug

- For VS Code user, copy `launch.json` below

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug",
            "type": "go",
            "request": "attach",
            "mode": "remote",
            "port": 16517,
            "remotePath": "${workspaceFolder}",
            "host": "127.0.0.1",
            "cwd": "${workspaceFolder}"
        }
    ]
}

```

- Just use Debug Panel as usual

##  Linux trace error

`.... caused by a kernel security setting, try to writing "0" to /proc/sys/kernel/yama/ptrace_scope`

- Edit `/etc/sysctl.d/10-ptrace.conf` file and set `kernel.yama.ptrace_scope = 0`

- Restart
