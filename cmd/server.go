package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"sshcli/internal/config"
	"sshcli/internal/ssh"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Gestiona servidores SSH configurados",
	Long: `Permite agregar, editar, copiar, exportar e importar servidores SSH configurados.
Puedes tener múltiples servidores y cambiar entre ellos fácilmente.`,
}

var (
	addHost  string
	addPort  int
	addUser  string
	addPass  string
	addForce bool
	addTags  []string

	editHost       string
	editPort       int
	editUser       string
	editPass       string
	editTags       []string
	editAddTags    []string
	editRemoveTags []string
	editClearTags  bool

	copyHost string
	copyPort int
	copyUser string
	copyPass string
	copyTags []string

	removeYes bool

	listTags   []string
	listFormat string

	infoShowPass bool

	importForce   bool
	importReplace bool

	setPassValue string
)

type serverListItem struct {
	Name   string   `json:"name"`
	Host   string   `json:"host"`
	Port   int      `json:"port"`
	User   string   `json:"user"`
	Tags   []string `json:"tags,omitempty"`
	Active bool     `json:"active"`
}

type doctorCheck struct {
	Name   string `json:"name"`
	OK     bool   `json:"ok"`
	Detail string `json:"detail,omitempty"`
}

var serverAddCmd = &cobra.Command{
	Use:   "add [nombre]",
	Short: "Agrega un nuevo servidor SSH",
	Long: `Agrega un nuevo servidor SSH a la configuración.
Valida la conexión antes de guardar.`,
	Args: cobra.ExactArgs(1),
	RunE: runServerAdd,
}

var serverEditCmd = &cobra.Command{
	Use:   "edit [nombre]",
	Short: "Edita un servidor configurado",
	Long: `Edita host, puerto, usuario, contraseña o tags de un servidor existente.
Si cambias datos de conexión, valida el acceso antes de guardar.`,
	Args: cobra.ExactArgs(1),
	RunE: runServerEdit,
}

var serverCopyCmd = &cobra.Command{
	Use:   "copy [origen] [nuevo]",
	Short: "Copia la configuración de un servidor",
	Long: `Duplica un servidor configurado bajo un nuevo nombre.
Puedes sobrescribir host, puerto, usuario, contraseña o tags en la copia.`,
	Args: cobra.ExactArgs(2),
	RunE: runServerCopy,
}

var serverRemoveCmd = &cobra.Command{
	Use:   "remove [nombre]",
	Short: "Elimina un servidor de la configuración",
	Args:  cobra.ExactArgs(1),
	RunE:  runServerRemove,
}

var serverRenameCmd = &cobra.Command{
	Use:   "rename [nombre_actual] [nombre_nuevo]",
	Short: "Renombra un servidor configurado",
	Long: `Renombra un servidor SSH existente.
Valida que el servidor actual exista y que el nuevo nombre no esté en uso.`,
	Args: cobra.ExactArgs(2),
	RunE: runServerRename,
}

var serverListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista todos los servidores configurados",
	Long: `Muestra todos los servidores SSH configurados.
El servidor activo se marca con un asterisco (*).`,
	RunE: runServerList,
}

var serverSearchCmd = &cobra.Command{
	Use:   "search [texto]",
	Short: "Busca servidores por nombre, host, usuario o tag",
	Args:  cobra.ExactArgs(1),
	RunE:  runServerSearch,
}

var serverUseCmd = &cobra.Command{
	Use:   "use [nombre]",
	Short: "Selecciona el servidor activo",
	Args:  cobra.ExactArgs(1),
	RunE:  runServerUse,
}

var serverInfoCmd = &cobra.Command{
	Use:   "info [nombre]",
	Short: "Muestra información detallada de un servidor",
	Long: `Muestra los detalles de configuración de un servidor específico.
Si no se especifica nombre, muestra el servidor activo.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runServerInfo,
}

var serverTestCmd = &cobra.Command{
	Use:   "test [nombre]",
	Short: "Valida la conexión a un servidor configurado",
	Long: `Prueba la conexión SSH usando la configuración guardada.
Si no se especifica nombre, usa el servidor activo.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runServerTest,
}

