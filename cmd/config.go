package cmd

import (
	"fmt"
	"strings"

	"sshcli/internal/config"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Gestiona la configuración global de sshcli",
	Long: `Permite ver y modificar la configuración global del cliente SSH.

Ejemplos:
  sshcli config show                # Ver configuración actual
  sshcli config set tty true        # Habilitar -t por defecto
  sshcli config set tty false       # Deshabilitar -t por defecto

Opciones disponibles:
  tty     - Usar modo interactivo por defecto (true/false)`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Muestra la configuración actual",
	RunE:  runConfigShow,
}

var configSetCmd = &cobra.Command{
	Use:   "set [opcion] [valor]",
	Short: "Establece un valor de configuración",
	Long: `Establece un valor de configuración global.

Ejemplos:
  sshcli config set tty true     # Habilitar TTY por defecto
  sshcli config set tty false    # Deshabilitar TTY por defecto`,
	Args: cobra.ExactArgs(2),
	RunE: runConfigSet,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("No hay configuración guardada")
		fmt.Println("Usa 'sshcli server add' para comenzar")
		return nil
	}

	fmt.Println("=== Configuración Global ===")
	fmt.Printf("Servidor activo: %s\n", cfg.ActiveServer)
	fmt.Printf("TTY por defecto: %v\n", cfg.DefaultTTY)
	fmt.Printf("Servidores: %d\n", len(cfg.Servers))

	return nil
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	option := strings.ToLower(args[0])
	value := strings.ToLower(args[1])

	cfg := config.LoadOrCreate()

	switch option {
	case "tty":
		if value == "true" || value == "1" || value == "on" || value == "yes" {
			cfg.DefaultTTY = true
			fmt.Println("TTY habilitado por defecto")
			fmt.Println("Todos los comandos exec usarán modo interactivo")
			fmt.Println("Usa --no-tty para desactivar temporalmente")
		} else if value == "false" || value == "0" || value == "off" || value == "no" {
			cfg.DefaultTTY = false
			fmt.Println("TTY deshabilitado por defecto")
			fmt.Println("Usa -t para activar modo interactivo")
		} else {
			return fmt.Errorf("valor inválido: usa true/false")
		}
	default:
		return fmt.Errorf("opción desconocida: %s (opciones: tty)", option)
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("error al guardar: %v", err)
	}

	return nil
}
