# server — Gestión de servidores

## Descripción
Gestiona los servidores SSH configurados.

## Sintaxis
```bash
sshcli server list [--format text|json] [--tag prod]
sshcli server search <texto>
sshcli server add <nombre> --host <ip> --user <user> --pass <pass>
sshcli server add <nombre> --host <ip> --user <user> --pass <pass> --force
sshcli server use <nombre>
sshcli server rename <actual> <nuevo>
sshcli server edit <nombre> [--host ...] [--port ...] [--user ...] [--pass ...]
sshcli server copy <origen> <nuevo>
sshcli server sync-tags <origen> <destino>
sshcli server test [nombre]
sshcli server ping [nombre]
sshcli server doctor [nombre]
sshcli server export <archivo.json>
sshcli server import <archivo.json>
sshcli server set-pass <nombre> --pass <nueva_pass>
sshcli server info [nombre] [--show-pass]
sshcli server remove <nombre> --yes
```

## Ejemplos
```bash
# Ver servidores en texto o JSON
sshcli server list
sshcli server list --format json

# Filtrar por tags
sshcli server list --tag prod --tag eu

# Buscar por nombre, host, usuario o tag
sshcli server search "staging"
sshcli server search "10.0.0.20"

# Agregar servidor
sshcli server add "prod" --host "192.168.1.10" --user "root" --pass "secret" --tag prod --tag eu

# Sobrescribir servidor existente
sshcli server add "prod" --host "192.168.1.11" --user "root" --pass "secret" --force

# Renombrar servidor
sshcli server rename "prod" "prod-eu"

# Editar host y tags
sshcli server edit "prod-eu" --host "10.0.0.20" --add-tag api --remove-tag eu

# Copiar servidor
sshcli server copy "prod-eu" "staging" --host "10.0.0.30" --tag staging

# Sincronizar tags entre dos servidores
sshcli server sync-tags "prod-eu" "staging"

# Probar conexión y medir latencia
sshcli server test "staging"
sshcli server ping "staging"

# Diagnóstico remoto
sshcli server doctor "staging"

# Exportar e importar
sshcli server export "servers.json"
sshcli server import "servers.json"
sshcli server import "servers.json" --force
sshcli server import "servers.json" --replace

# Actualizar contraseña guardada
sshcli server set-pass "staging" --pass "Nueva_Clave123!"

# Ver información, incluso contraseña guardada
sshcli server info "staging" --show-pass

# Eliminar sin prompt
sshcli server remove "staging" --yes
```

## Notas
- `*` = servidor activo.
- `server add` falla si el nombre ya existe, salvo que uses `--force`.
- `server rename` valida que el nombre actual exista y que el nuevo nombre no esté en uso.
- `server edit` valida la conexión si cambias host, puerto, usuario o contraseña.
- `server test` usa la configuración guardada y verifica autenticación SSH real.
- `server ping` mide el tiempo de establecimiento de conexión SSH.
- `server doctor` revisa SSH y utilidades comunes como `sudo`, `git`, `docker`, `php`, `composer` y `caddy`.
- `server export` / `server import` usan JSON compatible con la configuración de `sshcli`.
- `server set-pass` actualiza **solo** la contraseña guardada localmente.
- `server sync-tags` copia exactamente los tags del origen al destino.
- `server list --format json` devuelve una lista lista para automatización.
- `server remove` pide confirmación; usa `--yes` para modo no interactivo.
- Los tags permiten clasificar servidores por entorno, región o rol.
- **Caracteres especiales en contraseñas:** si la contraseña tiene `!`, `$`, `` ` ``, `"`, `\` u otros caracteres especiales, usa comillas dobles.
