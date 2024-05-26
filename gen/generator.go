package gen

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/jzero-io/goctl-types/parser"
	"github.com/rinchsan/gosimports"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/tools/goctl/api/gogen"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
)

type GeneratedFile struct {
	Path    string
	Content bytes.Buffer
	Skip    bool
}

type Generator struct {
	p *plugin.Plugin

	filenameTemplate string
}

func (g *Generator) Generate() ([]*GeneratedFile, error) {
	var generatedFiles []*GeneratedFile

	groupSpecs, err := parser.Parse(g.p.Api)
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

	allTypesRawName := getAllTypesRawName(*g.p.Api)
	t1, _ := lo.Difference(allTypesRawName, allGroupTypesRawNames)
	for _, t := range t1 {
		for _, apiType := range g.p.Api.Types {
			if t == apiType.Name() {
				baseTypes = append(baseTypes, apiType)
			}
		}
	}

	if _, err := os.Stat(filepath.Join(g.p.Dir, "internal", "types")); err != nil {
		_ = os.MkdirAll(filepath.Join(g.p.Dir, "internal", "types"), 0o755)
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

		file, err := g.newGeneratedTypesGoFile(gs.GenTypes, gs.GroupName)
		if err != nil {
			return nil, err
		}
		generatedFiles = append(generatedFiles, file)
	}

	// generate base types
	baseFile, err := g.newGeneratedTypesGoFile(baseTypes, "")
	if err != nil {
		return nil, err
	}
	generatedFiles = append(generatedFiles, baseFile)

	return generatedFiles, nil
}

func (g *Generator) newGeneratedTypesGoFile(types []spec.Type, groupName string) (*GeneratedFile, error) {
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

	var styledGroup string
	if len(strings.Split(groupName, "/")) == 1 {
		styledGroup = ""
	} else {
		styledGroup, err = format.FileNamingFormat(g.p.Style, strings.ReplaceAll(groupName, "/", "_"))
		if err != nil {
			return nil, err
		}
	}

	typesGoFilePathBytes, err := ParseTemplate(map[string]interface{}{
		"group": styledGroup,
	}, []byte(g.filenameTemplate))
	if err != nil {
		return nil, err
	}

	typesGoFilePath := string(typesGoFilePathBytes)
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
