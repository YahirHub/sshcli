---
name: sshcli
description: >
  Gestión profesional de servidores remotos Linux mediante `sshcli`. Activa este skill para 
  administración de sistemas, despliegues, edición de código remota, monitoreo de salud, 
  gestión de servicios (systemd), contenedores, procesos y flujos Git en servidores. 
  Optimizado para entornos Windows (Git Bash/MSYS2) -> Linux remoto.
---

# Skill: Gestión de Servidores Remotos via SSHCLI

> Validado y corregido contra Windows Git Bash/MSYS2 + servidor Debian 12 remoto.

## 🛠 Configuración e Inspección Inicial

Antes de operar, valida contexto y conexión:

```bash
sshcli server list
sshcli status
```

Para agregar un servidor:

```bash
sshcli server add <nombre> --host <ip> --user <user> --pass "<password>"
```

Si la contraseña tiene caracteres especiales, usa comillas dobles.

---

## 🛣 Manejo de Rutas (CRÍTICO para Windows)

1. **Rutas locales:** usa formato MSYS2/Unix:
   - `/c/Users/Admin/Desktop/file.txt`
2. **Rutas remotas:** usa rutas absolutas Linux:
   - `/var/www/app`
   - `/home/root/script.sh`

La normalización de rutas quedó validada en `upload`, `download`, `diff` y escritura remota.

---

## 📂 Protocolo de Gestión de Archivos

### 1. Exploración y lectura
- **Mapear:** `sshcli tree /ruta -d 2 --dirs`
- **Buscar archivos:** `sshcli find /ruta --name "*.conf"`
- **Leer:**
  - `sshcli cat-lines /ruta/archivo 1 100 -n`
  - `sshcli read /ruta/archivo`
  - `sshcli head /ruta/archivo -n 20`
  - `sshcli tail /ruta/archivo -n 20`
  - `sshcli wc /ruta/archivo`

### 2. Creación y edición
| Acción | Comando recomendado | Nota |
|---|---|---|
| Crear archivo simple | `sshcli write /ruta/file.txt "contenido"` | OK |
| Crear archivo multilínea | `sshcli write /ruta/file.txt "linea1\nlinea2\n"` | `\n` ya funciona |
| Crear desde stdin | `printf 'a\nb\n' \| sshcli write /ruta/file.txt` | También funciona |
| Crear script ejecutable | `sshcli write /ruta/script.sh "#!/bin/bash\necho ok\n" -x` | `-x` deja 755 |
| Reemplazar texto | `sshcli search-replace /ruta/file "viejo" "nuevo" -a` | OK |
| Insertar al inicio | `sshcli insert-line /ruta/file 1 "import os"` | Línea 1 = inicio |
| Insertar con alias al inicio | `sshcli insert-line /ruta/file 0 "import os"` | También válido |
| Reemplazar línea | `sshcli replace-line /ruta/file 5 "nuevo contenido"` | OK |
| Eliminar línea(s) | `sshcli delete-line /ruta/file 10 20` | OK |
| Agregar al final | `sshcli append /ruta/file "\ntexto"` | Soporta `\n` |

### 3. Validación post-edición
Después de **cada cambio**:

```bash
sshcli syntax-check /ruta/archivo
sshcli cat-lines /ruta/archivo 1 80 -n
```

---

## 🚀 Ejecución de Comandos y Procesos

- **Modo texto:** `sshcli exec "comando"`
- **Modo interactivo / PTY:** `sshcli exec -t "comando"`
- **Shell completa:** `sshcli connect`
- **Scripts:** `sshcli run /ruta/script.py`
- **Segundo plano:** `sshcli daemon "python app.py" --log /var/log/app.log`
- **Procesos:** `sshcli ps`, `sshcli kill <pid|nombre>`

`exec -t` quedó corregido para solicitar PTY. `connect` funciona, pero requiere una terminal real para uso cómodo.

---

## 🧰 Git remoto

Comandos validados:

```bash
sshcli git-status /app
sshcli git-diff /app
sshcli git-add /app archivo.py
sshcli git-commit /app -m "mensaje"
sshcli git-log /app --oneline -n 10
sshcli git-branch /app -c feature/x
sshcli git-checkout /app feature/x
sshcli git-clone /repo/origen.git /app-clone
sshcli git-pull /app
sshcli git-push /app
```

Notas:
- `git-push` requiere un remoto configurado.
- `git-branch` y `git-checkout` pueden imprimir poca salida, pero funcionan.

---

## 📊 Monitoreo y sistema

| Comando | Estado |
|---|---|
| `sshcli disk /` | OK |
| `sshcli memory` | OK |
| `sshcli ports -l` | OK |
| `sshcli service <nombre> status` | OK |
| `sshcli env-set /path/.env "K=V"` | OK |
| `sshcli docker ps` | Requiere Docker instalado en el servidor |
| `sshcli project-snapshot /app` | OK |
| `sshcli find-code "pattern" /app` | OK |

---

## 🛡 Reglas de seguridad para el agente

1. Usa siempre rutas absolutas.
2. Después de editar, valida con `syntax-check` y vuelve a leer el archivo.
3. Para scripts, usa `write -x` o `sshcli chmod +x /ruta/script.sh`.
4. Antes de `remove -rf`, ejecuta `list`.
5. Para transferencias Windows -> Linux, usa `upload "/c/..." "/ruta/remota"`.
6. `exists` devuelve `NO` y sale con código de error si la ruta no existe; úsalo considerando ese comportamiento.

---
*Este skill requiere el binario `sshcli` disponible en el PATH.*
