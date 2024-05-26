package parser

import (
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	apiSpec, err := parser.Parse(filepath.Join("testdata", "desc", "main.api"))
	if err != nil {
		t.Fatal(err)
	}

	parse, err := Parse(apiSpec)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(parse)
}
