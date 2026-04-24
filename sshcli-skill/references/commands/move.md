# move — Mover o renombrar

## Descripción
Mueve o renombra archivos y directorios en el servidor remoto.

## Sintaxis
```bash
sshcli move [origen_remoto] [destino_remoto] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Renombrar archivo
sshcli move "/home/user/old.txt" "/home/user/new.txt"

# Mover archivo
sshcli move "/tmp/file.txt" "/home/user/file.txt"
```

## Notas
- Equivalente a `mv` en Linux
- Normaliza rutas automáticamente
