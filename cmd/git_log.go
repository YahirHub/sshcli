package cmd

import (
	"fmt"
	"sshcli/internal/paths"

	"github.com/spf13/cobra"
)

var (
	gitLogCount   int
	gitLogOneline bool
	gitLogServer  string
)

var gitLogCmd = &cobra.Command{
	Use:   "git-log [directorio]",
	Short: "Historial de commits",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runGitLog,
}

func init() {
	rootCmd.AddCommand(gitLogCmd)
	gitLogCmd.Flags().IntVarP(&gitLogCount, "number", "n", 10, "Número de commits")
	gitLogCmd.Flags().BoolVar(&gitLogOneline, "oneline", false, "Formato compacto")
	gitLogCmd.Flags().StringVarP(&gitLogServer, "server", "s", "", "Servidor específico")
}

func runGitLog(cmd *cobra.Command, args []string) error {
	dir := "/"
	if len(args) > 0 {
		dir = paths.ToRemote(args[0])
	}

	client, _, err := getClient(gitLogServer)
	if err != nil {
		return err
	}
	defer client.Close()

	format := ""
	if gitLogOneline { format = "--oneline" }
	
	gitCmd := fmt.Sprintf("cd '%s' && git log -n %d %s", dir, gitLogCount, format)
	output, err := client.Run(gitCmd)
	if err != nil {
		return fmt.Errorf("error al obtener log: %v", err)
	}

	fmt.Print(output)
	return nil
}