package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const configFileName = ".sshcli.conf"

// Server representa un servidor SSH configurado
type Server struct {
	Name     string   `json:"name"`
	Host     string   `json:"host"`
	Port     int      `json:"port"`
	User     string   `json:"user"`
	Password string   `json:"password"`
	Tags     []string `json:"tags,omitempty"`
}

// Clone crea una copia profunda del servidor
func (s *Server) Clone() *Server {
	if s == nil {
		return nil
	}
	clone := *s
	clone.Tags = append([]string(nil), s.Tags...)
	clone.Tags = normalizeTags(clone.Tags)
	return &clone
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
	for name, server := range cfg.Servers {
		if server == nil {
			delete(cfg.Servers, name)
			continue
		}
		server.Name = name
		server.Tags = normalizeTags(server.Tags)
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
func (c *Config) AddServer(server *Server, force bool) error {
	if server == nil {
		return fmt.Errorf("servidor inválido")
	}
	if server.Name == "" {
		return fmt.Errorf("el nombre del servidor no puede estar vacío")
	}
	if _, exists := c.Servers[server.Name]; exists && !force {
		return fmt.Errorf("ya existe un servidor con el nombre '%s'", server.Name)
	}
	server.Tags = normalizeTags(server.Tags)
	c.Servers[server.Name] = server.Clone()
	if c.ActiveServer == "" {
		c.ActiveServer = server.Name
	}
	return nil
}

// RemoveServer elimina un servidor de la configuración
func (c *Config) RemoveServer(name string) error {
	if _, exists := c.Servers[name]; !exists {
		return fmt.Errorf("servidor '%s' no encontrado", name)
	}
	delete(c.Servers, name)

	if c.ActiveServer == name {
		c.ActiveServer = ""
		for _, n := range c.ListServers() {
			c.ActiveServer = n
			break
		}
	}
	return nil
}

// RenameServer renombra un servidor existente validando que el nuevo nombre no exista
func (c *Config) RenameServer(oldName, newName string) error {
	if oldName == "" || newName == "" {
		return fmt.Errorf("los nombres de servidor no pueden estar vacíos")
	}
	if oldName == newName {
		return fmt.Errorf("el nombre nuevo debe ser diferente al actual")
	}

	server, exists := c.Servers[oldName]
	if !exists {
		return fmt.Errorf("servidor '%s' no encontrado", oldName)
	}
	if _, exists := c.Servers[newName]; exists {
		return fmt.Errorf("ya existe un servidor con el nombre '%s'", newName)
	}

	delete(c.Servers, oldName)
	server.Name = newName
	c.Servers[newName] = server

	if c.ActiveServer == oldName {
		c.ActiveServer = newName
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

func normalizeTags(tags []string) []string {
	if len(tags) == 0 {
		return nil
	}
	seen := map[string]struct{}{}
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		result = append(result, tag)
	}
	sort.Strings(result)
	if len(result) == 0 {
		return nil
	}
	return result
}
