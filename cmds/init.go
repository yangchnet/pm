package cmds

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/yangchnet/pm/utils"
)

func InitCmd() *cobra.Command {
	var (
		storeType  string
		remoteType string
	)

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "init config",
		Run: func(cmd *cobra.Command, args []string) {
			home, err := os.UserHomeDir()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			pmDir := filepath.Join(home, ".pm")
			storeDir := filepath.Join(pmDir, "store")
			userKeyPath := filepath.Join(pmDir, "user.key")
			scriptPath := filepath.Join(pmDir, "completion")

			utils.CreateDirIfNotExist(pmDir)    // 创建.pm主目录
			utils.CreateDirIfNotExist(storeDir) // 创建存储目录

			remote, err := NewRemote(cmd.Context(), remoteType, nil)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			// remote 初始化
			remoteConfigStr, err := remote.Init(cmd.Context())
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			localStore, err := NewStore(cmd.Context(), storeType, nil)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			storeConfigStr, err := localStore.Init(cmd.Context())
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			// 写入配置文件
			params := map[string]string{
				"remoteConfig": remoteConfigStr,
				"storeConfig":  storeConfigStr,
				"localPath":    filepath.Join(home, ".pm/store"),
				"userKeyPath":  userKeyPath,
				"latency":      "24h",
			}
			tmpl, err := template.New("conf").Parse(confTmpl)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			buf := &bytes.Buffer{}
			if err := tmpl.Execute(buf, params); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			confPath := filepath.Join(pmDir, "conf.yaml")
			confFile, err := os.Create(confPath)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			defer confFile.Close()

			confFile.Write(buf.Bytes())

			writeAutoCompletionScript(scriptPath)
		},
	}

	initCmd.Flags().StringVarP(&storeType, "store", "s", "file", "store type: [file, sqlite], default file")

	initCmd.Flags().StringVarP(&remoteType, "remote", "r", "empty", "remote type: [git, empty], default empty")

	return initCmd
}

var confTmpl = `
{{.remoteConfig}}

{{.storeConfig}}

local:
  path: {{.localPath}}

user_key_path: {{.userKeyPath}}
`

func writeAutoCompletionScript(scriptPath string) {
	if runtime.GOOS == "windows" {
		windowsWriteAutoCompleteScript(scriptPath + ".ps1")
		return
	}
	scriptPath = scriptPath + ".sh"

	shell := filepath.Base(os.Getenv("SHELL"))
	switch shell {
	case "bash":
		profile := filepath.Join(os.Getenv("HOME"), ".bashrc")
		if err := RootCmd.GenBashCompletionFileV2(scriptPath, true); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		utils.AppendLineToFile(profile, fmt.Sprintf(". %s", scriptPath))

	case "zsh":
		profile := filepath.Join(os.Getenv("HOME"), ".zshrc")
		if err := RootCmd.GenZshCompletionFile(scriptPath); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		utils.AppendLineToFile(profile, fmt.Sprintf(". %s", scriptPath))

	case "fish":
		profile := filepath.Join(os.Getenv("HOME"), ".config", "fish", "config.fish")
		if err := RootCmd.GenFishCompletionFile(scriptPath, true); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		utils.AppendLineToFile(profile, fmt.Sprintf(". %s", scriptPath))
	default:
		fmt.Println("不支持的shell:", shell)
	}
}

func windowsWriteAutoCompleteScript(path string) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	profile := filepath.Join(home, "Documents", "WindowsPowerShell", "Microsoft.PowerShell_profile.ps1")
	utils.CreateDirIfNotExist(filepath.Dir(profile))
	if err := RootCmd.GenPowerShellCompletionFile(path); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	utils.AppendLineToFile(profile, fmt.Sprintf(". %s", path))
}
