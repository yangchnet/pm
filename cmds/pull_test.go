package cmds

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/atotto/clipboard"
	"github.com/yangchnet/pm/config"
)

func Test_pull(t *testing.T) {
	ctx := context.Background()
	config.InitConfig()

	service, err := NewService(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := service.remote.Init(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// if err := service.remote.Pull(ctx); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
}

func Test_clip(t *testing.T) {
	clipboard.WriteAll("123")
}
