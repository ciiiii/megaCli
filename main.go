package main

import (
	"github.com/manifoldco/promptui"
	"github.com/t3rm1n4l/go-mega"
	"github.com/spf13/cobra"
	"fmt"
	"os"
)

const prev = "-parent-"

type Item struct {
	Name string
	Type int
	Size string
	Node *mega.Node
}

type Config struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}

var (
	chooseTemplate = &promptui.SelectTemplates{
		Label:    " 🗂 {{ . | white}}",
		Active:   "-> {{ .Name | cyan }} ({{ .Size | red }}) [{{ .Type | green}}]",
		Inactive: "   {{ .Name | red }} ({{ .Size | red }}) [{{ .Type | green}}]",
	}
	m    = mega.New()
	root *mega.Node
)

func main() {
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Init/reset your account info",
		Long:  "init/reset your account info in $HOME/.config/megaCli/mega.toml",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			initConf(setConf())
		},
	}

	var uploadCmd = &cobra.Command{
		Use:   "upload <remotePath> <localPath>",
		Short: "Upload file from local",
		Long:  "upload 'localPath' to 'remotePath' on mega.nz cloud",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			auth()
			fileName := args[0]
			filePath := args[1]
			fmt.Printf("upload '%s' to '%s'!\n", fileName, filePath)
			f, err := os.Stat(filePath)
			if err != nil {
				fmt.Printf("get '%s' size failed!\n", filePath)
				return
			}
			uploadFile(fileName, filePath, f.Size())
		},
	}

	var discoverCmd = &cobra.Command{
		Use:   "discover",
		Short: "List/download/delete file in mega.nz",
		Long:  "List or download or delete file in mega.nz",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			auth()
			for ; ; {
				selectedFile := choose(nil, getChildren(root))
				if selectedFile != nil {
					operate(selectedFile)
				}
			}
		},
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version of megaCli",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("megaCli v0.1.0 -- HEAD")
		},
	}

	var rootCmd = &cobra.Command{Use: "megaCli"}
	rootCmd.AddCommand(initCmd, uploadCmd, discoverCmd, versionCmd)
	rootCmd.Execute()
}
