package schemaParser

import (
	"fmt"
	"github.com/graphql-go/graphql"
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

func Parsing(sdlContent string) (*GraphQLType, error) {

	namespaces := strings.Split(sdlContent, "}")
	types := []GraphQLParsingType{}
	for _, ns := range namespaces {
		// 解析类型定义
		ns = strings.TrimSpace(ns)
		if strings.HasPrefix(ns, "type ") {
			parsingType := GraphQLParsingType{}
			typep, fields := splitStringPairs(ns, "{")
			typeDef, typeName := splitSpaceAndTabPairs(typep)
			parsingType.name = typeName
			parsingType.typeName = "type"
			fmt.Printf("解析类型：%s\n", typeName)
			fmt.Printf("类型定义：%s\n", typeDef)
			fs := splitNewLineAndTabPairs(fields)
			for _, fline := range fs {
				fieldName, fieldType := splitLastSpacePairs(fline)
				fieldName = strings.TrimSpace(fieldName)
				fieldType = strings.TrimSpace(fieldType)
				fmt.Printf("字段名称：%s\n", fieldName)
				fmt.Printf("字段类型：%s\n", fieldType)
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
	return nil, nil
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
		return name, gpp
	} else {
		panic("field format error" + str)
	}

}

func splitStringPairs(s, sep string) (string, string) {
	parts := strings.Split(s, sep)
	if len(parts) != 2 {
		panic("the string can not split a pairs")
	}
	return parts[0], parts[1]
}

func splitSpaceAndTabPairs(s string) (string, string) {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == '\t'
	})
	if len(parts) != 2 {
		panic("the string can not split a pairs")
	}
	return parts[0], parts[1]
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
		return s[:lastSpaceIndex], s[lastSpaceIndex+1:]
	} else {
		panic("the string can not split a pairs")
	}
}
