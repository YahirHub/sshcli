# disk — Uso de disco

## Descripción
Muestra el uso de disco del sistema.

## Sintaxis
```bash
sshcli disk [ruta] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Disco raíz
sshcli disk /

# Directorio específico
sshcli disk /home
```

## Output Example
```
Filesystem      Size  Used Avail Use% Mounted on
/dev/sda2        46G   26G   18G  61% /
```

## Notas
- Si no se especifica ruta, usa `/`
- Útil para verificar espacio antes de instalaciones
