package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const configFileName = ".sshcli.conf"

// Server representa un servidor SSH configurado
type Server struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// Config representa la configuración completa con múltiples servidores
type Config struct {
	ActiveServer string             `json:"active_server"`
	Servers      map[string]*Server `json:"servers"`
	DefaultTTY   bool               `json:"default_tty"`
}

// getConfigPath devuelve la ruta completa del archivo de configuración
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener directorio home: %v", err)
	}
	return filepath.Join(home, configFileName), nil
}

// NewConfig crea una nueva configuración vacía
func NewConfig() *Config {
	return &Config{
		Servers: make(map[string]*Server),
	}
}

// Save guarda la configuración en el archivo con permisos restrictivos
func Save(cfg *Config) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("error al serializar configuración: %v", err)
	}

	// 0600 asegura que solo el usuario actual pueda leer/escribir las contraseñas
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("error al escribir configuración: %v", err)
	}

	return nil
}

// Load carga la configuración desde el archivo
func Load() (*Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("archivo de configuración no encontrado")
		}
		return nil, fmt.Errorf("error al leer configuración: %v", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error al parsear configuración: %v", err)
	}

	if cfg.Servers == nil {
		cfg.Servers = make(map[string]*Server)
	}

	return &cfg, nil
}

// LoadOrCreate carga la configuración o crea una nueva si no existe
func LoadOrCreate() *Config {
	cfg, err := Load()
	if err != nil {
		return NewConfig()
	}
	return cfg
}

// AddServer agrega un servidor a la configuración
func (c *Config) AddServer(server *Server) {
	c.Servers[server.Name] = server
	if c.ActiveServer == "" {
		c.ActiveServer = server.Name
	}
}

// RemoveServer elimina un servidor de la configuración
func (c *Config) RemoveServer(name string) error {
	if _, exists := c.Servers[name]; !exists {
		return fmt.Errorf("servidor '%s' no encontrado", name)
	}
	delete(c.Servers, name)
	
	if c.ActiveServer == name {
		c.ActiveServer = ""
		for n := range c.Servers {
			c.ActiveServer = n
			break
		}
	}
	return nil
}

// GetServer obtiene un servidor por nombre
func (c *Config) GetServer(name string) (*Server, error) {
	server, exists := c.Servers[name]
	if !exists {
		return nil, fmt.Errorf("servidor '%s' no encontrado", name)
	}
	return server, nil
}

// GetActiveServer obtiene el servidor activo
func (c *Config) GetActiveServer() (*Server, error) {
	if c.ActiveServer == "" {
		return nil, fmt.Errorf("no hay servidor activo configurado")
	}
	return c.GetServer(c.ActiveServer)
}

// SetActiveServer establece el servidor activo
func (c *Config) SetActiveServer(name string) error {
	if _, exists := c.Servers[name]; !exists {
		return fmt.Errorf("servidor '%s' no encontrado", name)
	}
	c.ActiveServer = name
	return nil
}

// ListServers devuelve lista ordenada de nombres de servidores
func (c *Config) ListServers() []string {
	names := make([]string, 0, len(c.Servers))
	for name := range c.Servers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Exists verifica si existe el archivo de configuración
func Exists() bool {
	path, err := getConfigPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}

// Delete elimina el archivo de configuración
func Delete() error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}
	return os.Remove(path)
}