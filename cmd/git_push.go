package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var gitPushServer string

var gitPushCmd = &cobra.Command{
	Use:   "git-push [directorio]",
	Short: "Envía commits al repositorio remoto",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runGitPush,
}

func init() {
	rootCmd.AddCommand(gitPushCmd)
	gitPushCmd.Flags().StringVarP(&gitPushServer, "server", "s", "", "Servidor específico a usar")
}

func runGitPush(cmd *cobra.Command, args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = cleanRemotePath(args[0])
	}

	client, _, err := getClient(gitPushServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	gitCmd := fmt.Sprintf("cd '%s' && git push", dir)
	output, err := client.Run(gitCmd)
	if err != nil {
		return fmt.Errorf("error al hacer push: %v\n%s", err, output)
	}

	fmt.Print(output)
	return nil
}