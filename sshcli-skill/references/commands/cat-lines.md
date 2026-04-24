# cat-lines — Leer rango de líneas

## Descripción
Lee un rango específico de líneas de un archivo.

## Sintaxis
```bash
sshcli cat-lines [archivo] [inicio] [fin] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-n, --numbers` | Mostrar números de línea |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Leer líneas 1-10
sshcli cat-lines "/var/log/syslog" 1 10

# Con números de línea
sshcli cat-lines "/app/main.py" 50 100 -n
```

## Notas
- Es útil para archivos grandes
- Evita cargar el archivo completo