var serverPingCmd = &cobra.Command{
	Use:   "ping [nombre]",
	Short: "Mide latencia de conexión SSH",
	Long: `Mide el tiempo de establecimiento de conexión SSH.
Si no se especifica nombre, usa el servidor activo.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runServerPing,
}

var serverDoctorCmd = &cobra.Command{
	Use:   "doctor [nombre]",
	Short: "Diagnostica utilidades comunes del servidor",
	Long: `Valida conectividad SSH y la presencia de herramientas comunes como sudo, git, docker, php, composer y caddy.
Si no se especifica nombre, usa el servidor activo.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runServerDoctor,
}

var serverSyncTagsCmd = &cobra.Command{
	Use:   "sync-tags [origen] [destino]",
	Short: "Copia los tags de un servidor a otro",
	Args:  cobra.ExactArgs(2),
	RunE:  runServerSyncTags,
}

var serverExportCmd = &cobra.Command{
	Use:   "export [archivo]",
	Short: "Exporta la configuración de servidores a un archivo JSON",
	Args:  cobra.ExactArgs(1),
	RunE:  runServerExport,
}

var serverImportCmd = &cobra.Command{
	Use:   "import [archivo]",
	Short: "Importa servidores desde un archivo JSON",
	Long: `Importa servidores desde un archivo JSON exportado previamente.
Por defecto fusiona con la configuración actual. Usa --replace para reemplazarla completa.`,
	Args: cobra.ExactArgs(1),
	RunE: runServerImport,
}

var serverSetPassCmd = &cobra.Command{
	Use:   "set-pass [nombre]",
	Short: "Actualiza la contraseña guardada de un servidor",
	Long: `Actualiza solo la contraseña almacenada en la configuración local.
Útil después de rotar la contraseña real del usuario SSH en el servidor remoto.`,
	Args: cobra.ExactArgs(1),
	RunE: runServerSetPass,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.AddCommand(serverAddCmd)
	serverAddCmd.Flags().StringVar(&addHost, "host", "", "Dirección del servidor SSH (requerido)")
	serverAddCmd.Flags().IntVar(&addPort, "port", 22, "Puerto SSH (por defecto: 22)")
	serverAddCmd.Flags().StringVar(&addUser, "user", "", "Nombre de usuario SSH (requerido)")
	serverAddCmd.Flags().StringVar(&addPass, "pass", "", "Contraseña SSH (requerido)")
	serverAddCmd.Flags().BoolVar(&addForce, "force", false, "Sobrescribir si el nombre ya existe")
	serverAddCmd.Flags().StringSliceVar(&addTags, "tag", nil, "Tag(s) del servidor (repetible o separadas por coma)")
	serverAddCmd.MarkFlagRequired("host")
	serverAddCmd.MarkFlagRequired("user")
	serverAddCmd.MarkFlagRequired("pass")

	serverCmd.AddCommand(serverEditCmd)
	serverEditCmd.Flags().StringVar(&editHost, "host", "", "Nuevo host")
	serverEditCmd.Flags().IntVar(&editPort, "port", 0, "Nuevo puerto")
	serverEditCmd.Flags().StringVar(&editUser, "user", "", "Nuevo usuario SSH")
	serverEditCmd.Flags().StringVar(&editPass, "pass", "", "Nueva contraseña SSH")
	serverEditCmd.Flags().StringSliceVar(&editTags, "tag", nil, "Reemplazar tags por esta lista")
	serverEditCmd.Flags().StringSliceVar(&editAddTags, "add-tag", nil, "Agregar tag(s)")
	serverEditCmd.Flags().StringSliceVar(&editRemoveTags, "remove-tag", nil, "Eliminar tag(s)")
	serverEditCmd.Flags().BoolVar(&editClearTags, "clear-tags", false, "Eliminar todos los tags")

	serverCmd.AddCommand(serverCopyCmd)
	serverCopyCmd.Flags().StringVar(&copyHost, "host", "", "Host de la copia")
	serverCopyCmd.Flags().IntVar(&copyPort, "port", 0, "Puerto de la copia")
	serverCopyCmd.Flags().StringVar(&copyUser, "user", "", "Usuario SSH de la copia")
	serverCopyCmd.Flags().StringVar(&copyPass, "pass", "", "Contraseña SSH de la copia")
	serverCopyCmd.Flags().StringSliceVar(&copyTags, "tag", nil, "Tags de la copia (reemplaza los originales)")

	serverCmd.AddCommand(serverRemoveCmd)
	serverRemoveCmd.Flags().BoolVarP(&removeYes, "yes", "y", false, "Confirmar eliminación sin preguntar")
	serverCmd.AddCommand(serverRenameCmd)
	serverCmd.AddCommand(serverListCmd)
	serverListCmd.Flags().StringSliceVar(&listTags, "tag", nil, "Filtrar por tag(s)")
	serverListCmd.Flags().StringVar(&listFormat, "format", "text", "Formato de salida: text|json")
	serverCmd.AddCommand(serverSearchCmd)
	serverCmd.AddCommand(serverUseCmd)
	serverCmd.AddCommand(serverInfoCmd)
	serverInfoCmd.Flags().BoolVar(&infoShowPass, "show-pass", false, "Mostrar la contraseña guardada")
	serverCmd.AddCommand(serverTestCmd)
	serverCmd.AddCommand(serverPingCmd)
	serverCmd.AddCommand(serverDoctorCmd)
	serverCmd.AddCommand(serverSyncTagsCmd)
	serverCmd.AddCommand(serverExportCmd)
	serverCmd.AddCommand(serverImportCmd)
	serverImportCmd.Flags().BoolVar(&importForce, "force", false, "Sobrescribir servidores existentes al importar")
	serverImportCmd.Flags().BoolVar(&importReplace, "replace", false, "Reemplazar toda la configuración actual por la importada")
	serverCmd.AddCommand(serverSetPassCmd)
	serverSetPassCmd.Flags().StringVar(&setPassValue, "pass", "", "Nueva contraseña guardada")
	serverSetPassCmd.MarkFlagRequired("pass")
}

