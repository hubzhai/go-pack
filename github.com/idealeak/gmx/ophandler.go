package gmx

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func watchRoot(oe *OpEvent) *ApiResult {
	ret := make(map[string]string)
	type OpParam struct {
		Name string
	}

	var filter OpParam
	if oe.param != "" {
		err := json.Unmarshal([]byte(oe.param), &filter)
		if err != nil {
			return &ApiResult{Tag: Tag_ParamIllegal}
		}
	}
	var vals map[string]reflect.Value
	if filter.Name != "" {
		vals = ValuesByNameFunc(func(str string) bool {
			return strings.Contains(str, filter.Name)
		})
	} else {
		vals = Values()
	}
	for k, v := range vals {
		ret[k] = FullNameOfType(v.Type())
	}

	return &ApiResult{Tag: Tag_OK, Msg: ret}
}

func watchVal(oe *OpEvent) *ApiResult {
	type OpParam struct {
		Name  string
		Filed string
	}

	var param OpParam
	if oe.param != "" {
		err := json.Unmarshal([]byte(oe.param), &param)
		if err != nil {
			return &ApiResult{Tag: Tag_ParamIllegal}
		}
	}

	v := FindValue(param.Name)
	if !v.IsValid() {
		return &ApiResult{Tag: Tag_NoFindValue}
	}

	s := GetSession(oe.token)
	if s == nil {
		return &ApiResult{Tag: Tag_AuthNeed}
	}

	f := FieldByNameOfValue(v, param.Name, param.Filed, s)
	if !f.IsValid() {
		return &ApiResult{Tag: Tag_NoFindField}
	}

	desc := DescOfValue(f)
	return &ApiResult{Tag: Tag_OK, Msg: desc}
}

func watchMap(oe *OpEvent) *ApiResult {
	type OpParam struct {
		Name         string
		Filed        string
		PageNo       int
		NumOfPerPage int
	}

	var param OpParam
	if oe.param != "" {
		err := json.Unmarshal([]byte(oe.param), &param)
		if err != nil {
			return &ApiResult{Tag: Tag_ParamIllegal}
		}
	}

	v := FindValue(param.Name)
	if !v.IsValid() {
		return &ApiResult{Tag: Tag_NoFindValue}
	}

	s := GetSession(oe.token)
	if s == nil {
		return &ApiResult{Tag: Tag_AuthNeed}
	}

	f := FieldByNameOfValue(v, param.Name, param.Filed, s)
	if !f.IsValid() {
		return &ApiResult{Tag: Tag_NoFindField}
	}

	f = indirect(f)
	if !f.IsValid() || f.Kind() != reflect.Map || f.IsNil() {
		return &ApiResult{Tag: Tag_MapInvalid}
	}

	var mapKeys []reflect.Value
	sKey := param.Name + "." + param.Filed
	if mk, ok := s.Get(sKey); ok {
		if keys, ok := mk.([]reflect.Value); ok {
			mapKeys = keys
		}
	} else {
		mapKeys = f.MapKeys()
		s.Set(sKey, mapKeys)
	}

	start := (param.PageNo - 1) * param.NumOfPerPage
	cnt := len(mapKeys)
	if cnt < start {
		return &ApiResult{Tag: Tag_NoMoreData}
	}

	end := param.PageNo * param.NumOfPerPage
	if cnt < end {
		end = cnt
	}

	retKeys := mapKeys[start:end]
	var desc MapDesc
	desc.Name = param.Name
	desc.Type = FullNameOfValue(f)
	desc.NumOfPerPage = param.NumOfPerPage
	desc.PageNo = param.PageNo
	desc.TotalCount = cnt
	for _, k := range retKeys {
		desc.Keys = append(desc.Keys, DescOfValue(k))
		kv := f.MapIndex(k)
		desc.Vals = append(desc.Vals, DescOfValue(kv))
	}
	return &ApiResult{Tag: Tag_OK, Msg: desc}
}

func unwatchMap(oe *OpEvent) *ApiResult {
	type OpParam struct {
		Name  string
		Filed string
	}

	var param OpParam
	if oe.param != "" {
		err := json.Unmarshal([]byte(oe.param), &param)
		if err != nil {
			return &ApiResult{Tag: Tag_ParamIllegal}
		}
	}

	s := GetSession(oe.token)
	if s == nil {
		return &ApiResult{Tag: Tag_AuthNeed}
	}

	sKey := param.Name + "." + param.Filed
	s.Delete(sKey)
	return &ApiResult{Tag: Tag_OK}
}

