# connect — Terminal interactiva

## Descripción
Abre una terminal SSH interactiva con el servidor activo.

## Sintaxis
```bash
sshcli connect
sshcli connect --server <nombre>
```

## Ejemplos
```bash
sshcli connect
sshcli connect --server prod
```

## Notas
- Requiere servidor activo o `--server`.
- Solicita PTY remoto.
- Requiere una terminal local real para uso cómodo.
- Ctrl+D o `exit` para salir.
