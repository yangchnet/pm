package gitremote

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/yangchnet/pm/config"
	"github.com/yangchnet/pm/utils"

	"github.com/go-git/go-git/v5/plumbing/object"
)

type GitRemote struct {
	LocalPath string
	*GitRemoteConfig
}

type GitRemoteConfig struct {
	Url            string `json:"url" yaml:"url"`
	Name           string `json:"name" yaml:"name"`
	Email          string `json:"email" yaml:"email"`
	PrivateKeyPath string `json:"privateKeyPath" yaml:"privateKeyPath"`
}

func NewGitRemote(ctx context.Context, gitConfigMap map[string]any) *GitRemote {
	if gitConfigMap == nil {
		return &GitRemote{}
	}

	return &GitRemote{
		LocalPath: config.GetString("local.path"),
		GitRemoteConfig: &GitRemoteConfig{
			Url:            gitConfigMap["url"].(string),
			Name:           gitConfigMap["name"].(string),
			Email:          gitConfigMap["email"].(string),
			PrivateKeyPath: gitConfigMap["private_key_path"].(string),
		},
	}

}

// Init 初始化remote
func (gr *GitRemote) Init(ctx context.Context) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home dir: %w", err)
	}

	var username, userEmail, gitUrl, privateKeyPath string
	fmt.Printf("Please input your git username: ")
	fmt.Scanln(&username)
	if username == "" {
		return "", fmt.Errorf("username can't be empty")
	}
	fmt.Printf("Please input your git email: ")
	fmt.Scanln(&userEmail)
	if userEmail == "" {
		return "", fmt.Errorf("user email can't be empty")
	}
	fmt.Printf("Please input your git url: ")
	fmt.Scanln(&gitUrl)
	if gitUrl == "" {
		return "", fmt.Errorf("git url can't be empty")
	}
	_, _, _, err = utils.ExtractGitInfo(gitUrl)
	if err != nil {
		return "", err
	}
	fmt.Printf("Please input your private key path(leave blank to generate a new key): ")
	fmt.Scanln(&privateKeyPath)

	sshDir := filepath.Join(home, ".ssh")
	if privateKeyPath == "" {
		utils.CreateDirIfNotExist(sshDir)

		privateKeyPath = filepath.Join(sshDir, "pm")
		sshCmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-f", privateKeyPath, "-N", "", "-C", "")
		sshCmd.Run()

		pubKey, err := os.ReadFile(filepath.Join(sshDir, "pm.pub"))
		if err != nil {
			return "", fmt.Errorf("failed to read public key: %w", err)
		}

		defer func() {
			fmt.Println("Please add the following public key to your git server:")
			fmt.Println(string(pubKey))
		}()

	}

	writeKnowHosts(gitUrl)

	params := map[string]string{
		"gitUserName":    username,
		"gitEmail":       userEmail,
		"gitUrl":         gitUrl,
		"privateKeyPath": privateKeyPath,
	}
	tmpl, err := template.New("conf").Parse(gitRemoteConfigTemplate)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, params); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return buf.String(), nil
}

var gitRemoteConfigTemplate = `remote:
  type: git
  name: {{.gitUserName}}
  email: {{.gitEmail}}
  url: {{.gitUrl}}
  private_key_path: {{.privateKeyPath}}
`

func writeKnowHosts(gitUrl string) {
	gitServer, _, _, err := utils.ExtractGitInfo(gitUrl)
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

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", fmt.Sprintf("ssh-keyscan -t rsa %s >> %s", gitServer, knownHostsFile))
	} else {
		cmd = exec.Command("sh", "-c", fmt.Sprintf("ssh-keyscan -t rsa %s >> %s", gitServer, knownHostsFile))
	}

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

// Push 将本地的store文件推送到远端
func (gr *GitRemote) Push(ctx context.Context, msg ...string) error {

	auth, err := auth(gr.PrivateKeyPath, "")
	if err != nil {
		return fmt.Errorf("failed to create git auth credentials: %w", err)
	}

	r, err := git.PlainOpen(gr.LocalPath)
	if err != nil {
		return err
	}

	if err := gr.commit(ctx, r); err != nil {
		return err
	}

	return r.Push(&git.PushOptions{
		Auth: auth,
	})
}

// Pull 将远端的store文件拉取到本地
func (gr *GitRemote) Pull(ctx context.Context) error {
	utils.CreateDirIfNotExist(gr.LocalPath)

	auth, err := auth(gr.PrivateKeyPath, "")
	if err != nil {
		return fmt.Errorf("failed to create git auth credentials: %w", err)
	}

	var r *git.Repository

	r, err = git.PlainOpen(gr.LocalPath)
	if err != nil && !errors.Is(err, git.ErrRepositoryNotExists) {
		return err
	}

	if errors.Is(err, git.ErrRepositoryNotExists) {
		_, err = git.PlainClone(config.GetString("local.path"), false, &git.CloneOptions{
			URL:               config.GetString("remote.url"),
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			Auth:              auth,
		})
		if err != nil && !errors.Is(err, transport.ErrEmptyRemoteRepository) {
			return fmt.Errorf("failed to clone git repository: %w", err)
		}
	}

	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get git worktree: %w", err)
	}

	// Pull the latest changes from the origin remote and merge into the current branch
	if err := w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       auth,
	}); err != nil && !errors.Is(git.NoErrAlreadyUpToDate, err) && !errors.Is(err, transport.ErrEmptyRemoteRepository) {
		return fmt.Errorf("failed to pull latest changes from origin remote: %w", err)
	}

	return nil
}

func (gr *GitRemote) commit(ctx context.Context, g *git.Repository, msg ...string) error {
	w, err := g.Worktree()
	if err != nil {
		return err
	}
	_, err = w.Add(".")
	if err != nil {
		return err
	}

	commitMsg := fmt.Sprintf("pm: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	if len(msg) > 0 {
		commitMsg += strings.Join(msg, " ")
	}

	_, err = w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  gr.Name,
			Email: gr.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func auth(privateKey string, password string) (transport.AuthMethod, error) {
	return ssh.NewPublicKeysFromFile("git", privateKey, password)
}
