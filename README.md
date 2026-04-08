# SSHCLI - Cliente SSH para Agentes de IA

<p align="center">
  <b>Herramienta CLI diseñada para que modelos de IA y agentes autónomos puedan conectarse y operar en servidores Linux remotos.</b>
</p>

---

## Descripción

**sshcli** es un cliente SSH de línea de comandos escrito en Go, diseñado específicamente para ser utilizado por **agentes de IA** y **modelos de lenguaje** que necesitan:

- Ejecutar comandos en servidores Linux remotos
- Crear, editar y gestionar archivos de código
- Administrar repositorios Git
- Gestionar procesos y servicios
- Realizar operaciones de sistema

A diferencia de clientes SSH tradicionales, sshcli está optimizado para uso programático con comandos atómicos y salidas predecibles.

## Características Principales

- **51 comandos** especializados para operaciones de desarrollo
- **Multi-servidor** con cambio rápido entre servidores
- **Modo interactivo** (-t) para programas como htop, vim, apt
- **Sin dependencias** - binario único multiplataforma
- **Configuración persistente** en `~/.sshcli.conf`

## Instalación

### Desde código fuente

```bash
git clone https://github.com/tu-usuario/sshcli.git
cd sshcli
go build -o sshcli .

# Linux/Mac
sudo cp sshcli /usr/local/bin/

# Windows
copy sshcli.exe C:\Windows\System32\
```

### Compilación cruzada

```bash
# Para Linux desde Windows
set GOOS=linux
set GOARCH=amd64
go build -o sshcli .

# Para Windows desde Linux
GOOS=windows GOARCH=amd64 go build -o sshcli.exe .
```

## Inicio Rápido

### 1. Agregar un servidor

```bash
sshcli server add mi-servidor --host 192.168.1.100 --user root --pass clave
```

### 2. Ejecutar comandos

```bash
sshcli exec "ls -la"
sshcli exec "cat /etc/hostname"
```

### 3. Modo interactivo

```bash
sshcli exec -t htop
sshcli exec -t "apt install nginx"
```

## Comandos por Categoría

### Gestión de Servidores
```bash
sshcli server add <nombre> --host X --user Y --pass Z
sshcli server list
sshcli server use <nombre>
sshcli server remove <nombre>
sshcli connect                    # Sesión SSH interactiva
```

### Ejecución de Comandos
```bash
sshcli exec "comando"             # Modo normal
sshcli exec -t "comando"          # Modo interactivo (TTY)
sshcli run /ruta/script.py        # Ejecutar script
sshcli syntax-check /ruta/file.py # Verificar sintaxis
```

### Edición de Archivos
```bash
sshcli write /ruta/archivo        # Crear desde stdin
sshcli read /ruta/archivo         # Leer archivo
sshcli append /ruta/archivo "texto"
sshcli insert-line /ruta/archivo 10 "código"
sshcli delete-line /ruta/archivo 10 15
sshcli replace-line /ruta/archivo 5 "nuevo contenido"
sshcli search-replace /ruta/archivo "buscar" "reemplazar"
```

### Exploración
```bash
sshcli list /ruta
sshcli tree /ruta -d 3
sshcli cat-lines /ruta/archivo 50 100
sshcli head -n 20 /ruta/archivo
sshcli tail -n 50 /ruta/archivo
sshcli find /ruta -name "*.py"
sshcli grep "patrón" /ruta/archivo
```

### Git
```bash
sshcli git-status /repo
sshcli git-diff /repo
sshcli git-add /repo .
sshcli git-commit /repo -m "mensaje"
sshcli git-pull /repo
sshcli git-push /repo
sshcli git-branch /repo
sshcli git-checkout /repo main
sshcli git-clone URL /destino
```

### Sistema
```bash
sshcli ps --filter python
sshcli kill <pid_o_nombre>
sshcli service nginx restart
sshcli ports
sshcli disk /ruta
sshcli memory
sshcli env /app/.env
sshcli env-set /app/.env "KEY=value"
```

## Configuración

### Habilitar TTY por defecto

```bash
sshcli config set tty true
```

Esto hace que todos los `exec` usen modo interactivo. Para desactivar temporalmente:

```bash
sshcli exec --no-tty "ls -la"
```

### Archivo de configuración

`~/.sshcli.conf`:
```json
{
  "active_server": "produccion",
  "default_tty": false,
  "servers": {
    "produccion": {
      "name": "produccion",
      "host": "prod.example.com",
      "port": 22,
      "user": "deploy",
      "password": "***"
    }
  }
}
```

## Uso con Agentes de IA

### Ejemplo: Flujo de desarrollo

```bash
# 1. Explorar proyecto
sshcli tree /app -d 2

# 2. Leer código existente
sshcli cat-lines -n /app/main.py 1 50

# 3. Modificar código
sshcli insert-line /app/main.py 1 "import logging"
sshcli search-replace /app/main.py "print(" "logging.info("

# 4. Verificar sintaxis
sshcli syntax-check /app/main.py

# 5. Ejecutar
sshcli run /app/main.py

# 6. Commit cambios
sshcli git-add /app .
sshcli git-commit /app -m "refactor: usar logging"
```

### Ejemplo: Deploy

```bash
sshcli git-pull --server prod /var/www/app
sshcli exec --server prod "pip install -r requirements.txt"
sshcli service --server prod gunicorn restart
sshcli tail --server prod /var/log/gunicorn/error.log
```

## Solución de Problemas

### Terminal bugueado después de -t

Si el terminal queda en modo raw, escribe:
```
reset
```
Y presiona Enter (aunque no veas lo que escribes).

### Error de conexión

Verifica credenciales:
```bash
sshcli server info mi-servidor
sshcli status
```

## Licencia

MIT

## Contribuir

Las contribuciones son bienvenidas. Por favor abre un issue primero para discutir cambios mayores.
