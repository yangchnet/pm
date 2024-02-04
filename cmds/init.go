package cmds

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

func InitCmd() *cobra.Command {
	var (
		gitUrl         string
		gitUserName    string
		gitEmail       string
		privateKeyPath string
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

			sshDir := filepath.Join(home, ".ssh")
			pmDir := filepath.Join(home, ".pm")
			storeDir := filepath.Join(pmDir, "store")
			userKeyPath := filepath.Join(pmDir, "user.key")
			scriptPath := filepath.Join(pmDir, "completion.sh")

			// 1. 设置密钥
			if privateKeyPath == "" {
				if _, err := os.Stat(sshDir); os.IsNotExist(err) {
					os.MkdirAll(sshDir, 0755)
				}
				privateKeyPath = filepath.Join(sshDir, "pm")
				sshCmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-f", privateKeyPath, "-N", "")
				sshCmd.Run()
			}

			// 2. 创建.pm主目录
			if _, err := os.Stat(pmDir); os.IsNotExist(err) {
				os.MkdirAll(pmDir, 0755)
			}

			// 3. 创建存储目录
			if _, err := os.Stat(storeDir); os.IsNotExist(err) {
				os.MkdirAll(storeDir, 0755)
			}

			// 4. 在~/.pm中创建配置文件conf.yaml
			confPath := filepath.Join(pmDir, "conf.yaml")
			confFile, err := os.Create(confPath)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			defer confFile.Close()

			params := map[string]string{
				"gitUserName":    gitUserName,
				"gitEmail":       gitEmail,
				"gitUrl":         gitUrl,
				"privateKeyPath": privateKeyPath,
				"localPath":      filepath.Join(home, ".pm/store"),
				"userKeyPath":    userKeyPath,
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

			confFile.Write(buf.Bytes())

			writeKnowHosts(gitUrl)

			writeAutoCompletionScript(scriptPath)
		},
	}

	initCmd.Flags().StringVarP(&gitUrl, "git-url", "", "", "git url")
	initCmd.MarkFlagRequired("git-url")

	initCmd.Flags().StringVarP(&gitUserName, "git-username", "", "", "git user name")
	initCmd.MarkFlagRequired("git-username")

	initCmd.Flags().StringVarP(&gitEmail, "git-email", "", "", "git email")
	initCmd.MarkFlagRequired("git-email")

	initCmd.Flags().StringVarP(&privateKeyPath, "private-key-path", "", "", "private key path")

	return initCmd
}

var confTmpl = `remote:
  git:
    name: {{.gitUserName}}
    email: {{.gitEmail}}
    url: {{.gitUrl}}
    privateKeyPath: {{.privateKeyPath}}

local:
  path: {{.localPath}}

userKeyPath: {{.userKeyPath}}
`

func writeAutoCompletionScript(path string) {
	shell := filepath.Base(os.Getenv("SHELL"))
	switch shell {
	case "bash":
		profile := filepath.Join(os.Getenv("HOME"), ".bashrc")
		if err := RootCmd.GenBashCompletionFileV2(path, true); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		appendLineToFile(profile, fmt.Sprintf(". %s", path))

	case "zsh":
		profile := filepath.Join(os.Getenv("HOME"), ".zshrc")
		if err := RootCmd.GenZshCompletionFile(path); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		appendLineToFile(profile, fmt.Sprintf(". %s", path))

	case "fish":
		profile := filepath.Join(os.Getenv("HOME"), ".config", "fish", "config.fish")
		if err := RootCmd.GenFishCompletionFile(path, true); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		appendLineToFile(profile, fmt.Sprintf(". %s", path))

	// case "powershell":
	// 	profile := filepath.Join(os.Getenv("HOME"), "Documents", "WindowsPowerShell", "Microsoft.PowerShell_profile.ps1")
	default:
		fmt.Println("不支持的shell:", shell)
	}
}

func appendLineToFile(path, line string) error {
	// Open the file in append mode
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the line to the file
	_, err = fmt.Fprintln(file, line)
	return err
}

func writeKnowHosts(gitUrl string) {
	if strings.HasPrefix(gitUrl, "https") {
		fmt.Println("请使用git@开头的git地址")
		os.Exit(1)
	}

	gitServer, _, _, err := extractGitInfo(gitUrl)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	knownHostsFile := filepath.Join(home, ".ssh", "known_hosts")
	cmd := exec.Command("sh", "-c", fmt.Sprintf("ssh-keyscan -t rsa %s >> %s", gitServer, knownHostsFile))
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
