# service — Gestionar servicios

## Descripción
Gestiona servicios del sistema usando systemctl.

## Sintaxis
```bash
sshcli service [nombre] [accion] [flags]
```

## Acciones
| Acción | Descripción |
|--------|-------------|
| `start` | Iniciar servicio |
| `stop` | Detener servicio |
| `restart` | Reiniciar servicio |
| `status` | Ver estado |
| `enable` | Habilitar al inicio |
| `disable` | Deshabilitar al inicio |

## Flags
| Flag | Descripción |
|------|-------------|
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Ver estado
sshcli service nginx status

# Reiniciar
sshcli service nginx restart

# Iniciar
sshcli service postgresql start

# Habilitar al inicio
sshcli service nginx enable
```

## Notas
- Usa `systemctl` internamente
- Requiere permisos de root
