package schemaParser

import (
	"github.com/graphql-go/graphql"
	"reflect"
	"regexp"
	"strings"
)

type GraphQLType struct {
	Query *graphql.Object
}

type GraphQLParsingType struct {
	name     string
	typeName string
	fields   map[string]GraphQLParsingField
}

type GraphqlParingParam struct {
	name  string
	typeP string
}

type GraphQLParsingField struct {
	name     string
	typeName string
	params   []GraphqlParingParam
}

func Parsing(sdlContent string) (*graphql.Object, error) {

	namespaces := strings.Split(sdlContent, "}")
	types := []GraphQLParsingType{}
	for _, ns := range namespaces {
		// 解析类型定义
		ns = strings.TrimSpace(ns)
		if strings.HasPrefix(ns, "type ") {
			parsingType := GraphQLParsingType{}
			typep, fields := splitStringPairs(ns, "{")
			_, typeName := splitSpaceAndTabPairs(typep)
			parsingType.name = typeName
			parsingType.typeName = "type"
			fs := splitNewLineAndTabPairs(fields)
			for _, fline := range fs {
				fieldName, fieldType := splitLastSpacePairs(fline)
				fieldName = strings.TrimSpace(fieldName)
				fieldType = strings.TrimSpace(fieldType)
				f := parsingField(fieldName, fieldType)
				if parsingType.fields == nil {
					parsingType.fields = make(map[string]GraphQLParsingField)
				}
				parsingType.fields[f.name] = f
			}
			types = append(types, parsingType)
		}
	}
	println(types)
	graphqlObjMap := make(map[string]*graphql.Object)
	for _, t := range types {

		graphqlObj := createGraphqlObj(t, graphqlObjMap)
		graphqlObjMap[t.name] = graphqlObj

	}
	return graphqlObjMap["Query"], nil
}

func createGraphqlObj(t GraphQLParsingType, objMap map[string]*graphql.Object) *graphql.Object {

	gfs := make(graphql.Fields)

	for k, v := range t.fields {
		gfs[k] = &graphql.Field{
			Type: graphqlType(v.typeName, objMap),
			Args: graphQLArgs(v.params, objMap),
		}
	}

	ob := graphql.NewObject(graphql.ObjectConfig{
		Name:   t.name,
		Fields: gfs,
	})

	return ob
}

func graphQLArgs(params []GraphqlParingParam, objMap map[string]*graphql.Object) graphql.FieldConfigArgument {
	args := make(graphql.FieldConfigArgument)
	for _, p := range params {
		p.typeP, _ = extractArgs(p.typeP)
		args[p.name] = &graphql.ArgumentConfig{
			Type: graphqlType(p.typeP, objMap),
		}
	}
	return args
}

func extractArgs(s string) (string, bool) {
	if s[len(s)-1] == '!' {
		return s[:len(s)-1], true
	}
	return s, false

}

func graphqlType(name string, objMap map[string]*graphql.Object) graphql.Output {
	switch name {
	case "String":
		return graphql.String
	case "Int":
		return graphql.Int
	case "BigDecimal":
		return graphql.Float
	default:
		return graphql.NewList(getTypeByName(name, objMap))
	}

}

func getTypeByName(name string, objMap map[string]*graphql.Object) graphql.Type {
	name = strings.TrimSpace(name)
	tname, isArray := extractType(name)
	t, exist := objMap[tname]
	if !exist {
		panic("not find type " + name)
	}
	//ptrType := reflect.TypeOf(t)
	//fmt.Println("Pointer type:", ptrType)
	//listType, _ := reflect.TypeOf(t).Elem().(graphql.Type)

	structValue := reflect.New(reflect.TypeOf(t).Elem())
	structValue.Elem().Set(reflect.ValueOf(t).Elem())

	if isArray {
		return graphql.NewList(structValue.Interface().(graphql.Type))
	} else {
		return structValue.Interface().(graphql.Type)
	}
	return nil
}

func extractType(input string) (string, bool) {
	pattern := `\[(.*?)\]`
	regex := regexp.MustCompile(pattern)
	match := regex.FindStringSubmatch(input)
	if len(match) > 1 {
		return match[1], true
	}
	return strings.TrimSpace(input), false
}

func parsingField(s string, t string) GraphQLParsingField {
	f := GraphQLParsingField{}
	f.typeName = t
	p1, _ := splitLasLetterPairs(s, ":")
	f.name, f.params = parsingFiledParam(p1)
	return f

}

var parsingFiledParamRe = regexp.MustCompile(`^([^(]+)(\(([^)]+)\))?`)

func parsingFiledParam(str string) (string, []GraphqlParingParam) {
	match := parsingFiledParamRe.FindStringSubmatch(str)

	if len(match) > 1 {
		name := match[1]
		bbb := match[3]
		gpp := []GraphqlParingParam{}
		if bbb != "" {
			ps := strings.Split(bbb, ",")
			for _, p := range ps {
				n, t := splitStringPairs(p, ":")
				param := GraphqlParingParam{}
				param.name = n
				param.typeP = t
				gpp = append(gpp, param)
			}
		}
		return strings.TrimSpace(name), gpp
	} else {
		panic("field format error" + str)
	}

}

func splitStringPairs(s, sep string) (string, string) {
	parts := strings.Split(s, sep)
	if len(parts) != 2 {
		panic("the string can not split a pairs")
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

func splitSpaceAndTabPairs(s string) (string, string) {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == '\t'
	})
	if len(parts) != 2 {
		panic("the string can not split a pairs")
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

func splitNewLineAndTabPairs(s string) []string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '\r' || r == '\n'
	})
	return parts
}

func splitLastSpacePairs(s string) (string, string) {
	return splitLasLetterPairs(s, " ")
}

func splitLasLetterPairs(s string, sep string) (string, string) {
	lastSpaceIndex := strings.LastIndex(s, sep)
	if lastSpaceIndex != -1 {
		return strings.TrimSpace(s[:lastSpaceIndex]), strings.TrimSpace(s[lastSpaceIndex+1:])
	} else {
		panic("the string can not split a pairs")
	}
}
