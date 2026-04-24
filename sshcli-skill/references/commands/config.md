# config — Configuración

## Descripción
Muestra y modifica la configuración de sshcli.

## Sintaxis
```bash
sshcli config show              # Ver toda la config
sshcli config set <key> <valor> # Modificar valor
```

## Ejemplos
```bash
# Ver configuración
sshcli config show

# Establecer valor
sshcli config set "timeout" "30"
```

## Notas
- Guardado en archivo de configuración local
- Affecta comportamiento global de sshcli