func runServerAdd(cmd *cobra.Command, args []string) error {
	name := args[0]
	cfg := config.LoadOrCreate()

	server := &config.Server{
		Name:     name,
		Host:     addHost,
		Port:     addPort,
		User:     addUser,
		Password: addPass,
		Tags:     addTags,
	}

	if err := validateServerConnection(server); err != nil {
		return err
	}
	if err := cfg.AddServer(server, addForce); err != nil {
		return err
	}
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("error al guardar configuración: %v", err)
	}

	fmt.Printf("Servidor '%s' agregado exitosamente\n", name)
	if cfg.ActiveServer == name {
		fmt.Printf("Servidor '%s' establecido como activo\n", name)
	}
	return nil
}

func runServerEdit(cmd *cobra.Command, args []string) error {
	name := args[0]
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("no hay configuración: %v", err)
	}

	current, err := cfg.GetServer(name)
	if err != nil {
		return err
	}
	updated := current.Clone()

	changed := false
	validateConn := false

	if cmd.Flags().Changed("host") {
		updated.Host = editHost
		changed = true
		validateConn = true
	}
	if cmd.Flags().Changed("port") {
		updated.Port = editPort
		changed = true
		validateConn = true
	}
	if cmd.Flags().Changed("user") {
		updated.User = editUser
		changed = true
		validateConn = true
	}
	if cmd.Flags().Changed("pass") {
		updated.Password = editPass
		changed = true
		validateConn = true
	}
	if editClearTags {
		updated.Tags = nil
		changed = true
	}
	if cmd.Flags().Changed("tag") {
		updated.Tags = append([]string(nil), editTags...)
		changed = true
	}
	if cmd.Flags().Changed("add-tag") {
		updated.Tags = append(updated.Tags, editAddTags...)
		changed = true
	}
	if cmd.Flags().Changed("remove-tag") {
		remove := map[string]struct{}{}
		for _, tag := range editRemoveTags {
			remove[strings.TrimSpace(tag)] = struct{}{}
		}
		filtered := make([]string, 0, len(updated.Tags))
		for _, tag := range updated.Tags {
			if _, exists := remove[tag]; !exists {
				filtered = append(filtered, tag)
			}
		}
		updated.Tags = filtered
		changed = true
	}

	if !changed {
		return fmt.Errorf("no se especificaron cambios")
	}
	if validateConn {
		if err := validateServerConnection(updated); err != nil {
			return err
		}
	}

	cfg.Servers[name] = updated
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("error al guardar configuración: %v", err)
	}
	fmt.Printf("Servidor '%s' actualizado\n", name)
	return nil
}

