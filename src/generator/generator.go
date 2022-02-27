package generator

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/anc95/golang-enum-to-ts/src/ast"
)

type EnumValue struct {
	Comments ast.BaseDeclaration
	Value    interface{}
}

type EnumValueMap map[string]map[string]EnumValue

func getEnum(f ast.File) EnumValueMap {
	body := f.Body
	enumType := make(map[string]string)
	enumValueMap := make(EnumValueMap)

	for i := 0; i < len(body); i++ {
		decl := body[i]

		switch decl.(type) {
		case ast.TypeDeclaration:
			decl := decl.(ast.TypeDeclaration)
			enumType[decl.Id] = string(decl.Kind)
		case ast.ConstDeclaration:
			decl := decl.(ast.ConstDeclaration)
			kind := ""
			iotaValue := 0
			useIota := false
			iotaFlag := false
			var prevValue interface{}

			for _, x := range decl.Declarators {
				if x.Kind != "" {
					kind = x.Kind

					if enumValueMap[kind] == nil {
						enumValueMap[kind] = map[string]EnumValue{}
					}
				}

				enumValue := EnumValue{Comments: x.BaseDeclaration}

				if iotaFlag {
					iotaValue += 1
				}

				if x.Value != "" {
					if x.Value == "iota" {
						useIota = true
						iotaFlag = true
						enumValue.Value = iotaValue
					} else {
						useIota = false

						if kind == "int" || enumType[kind] == "int" {
							val, _ := strconv.Atoi(x.Value)
							enumValue.Value = val
							prevValue = val
						} else {
							enumValue.Value = x.Value
							prevValue = x.Value
						}
					}
				} else {
					if useIota {
						enumValue.Value = iotaValue
					} else {
						enumValue.Value = prevValue
					}
				}

				enumValueMap[kind][x.Id] = enumValue
			}
		}
	}

	return enumValueMap
}

func GenerateTS(f ast.File) string {
	enumValueMap := getEnum(f)
	result := ""

	keys := make([]string, 0)
	for k, _ := range enumValueMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		ret := ""

		keys := make([]string, 0)
		for k, _ := range enumValueMap[key] {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := enumValueMap[key][k]

			if ret != "" {
				ret += "\n"
			}
			leadingComments, trailingComments := v.Comments.LeadingComments, v.Comments.TrailingComments
			statement := ""

			switch v.Value.(type) {
			case int:
				statement += fmt.Sprintf("    %s = %d,", k, v.Value)
			case string:
				statement += fmt.Sprintf("    %s = '%s',", k, v.Value)
			}

			for i := 0; leadingComments != nil && i < len(leadingComments); i++ {
				statement = fmt.Sprintf("    //%s", leadingComments[i].Value) + statement
			}

			for i := 0; trailingComments != nil && i < len(trailingComments); i++ {
				statement += fmt.Sprintf(" //%s", trailingComments[i].Value)
			}

			ret += statement
		}

		ret = fmt.Sprintf("  export enum %s {\n", key) + ret
		ret += "\n  }\n"

		result += ret
	}

	result = fmt.Sprintf("namespace %s {\n", f.Name) + result
	result += "}\n"

	return result
}
