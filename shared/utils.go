package shared

import (
	"fmt"
	"reflect"
	"strings"
)

func ToString(v interface{}) string {
	var sb strings.Builder
	visited := make(map[uintptr]bool) // 순환 참조 방지
	toStringRecursive(reflect.ValueOf(v), &sb, 0, visited)
	return sb.String()
}

func toStringRecursive(v reflect.Value, sb *strings.Builder, indent int, visited map[uintptr]bool) {
	if !v.IsValid() {
		sb.WriteString("nil")
		return
	}

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			sb.WriteString("nil")
			return
		}
		ptr := v.Pointer()
		if visited[ptr] {
			sb.WriteString(fmt.Sprintf("&<circular reference: %p>", v.Pointer()))
			return
		}
		visited[ptr] = true
		sb.WriteString("&")
		toStringRecursive(v.Elem(), sb, indent, visited)

	case reflect.Struct:
		t := v.Type()
		sb.WriteString(t.Name() + " {\n")
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			fieldValue := v.Field(i)
			if field.PkgPath != "" { // unexported
				continue
			}
			sb.WriteString(strings.Repeat("  ", indent+1))
			sb.WriteString(field.Name + ": ")
			toStringRecursive(fieldValue, sb, indent+1, visited)
			sb.WriteString("\n")
		}
		sb.WriteString(strings.Repeat("  ", indent) + "}")

	case reflect.Slice, reflect.Array:
		sb.WriteString("[")
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				sb.WriteString(", ")
			}
			toStringRecursive(v.Index(i), sb, indent, visited)
		}
		sb.WriteString("]")

	case reflect.Map:
		sb.WriteString("{")
		for i, key := range v.MapKeys() {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%v: ", key))
			toStringRecursive(v.MapIndex(key), sb, indent, visited)
		}
		sb.WriteString("}")

	case reflect.String:
		sb.WriteString(fmt.Sprintf("\"%s\"", v.String()))

	case reflect.Interface:
		toStringRecursive(v.Elem(), sb, indent, visited)

	default:
		sb.WriteString(fmt.Sprintf("%v", v.Interface()))
	}
}
