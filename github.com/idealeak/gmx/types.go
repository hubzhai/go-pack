package gmx

import (
	"reflect"
	"sync"
	"unicode"
	"unicode/utf8"
)

var typesMap = new(sync.Map)

func RegisteType(ud interface{}) {
	v := reflect.Indirect(reflect.ValueOf(ud))
	t := v.Type()
	typesMap.Store(FullNameOfType(t), t)
}

func Types() map[string]reflect.Type {
	m := make(map[string]reflect.Type)
	typesMap.Range(func(k, v interface{}) bool {
		m[k.(string)] = v.(reflect.Type)
		return true
	})
	return m
}

func TypesByNameFunc(match func(string) bool) map[string]reflect.Type {
	m := make(map[string]reflect.Type)
	typesMap.Range(func(k, v interface{}) bool {
		kn := k.(string)
		if match(kn) {
			m[kn] = v.(reflect.Type)
		}
		return true
	})
	return m
}

func FieldsOfType(t reflect.Type) []reflect.StructField {
	var m []reflect.StructField
	n := t.NumField()
	for i := 0; i < n; i++ {
		m = append(m, t.Field(i))
	}
	return m
}

func MethodsOfType(t reflect.Type) []reflect.Method {
	var m []reflect.Method
	// To help the user, see if a pointer receiver would work.
	t = reflect.PtrTo(t)
	n := t.NumMethod()
	for i := 0; i < n; i++ {
		m = append(m, t.Method(i))
	}
	return m
}

func FullNameOfType(t reflect.Type) string {
	var name string
	pkgPath := t.PkgPath()
	if pkgPath != "" {
		name = pkgPath + "/"
	}
	name = name + t.Name()
	return name
}

func MetaOfType(t reflect.Type) *TypeDesc {
	var td TypeDesc
	td.Name = t.Name()
	td.Type = FullNameOfType(t)
	td.IsExported = isExported(td.Name)
	n := t.NumField()
	for i := 0; i < n; i++ {
		var fd FieldTypeDesc
		f := t.Field(i)
		fd.Name = f.Name
		ft := f.Type
		var p string
		for ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
			p = p + "*"
		}
		name := FullNameOfType(ft)
		if name == "" {
			name = ft.String()
		}
		fd.Type = p + name
		fd.Anonymous = f.Anonymous
		fd.IsExported = isExported(f.Name)
		fd.Index = f.Index
		fd.Offset = int(f.Offset)
		fd.Tag = string(f.Tag)
		td.Fields = append(td.Fields, fd)
	}

	existMethod := make(map[string]struct{})
	// method receive is value
	n = t.NumMethod()
	for i := 0; i < n; i++ {
		var md MethodDesc
		m := t.Method(i)
		md.Name = m.Name
		md.Index = m.Index
		md.RcvrIsPtr = false
		numIn := m.Type.NumIn()
		for j := 0; j < numIn; j++ {
			var p string
			in := m.Type.In(j)
			for in.Kind() == reflect.Ptr {
				in = in.Elem()
				p = p + "*"
			}
			name := FullNameOfType(in)
			if name == "" {
				name = in.String()
			}
			md.In = append(md.In, MethodParam{
				Name:  in.Name(),
				Type:  p + name,
				Index: j,
			})
		}

		numOut := m.Type.NumOut()
		for j := 0; j < numOut; j++ {
			var p string
			out := m.Type.Out(j)
			for out.Kind() == reflect.Ptr {
				out = out.Elem()
				p = p + "*"
			}
			name := FullNameOfType(out)
			if name == "" {
				name = out.String()
			}
			md.Out = append(md.Out, MethodParam{
				Name:  out.Name(),
				Type:  p + name,
				Index: j,
			})
		}
		td.Methods = append(td.Methods, md)
		existMethod[md.Name] = struct{}{}
	}

	// method receive is pointer
	t = reflect.PtrTo(t)
	n = t.NumMethod()
	for i := 0; i < n; i++ {
		m := t.Method(i)
		if _, exist := existMethod[m.Name]; exist {
			continue
		}
		var md MethodDesc
		md.Name = m.Name
		md.Index = m.Index
		md.RcvrIsPtr = true
		numIn := m.Type.NumIn()
		for j := 0; j < numIn; j++ {
			var p string
			in := m.Type.In(j)
			for in.Kind() == reflect.Ptr {
				in = in.Elem()
				p = p + "*"
			}
			name := FullNameOfType(in)
			if name == "" {
				name = in.String()
			}
			md.In = append(md.In, MethodParam{
				Name:  in.Name(),
				Type:  p + name,
				Index: j,
			})
		}

		numOut := m.Type.NumOut()
		for j := 0; j < numOut; j++ {
			var p string
			out := m.Type.Out(j)
			for out.Kind() == reflect.Ptr {
				out = out.Elem()
				p = p + "*"
			}
			name := FullNameOfType(out)
			if name == "" {
				name = out.String()
			}
			md.Out = append(md.Out, MethodParam{
				Name:  out.Name(),
				Type:  p + name,
				Index: j,
			})
		}
		td.Methods = append(td.Methods, md)
	}
	return &td
}

// Is this an exported - upper case - name?
func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}
