# memory — Uso de memoria

## Descripción
Muestra el uso de memoria RAM del sistema.

## Sintaxis
```bash
sshcli memory [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
sshcli memory
```

## Output Example
```
               total        used        free      shared  buff/cache   available
Mem:            23Gi       4.0Gi       5.5Gi       5.0Mi        14Gi        19Gi
Swap:             0B          0B          0B          0B
```

## Notas
- Muestra: total, usado, libre, shared, buff/cache, available
- Available = memoria real disponible para aplicaciones
