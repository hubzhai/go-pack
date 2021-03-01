package gmx

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

var valuesMap = new(sync.Map)

func RegisteValue(name string, ud interface{}) {
	v := reflect.Indirect(reflect.ValueOf(ud))
	valuesMap.Store(FullNameOfValue(v)+"/"+name, v)
}

func UnregisteValue(name string, ud interface{}) {
	v := reflect.Indirect(reflect.ValueOf(ud))
	valuesMap.Delete(FullNameOfValue(v) + "/" + name)
}

func Values() map[string]reflect.Value {
	m := make(map[string]reflect.Value)
	valuesMap.Range(func(k, v interface{}) bool {
		m[k.(string)] = v.(reflect.Value)
		return true
	})
	return m
}

func ValuesByNameFunc(match func(string) bool) map[string]reflect.Value {
	m := make(map[string]reflect.Value)
	valuesMap.Range(func(k, v interface{}) bool {
		kn := k.(string)
		if match(kn) {
			m[kn] = v.(reflect.Value)
		}
		return true
	})
	return m
}

func FindValue(name string) reflect.Value {
	if v, ok := valuesMap.Load(name); ok {
		return v.(reflect.Value)
	}
	return reflect.Value{}
}

func FieldByIndexOfValue(v reflect.Value, idx []int) reflect.Value {
	v = reflect.Indirect(v)
	return v.FieldByIndex(idx)
}

func FieldByNameOfValue(v reflect.Value, vname, fname string, s *Session) reflect.Value {
	fname = strings.TrimSpace(fname)
	if fname == "" {
		return v
	}
	ns := strings.Split(fname, ".")
	for i, name := range ns {
		v = indirectEnsure(v, false)
		if !v.IsValid() {
			return reflect.Value{}
		}
		n := len(name)
		if n > 2 && name[0] == '[' && name[n-1] == ']' { // array,slice,map:format: a.[x].[k:x].[x]
			kind := v.Kind()
			switch kind {
			case reflect.Map:
				strIter := name[1 : n-1]
				if len(strIter) > 2 && strIter[1] == ':' && strIter[0] == 'k' || strIter[0] == 'K' { //按键值找,只支持有限的类型
					strkey := strIter[2:]
					iter := v.MapRange()
					if iter.Next() {
						iKey := iter.Key()
						kc := reflect.New(iKey.Type())
						indirectK := indirectEnsure(kc, true)
						switch indirectK.Type().Kind() {
						case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
							kval, err := strconv.ParseInt(strkey, 10, 64)
							if err != nil {
								return reflect.Value{}
							}
							SetInt(kc, kval)
							v = v.MapIndex(kc.Elem())
						case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
							kval, err := strconv.ParseUint(strkey, 10, 64)
							if err != nil {
								return reflect.Value{}
							}
							SetUint(kc, kval)
							v = v.MapIndex(kc.Elem())
						case reflect.Float32, reflect.Float64:
							kval, err := strconv.ParseFloat(strkey, 64)
							if err != nil {
								return reflect.Value{}
							}
							SetFloat(kc, kval)
							v = v.MapIndex(kc.Elem())
						case reflect.String:
							SetString(kc, strkey)
							v = v.MapIndex(kc.Elem())
						default:
							return reflect.Value{}
						}
					}
				} else { //按快照索引找,需要提前将keys关联进session
					if len(strIter) > 2 && strIter[1] == ':' {
						strIter = strIter[2:]
					}
					idx, err := strconv.ParseInt(strIter, 10, 64)
					if err != nil {
						return reflect.Value{}
					}
					sKey := vname + "." + strings.Join(ns[:i], ".")
					if mk, ok := s.Get(sKey); ok {
						if keys, ok := mk.([]reflect.Value); ok {
							if idx < 0 || int(idx) > len(keys) {
								return reflect.Value{}
							}
							k := keys[int(idx)]
							v = v.MapIndex(k)
							continue
						}
					}
				}
			case reflect.Slice, reflect.Array:
				strIter := name[1 : n-1]
				if len(strIter) > 2 {
					if strIter[1] == ':' {
						strIter = strIter[2:]
					}
				}
				idx, err := strconv.ParseInt(strIter, 10, 64)
				if err != nil {
					return reflect.Value{}
				}
				if idx < 0 || int(idx) > v.Len() {
					return reflect.Value{}
				}
				v = v.Index(int(idx))
				continue
			default:
				return reflect.Value{}
			}
		} else {
			v = v.FieldByName(name)
			if !v.IsValid() {
				return v
			}
		}
	}
	return v
}

