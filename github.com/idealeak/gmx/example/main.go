package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/idealeak/gmx"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

type innerType struct {
	i       int
	ip      *int
	Int     int
	Int64   int64
	Int32   int32
	Int16   int16
	Int8    int8
	Uint64  uint64
	Uint32  uint32
	Uint16  uint16
	Uint8   uint8
	Str     string
	IntP    *int
	Int64P  ***int64
	Int32P  *int32
	Int16P  *int16
	Int8P   *int8
	Uint64P **uint64
	Uint32P ***uint32
	Uint16P ****uint16
	Uint8P  *****uint8
	StrP    ***string
	Slice   []*innerType
	Slice2  [][]*innerType
	sssMap  map[int]*innerType
	psssMap map[*int]*innerType
	ssssMap map[string]*innerType
}

func (it innerType) Func2(pi **int, i2 int32, s1 string, ps2 *string, b1 bool, pb2 *bool, f1 float32, pf2 *float64, it2 innerType) (int, innerType) {
	var tem int
	if pi != nil {
		ppi := *pi
		if ppi != nil {
			tem = *ppi
		}
	}
	it.Int32 = i2
	fmt.Println(it.Int, it.Str, pi, tem, i2, s1, ps2, b1, pb2, f1, pf2, it2)
	return it.i, it
}

func (it innerType) String() string {
	return "innerType"
}

type Test struct {
	IT    *innerType
	itt   innerType
	m     map[string]struct{}
	mi    map[*innerType]*gmx.Test
	l     list.List
	sm    sync.Map
	mp    *map[string]struct{}
	c     chan int
	cp    *chan int
	a     []int
	ap    *[]int
	arr   [8]int
	i1    interface{}
	i2    interface{}
	slice []innerType
	ssMap map[int]*innerType
	sss   []*innerType
}

func (t Test) Func1() {
	fmt.Println(t.IT)
}

func (t *Test) Func3() {
	fmt.Println(t.IT)
}

