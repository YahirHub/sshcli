package cmd

import (
	"fmt"
	"sort"
	"strings"

	"sshcli/internal/paths"

	"github.com/spf13/cobra"
)

var (
	treeDepth  int
	treeDirs   bool
	treeServer string
)

type treeNode struct {
	children map[string]*treeNode
}

var treeCmd = &cobra.Command{
	Use:   "tree [directorio]",
	Short: "Muestra la estructura de directorios en forma de árbol",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runTree,
}

func init() {
	rootCmd.AddCommand(treeCmd)
	treeCmd.Flags().IntVarP(&treeDepth, "depth", "d", 3, "Profundidad máxima")
	treeCmd.Flags().BoolVar(&treeDirs, "dirs", false, "Solo directorios")
	treeCmd.Flags().StringVarP(&treeServer, "server", "s", "", "Servidor específico a usar")
}

func runTree(cmd *cobra.Command, args []string) error {
	remotePath := "/"
	if len(args) > 0 {
		remotePath = paths.ToRemote(args[0])
	}

	client, _, err := getClient(treeServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	var nativeTreeCmd string
	if treeDirs {
		nativeTreeCmd = fmt.Sprintf("tree -d -L %d '%s' 2>/dev/null", treeDepth, remotePath)
	} else {
		nativeTreeCmd = fmt.Sprintf("tree -L %d '%s' 2>/dev/null", treeDepth, remotePath)
	}

	output, err := client.Run(nativeTreeCmd)
	if err == nil && strings.TrimSpace(output) != "" {
		fmt.Print(output)
		return nil
	}

	findCmd := fmt.Sprintf("find '%s' -maxdepth %d", remotePath, treeDepth)
	if treeDirs {
		findCmd += " -type d"
	}
	findCmd += " 2>/dev/null | sort | head -200"

	output, err = client.Run(findCmd)
	if err != nil {
		return fmt.Errorf("error al obtener estructura: %v", err)
	}

	lines := []string{}
	for _, line := range strings.Split(strings.ReplaceAll(output, "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}

	if len(lines) == 0 {
		fmt.Println(remotePath)
		return nil
	}

	fmt.Print(renderTreeFromPaths(remotePath, lines))
	return nil
}

func renderTreeFromPaths(rootPath string, fullPaths []string) string {
	root := &treeNode{children: map[string]*treeNode{}}
	for _, full := range fullPaths {
		addTreePath(root, rootPath, full)
	}

	var b strings.Builder
	b.WriteString(rootPath)
	b.WriteString("\n")
	renderTreeChildren(&b, root, "")
	return b.String()
}

func addTreePath(root *treeNode, rootPath, fullPath string) {
	if fullPath == rootPath {
		return
	}

	rel := strings.TrimPrefix(fullPath, rootPath)
	rel = strings.TrimPrefix(rel, "/")
	if rel == "" {
		return
	}

	parts := strings.Split(rel, "/")
	cur := root
	for _, part := range parts {
		if part == "" {
			continue
		}
		if cur.children == nil {
			cur.children = map[string]*treeNode{}
		}
		next, ok := cur.children[part]
		if !ok {
			next = &treeNode{children: map[string]*treeNode{}}
			cur.children[part] = next
		}
		cur = next
	}
}

func renderTreeChildren(b *strings.Builder, node *treeNode, prefix string) {
	if node == nil || len(node.children) == 0 {
		return
	}

	keys := make([]string, 0, len(node.children))
	for k := range node.children {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		last := i == len(keys)-1
		branch := "├── "
		nextPrefix := prefix + "│   "
		if last {
			branch = "└── "
			nextPrefix = prefix + "    "
		}

		b.WriteString(prefix)
		b.WriteString(branch)
		b.WriteString(k)
		b.WriteString("\n")
		renderTreeChildren(b, node.children[k], nextPrefix)
	}
}
