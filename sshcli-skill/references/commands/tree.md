# tree — Árbol de directorios

## Descripción
Muestra la estructura de directorios en forma de árbol.

## Sintaxis
```bash
sshcli tree [ruta] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-d, --depth` | Profundidad máxima |
| `--dirs` | Solo directorios |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
sshcli tree "/var/www"
sshcli tree "/home" -d 2
sshcli tree "/var" -d 1 --dirs
```

## Notas
- Si el binario `tree` no existe en el servidor, `sshcli` genera un árbol ASCII de fallback.
- La profundidad `-d` limita la recursión.
