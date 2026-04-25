package ssh

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Client representa una conexión SSH con capacidades SFTP
type Client struct {
	sshClient  *ssh.Client
	sftpClient *sftp.Client
}

// Connect establece una conexión SSH al servidor remoto
func Connect(host string, port int, user, password string) (*Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth:[]ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         20 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("fallo al conectar a %s: %v", addr, err)
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		sshClient.Close()
		return nil, fmt.Errorf("fallo al iniciar SFTP: %v", err)
	}

	return &Client{
		sshClient:  sshClient,
		sftpClient: sftpClient,
	}, nil
}

// Close cierra todas las conexiones de forma segura
func (c *Client) Close() error {
	if c.sftpClient != nil {
		c.sftpClient.Close()
	}
	if c.sshClient != nil {
		return c.sshClient.Close()
	}
	return nil
}

// NewSession crea una nueva sesión SSH raw para uso interactivo
func (c *Client) NewSession() (*ssh.Session, error) {
	return c.sshClient.NewSession()
}

// Run ejecuta un comando en el servidor remoto y retorna la salida
func (c *Client) Run(command string) (string, error) {
	session, err := c.sshClient.NewSession()
	if err != nil {
		return "", fmt.Errorf("fallo al crear sesion: %v", err)
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(command)
	
	output := stdout.String()
	if err != nil {
		errStr := strings.TrimSpace(stderr.String())
		if errStr != "" {
			return output, fmt.Errorf("%s", errStr)
		}
		return output, err
	}

	return output, nil
}

// WriteFile escribe contenido a un archivo remoto
func (c *Client) WriteFile(remotePath string, data[]byte, perm os.FileMode) error {
	// SFTP Create falla si el directorio no existe.
	file, err := c.sftpClient.OpenFile(remotePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return fmt.Errorf("fallo al abrir/crear archivo remoto (existe la carpeta?): %v", err)
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("fallo al escribir datos en SFTP: %v", err)
	}

	return c.sftpClient.Chmod(remotePath, perm)
}

// ReadFile lee el contenido de un archivo remoto
func (c *Client) ReadFile(remotePath string) ([]byte, error) {
	file, err := c.sftpClient.Open(remotePath)
	if err != nil {
		return nil, fmt.Errorf("fallo al abrir archivo remoto: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("fallo al leer archivo: %v", err)
	}

	return data, nil
}

// FileExists verifica si un archivo existe en el servidor remoto
func (c *Client) FileExists(remotePath string) bool {
	_, err := c.sftpClient.Stat(remotePath)
	return err == nil
}

// IsDir verifica si una ruta remota es un directorio
func (c *Client) IsDir(remotePath string) bool {
	info, err := c.sftpClient.Stat(remotePath)
	if err != nil {
		return false
	}
	return info.IsDir()
}