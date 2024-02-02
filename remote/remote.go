package remote

import (
	"context"
	"fmt"

	"github.com/yangchnet/pm/config"
	gitremote "github.com/yangchnet/pm/remote/git"
)

type Remote interface {
	// Push 将本地的store文件推送到远端
	Push(ctx context.Context) error

	// Pull 将远端的store文件拉取到本地
	Pull(ctx context.Context) error
}

func NewRemote(ctx context.Context) (Remote, error) {
	for name, _ := range config.GetStringMap("remote") {
		switch name {
		case "git":
			return gitremote.NewGitRemote(
				ctx,
				config.GetString("local.path"),
				&gitremote.GitRemoteConfig{
					Name:           config.GetString("remote.git.name"),
					Email:          config.GetString("remote.git.email"),
					Url:            config.GetString("remote.git.url"),
					PrivateKeyPath: config.GetString("remote.git.privateKeyPath"),
					PublicKeyPath:  config.GetString("remote.git.publicKeyPath"),
				},
			)
		}
	}

	return nil, fmt.Errorf("未配置remote")
}

// func isGit(urlStr string) bool {
// 	parsedUrl, err := url.Parse(urlStr)
// 	if err != nil {
// 		fmt.Println(err)
// 		return false
// 	}

// 	if parsedUrl.Scheme != "git" && parsedUrl.Scheme != "https" {
// 		return false
// 	}

// 	// If the URL is an HTTPS URL, we also check if the path ends with .git
// 	if parsedUrl.Scheme == "https" && !strings.HasSuffix(parsedUrl.Path, ".git") {
// 		return false
// 	}

// 	return true
// }
