package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	dockerServer string
	dockerTail   int
)

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Gestiona contenedores Docker en el servidor",
}

var dockerPsCmd = &cobra.Command{
	Use:   "ps",
	Short: "Lista contenedores activos",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := getClient(dockerServer)
		if err != nil {
			return err
		}
		defer client.Close()
		// Formato optimizado para agentes
		output, err := client.Run("docker ps --format 'table {{.ID}}\t{{.Names}}\t{{.Status}}\t{{.Ports}}'")
		if err != nil {
			return fmt.Errorf("error: %v\n¿Está Docker instalado?", err)
		}
		fmt.Print(output)
		return nil
	},
}

var dockerLogsCmd = &cobra.Command{
	Use:   "logs [container_id_o_nombre]",
	Short: "Muestra logs de un contenedor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := getClient(dockerServer)
		if err != nil {
			return err
		}
		defer client.Close()
		dockerCommand := fmt.Sprintf("docker logs --tail %d %s", dockerTail, args[0])
		output, err := client.Run(dockerCommand)
		if err != nil {
			return err
		}
		fmt.Print(output)
		return nil
	},
}

var dockerExecCmd = &cobra.Command{
	Use:   "exec [container_id] [comando]",
	Short: "Ejecuta un comando dentro de un contenedor",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := getClient(dockerServer)
		if err != nil {
			return err
		}
		defer client.Close()
		dockerCommand := fmt.Sprintf("docker exec %s %s", args[0], args[1])
		output, err := client.Run(dockerCommand)
		if err != nil {
			return err
		}
		fmt.Print(output)
		return nil
	},
}

var dockerStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Muestra consumo de recursos de contenedores",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := getClient(dockerServer)
		if err != nil {
			return err
		}
		defer client.Close()
		output, err := client.Run("docker stats --no-stream --format 'table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}'")
		if err != nil {
			return err
		}
		fmt.Print(output)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dockerCmd)
	dockerCmd.AddCommand(dockerPsCmd)
	dockerCmd.AddCommand(dockerLogsCmd)
	dockerCmd.AddCommand(dockerExecCmd)
	dockerCmd.AddCommand(dockerStatsCmd)

	dockerCmd.PersistentFlags().StringVarP(&dockerServer, "server", "s", "", "Servidor específico")
	dockerLogsCmd.Flags().IntVarP(&dockerTail, "tail", "n", 50, "Número de líneas de logs")
}