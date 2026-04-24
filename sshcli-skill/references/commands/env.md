# env — Variables de entorno

## Descripción
Muestra las variables de entorno del sistema.

## Sintaxis
```bash
sshcli env [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
sshcli env
sshcli env | grep PATH
```

## Output Example
```
HOME=/root
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin
USER=root
...
```

## Notas
- Equivalente a `env` en Linux
- Útil para verificar configuraciones
