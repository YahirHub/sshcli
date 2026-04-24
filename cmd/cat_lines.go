package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"sshcli/internal/paths"

	"github.com/spf13/cobra"
)

var (
	catLinesNumbers bool
	catLinesServer  string
)

var catLinesCmd = &cobra.Command{
	Use:   "cat-lines [archivo] [inicio] [fin]",
	Short: "Lee rango de líneas",
	Args:  cobra.ExactArgs(3),
	RunE:  runCatLines,
}

func init() {
	rootCmd.AddCommand(catLinesCmd)
	catLinesCmd.Flags().BoolVarP(&catLinesNumbers, "numbers", "n", false, "Números de línea")
	catLinesCmd.Flags().StringVarP(&catLinesServer, "server", "s", "", "Servidor específico")
}

func runCatLines(cmd *cobra.Command, args []string) error {
	remotePath := paths.ToRemote(args[0])
	start, _ := strconv.Atoi(args[1])
	end, _ := strconv.Atoi(args[2])

	client, _, err := getClient(catLinesServer)
	if err != nil {
		return err
	}
	defer client.Close()

	data, err := client.ReadFile(remotePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	for i := start - 1; i < end && i < len(lines); i++ {
		if catLinesNumbers {
			fmt.Printf("%4d | %s\n", i+1, lines[i])
		} else {
			fmt.Println(lines[i])
		}
	}
	return nil
}