func runServerCopy(cmd *cobra.Command, args []string) error {
	sourceName := args[0]
	newName := args[1]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("no hay configuración: %v", err)
	}
	if _, exists := cfg.Servers[newName]; exists {
		return fmt.Errorf("ya existe un servidor con el nombre '%s'", newName)
	}

	source, err := cfg.GetServer(sourceName)
	if err != nil {
		return err
	}
	copyServer := source.Clone()
	copyServer.Name = newName
	if cmd.Flags().Changed("host") {
		copyServer.Host = copyHost
	}
	if cmd.Flags().Changed("port") {
		copyServer.Port = copyPort
	}
	if cmd.Flags().Changed("user") {
		copyServer.User = copyUser
	}
	if cmd.Flags().Changed("pass") {
		copyServer.Password = copyPass
	}
	if cmd.Flags().Changed("tag") {
		copyServer.Tags = append([]string(nil), copyTags...)
	}

	if err := validateServerConnection(copyServer); err != nil {
		return err
	}
	if err := cfg.AddServer(copyServer, false); err != nil {
		return err
	}
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("error al guardar configuración: %v", err)
	}
	fmt.Printf("Servidor copiado: %s -> %s\n", sourceName, newName)
	return nil
}

func runServerRemove(cmd *cobra.Command, args []string) error {
	name := args[0]
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("no hay configuración: %v", err)
	}
	if !removeYes {
		if !term.IsTerminal(int(os.Stdin.Fd())) {
			return fmt.Errorf("usa --yes para confirmar la eliminación en modo no interactivo")
		}
		fmt.Printf("Confirma la eliminación escribiendo el nombre '%s': ", name)
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		if strings.TrimSpace(input) != name {
			return fmt.Errorf("confirmación inválida; eliminación cancelada")
		}
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

func runServerRename(cmd *cobra.Command, args []string) error {
	oldName := args[0]
	newName := args[1]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("no hay configuración: %v", err)
	}
	if err := cfg.RenameServer(oldName, newName); err != nil {
		return err
	}
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("error al guardar configuración: %v", err)
	}

	fmt.Printf("Servidor renombrado: %s -> %s\n", oldName, newName)
	if cfg.ActiveServer == newName {
		fmt.Printf("Servidor activo: %s\n", newName)
	}
	return nil
}

func runServerList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		if listFormat == "json" {
			fmt.Println("[]")
			return nil
		}
		fmt.Println("No hay servidores configurados")
		fmt.Println("Usa 'sshcli server add' para agregar uno")
		return nil
	}

	items := collectServerItems(cfg, "", listTags)
	if strings.EqualFold(listFormat, "json") {
		data, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			return fmt.Errorf("error al serializar salida: %v", err)
		}
		fmt.Println(string(data))
		return nil
	}
	if !strings.EqualFold(listFormat, "text") {
		return fmt.Errorf("formato inválido: %s (usa text o json)", listFormat)
	}
	printServerItemsText(items)
	return nil
}

func runServerSearch(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("no hay configuración: %v", err)
	}
	items := collectServerItems(cfg, args[0], nil)
	if len(items) == 0 {
		fmt.Println("Sin resultados")
		return nil
	}
	printServerItemsText(items)
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
	if len(server.Tags) > 0 {
		fmt.Printf("Tags: %s\n", strings.Join(server.Tags, ", "))
	} else {
		fmt.Println("Tags: (sin tags)")
	}
	if infoShowPass {
		fmt.Printf("Contraseña: %s\n", server.Password)
	}
	if name == cfg.ActiveServer {
		fmt.Println("Estado: ACTIVO")
	}
	return nil
}

func runServerTest(cmd *cobra.Command, args []string) error {
	server, err := resolveServerArg(args)
	if err != nil {
		return err
	}
	if err := validateServerConnection(server); err != nil {
		return err
	}
	fmt.Printf("Servidor '%s': OK\n", server.Name)
	return nil
}

func runServerPing(cmd *cobra.Command, args []string) error {
	server, err := resolveServerArg(args)
	if err != nil {
		return err
	}
	started := time.Now()
	client, err := ssh.Connect(server.Host, server.Port, server.User, server.Password)
	if err != nil {
		return fmt.Errorf("error de conexión: %v", err)
	}
	defer client.Close()
	latency := time.Since(started)
	fmt.Printf("Servidor '%s': %d ms\n", server.Name, latency.Milliseconds())
	return nil
}

