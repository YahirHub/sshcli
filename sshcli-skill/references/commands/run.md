# run — Ejecutar scripts

## Descripción
Ejecuta un archivo de código con el intérprete apropiado.

## Sintaxis
```bash
sshcli run [archivo] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-a, --args` | Argumentos para el script |
| `-e, --env` | Variables de entorno |
| `-w, --workdir` | Directorio de trabajo |
| `-s, --server` | Servidor específico |

## Intérpretes Soportados
| Extensión | Intérprete |
|-----------|-------------|
| `.py` | python3 |
| `.js` | node |
| `.sh` | bash |
| `.go` | go run |
| `.rb` | ruby |
| `.php` | php |

## Ejemplos
```bash
# Ejecutar script Python
sshcli run "/app/main.py"

# Con argumentos
sshcli run "/app/script.py" -a "--port 3000 --debug"

# Con variables de entorno
sshcli run "/app/server.py" -e "DEBUG=1 API_KEY=xxx"
```

## Notas
- Auto-detecta el intérprete por la extensión
- El archivo debe tener permisos de ejecución
