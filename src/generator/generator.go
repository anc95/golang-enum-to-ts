package generator

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/anc95/golang-enum-to-ts/src/ast"
)

func getEnum(f ast.File) map[string]map[string]interface{} {
	body := f.Body
	enumType := make(map[string]string)
	enumValue := make(map[string]map[string]interface{})

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

					if enumValue[kind] == nil {
						enumValue[kind] = map[string]interface{}{}
					}
				}

				if iotaFlag {
					iotaValue += 1
				}

				if x.Value != "" {
					if x.Value == "iota" {
						useIota = true
						iotaFlag = true
						enumValue[kind][x.Id] = iotaValue
					} else {
						useIota = false

						if kind == "int" || enumType[kind] == "int" {
							val, _ := strconv.Atoi(x.Value)
							enumValue[kind][x.Id] = val
							prevValue = val
						}

						enumValue[kind][x.Id] = x.Value
						prevValue = x.Value
					}
				} else {
					if useIota {
						enumValue[kind][x.Id] = iotaValue
					} else {
						enumValue[kind][x.Id] = prevValue
					}
				}
			}
		}
	}

	return enumValue
}

func GenerateTS(f ast.File) string {
	enumValue := getEnum(f)
	result := ""

	keys := make([]string, 0)
	for k, _ := range enumValue {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		ret := ""

		for k, v := range enumValue[key] {
			if ret != "" {
				ret += "\n"
			}

			switch v.(type) {
			case int:
				ret += fmt.Sprintf("    %s = %d,", k, v)
			case string:
				ret += fmt.Sprintf("    %s = '%s',", k, v)
			}
		}

		ret = fmt.Sprintf("  export enum %s {\n", key) + ret
		ret += "\n  }\n"

		result += ret
	}

	result = fmt.Sprintf("namespace %s {\n", f.Name) + result
	result += "}\n"

	return result
}
