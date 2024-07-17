package parser

import (
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

type GroupSpec struct {
	GroupName string
	Types     []spec.Type // group all types
	GenTypes  []spec.Type // generated in group types if level == 0
}

type MergeGroupSpec struct {
	group     spec.Group
	groupName string
}

func Parse(apiSpec *spec.ApiSpec) ([]GroupSpec, error) {
	var groupSpecs []GroupSpec

	// 相同 group 名就合并 routes
	var mergeGroupSpecs []MergeGroupSpec
	for _, group := range apiSpec.Service.Groups {
		var groupName string
		if name, ok := group.Annotation.Properties["group"]; ok {
			groupName = name
		} else {
			continue // 如果没有 group 属性，跳过这个 group
		}

		// 查找是否已经存在相同组名的 MergeGroupSpec
		found := false
		for i, mergeGroupSpec := range mergeGroupSpecs {
			if groupName == mergeGroupSpec.groupName {
				mergeGroupSpecs[i].group.Routes = append(mergeGroupSpecs[i].group.Routes, group.Routes...)
				found = true
				break
			}
		}

		if !found {
			mergeGroupSpecs = append(mergeGroupSpecs, MergeGroupSpec{group: group, groupName: groupName})
		}
	}

	for _, mergeGroup := range mergeGroupSpecs {
		var groupName string
		var types []spec.Type
		var groupSpec GroupSpec
		if _, ok := mergeGroup.group.Annotation.Properties["group"]; ok {
			groupName = mergeGroup.group.Annotation.Properties["group"]
		}

		for _, route := range mergeGroup.group.Routes {
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
	var types []spec.Type

	switch t := handlerType.(type) {
	case spec.PointerType:
		if tt, ok := t.Type.(spec.DefineStruct); ok {
			defineStruct := findDefineStructFromPointerTypeRawName(apiSpec, tt.RawName)
			types = append(types, defineStruct)
			for _, m := range defineStruct.Members {
				types = append(types, getHandlerTypes(apiSpec, m.Type)...)
			}
		}
	case spec.DefineStruct:
		types = append(types, t)
		for _, m := range t.Members {
			types = append(types, getHandlerTypes(apiSpec, m.Type)...)
		}
	case spec.ArrayType:
		tt, ok := t.Value.(spec.DefineStruct)
		if ok {
			for _, x := range apiSpec.Types {
				if x.Name() == tt.RawName {
					types = append(types, x)
					//if ds, ok := x.(spec.DefineStruct); ok {
					//	for _, m := range ds.Members {
					//		requestTypes = append(requestTypes, getHandlerTypes(apiSpec, m.Type)...)
					//	}
					//}
				}

			}
		}
	}
	return types
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

func findDefineStructFromPointerTypeRawName(apiSpec *spec.ApiSpec, rawName string) spec.DefineStruct {
	for _, s := range apiSpec.Types {
		if ds, ok := s.(spec.DefineStruct); ok {
			if ds.RawName == rawName {
				return ds
			}
		}
	}
	return spec.DefineStruct{}
}
