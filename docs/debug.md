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
