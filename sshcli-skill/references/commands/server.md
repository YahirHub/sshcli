# server — Gestión de servidores

## Descripción
Gestiona los servidores SSH configurados.

## Sintaxis
```bash
sshcli server list         # Listar servidores
sshcli server add <nombre> --host <ip> --user <user> --pass <pass>
sshcli server use <nombre>  # Activar servidor
sshcli server info         # Info del servidor activo
sshcli server remove <nombre>  # Eliminar servidor
```

## Flags
| Flag | Descripción |
|------|-------------|
| `--host` | Dirección IP o hostname |
| `--user` | Usuario SSH |
| `--pass` | Contraseña |

## Ejemplos
```bash
# Ver servidores
sshcli server list

# Agregar servidor
sshcli server add "prod" --host "192.168.1.10" --user "root" --pass "secret"

# Agregar servidor con contraseña con caracteres especiales (importante: usar comillas)
sshcli server add "pangolin" --host "pangolin.yahirex.us.kg" --user "root" --pass "Linux_0145!"

# Usar servidor
sshcli server use "prod"

# Info del servidor activo
sshcli server info
```

## Notas
- `*` = servidor activo
- Requiere conexión SSH válida
- **Caracteres especiales en contraseñas:** Si la contraseña tiene `!`, `$`, `` ` ``, `"`, `\` u otros caracteres especiales, DEBE ir entre comillas dobles. De lo contrario el shell los interpretará y la conexión fallará.