func watchSlice(oe *OpEvent) *ApiResult {
	type OpParam struct {
		Name         string
		Filed        string
		PageNo       int
		NumOfPerPage int
	}

	var param OpParam
	if oe.param != "" {
		err := json.Unmarshal([]byte(oe.param), &param)
		if err != nil {
			return &ApiResult{Tag: Tag_ParamIllegal}
		}
	}

	v := FindValue(param.Name)
	if !v.IsValid() {
		return &ApiResult{Tag: Tag_NoFindValue}
	}

	s := GetSession(oe.token)
	if s == nil {
		return &ApiResult{Tag: Tag_AuthNeed}
	}

	f := FieldByNameOfValue(v, param.Name, param.Filed, s)
	if !f.IsValid() {
		return &ApiResult{Tag: Tag_NoFindField}
	}

	f = indirect(f)
	if !f.IsValid() || f.Kind() != reflect.Slice || f.IsNil() {
		return &ApiResult{Tag: Tag_SliceIsNil}
	}

	start := (param.PageNo - 1) * param.NumOfPerPage
	cnt := f.Len()
	if cnt < start {
		return &ApiResult{Tag: Tag_NoMoreData}
	}

	end := param.PageNo * param.NumOfPerPage
	if cnt < end {
		end = cnt
	}

	slice := f.Slice(start, end)
	var desc SliceDesc
	desc.Name = param.Name
	desc.Type = FullNameOfValue(f)
	desc.ElemType = FullNameOfType(f.Elem().Type())
	desc.NumOfPerPage = param.NumOfPerPage
	desc.PageNo = param.PageNo
	desc.TotalCount = cnt
	sliceKind := f.Type().Kind()
	elemKind := f.Type().Elem().Kind()
	switch sliceKind {
	case reflect.String:
		desc.Elems = slice.String()
	case reflect.Slice, reflect.Array:
		switch elemKind {
		case reflect.Uint8:
			desc.Elems = slice.Bytes()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			av := make([]int64, 0, slice.Len())
			for i := 0; i < slice.Len(); i++ {
				av = append(av, slice.Index(i).Int())
			}
			desc.Elems = av
		case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			av := make([]uint64, 0, slice.Len())
			for i := 0; i < slice.Len(); i++ {
				av = append(av, slice.Index(i).Uint())
			}
			desc.Elems = av
		case reflect.Float32, reflect.Float64:
			av := make([]float64, 0, slice.Len())
			for i := 0; i < slice.Len(); i++ {
				av = append(av, slice.Index(i).Float())
			}
			desc.Elems = av
		case reflect.Map:
			fallthrough
		case reflect.Ptr:
			fallthrough
		case reflect.Interface:
			fallthrough
		case reflect.Slice, reflect.Array:
			fallthrough
		case reflect.Struct:
			elems := make([]*ValueDesc, 0, slice.Len())
			for i := 0; i < slice.Len(); i++ {
				elems = append(elems, DescOfValue(slice.Index(i)))
			}
			desc.Elems = elems
		}
	}

	return &ApiResult{Tag: Tag_OK, Msg: desc}
}

