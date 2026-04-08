# Referencia de Comandos SSHCLI

## Índice

1. [Gestión de Servidores](#gestión-de-servidores)
2. [Configuración](#configuración)
3. [Conexión](#conexión)
4. [Ejecución](#ejecución)
5. [Edición de Archivos](#edición-de-archivos)
6. [Lectura de Archivos](#lectura-de-archivos)
7. [Gestión de Archivos](#gestión-de-archivos)
8. [Transferencia](#transferencia)
9. [Búsqueda](#búsqueda)
10. [Git](#git)
11. [Procesos y Sistema](#procesos-y-sistema)
12. [Información](#información)

---

## Gestión de Servidores

### server add
Agrega un nuevo servidor SSH.

```bash
sshcli server add <nombre> --host <ip> --user <usuario> --pass <contraseña> [--port <puerto>]
```

**Ejemplo:**
```bash
sshcli server add produccion --host 192.168.1.100 --user deploy --pass secreto
sshcli server add dev --host dev.local --port 2222 --user root --pass clave
```

### server list
Lista todos los servidores configurados.

```bash
sshcli server list
```

El servidor activo se marca con `*`.

### server use
Cambia el servidor activo.

```bash
sshcli server use <nombre>
```

### server remove
Elimina un servidor de la configuración.

```bash
sshcli server remove <nombre>
```

### server info
Muestra información detallada de un servidor.

```bash
sshcli server info [nombre]
```

---

## Configuración

### config show
Muestra la configuración actual.

```bash
sshcli config show
```

### config set
Establece un valor de configuración.

```bash
sshcli config set <opcion> <valor>
```

**Opciones:**
- `tty` - Habilitar/deshabilitar TTY por defecto (true/false)

**Ejemplo:**
```bash
sshcli config set tty true
```

---

## Conexión

### connect
Abre una sesión SSH interactiva.

```bash
sshcli connect [-s servidor]
```

Usa `exit` o `Ctrl+D` para salir.

---

## Ejecución

### exec
Ejecuta un comando en el servidor remoto.

```bash
sshcli exec [flags] "<comando>"
```

**Flags:**
- `-t, --tty` - Modo interactivo (para htop, vim, apt, etc.)
- `--no-tty` - Forzar modo no interactivo
- `-s, --server` - Servidor específico

**Ejemplos:**
```bash
sshcli exec "ls -la"
sshcli exec -t htop
sshcli exec -t "apt install nginx"
sshcli exec --server prod "systemctl status nginx"
```

### run
Ejecuta un script con su intérprete apropiado.

```bash
sshcli run <archivo> [flags]
```

**Flags:**
- `-a, --args` - Argumentos para el script
- `-e, --env` - Variables de entorno
- `-w, --workdir` - Directorio de trabajo

**Intérpretes soportados:**
- `.py` → python3
- `.js` → node
- `.ts` → npx tsx
- `.go` → go run
- `.sh` → bash
- `.rb` → ruby
- `.php` → php

**Ejemplo:**
```bash
sshcli run /app/main.py --args "--port 8000"
```

### syntax-check
Verifica la sintaxis de un archivo de código.

```bash
sshcli syntax-check <archivo>
```

**Lenguajes soportados:** Python, JavaScript, TypeScript, Go, Bash, Ruby, PHP, JSON, YAML

---

## Edición de Archivos

### write
Escribe contenido desde stdin a un archivo remoto.

```bash
echo "contenido" | sshcli write <ruta_remota>
cat archivo.txt | sshcli write <ruta_remota>
```

### append
Agrega contenido al final de un archivo.

```bash
sshcli append <ruta_remota> "<contenido>"
```

### insert-line
Inserta una línea en una posición específica.

```bash
sshcli insert-line <archivo> <numero_linea> "<contenido>"
```

**Ejemplo:**
```bash
sshcli insert-line /app/main.py 1 "import os"
```

### delete-line
Elimina líneas de un archivo.

```bash
sshcli delete-line <archivo> <linea_inicio> [linea_fin]
```

**Ejemplo:**
```bash
sshcli delete-line /app/main.py 10 15    # Elimina líneas 10-15
sshcli delete-line /app/main.py 5 5      # Elimina solo línea 5
```

### replace-line
Reemplaza una línea completa.

```bash
sshcli replace-line <archivo> <numero_linea> "<nuevo_contenido>"
```

### search-replace
Busca y reemplaza texto en un archivo.

```bash
sshcli search-replace <archivo> "<buscar>" "<reemplazar>" [flags]
```

**Flags:**
- `-a, --all` - Reemplazar todas las ocurrencias

---

## Lectura de Archivos

### read
Lee el contenido completo de un archivo.

```bash
sshcli read <ruta_remota>
```

### cat-lines
Lee un rango específico de líneas.

```bash
sshcli cat-lines <archivo> <linea_inicio> <linea_fin> [flags]
```

**Flags:**
- `-n, --numbers` - Mostrar números de línea

### head
Muestra las primeras líneas de un archivo.

```bash
sshcli head <archivo> [-n <lineas>]
```

### tail
Muestra las últimas líneas de un archivo.

```bash
sshcli tail <archivo> [-n <lineas>] [-f]
```

**Flags:**
- `-n, --lines` - Número de líneas (default: 20)
- `-f, --follow` - Seguir actualizaciones

---

## Gestión de Archivos

### list
Lista archivos y directorios.

```bash
sshcli list [ruta] [flags]
```

**Flags:**
- `-l, --long` - Formato largo
- `-a, --all` - Incluir ocultos

### tree
Muestra estructura de directorios en árbol.

```bash
sshcli tree [ruta] [flags]
```

**Flags:**
- `-d, --depth` - Profundidad máxima (default: 3)
- `--dirs` - Solo directorios

### mkdir
Crea un directorio.

```bash
sshcli mkdir <ruta> [-p]
```

### copy
Copia archivos dentro del servidor.

```bash
sshcli copy <origen> <destino>
```

### move
Mueve o renombra archivos.

```bash
sshcli move <origen> <destino>
```

### remove
Elimina archivos o directorios.

```bash
sshcli remove <ruta> [-r] [-f]
```

### chmod
Cambia permisos.

```bash
sshcli chmod <permisos> <ruta> [-r]
```

---

## Transferencia

### upload
Sube archivo o carpeta al servidor.

```bash
sshcli upload <local> <remoto>
```

### download
Descarga archivo o carpeta del servidor.

```bash
sshcli download <remoto> <local>
```

---

## Búsqueda

### find
Busca archivos por nombre.

```bash
sshcli find <ruta> [-n nombre] [-t tipo]
```

### grep
Busca texto en archivos.

```bash
sshcli grep "<patron>" <ruta> [-r] [-i]
```

### wc
Cuenta líneas, palabras y bytes.

```bash
sshcli wc <archivo> [-l] [-w] [-c]
```

### diff
Compara archivo local con remoto.

```bash
sshcli diff <local> <remoto>
```

---

## Git

| Comando | Descripción |
|---------|-------------|
| `git-status` | Estado del repositorio |
| `git-diff` | Ver cambios |
| `git-log` | Historial de commits |
| `git-add` | Agregar al staging |
| `git-commit` | Crear commit |
| `git-pull` | Actualizar desde remoto |
| `git-push` | Enviar al remoto |
| `git-branch` | Gestionar ramas |
| `git-checkout` | Cambiar rama |
| `git-clone` | Clonar repositorio |

---

## Procesos y Sistema

### ps
Lista procesos.

```bash
sshcli ps [--filter nombre] [--all]
```

### kill
Termina un proceso.

```bash
sshcli kill <pid_o_nombre> [--signal señal]
```

### daemon
Ejecuta comando en background.

```bash
sshcli daemon "<comando>" [--name nombre] [--log archivo]
```

### service
Gestiona servicios systemd.

```bash
sshcli service <nombre> <accion>
```

Acciones: start, stop, restart, status, enable, disable, reload

### ports
Lista puertos en uso.

```bash
sshcli ports [-l]
```

### disk
Muestra uso de disco.

```bash
sshcli disk [ruta]
```

### memory
Muestra uso de memoria.

```bash
sshcli memory
```

### env
Muestra variables de entorno o archivo .env.

```bash
sshcli env [archivo]
```

### env-set
Establece variable en archivo .env.

```bash
sshcli env-set <archivo> "KEY=value"
```

---

## Información

### status
Estado de conexión al servidor.

```bash
sshcli status
```

### exists
Verifica si archivo/directorio existe.

```bash
sshcli exists <ruta>
```

Retorna código 0 si existe, 1 si no.

### info
Muestra información detallada de archivo.

```bash
sshcli info <ruta>
```