func main() {
	type result struct {
		D struct {
			Code int
			Url  string
		}
	}
	res := &result{}
	res.D.Code = 1
	res.D.Url = "123"
	if data, err := json.Marshal(res); err == nil {
		r := new(result)
		r = nil
		err = json.Unmarshal(data, &r)
		if err == nil {
			fmt.Println(r)
		}
	}

	it := &innerType{
		//ip:     new(int),
		Int:     1,
		Int8:    2,
		Int16:   3,
		Int64:   4,
		Int8P:   new(int8),
		Int16P:  new(int16),
		Int32P:  new(int32),
		Str:     "hello worldCCCCCCCCCCCCCCCCCCCCCCC",
		sssMap:  make(map[int]*innerType),
		psssMap: make(map[*int]*innerType),
		ssssMap: make(map[string]*innerType),
	}
	for i := 0; i < 10; i++ {
		pk := new(int)
		*pk = i
		it.sssMap[i] = &innerType{Int: 10000 + i, IntP: pk}
		it.psssMap[pk] = &innerType{Int: 10000 + i}
		it.ssssMap[strconv.Itoa(i)] = &innerType{Int: 10000 + i}
		it.Slice = append(it.Slice, &innerType{Int: 20000 + i})
		var ss2 []*innerType
		for j := 0; j < 5; j++ {
			ss2 = append(ss2, &innerType{Int: 30000 + i*100 + j})
		}
		it.Slice2 = append(it.Slice2, ss2)
	}
	v := reflect.ValueOf(it)
	v = reflect.Indirect(v)

	//typ := reflect.TypeOf(it)
	typ := v.Type()
	if method, ok := typ.MethodByName("Func2"); ok {
		var inArgs []reflect.Value
		mtype := method.Type
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

		param := []string{"1", "2", "louhaoyu", "zhangjinghua", "true", "false", "1.0", "99.8", `{"Int32":999}`}
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
			if len(param) >= i {
				strArg := param[i-1]
				strArg = strings.ToLower(strings.TrimSpace(strArg))
				switch argTyp.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					i64, err := strconv.ParseInt(strArg, 10, 64)
					if err == nil {
						gmx.UnsafeSetInt(argPV, i64)
					}
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
					u64, err := strconv.ParseUint(strArg, 10, 64)
					if err == nil {
						gmx.UnsafeSetUint(argPV, u64)
					}
				case reflect.Bool:
					switch strArg {
					case "true", "ok", "yes", "1":
						gmx.UnsafeSetBool(argPV, true)
					default:
						gmx.UnsafeSetBool(argPV, false)
					}
				case reflect.Float32, reflect.Float64:
					f64, err := strconv.ParseFloat(strArg, 64)
					if err == nil {
						gmx.UnsafeSetFloat(argPV, f64)
					}
				case reflect.String:
					gmx.UnsafeSetString(argPV, strArg)
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
		for k := 0; k < len(outArgs); k++ {
			fmt.Println(gmx.DescOfValue(outArgs[k]))
		}
		fmt.Println(outArgs)
	}
	f := v.FieldByName("Int")
	if f.IsValid() && f.CanSet() {
		gmx.SetInt(f, 11)
	}
	f = v.FieldByName("Int8P")
	if f.IsValid() && f.CanSet() {
		gmx.SetInt(f, 12)
	}
	f = v.FieldByName("Int64P")
	if f.IsValid() && f.CanSet() {
		gmx.SetInt(f, 6422)
	}
	f = v.FieldByName("Uint32P")
	if f.IsValid() && f.CanSet() {
		gmx.SetUint(f, 322)
	}
	f = v.FieldByName("StrP")
	if f.IsValid() && f.CanSet() {
		gmx.SetString(f, "lyk")
	}

	vd := gmx.DescOfValue(v)
	data, err := json.MarshalIndent(vd, "", "\t")
	if err == nil {
		fmt.Println(string(data))
	}

	t := &Test{
		IT: it,
		m:  make(map[string]struct{}),
		mi: make(map[*innerType]*gmx.Test),
	}
	t.i1 = t.IT
	t.i2 = 1
	for i := 0; i < 1000; i++ {
		t.m[strconv.Itoa(i)] = struct{}{}
		t.mi[&innerType{i: i}] = &gmx.Test{T: strconv.Itoa(i)}
	}
	for i := 0; i < 10; i++ {
		t.slice = append(t.slice, innerType{})
	}
	v = reflect.ValueOf(t)
	//if v.CanAddr(){
	//	v=v.Addr()
	//}

	s := gmx.NewSession("lyk")
	gmx.SetSession(s)
	f = gmx.FieldByNameOfValue(v, "", "IT.sssMap", nil)
	if f.IsValid() {
		mapKeys := f.MapKeys()
		s.Set(".IT.sssMap", mapKeys)
	}

	iter := f.MapRange()
	if iter.Next() {
		k := iter.Key()
		kc := reflect.New(k.Type())
		kc.Elem().SetInt(1)
		v := f.MapIndex(kc.Elem())
		if v.IsValid() {
			f = gmx.FieldByNameOfValue(v, "", "Int64", nil)
			if f.IsValid() {
				gmx.UnsafeSetInt(f, 51)
			}
		}
	}

	f = gmx.FieldByNameOfValue(v, "", "IT.Slice2.[0].[0].Int64", nil)
	if f.IsValid() {
		gmx.UnsafeSetInt(f, 50)
	}

	f = gmx.FieldByNameOfValue(v, "", "IT.sssMap.[0].Int64", s)
	if f.IsValid() {
		gmx.UnsafeSetInt(f, 500)
		gmx.UnsafeSetPointer(f, 0)
	}

	f = gmx.FieldByNameOfValue(v, "", "IT.sssMap.[0].IntP", s)
	if f.IsValid() {
		gmx.UnsafeSetInt(f, 501)
	}

	f = gmx.FieldByNameOfValue(v, "", "IT.sssMap.[k:2].Int64", s)
	if f.IsValid() {
		gmx.UnsafeSetInt(f, 250)
	}

	f = gmx.FieldByNameOfValue(v, "", "IT.sssMap.[k:2].IntP", s)
	if f.IsValid() {
		gmx.UnsafeSetInt(f, 251)
		gmx.UnsafeSetPointer(f, 0)
	}

	f = gmx.FieldByNameOfValue(v, "", "IT.sssMap.[k:2].Int32P", s)
	if f.IsValid() {
		gmx.UnsafeSetInt(f, 251)
		gmx.UnsafeSetPointer(f, 0)
	}

	f = gmx.FieldByNameOfValue(v, "", "IT.psssMap.[k:2].Int64", s)
	if f.IsValid() {
		gmx.UnsafeSetInt(f, 250)
	}

	f = gmx.FieldByNameOfValue(v, "", "IT.ssssMap.[k:2].Int64", s)
	if f.IsValid() {
		gmx.UnsafeSetInt(f, 250)
	}

	p := v.Pointer()
	fmt.Println(p)
	itp := *(**innerType)(unsafe.Pointer(uintptr(p) + unsafe.Offsetof(t.IT)))
	itp.i = 99
	ittp := (*innerType)(unsafe.Pointer(uintptr(p) + unsafe.Offsetof(t.itt)))
	ittp.i = 100

	reflect.PtrTo(v.Type())
	f = gmx.FieldByNameOfValue(v, "", "mi", nil)
	if f.IsValid() && !f.IsNil() {
		keys1 := f.MapKeys()
		keys2 := f.MapKeys()
		for n, k := range keys1 {
			kk := reflect.NewAt(k.Elem().Type(), unsafe.Pointer(k.Pointer()))
			kki := gmx.FieldByNameOfValue(kk, "", "i", nil)
			if kki.IsValid() {
				fmt.Print(kki.Int(), "=")
			}
			//kkp:=reflect.NewAt(k.Type(), unsafe.Pointer(kk.UnsafeAddr()))

			k2 := keys2[n]
			val := f.MapIndex(kk)
			kval := gmx.FieldByNameOfValue(k, "", "i", nil)
			if kval.IsValid() {
				//gmx.UnsafeSetInt(kval, 11)
				fmt.Print(kval.Int(), "=")
			}
			kval2 := gmx.FieldByNameOfValue(k2, "", "i", nil)
			if kval2.IsValid() {
				//gmx.UnsafeSetInt(kval, 11)
				fmt.Print(kval2.Int(), "=")
			}
			tval := gmx.FieldByNameOfValue(val, "", "T", nil)
			if tval.IsValid() {
				//gmx.UnsafeSetString(tval, "lyk1")
				fmt.Println(tval.String())
			}
		}
	}

	f = gmx.FieldByNameOfValue(v, "", "slice", nil)
	if f.IsValid() && !f.IsNil() {
		n := f.Len()
		for i := 0; i < n; i++ {
			fv := f.Index(i)
			elem := gmx.FieldByNameOfValue(fv, "", "i", nil)
			if elem.IsValid() {
				gmx.UnsafeSetInt(elem, int64(i+1))
			}
		}
	}
	f = gmx.FieldByNameOfValue(v, "", "IT", nil)
	f = gmx.FieldByNameOfValue(v, "", "IT.StrP", nil)
	fmt.Println(f.String())
	f = gmx.FieldByNameOfValue(v, "", "IT.ip", nil)
	if f.IsValid() {
		gmx.UnsafeSetInt(f, 50)
	}
	f = gmx.FieldByNameOfValue(v, "", "IT.i", nil)
	if f.IsValid() {
		gmx.UnsafeSetInt(f, 60)
	}
	f = gmx.FieldByNameOfValue(v, "", "itt.ip", nil)
	if f.IsValid() {
		gmx.UnsafeSetInt(f, 50)
	}
	f = gmx.FieldByNameOfValue(v, "", "itt.i", nil)
	if f.IsValid() {
		gmx.UnsafeSetInt(f, 60)
	}
	//ip := (*int)(unsafe.Pointer(f.UnsafeAddr()))
	//*ip = 50
	fmt.Println(it)

	vd = gmx.DescOfValue(v)
	data, err = json.MarshalIndent(vd, "", "\t")
	if err == nil {
		fmt.Println(string(data))
	}
	fmt.Println(vd)
	//json.Unmarshal("", nil)
}

func init() {
	gmx.RegisteType(innerType{})
	gmx.RegisteType(&Test{})
}