func invoke(oe *OpEvent) *ApiResult {
	type OpParam struct {
		Name     string
		FuncName string
		Args     []string
	}

	var param OpParam
	if oe.param != "" {
		err := json.Unmarshal([]byte(oe.param), &param)
		if err != nil {
			return &ApiResult{Tag: Tag_ParamIllegal}
		}
	}

	v := FindValue(param.Name)
	if !v.IsValid() {
		return &ApiResult{Tag: Tag_NoFindValue}
	}

	var method *reflect.Method
	typ := v.Type()
	if m, ok := typ.MethodByName(param.FuncName); ok {
		method = &m
	} else {
		if typ.Kind() != reflect.Ptr {
			typ = reflect.PtrTo(typ)
		}
		if m, ok := typ.MethodByName(param.FuncName); ok {
			method = &m
		}
	}

	if method == nil {
		return &ApiResult{Tag: Tag_NoFindMethod}
	}

	var inArgs []reflect.Value
	mtype := method.Type
	//first arg mustbe receiver
	rcvTyp := mtype.In(0)
	if rcvTyp.Kind() == v.Type().Kind() {
		inArgs = append(inArgs, v)
	} else {
		if rcvTyp.Kind() == reflect.Ptr {
			pv := v.Addr()
			inArgs = append(inArgs, pv)
		} else {
			inArgs = append(inArgs, v.Elem())
		}
	}
	for i := 1; i < mtype.NumIn(); i++ {
		var argV reflect.Value
		var argPV reflect.Value
		var argIndirV reflect.Value
		argTyp := mtype.In(i)
		for {
			if argTyp.Kind() != reflect.Ptr {
				break
			}
			if !argV.IsValid() {
				argV = reflect.New(argTyp.Elem())
				argPV = argV
			} else {
				argV.Set(reflect.New(argTyp.Elem()))
			}
			argIndirV = argV
			argV = argV.Elem()
			argTyp = argV.Type()
		}
		if !argPV.IsValid() {
			argV = reflect.New(argTyp)
			argIndirV = argV
			argV = argV.Elem()
			argPV = argV
		}
		if len(param.Args) >= i {
			strArg := param.Args[i-1]
			strArg = strings.ToLower(strings.TrimSpace(strArg))
			switch argTyp.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				i64, err := strconv.ParseInt(strArg, 10, 64)
				if err == nil {
					UnsafeSetInt(argPV, i64)
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				u64, err := strconv.ParseUint(strArg, 10, 64)
				if err == nil {
					UnsafeSetUint(argPV, u64)
				}
			case reflect.Bool:
				switch strArg {
				case "true", "ok", "yes", "1":
					UnsafeSetBool(argPV, true)
				default:
					UnsafeSetBool(argPV, false)
				}
			case reflect.Float32, reflect.Float64:
				f64, err := strconv.ParseFloat(strArg, 64)
				if err == nil {
					UnsafeSetFloat(argPV, f64)
				}
			case reflect.String:
				UnsafeSetString(argPV, strArg)
			case reflect.Struct:
				err := json.Unmarshal([]byte(strArg), argIndirV.Interface())
				if err != nil {
					fmt.Errorf("Args Struct Unmarshal err:", err)
				}
			case reflect.Complex64:
				fmt.Errorf("Args not support reflect.Complex64")
			case reflect.Complex128:
				fmt.Errorf("Args not support reflect.Complex128")
			case reflect.Array:
				fmt.Errorf("Args not support Array")
			case reflect.Slice:
				fmt.Errorf("Args not support Slice")
			case reflect.Chan:
				fmt.Errorf("Args not support Chan")
			case reflect.Interface:
				fmt.Errorf("Args not support Interface")
			case reflect.Ptr:
				fmt.Errorf("Args not support Ptr")
			case reflect.Map:
				fmt.Errorf("Args not support Map")
			case reflect.Func:
				fmt.Errorf("Args not support Func")
			case reflect.UnsafePointer:
				fmt.Errorf("Args not support UnsafePointer")
			}
		}
		inArgs = append(inArgs, argPV)
	}
	outArgs := method.Func.Call(inArgs)
	ret := &InvokeResult{
		Name:     param.Name,
		FuncName: param.FuncName,
	}
	for i := 0; i < len(outArgs); i++ {
		ret.RetVals = append(ret.RetVals, DescOfValue(outArgs[i]))
	}
	return &ApiResult{Tag: Tag_OK, Msg: ret}
}

func setValue(oe *OpEvent) *ApiResult {
	type OpParam struct {
		Name  string
		Filed string
		Val   string
	}

	var param OpParam
	if oe.param != "" {
		err := json.Unmarshal([]byte(oe.param), &param)
		if err != nil {
			return &ApiResult{Tag: Tag_ParamIllegal}
		}
	}

	v := FindValue(param.Name)
	if !v.IsValid() {
		return &ApiResult{Tag: Tag_NoFindValue}
	}

	s := GetSession(oe.token)
	if s == nil {
		return &ApiResult{Tag: Tag_AuthNeed}
	}

	f := FieldByNameOfValue(v, param.Name, param.Filed, s)
	if !f.IsValid() {
		return &ApiResult{Tag: Tag_NoFindField}
	}

	f = indirectEnsure(f, true)
	if !f.IsValid() {

	}
	return &ApiResult{Tag: Tag_OK, Msg: ""}
}

func init() {
	RegisteOpHandler("watchRoot", HandlerWrapper(watchRoot))
	RegisteOpHandler("watchVal", HandlerWrapper(watchVal))
	RegisteOpHandler("watchMap", HandlerWrapper(watchMap))
	RegisteOpHandler("unwatchMap", HandlerWrapper(unwatchMap))
	RegisteOpHandler("watchSlice", HandlerWrapper(watchSlice))
	RegisteOpHandler("invoke", HandlerWrapper(invoke))
	RegisteOpHandler("setVal", HandlerWrapper(setValue))
}