func MethodsOfValue(v reflect.Value) []reflect.Value {
	var m []reflect.Value
	n := v.NumMethod()
	for i := 0; i < n; i++ {
		m = append(m, v.Method(i))
	}
	return m
}

func FieldsOfValue(v reflect.Value) []reflect.Value {
	var m []reflect.Value
	n := v.NumField()
	for i := 0; i < n; i++ {
		m = append(m, v.Field(i))
	}
	return m
}

func FullNameOfValue(v reflect.Value) string {
	var name string
	t := v.Type()
	pkgPath := t.PkgPath()
	if pkgPath != "" {
		name = pkgPath + "/"
	}
	name = name + t.Name()
	return name
}

func indirect(v reflect.Value) reflect.Value {
	for {
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() {
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if v.IsNil() {
			return v
		}
		v = v.Elem()
	}

	return v
}

func indirectEnsure(v reflect.Value, newWhenNil bool) reflect.Value {
	if !v.IsValid() {
		return v
	}
	// If v is a named type and is addressable,
	// start with its address, so that if the type has pointer methods,
	// we find them.
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}

	for {
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() {
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if v.IsNil() {
			if newWhenNil {
				if v.CanSet() {
					v.Set(reflect.New(v.Type().Elem()))
				} else {
					return reflect.Value{}
				}
			} else {
				return reflect.Value{}
			}
		}
		v = v.Elem()
	}

	return v
}

func SetBool(v reflect.Value, x bool) bool {
	pv := indirectEnsure(v, true)
	if pv.IsValid() {
		ret := true
		defer func() {
			if err := recover(); err != nil {
				ret = false
			}
		}()
		pv.SetBool(x)
		return ret
	}
	return false
}

func SetBytes(v reflect.Value, x []byte) bool {
	pv := indirectEnsure(v, true)
	if pv.IsValid() {
		ret := true
		defer func() {
			if err := recover(); err != nil {
				ret = false
			}
		}()
		pv.SetBytes(x)
		return ret
	}
	return false
}

func SetComplex(v reflect.Value, x complex128) bool {
	pv := indirectEnsure(v, true)
	if pv.IsValid() && !pv.OverflowComplex(x) {
		ret := true
		defer func() {
			if err := recover(); err != nil {
				ret = false
			}
		}()
		pv.SetComplex(x)
		return ret
	}
	return false
}

func SetInt(v reflect.Value, x int64) bool {
	pv := indirectEnsure(v, true)
	if pv.IsValid() && !pv.OverflowInt(x) {
		ret := true
		defer func() {
			if err := recover(); err != nil {
				ret = false
			}
		}()
		pv.SetInt(x)
		return ret
	}
	return false
}

func SetUint(v reflect.Value, x uint64) bool {
	pv := indirectEnsure(v, true)
	if pv.IsValid() && !pv.OverflowUint(x) {
		ret := true
		defer func() {
			if err := recover(); err != nil {
				ret = false
			}
		}()
		pv.SetUint(x)
		return ret
	}
	return false
}

func SetFloat(v reflect.Value, x float64) bool {
	pv := indirectEnsure(v, true)
	if pv.IsValid() && !pv.OverflowFloat(x) {
		ret := true
		defer func() {
			if err := recover(); err != nil {
				ret = false
			}
		}()
		pv.SetFloat(x)
		return ret
	}
	return false
}

func SetString(v reflect.Value, x string) bool {
	pv := indirectEnsure(v, true)
	if pv.IsValid() {
		ret := true
		defer func() {
			if err := recover(); err != nil {
				ret = false
			}
		}()
		pv.SetString(x)
		return ret
	}
	return false
}

func SetPointer(v reflect.Value, x unsafe.Pointer) bool {
	pv := indirectEnsure(v, true)
	if pv.IsValid() {
		ret := true
		defer func() {
			if err := recover(); err != nil {
				ret = false
			}
		}()
		pv.SetPointer(x)
		return ret
	}
	return false
}

func UnsafeSetBool(v reflect.Value, x bool) bool {
	pv := indirectEnsure(v, true)
	if pv.IsValid() && pv.CanAddr() && pv.UnsafeAddr() != 0 {
		if v.Kind() != reflect.Bool {
			return false
		}
		*(*bool)(unsafe.Pointer(pv.UnsafeAddr())) = x
		return true
	}
	return false
}

func UnsafeSetBytes(v reflect.Value, x []byte) bool {
	pv := indirectEnsure(v, true)
	if pv.IsValid() && pv.CanAddr() && pv.UnsafeAddr() != 0 {
		if v.Kind() != reflect.Slice {
			return false
		}
		if pv.Type().Elem().Kind() != reflect.Uint8 {
			return false
		}
		*(*[]byte)(unsafe.Pointer(pv.UnsafeAddr())) = x
		return true
	}
	return false
}

func UnsafeSetComplex(v reflect.Value, x complex128) bool {
	pv := indirectEnsure(v, true)
	if pv.IsValid() && pv.CanAddr() && pv.UnsafeAddr() != 0 {
		switch pv.Kind() {
		case reflect.Complex64:
			if pv.OverflowComplex(x) {
				return false
			}
			*(*complex64)(unsafe.Pointer(pv.UnsafeAddr())) = complex64(x)
		case reflect.Complex128:
			if pv.OverflowComplex(x) {
				return false
			}
			*(*complex128)(unsafe.Pointer(pv.UnsafeAddr())) = x
		default:
			return false
		}
		return true
	}
	return false
}

func UnsafeSetInt(v reflect.Value, x int64) bool {
	pv := indirectEnsure(v, true)
	if pv.IsValid() && pv.CanAddr() && pv.UnsafeAddr() != 0 {
		switch pv.Kind() {
		case reflect.Int:
			if pv.OverflowInt(x) {
				return false
			}
			*(*int)(unsafe.Pointer(pv.UnsafeAddr())) = int(x)
		case reflect.Int8:
			if pv.OverflowInt(x) {
				return false
			}
			*(*int8)(unsafe.Pointer(pv.UnsafeAddr())) = int8(x)
		case reflect.Int16:
			if pv.OverflowInt(x) {
				return false
			}
			*(*int16)(unsafe.Pointer(pv.UnsafeAddr())) = int16(x)
		case reflect.Int32:
			if pv.OverflowInt(x) {
				return false
			}
			*(*int32)(unsafe.Pointer(pv.UnsafeAddr())) = int32(x)
		case reflect.Int64:
			if pv.OverflowInt(x) {
				return false
			}
			*(*int64)(unsafe.Pointer(pv.UnsafeAddr())) = x
		default:
			return false
		}
		return true
	}
	return false
}

func UnsafeSetUint(v reflect.Value, x uint64) bool {
	pv := indirectEnsure(v, true)
	if pv.IsValid() && pv.CanAddr() && pv.UnsafeAddr() != 0 {
		switch pv.Kind() {
		case reflect.Uint:
			if pv.OverflowUint(x) {
				return false
			}
			*(*uint)(unsafe.Pointer(pv.UnsafeAddr())) = uint(x)
		case reflect.Uint8:
			if pv.OverflowUint(x) {
				return false
			}
			*(*uint8)(unsafe.Pointer(pv.UnsafeAddr())) = uint8(x)
		case reflect.Uint16:
			if pv.OverflowUint(x) {
				return false
			}
			*(*uint16)(unsafe.Pointer(pv.UnsafeAddr())) = uint16(x)
		case reflect.Uint32:
			if pv.OverflowUint(x) {
				return false
			}
			*(*uint32)(unsafe.Pointer(pv.UnsafeAddr())) = uint32(x)
		case reflect.Uint64:
			if pv.OverflowUint(x) {
				return false
			}
			*(*uint64)(unsafe.Pointer(pv.UnsafeAddr())) = x
		case reflect.Uintptr:
			if pv.OverflowUint(x) {
				return false
			}
			*(*uintptr)(unsafe.Pointer(pv.UnsafeAddr())) = uintptr(x)
		default:
			return false
		}
		return true
	}
	return false
}

func UnsafeSetFloat(v reflect.Value, x float64) bool {
	pv := indirectEnsure(v, true)
	if pv.IsValid() && pv.CanAddr() && pv.UnsafeAddr() != 0 {
		switch pv.Kind() {
		case reflect.Float32:
			if pv.OverflowFloat(x) {
				return false
			}
			*(*float32)(unsafe.Pointer(pv.UnsafeAddr())) = float32(x)
		case reflect.Float64:
			if pv.OverflowFloat(x) {
				return false
			}
			*(*float64)(unsafe.Pointer(pv.UnsafeAddr())) = x
		default:
			return false
		}
		return true
	}
	return false
}

func UnsafeSetString(v reflect.Value, x string) bool {
	pv := indirectEnsure(v, true)
	if pv.IsValid() && pv.CanAddr() && pv.UnsafeAddr() != 0 {
		if v.Kind() != reflect.String {
			return false
		}
		*(*string)(unsafe.Pointer(pv.UnsafeAddr())) = x
		return true
	}
	return false
}

func UnsafeSetPointer(v reflect.Value, x uintptr) bool {
	if v.IsValid() && v.CanAddr() && v.UnsafeAddr() != 0 {
		if v.Kind() != reflect.Ptr {
			return false
		}
		*(*uintptr)(unsafe.Pointer(v.UnsafeAddr())) = x
		return true
	}
	return false
}

func outlineOfValue(fv reflect.Value) interface{} {
	var i interface{}
	if fv.IsValid() {
		switch fv.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			i = fv.Int()
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			i = fv.Uint()
		case reflect.Float32, reflect.Float64:
			i = fv.Float()
		case reflect.Bool:
			i = fv.Bool()
		case reflect.Complex64, reflect.Complex128:
			i = fv.Complex()
		case reflect.Array:
			i = fmt.Sprintf("len:%d,cap:%d", fv.Len(), fv.Cap())
		case reflect.Slice:
			if !fv.IsNil() {
				i = fmt.Sprintf("len:%d,cap:%d", fv.Len(), fv.Cap())
			}
		case reflect.Chan:
			if !fv.IsNil() {
				i = fmt.Sprintf("len:%d,cap:%d", fv.Len(), fv.Cap())
			}
		case reflect.Interface, reflect.Ptr:
			if !fv.IsNil() && fv.CanInterface() {
				i = outlineOfValue(fv.Elem())
			}
		case reflect.String:
			n := fv.Len()
			var str string
			if n <= 256 {
				str = fv.String()
			} else {
				str = fv.String()[:256] + "..."
			}
			i = str
		case reflect.Struct:
			if fv.CanInterface() {
				d, err := json.Marshal(fv.Interface())
				if err == nil {
					if len(d) <= 128 {
						i = string(d)
					} else {
						i = string(d[:128]) + "..."
					}
				}
			} else {
				i = "[invisable]"
			}
		case reflect.Map:
			if !fv.IsNil() {
				i = fmt.Sprintf("len:%d", fv.Len())
			}
		case reflect.Func, reflect.UnsafePointer:
			if !fv.IsNil() {
				i = fmt.Sprintf("0x%x", fv.Pointer())
			}
		}
	}
	return i
}

func DescOfValue(v reflect.Value) *ValueDesc {
	var vd ValueDesc
	v = reflect.Indirect(v)
	t := v.Type()
	vd.Name = t.Name()
	vd.Type = FullNameOfType(t)
	vd.CanSet = v.CanSet()
	vd.Value = outlineOfValue(v)
	if v.Kind() == reflect.Struct {
		n := t.NumField()
		for i := 0; i < n; i++ {
			var fd FieldValueDesc
			f := t.Field(i)
			fd.Name = f.Name
			fd.Kind = uint(f.Type.Kind())
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
			fv := v.Field(i)
			fv = indirect(fv)
			if fv.IsValid() {
				fd.Value = outlineOfValue(fv)
			}
			vd.Fields = append(vd.Fields, fd)
		}

		// force method receive is pointer
		t = reflect.PtrTo(t)
		n = t.NumMethod()
		for i := 0; i < n; i++ {
			m := t.Method(i)
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
			vd.Methods = append(vd.Methods, md)
		}
	}
	return &vd
}
