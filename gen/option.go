package gen

import (
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
)

func NewGenerator(p *plugin.Plugin, ops ...Opt) *Generator {
	g := &Generator{p: p}
	for _, op := range ops {
		op(g)
	}
	return g
}

type Opt func(ops *Generator)

func WithFilenameTemplate(filenameTemplate string) Opt {
	return func(ops *Generator) {
		ops.filenameTemplate = filenameTemplate
	}
}
