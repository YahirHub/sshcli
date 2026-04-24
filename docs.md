# Documentación Oficial de SSHCLI

**SSHCLI** es un cliente SSH avanzado diseñado específicamente para ser la interfaz principal de agentes de IA y desarrolladores que administran servidores Linux de forma remota, atómica y profesional.

---

## 1. Instalación y Configuración

### Instalación
1. Clona el repositorio: `git clone <url-del-repo>`
2. Compila el binario: `go build -o sshcli .`
3. (Opcional) Mueve el binario a tu `PATH` (`/usr/local/bin` en Unix, `System32` en Windows).

### Configuración de Servidores
Para operar, debes registrar al menos un servidor:
```bash
sshcli server add <nombre> --host <host> --user <usuario> --pass <contraseña> [--port <puerto>]
```
*   **Listar:** `sshcli server list`
*   **Seleccionar:** `sshcli server use <nombre>`
*   **Eliminar:** `sshcli server remove <nombre>`

### Configuración Global
```bash
sshcli config set tty true/false  # Habilita o deshabilita TTY interactivo por defecto
sshcli config show                # Muestra el estado actual
```

---

## 2. Ejecución y Control

### Ejecución de Comandos
*   `exec "<comando>"`: Ejecuta comandos simples.
*   `exec -t "<comando>"`: Ejecuta en modo **PTY** (interactivo), necesario para comandos como `htop`, `vim`, o procesos que piden confirmación.
*   `--no-tty`: Forzar modo no interactivo aunque la configuración global diga lo contrario.

### Ejecución de Scripts (`run`)
Detecta automáticamente el intérprete según la extensión:
*   `.py` (python3), `.js` (node), `.ts` (npx tsx), `.go` (go run), `.sh` (bash), `.rb` (ruby), `.php` (php).
*   **Flags:** `--args`, `--env`, `--workdir`.

### Verificación de Sintaxis
`sshcli syntax-check <archivo>` valida errores antes de ejecutar. Soporta: JS, TS, Go, Python, Bash, Ruby, PHP, JSON y YAML.

---

## 3. Edición Quirúrgica de Archivos

Estos comandos permiten modificar archivos sin necesidad de abrir un editor interactivo:

*   **Lectura:** `read`, `cat-lines` (rango), `head`, `tail`.
*   **Creación:** `write` (desde `stdin` o argumento).
*   **Edición:**
    *   `insert-line`: Añade una línea en una posición específica.
    *   `delete-line`: Borra líneas específicas o rangos.
    *   `replace-line`: Cambia una línea completa.
    *   `search-replace`: Reemplazo de texto (flag `-a` para todas las ocurrencias).
    *   `append`: Añade texto al final del archivo.

---

## 4. Gestión de Archivos y Directorios

*   **Listado:** `list` (`-l` formato largo, `-a` ocultos).
*   **Estructura:** `tree` (muestra el árbol de archivos, soporte para profundidad `-d`).
*   **Manipulación:** `mkdir`, `copy`, `move`, `remove` (`-r` recursivo, `-f` fuerza).
*   **Permisos:** `chmod <modo> <archivo>`.
*   **Información:** `info` (muestra metadatos del archivo).
*   **Existencia:** `exists <ruta>` (retorna 0 si existe, 1 si no).

---

## 5. Transferencia de Archivos

*   **Subir:** `upload <local> <remoto>`
*   **Descargar:** `download <remoto> <local>`
    *   Ambos comandos soportan rutas absolutas y relativas. La herramienta normaliza automáticamente las rutas estilo Windows (`C:/`) a rutas Linux.

---

## 6. Git, Procesos y Sistema

### Integración Git
Comandos equivalentes a la CLI de git, optimizados para ejecutar remotamente:
`git-status`, `git-log`, `git-diff`, `git-add`, `git-commit`, `git-pull`, `git-push`, `git-branch`, `git-checkout`, `git-clone`.

### Diagnóstico de Sistema
*   **Procesos:** `ps` (filtro por nombre), `kill`.
*   **Recursos:** `memory` (uso de RAM), `disk` (uso de espacio), `ports` (puertos escuchando).
*   **Servicios:** `service <nombre> <accion>` (start, stop, restart, status, enable, disable, reload).
*   **Snapshot:** `project-snapshot <ruta>` (genera un reporte de todo el proyecto: estructura, commits, recursos, docker).

### Docker
*   `docker ps`: Lista contenedores.
*   `docker logs <ID>`: Muestra logs.
*   `docker exec <ID> <cmd>`: Ejecuta comando dentro del contenedor.
*   `docker stats`: Consumo de recursos.

---

## 7. Notas para Agentes de IA

1.  **Normalización de Rutas:** No te preocupes por el sistema operativo desde el que llamas a `sshcli`. El motor interno (`internal/paths`) convierte automáticamente las rutas de Git Bash/Windows a formatos válidos para el servidor Linux remoto.
2.  **Modo TTY:** Si el agente necesita interactuar con un proceso, usa siempre `-t`. Si el comando es un script de automatización, no uses `-t` para evitar caracteres de control innecesarios en el output.
3.  **Encadenamiento:** Como todas las salidas son texto plano, puedes redirigir la salida a otros archivos o analizarlas con un LLM fácilmente.
4.  **Seguridad:** Toda la configuración se guarda con permisos `0600` para proteger las contraseñas guardadas en `~/.sshcli.conf`.

---
*Para más detalles sobre un comando específico, usa `sshcli <comando> --help`.*