package parser

import "github.com/zeromicro/go-zero/tools/goctl/api/spec"

type GroupSpec struct {
	GroupName string
	Types     []spec.Type // group all types
	GenTypes  []spec.Type // generated in group types if level == 0
	BaseTypes []spec.Type // generated in base types if level == 0
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
			types = append(types, getRequestTypes(route.RequestType)...)
			types = append(types, getResponseTypes(route.ResponseType)...)
		}
		groupSpec.GroupName = groupName
		groupSpec.Types = types
		groupSpecs = append(groupSpecs, groupSpec)
	}

	return groupSpecs, nil
}

func getRequestTypes(requestType spec.Type) []spec.Type {
	var requestTypes []spec.Type

	switch requestType.(type) {
	case spec.DefineStruct:
		t := requestType.(spec.DefineStruct)
		requestTypes = append(requestTypes, t)
		for _, m := range t.Members {
			requestTypes = append(requestTypes, getRequestTypes(m.Type)...)
		}
	}

	return requestTypes
}

func getResponseTypes(responseType spec.Type) []spec.Type {
	var responseTypes []spec.Type

	switch responseType.(type) {
	case spec.DefineStruct:
		t := responseType.(spec.DefineStruct)
		responseTypes = append(responseTypes, t)
		for _, m := range t.Members {
			responseTypes = append(responseTypes, getRequestTypes(m.Type)...)
		}
	}

	return responseTypes
}