func runServerDoctor(cmd *cobra.Command, args []string) error {
	server, err := resolveServerArg(args)
	if err != nil {
		return err
	}

	fmt.Printf("Diagnóstico de %s (%s@%s:%d)\n", server.Name, server.User, server.Host, server.Port)
	start := time.Now()
	client, err := ssh.Connect(server.Host, server.Port, server.User, server.Password)
	if err != nil {
		return fmt.Errorf("error de conexión: %v", err)
	}
	defer client.Close()
	latency := time.Since(start)

	checks := []doctorCheck{{Name: "ssh", OK: true, Detail: fmt.Sprintf("conectado en %d ms", latency.Milliseconds())}}
	checks = append(checks, remoteCommandCheck(client, "sudo", "command -v sudo >/dev/null 2>&1 && echo ok || echo missing")...)
	checks = append(checks, remoteCommandCheck(client, "git", "command -v git >/dev/null 2>&1 && echo ok || echo missing")...)
	checks = append(checks, remoteCommandCheck(client, "docker", "command -v docker >/dev/null 2>&1 && echo ok || echo missing")...)
	checks = append(checks, remoteCommandCheck(client, "php", "command -v php >/dev/null 2>&1 && php -v | head -1 || echo missing")...)
	checks = append(checks, remoteCommandCheck(client, "composer", "command -v composer >/dev/null 2>&1 && composer --version | head -1 || echo missing")...)
	checks = append(checks, remoteCommandCheck(client, "caddy", "command -v caddy >/dev/null 2>&1 && caddy version || echo missing")...)
	checks = append(checks, remoteCommandCheck(client, "uname", "uname -a")...)

	for _, c := range checks {
		state := "OK"
		if !c.OK {
			state = "MISSING"
		}
		if c.Name == "uname" {
			state = "INFO"
		}
		fmt.Printf("- %-8s %s", c.Name, state)
		if c.Detail != "" {
			fmt.Printf(" - %s", c.Detail)
		}
		fmt.Println()
	}
	return nil
}

func runServerSyncTags(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("no hay configuración: %v", err)
	}
	origin, err := cfg.GetServer(args[0])
	if err != nil {
		return err
	}
	dest, err := cfg.GetServer(args[1])
	if err != nil {
		return err
	}
	dest.Tags = append([]string(nil), origin.Tags...)
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("error al guardar configuración: %v", err)
	}
	fmt.Printf("Tags sincronizados: %s -> %s\n", args[0], args[1])
	return nil
}

func runServerExport(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("no hay configuración: %v", err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("error al serializar exportación: %v", err)
	}
	if err := os.WriteFile(args[0], data, 0600); err != nil {
		return fmt.Errorf("error al escribir archivo de exportación: %v", err)
	}
	fmt.Printf("Configuración exportada a: %s\n", args[0])
	return nil
}

func runServerImport(cmd *cobra.Command, args []string) error {
	imported, err := loadConfigFromFile(args[0])
	if err != nil {
		return err
	}

	if importReplace {
		if imported.ActiveServer == "" {
			for _, name := range imported.ListServers() {
				imported.ActiveServer = name
				break
			}
		}
		if err := config.Save(imported); err != nil {
			return fmt.Errorf("error al guardar configuración importada: %v", err)
		}
		fmt.Printf("Configuración importada desde: %s\n", args[0])
		return nil
	}

	cfg := config.LoadOrCreate()
	for _, name := range imported.ListServers() {
		server := imported.Servers[name]
		if err := cfg.AddServer(server, importForce); err != nil {
			return fmt.Errorf("error al importar '%s': %v", name, err)
		}
	}
	if cfg.ActiveServer == "" && imported.ActiveServer != "" {
		if _, exists := cfg.Servers[imported.ActiveServer]; exists {
			cfg.ActiveServer = imported.ActiveServer
		}
	}
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("error al guardar configuración importada: %v", err)
	}
	fmt.Printf("Configuración importada desde: %s\n", args[0])
	return nil
}

