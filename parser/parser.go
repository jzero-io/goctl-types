package parser

import (
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

type GroupSpec struct {
	GroupName string
	Types     []spec.Type // group all types
	GenTypes  []spec.Type // generated in group types if level == 0
}

func Parse(apiSpec *spec.ApiSpec) ([]GroupSpec, error) {
	var groupSpecs []GroupSpec

	for _, group := range apiSpec.Service.Groups {
		var groupName string
		var types []spec.Type
		var groupSpec GroupSpec
		if _, ok := group.Annotation.Properties["group"]; ok {
			groupName = group.Annotation.Properties["group"]
		}

		for _, route := range group.Routes {
			types = append(types, getHandlerTypes(apiSpec, route.RequestType)...)
			types = append(types, getHandlerTypes(apiSpec, route.ResponseType)...)
		}
		groupSpec.GroupName = groupName
		groupSpec.Types = types
		groupSpecs = append(groupSpecs, groupSpec)
	}

	for i := range groupSpecs {
		groupSpecs[i].Types = removeDuplicateTypes(groupSpecs[i].Types)
	}

	// separate group types
	groupTypesRawName := make([][]string, 0)
	for _, group := range groupSpecs {
		var typesRawName []string
		for _, name := range group.Types {
			typesRawName = append(typesRawName, name.Name())
		}
		groupTypesRawName = append(groupTypesRawName, typesRawName)
	}

	elements := separateCommonElements(groupTypesRawName...)

	for i := range groupSpecs {
		var genTypes []spec.Type
		elementArray := elements[i]
		for _, e := range elementArray {
			for _, t := range groupSpecs[i].Types {
				if t.Name() == e {
					genTypes = append(genTypes, t)
				}
			}
		}
		groupSpecs[i].GenTypes = genTypes
	}
	return groupSpecs, nil
}

func getHandlerTypes(apiSpec *spec.ApiSpec, handlerType spec.Type) []spec.Type {
	var requestTypes []spec.Type

	switch t := handlerType.(type) {
	case spec.DefineStruct:
		requestTypes = append(requestTypes, t)
		for _, m := range t.Members {
			requestTypes = append(requestTypes, getHandlerTypes(apiSpec, m.Type)...)
		}
	case spec.ArrayType:
		tt, ok := t.Value.(spec.DefineStruct)
		if ok {
			for _, x := range apiSpec.Types {
				if x.Name() == tt.Name() {
					requestTypes = append(requestTypes, x)
				}
				if ds, ok := x.(spec.DefineStruct); ok {
					for _, m := range ds.Members {
						requestTypes = append(requestTypes, getHandlerTypes(apiSpec, m.Type)...)
					}
				}
			}
		}
	}
	return requestTypes
}

func removeDuplicateTypes(types []spec.Type) []spec.Type {
	var newTypes []spec.Type

	var existMap = make(map[string]struct{})
	for _, t := range types {
		if _, ok := existMap[t.Name()]; !ok {
			newTypes = append(newTypes, t)
			existMap[t.Name()] = struct{}{}
		}
	}

	return newTypes
}
