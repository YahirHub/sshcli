package ssh

import (
	"bytes"
	"fmt"
	"io"
	"os"
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
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
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

// Close cierra todas las conexiones
func (c *Client) Close() error {
	if c.sftpClient != nil {
		c.sftpClient.Close()
	}
	if c.sshClient != nil {
		return c.sshClient.Close()
	}
	return nil
}

// Run ejecuta un comando en el servidor remoto y retorna la salida
func (c *Client) Run(command string) (string, error) {
	session, err := c.sshClient.NewSession()
	if err != nil {
		return "", fmt.Errorf("fallo al crear sesión: %v", err)
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(command)
	if err != nil {
		if stderr.Len() > 0 {
			return "", fmt.Errorf("%s: %v", stderr.String(), err)
		}
		return "", err
	}

	return stdout.String(), nil
}

// WriteFile escribe contenido a un archivo remoto
func (c *Client) WriteFile(remotePath string, data []byte, perm os.FileMode) error {
	file, err := c.sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("fallo al crear archivo remoto: %v", err)
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("fallo al escribir archivo: %v", err)
	}

	if err := c.sftpClient.Chmod(remotePath, perm); err != nil {
		// No es crítico si falla chmod
		return nil
	}

	return nil
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

// ListDir lista el contenido de un directorio remoto
func (c *Client) ListDir(remotePath string) ([]os.FileInfo, error) {
	return c.sftpClient.ReadDir(remotePath)
}