func runServerSetPass(cmd *cobra.Command, args []string) error {
	name := args[0]
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("no hay configuración: %v", err)
	}
	server, err := cfg.GetServer(name)
	if err != nil {
		return err
	}
	server.Password = setPassValue
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("error al guardar configuración: %v", err)
	}
	fmt.Printf("Contraseña actualizada para el servidor '%s'\n", name)
	return nil
}

func validateServerConnection(server *config.Server) error {
	fmt.Printf("Validando conexión a %s@%s:%d...\n", server.User, server.Host, server.Port)
	client, err := ssh.Connect(server.Host, server.Port, server.User, server.Password)
	if err != nil {
		return fmt.Errorf("error de conexión: %v", err)
	}
	defer client.Close()

	output, err := client.Run("echo 'Conexión exitosa'")
	if err != nil {
		return fmt.Errorf("error al validar conexión: %v", err)
	}
	fmt.Print(output)
	return nil
}

func loadConfigFromFile(path string) (*config.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error al leer archivo de importación: %v", err)
	}
	var cfg config.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error al parsear archivo de importación: %v", err)
	}
	if cfg.Servers == nil {
		cfg.Servers = make(map[string]*config.Server)
	}
	for name, server := range cfg.Servers {
		if server == nil {
			delete(cfg.Servers, name)
			continue
		}
		server.Name = name
		cfg.Servers[name] = server.Clone()
	}
	return &cfg, nil
}

func matchesAllTags(serverTags, requiredTags []string) bool {
	if len(requiredTags) == 0 {
		return true
	}
	have := map[string]struct{}{}
	for _, tag := range serverTags {
		have[strings.TrimSpace(tag)] = struct{}{}
	}
	for _, tag := range requiredTags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		if _, ok := have[tag]; !ok {
			return false
		}
	}
	return true
}

func resolveServerArg(args []string) (*config.Server, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("no hay configuración: %v", err)
	}
	if len(args) > 0 {
		return cfg.GetServer(args[0])
	}
	return cfg.GetActiveServer()
}

func collectServerItems(cfg *config.Config, query string, requiredTags []string) []serverListItem {
	query = strings.ToLower(strings.TrimSpace(query))
	items := make([]serverListItem, 0)
	for _, name := range cfg.ListServers() {
		server := cfg.Servers[name]
		if !matchesAllTags(server.Tags, requiredTags) {
			continue
		}
		if query != "" && !serverMatchesQuery(name, server, query) {
			continue
		}
		items = append(items, serverListItem{
			Name:   name,
			Host:   server.Host,
			Port:   server.Port,
			User:   server.User,
			Tags:   append([]string(nil), server.Tags...),
			Active: name == cfg.ActiveServer,
		})
	}
	return items
}

func serverMatchesQuery(name string, server *config.Server, query string) bool {
	if strings.Contains(strings.ToLower(name), query) ||
		strings.Contains(strings.ToLower(server.Host), query) ||
		strings.Contains(strings.ToLower(server.User), query) {
		return true
	}
	for _, tag := range server.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}
	return false
}

func printServerItemsText(items []serverListItem) {
	if len(items) == 0 {
		fmt.Println("Sin resultados")
		return
	}
	fmt.Println("Servidores configurados:")
	fmt.Println()
	for _, item := range items {
		marker := "  "
		if item.Active {
			marker = "* "
		}
		tagSuffix := ""
		if len(item.Tags) > 0 {
			tags := append([]string(nil), item.Tags...)
			sort.Strings(tags)
			tagSuffix = fmt.Sprintf(" [%s]", strings.Join(tags, ", "))
		}
		fmt.Printf("%s%-15s %s@%s:%d%s\n", marker, item.Name, item.User, item.Host, item.Port, tagSuffix)
	}
	fmt.Println()
	fmt.Println("(*) = servidor activo")
}

func remoteCommandCheck(client *ssh.Client, name, command string) []doctorCheck {
	output, err := client.Run(command)
	if err != nil {
		return []doctorCheck{{Name: name, OK: false, Detail: err.Error()}}
	}
	output = strings.TrimSpace(output)
	if output == "missing" {
		return []doctorCheck{{Name: name, OK: false, Detail: "no instalado"}}
	}
	return []doctorCheck{{Name: name, OK: true, Detail: output}}
}
