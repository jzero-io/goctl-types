package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
	"time"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "goctl-types gen",
	Long:  `goctl-types gen`,
	RunE:  gen,
}

func gen(_ *cobra.Command, _ []string) error {
	time.Sleep(time.Second * 15)
	p, err := plugin.NewPlugin()
	if err != nil {
		return err
	}

	fmt.Println(p)

	return nil
}

func init() {
	rootCmd.AddCommand(genCmd)
}
