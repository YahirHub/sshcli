Ya tienes un archivo `README.md` muy completo y profesional en el contenido que proporcionaste. Sin embargo, si deseas **optimizarlo para GitHub** (añadiendo insignias, mejorando el formato visual y haciendo la lectura más ágil), aquí tienes una versión mejorada basada en tu contexto:

---

# 🚀 SSHCLI - Cliente SSH Avanzado para Agentes de IA

![Go Version](https://img.shields.io/badge/go-1.25.0-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Platform](https://img.shields.io/badge/platform-linux%20|%20macos%20|%20windows-lightgrey.svg)

**sshcli** es un cliente SSH de línea de comandos diseñado específicamente para que **agentes de IA** y **modelos de lenguaje (LLMs)** interactúen de forma atómica, segura y eficiente con servidores Linux remotos.

---

## 💡 ¿Por qué SSHCLI?

A diferencia de clientes SSH tradicionales, **sshcli** está diseñado para ser "amigable con las máquinas":
*   **Comandos Atómicos:** Operaciones de archivo, Git y sistema con salidas limpias.
*   **Multi-Servidor:** Gestione flotas de servidores desde una configuración centralizada.
*   **Tooling para IA:** Especialmente útil para entornos donde una IA necesita leer código, verificar sintaxis, editar archivos mediante reemplazos programáticos y gestionar contenedores Docker.
*   **Multiplataforma:** Funciona en Windows (Git Bash/CMD/PS), Linux y macOS con un binario único.

---

## 🛠️ Instalación

### Desde el código fuente
```bash
git clone https://github.com/YahirHub/sshcli
cd sshcli
go build -o sshcli .
```

### Binarios
Puedes descargar los binarios precompilados desde la sección de **Releases** de este repositorio.

---

## 🚀 Inicio Rápido

### 1. Configuración
```bash
# Agregar un servidor
sshcli server add produccion --host 192.168.1.100 --user root --pass clave

# Listar servidores
sshcli server list
```

### 2. Comandos Comunes
```bash
# Ejecutar un comando
sshcli exec "uptime"

# Leer un archivo
sshcli read /var/www/app/.env

# Editar código quirúrgicamente
sshcli search-replace /app/config.py "DEBUG = True" "DEBUG = False"

# Gestionar Docker
sshcli docker ps
sshcli docker logs <container_id>
```

---

## 📋 Categorías de Comandos

| Categoría | Ejemplos |
| :--- | :--- |
| **Gestión** | `server add`, `server list`, `server use` |
| **Ejecución** | `exec`, `run` (py, js, go, sh), `syntax-check` |
| **Archivos** | `read`, `write`, `search-replace`, `insert-line`, `delete-line` |
| **Exploración** | `tree`, `list`, `find`, `grep`, `cat-lines` |
| **Git** | `git-status`, `git-commit`, `git-pull`, `git-push` |
| **Sistema** | `ps`, `kill`, `service`, `ports`, `disk`, `memory` |

---

## 🤖 Uso con Agentes de IA

**sshcli** es la herramienta ideal para flujos de trabajo autónomos. Ejemplo de una sesión de agente:

```bash
# 1. El agente explora el entorno
sshcli tree /app -d 2

# 2. El agente detecta un error de sintaxis y lo corrige
sshcli syntax-check /app/main.py
sshcli replace-line /app/main.py 5 "import logging"

# 3. Despliegue automático
sshcli git-commit /app -m "fix: importar logging"
sshcli git-push /app
sshcli service --server prod gunicorn restart
```

---

## ⚙️ Configuración

Puedes personalizar el comportamiento global:

```bash
# Habilitar TTY por defecto para todos los comandos 'exec'
sshcli config set tty true
```

Toda la configuración se guarda de forma segura en `~/.sshcli.conf`.

---

## 🏗️ Desarrollo

¿Quieres contribuir? ¡Las PRs son bienvenidas!

1. **Requisitos:** Go 1.25+
2. **Estructura:**
   - `cmd/`: Lógica de los comandos CLI (usando `cobra`).
   - `internal/ssh/`: Wrapper de conexión SSH/SFTP.
   - `internal/paths/`: Normalizador de rutas (crucial para Windows/Linux).
3. **Testing:** `go test ./...`

---

## 📜 Licencia
Este proyecto está bajo la licencia **MIT**.

---
*Desarrollado para facilitar la administración remota de sistemas mediante automatización.*