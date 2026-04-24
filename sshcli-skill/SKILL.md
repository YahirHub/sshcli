---
name: sshcli
description: >
  Gestión profesional de servidores remotos Linux mediante `sshcli`. Activa este skill para 
  administración de sistemas, despliegues, edición de código remota, monitoreo de salud, 
  gestión de servicios (systemd), contenedores, procesos y flujos Git en servidores. 
  Optimizado para entornos donde el cliente es Windows (Git Bash/MSYS2) y el servidor es Linux.
---

# Skill: Gestión de Servidores Remotos via SSHCLI

## 🛠 Configuración e Inspección Inicial

Antes de realizar cualquier operación, verifica el entorno:

```bash
sshcli server list    # Ver servidores configurados (* = activo)
sshcli status         # Verificar conexión y salud del servidor activo
```

Para agregar un nuevo servidor (solicitar Host, User, Pass al usuario):
```bash
sshcli server add <nombre> --host <ip> --user <user> --pass "<password>"
```

**IMPORTANTE:** Si la contraseña contiene caracteres especiales (`!`, `$`, `` ` ``, `"`, `\`, etc.), debe ir entre comillas dobles para evitar que el shell los interprete. Ejemplo: `--pass "Linux_0145!"`

---

## 🛣 Manejo de Rutas (CRÍTICO para Windows)

`sshcli` normaliza automáticamente las rutas para evitar conflictos entre Windows y Linux.

1.  **Rutas Locales (Tu entorno actual):**
    - **REGLA DE ORO:** Usa el formato MSYS2/Unix: `/c/Users/Admin/Desktop/file.txt`.
    - Esto evita que las barras invertidas `\` sean malinterpretadas por el shell.
2.  **Rutas Remotas (El Servidor):**
    - Usa siempre rutas absolutas de Linux: `/var/www/app/` o `/home/user/script.sh`.
    - **Protección Automática:** El binario detecta y elimina prefijos basura como `/Program Files/Git/` que Git Bash inyecta erróneamente.

---

## 📂 Protocolo de Gestión de Archivos

Sigue este flujo para garantizar la integridad del código:

### 1. Exploración y Lectura
- **Mapear:** `sshcli tree /ruta -d 2` para entender la jerarquía.
- **Buscar:** `sshcli find /ruta -name "*.conf"` para localizar archivos.
- **Leer:** `sshcli cat-lines /ruta/archivo 1 100` para archivos largos o `sshcli read /ruta/archivo` para cortos.

### 2. Creación y Edición (Atómica)
| Acción | Comando | Ventaja |
|--------|---------|---------|
| **Crear Script** | `sshcli write /ruta/file.sh "contenido" -x` | **`-x`** otorga permisos de ejecución (755) al instante. |
| **Crear Config** | `sshcli write /ruta/file.json "content"` | Crea directorios padres automáticamente (mkdir -p). |
| **Modificar** | `sshcli search-replace /ruta "viejo" "nuevo" -a` | Edición quirúrgica sin reescribir todo el archivo. |
| **Insertar** | `sshcli insert-line /ruta 1 "import os"` | Ideal para añadir dependencias o cabeceras. |

### 3. Validación (OBLIGATORIA)
Después de CUALQUIER edición, verifica la sintaxis antes de reiniciar servicios:
```bash
sshcli syntax-check /ruta/archivo
```

---

## 🚀 Ejecución de Comandos y Procesos

- **Modo Texto:** `sshcli exec "comando"` (Rápido, para `ls`, `cat`, `uptime`).
- **Modo Interactivo (TTY):** `sshcli exec -t "comando"` (Obligatorio para `apt install`, `htop`, `vim` o comandos que pidan confirmación).
- **Scripts:** `sshcli run /ruta/script.py --args "--port 80"`. (Detecta intérprete por extensión).
- **Background:** `sshcli daemon "python app.py" --log /var/log/app.log`. (El proceso persiste tras desconectarte).

---

## 📊 Monitoreo y Sistema

| Comando | Función |
|---------|---------|
| `sshcli disk /` | Verificar espacio en disco (evita fallos de escritura). |
| `sshcli memory` | Monitorear consumo de RAM. |
| `sshcli ports -l` | Ver puertos escuchando (LISTEN) y sus PIDs. |
| `sshcli service <n> <acc>` | Gestionar systemd (status, restart, stop, start). |
| `sshcli env-set /path/.env "K=V"` | Actualizar configuraciones de entorno de forma segura. |

---

## 🛡 Reglas de Seguridad para el Agente

1.  **Rutas Absolutas:** No confíes en rutas relativas; usa siempre el path completo.
2.  **Validación Post-Edit:** Nunca asumas que un `write` o `replace` fue perfecto. Ejecuta `syntax-check`.
3.  **Permisos:** Si un script falla con "Permission denied", usa la flag `-x` en `write` o ejecuta `sshcli chmod +x /ruta`.
4.  **Borrado Seguro:** Antes de un `remove -rf`, ejecuta `list` para confirmar que estás en la carpeta correcta.
5.  **Transferencias:** Para subir archivos locales a remoto, usa `sshcli upload "/c/ruta/local" "/ruta/remota"`.

---
*Este skill requiere el binario `sshcli` compilado y disponible en el PATH del sistema.*