package gen

import (
	"bytes"
	"github.com/jzero-io/goctl-types/parser"
	"github.com/rinchsan/gosimports"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/tools/goctl/api/gogen"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type GeneratedFile struct {
	Path    string
	Content bytes.Buffer
	Skip    bool
}

func Generate(plugin *plugin.Plugin) ([]*GeneratedFile, error) {
	var generatedFiles []*GeneratedFile

	groupSpecs, err := parser.Parse(plugin.Api)
	if err != nil {
		return nil, err
	}

	// get base types
	var baseTypes []spec.Type

	allGroupTypesRawNames := make([]string, 0)
	for _, groupSpec := range groupSpecs {
		for _, groupType := range groupSpec.GenTypes {
			if groupType.Name() != "" {
				allGroupTypesRawNames = append(allGroupTypesRawNames, groupType.Name())
			}
		}
	}

	allTypesRawName := getAllTypesRawName(*plugin.Api)
	t1, _ := lo.Difference(allTypesRawName, allGroupTypesRawNames)
	for _, t := range t1 {
		for _, apiType := range plugin.Api.Types {
			if t == apiType.Name() {
				baseTypes = append(baseTypes, apiType)
			}
		}
	}

	if _, err := os.Stat(filepath.Join(plugin.Dir, "internal", "types")); err != nil {
		_ = os.MkdirAll(filepath.Join(plugin.Dir, "internal", "types"), 0o755)
	}

	// generate group types
	for _, gs := range groupSpecs {
		if len(gs.GenTypes) == 0 {
			continue
		}
		if gs.GroupName == "" {
			baseTypes = append(baseTypes, gs.GenTypes...)
			continue
		}

		file, err := newGeneratedTypesGoFile(gs.GenTypes, gs.GroupName)
		if err != nil {
			return nil, err
		}
		generatedFiles = append(generatedFiles, file)
	}

	// generate base types
	baseFile, err := newGeneratedTypesGoFile(baseTypes, "")
	if err != nil {
		return nil, err
	}
	generatedFiles = append(generatedFiles, baseFile)

	return generatedFiles, nil
}

func newGeneratedTypesGoFile(types []spec.Type, groupName string) (*GeneratedFile, error) {
	typesGoString, err := gogen.BuildTypes(types)
	if err != nil {
		return nil, err
	}
	typesGoBytes, err := ParseTemplate(map[string]interface{}{
		"Types": typesGoString,
	}, []byte(TypesGoTpl))
	if err != nil {
		return nil, err
	}

	typesGoFormatBytes, err := gosimports.Process("", typesGoBytes, &gosimports.Options{
		FormatOnly: true,
		Comments:   true,
	})
	if err != nil {
		return nil, err
	}

	prefix := strings.ReplaceAll(filepath.Dir(groupName), "/", "_") + "_"
	if len(strings.Split(groupName, "/")) == 1 {
		prefix = ""
	}

	fileBase := filepath.Base(groupName)
	typesGoFilePath := prefix + fileBase[0:len(fileBase)-len(path.Ext(fileBase))] + ".types.go"
	if groupName == "" {
		typesGoFilePath = "types.go"
	}

	return &GeneratedFile{
		Path:    typesGoFilePath,
		Content: *bytes.NewBuffer(typesGoFormatBytes),
	}, nil
}

func getAllTypesRawName(spec spec.ApiSpec) []string {
	var allTypes []string
	for _, v := range spec.Types {
		if v.Name() != "" {
			allTypes = append(allTypes, v.Name())
		}
	}
	return allTypes
}
