package cmd

import (
	"sshcli/internal/config"
	"sshcli/internal/ssh"
)

// Variable global para especificar servidor en cualquier comando
var targetServer string

// getClient obtiene un cliente SSH para el servidor especificado o el activo
func getClient(serverName string) (*ssh.Client, *config.Server, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, err
	}

	var server *config.Server
	if serverName != "" {
		server, err = cfg.GetServer(serverName)
	} else {
		server, err = cfg.GetActiveServer()
	}
	if err != nil {
		return nil, nil, err
	}

	client, err := ssh.Connect(server.Host, server.Port, server.User, server.Password)
	if err != nil {
		return nil, nil, err
	}

	return client, server, nil
}
