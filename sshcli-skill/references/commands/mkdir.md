# mkdir — Crear directorios

## Descripción
Crea un directorio en el servidor remoto.

## Sintaxis
```bash
sshcli mkdir [ruta_remota] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-p, --parents` | Crear directorios padres si no existen |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Crear directorio
sshcli mkdir "/home/user/newdir"

# Crear ruta anidada
sshcli mkdir "/home/user/a/b/c" -p
```

## Notas
- La flag `-p` crea directorios padre automáticamente
- No falla si el directorio ya existe (con `-p`)
