package cmd

import (
	"os"
	"path/filepath"
	"time"

	"github.com/jzero-io/goctl-types/gen"
	"github.com/spf13/cobra"
	apiParser "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	Debug bool

	FilenameTemplate string
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "goctl-types gen",
	Long:  `goctl-types gen`,
	RunE:  do,
}

func do(_ *cobra.Command, _ []string) error {
	if Debug {
		time.Sleep(time.Second * 15)
	}

	p, err := plugin.NewPlugin()
	if err != nil {
		return err
	}
	if p.Style == "" {
		p.Style = "go_zero"
	}

	p.Api, err = apiParser.Parse(p.ApiFilePath, "")
	if err != nil {
		return err
	}

	var opts []gen.Opt
	opts = append(opts, gen.WithFilenameTemplate(FilenameTemplate))

	g := gen.NewGenerator(p, opts...)

	files, err := g.Generate()
	if err != nil {
		return err
	}

	typesDir := filepath.Join(p.Dir, "internal", "types")
	emptyTypesGoBytes, err := gen.ParseTemplate(nil, []byte(gen.TypesGoTpl))
	if err != nil {
		return err
	}
	if err = os.WriteFile(filepath.Join(typesDir, "types.go"), emptyTypesGoBytes, 0o644); err != nil {
		return err
	}

	for _, v := range files {
		typesGoFilePath := filepath.Join(typesDir, v.Path)

		if !pathx.FileExists(filepath.Dir(typesGoFilePath)) {
			if err = os.MkdirAll(filepath.Dir(typesGoFilePath), 0o755); err != nil {
				return err
			}
		}
		if pathx.FileExists(typesGoFilePath) && v.Skip {
			continue
		}
		if err = os.WriteFile(typesGoFilePath, v.Content.Bytes(), 0o644); err != nil {
			return err
		}
	}

	return err
}

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().BoolVarP(&Debug, "debug", "", false, "debug")
	genCmd.Flags().StringVarP(&FilenameTemplate, "filename-template", "", "{{.group}}.types.go", "filename template")
}
