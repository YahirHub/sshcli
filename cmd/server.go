package cmd

import (
	"fmt"

	"sshcli/internal/config"
	"sshcli/internal/ssh"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Gestiona servidores SSH configurados",
	Long: `Permite agregar, eliminar, listar y seleccionar servidores SSH.
Puedes tener múltiples servidores y cambiar entre ellos fácilmente.

Ejemplos:
  sshcli server add produccion --host prod.example.com --user deploy --pass secreto
  sshcli server add desarrollo --host dev.local --port 2222 --user admin --pass clave
  sshcli server list
  sshcli server use produccion
  sshcli server remove desarrollo`,
}

var (
	serverHost string
	serverPort int
	serverUser string
	serverPass string
)

var serverAddCmd = &cobra.Command{
	Use:   "add [nombre]",
	Short: "Agrega un nuevo servidor SSH",
	Long: `Agrega un nuevo servidor SSH a la configuración.
Valida la conexión antes de guardar.

Ejemplos:
  sshcli server add produccion --host prod.example.com --user deploy --pass secreto
  sshcli server add dev --host 192.168.1.100 --port 2222 --user root --pass clave`,
	Args: cobra.ExactArgs(1),
	RunE: runServerAdd,
}

var serverRemoveCmd = &cobra.Command{
	Use:   "remove [nombre]",
	Short: "Elimina un servidor de la configuración",
	Long: `Elimina un servidor SSH de la lista de servidores configurados.

Ejemplo:
  sshcli server remove desarrollo`,
	Args: cobra.ExactArgs(1),
	RunE: runServerRemove,
}

var serverListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista todos los servidores configurados",
	Long: `Muestra todos los servidores SSH configurados.
El servidor activo se marca con un asterisco (*).

Ejemplo:
  sshcli server list`,
	RunE: runServerList,
}

var serverUseCmd = &cobra.Command{
	Use:   "use [nombre]",
	Short: "Selecciona el servidor activo",
	Long: `Cambia el servidor activo que se usará por defecto en los comandos.

Ejemplo:
  sshcli server use produccion`,
	Args: cobra.ExactArgs(1),
	RunE: runServerUse,
}

var serverInfoCmd = &cobra.Command{
	Use:   "info [nombre]",
	Short: "Muestra información detallada de un servidor",
	Long: `Muestra los detalles de configuración de un servidor específico.
Si no se especifica nombre, muestra el servidor activo.

Ejemplos:
  sshcli server info
  sshcli server info produccion`,
	Args: cobra.MaximumNArgs(1),
	RunE: runServerInfo,
}

func init() {
	rootCmd.AddCommand(serverCmd)
	
	serverCmd.AddCommand(serverAddCmd)
	serverAddCmd.Flags().StringVar(&serverHost, "host", "", "Dirección del servidor SSH (requerido)")
	serverAddCmd.Flags().IntVar(&serverPort, "port", 22, "Puerto SSH (por defecto: 22)")
	serverAddCmd.Flags().StringVar(&serverUser, "user", "", "Nombre de usuario SSH (requerido)")
	serverAddCmd.Flags().StringVar(&serverPass, "pass", "", "Contraseña SSH (requerido)")
	serverAddCmd.MarkFlagRequired("host")
	serverAddCmd.MarkFlagRequired("user")
	serverAddCmd.MarkFlagRequired("pass")
	
	serverCmd.AddCommand(serverRemoveCmd)
	serverCmd.AddCommand(serverListCmd)
	serverCmd.AddCommand(serverUseCmd)
	serverCmd.AddCommand(serverInfoCmd)
}

func runServerAdd(cmd *cobra.Command, args []string) error {
	name := args[0]
	
	fmt.Printf("Validando conexión a %s@%s:%d...\n", serverUser, serverHost, serverPort)
	
	client, err := ssh.Connect(serverHost, serverPort, serverUser, serverPass)
	if err != nil {
		return fmt.Errorf("error de conexión: %v", err)
	}
	defer client.Close()
	
	output, err := client.Run("echo 'Conexión exitosa'")
	if err != nil {
		return fmt.Errorf("error al validar conexión: %v", err)
	}
	fmt.Print(output)
	
	cfg := config.LoadOrCreate()
	
	server := &config.Server{
		Name:     name,
		Host:     serverHost,
		Port:     serverPort,
		User:     serverUser,
		Password: serverPass,
	}
	
	cfg.AddServer(server)
	
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("error al guardar configuración: %v", err)
	}
	
	fmt.Printf("Servidor '%s' agregado exitosamente\n", name)
	if cfg.ActiveServer == name {
		fmt.Printf("Servidor '%s' establecido como activo\n", name)
	}
	
	return nil
}

func runServerRemove(cmd *cobra.Command, args []string) error {
	name := args[0]
	
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("no hay configuración: %v", err)
	}
	
	if err := cfg.RemoveServer(name); err != nil {
		return err
	}
	
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("error al guardar configuración: %v", err)
	}
	
	fmt.Printf("Servidor '%s' eliminado\n", name)
	return nil
}

func runServerList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("No hay servidores configurados")
		fmt.Println("Usa 'sshcli server add' para agregar uno")
		return nil
	}
	
	servers := cfg.ListServers()
	if len(servers) == 0 {
		fmt.Println("No hay servidores configurados")
		return nil
	}
	
	fmt.Println("Servidores configurados:")
	fmt.Println()
	for _, name := range servers {
		server := cfg.Servers[name]
		marker := "  "
		if name == cfg.ActiveServer {
			marker = "* "
		}
		fmt.Printf("%s%-15s %s@%s:%d\n", marker, name, server.User, server.Host, server.Port)
	}
	fmt.Println()
	fmt.Println("(*) = servidor activo")
	
	return nil
}

func runServerUse(cmd *cobra.Command, args []string) error {
	name := args[0]
	
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("no hay configuración: %v", err)
	}
	
	if err := cfg.SetActiveServer(name); err != nil {
		return err
	}
	
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("error al guardar configuración: %v", err)
	}
	
	fmt.Printf("Servidor activo: %s\n", name)
	return nil
}

func runServerInfo(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("no hay configuración: %v", err)
	}
	
	var server *config.Server
	var name string
	
	if len(args) > 0 {
		name = args[0]
		server, err = cfg.GetServer(name)
	} else {
		server, err = cfg.GetActiveServer()
		name = cfg.ActiveServer
	}
	
	if err != nil {
		return err
	}
	
	fmt.Printf("Servidor: %s\n", name)
	fmt.Printf("Host: %s\n", server.Host)
	fmt.Printf("Puerto: %d\n", server.Port)
	fmt.Printf("Usuario: %s\n", server.User)
	if name == cfg.ActiveServer {
		fmt.Println("Estado: ACTIVO")
	}
	
	return nil
}