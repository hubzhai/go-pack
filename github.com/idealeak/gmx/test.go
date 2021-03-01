package gmx

import "fmt"

type TestNest struct {
	i     int
	I     int
	pi    *int
	ppi   **int
	i8    int8
	I8    int8
	pi8   *int8
	ppi8  **int8
	i16   int16
	I16   int16
	pi16  *int16
	ppi16 **int16
	i32   int32
	I32   int32
	pi32  *int32
	ppi32 **int32
	i64   int64
	I64   int64
	pi64  *int64
	ppi64 **int64
	u     uint
	U     uint
	pu    *uint
	ppu   **uint
	u8    uint8
	U8    uint8
	pu8   *uint8
	ppu8  **uint8
	u16   uint16
	U16   uint16
	pu16  *uint16
	ppu16 **uint16
	u32   uint32
	U32   uint32
	pu32  *uint32
	ppu32 **uint32
	u64   uint64
	U64   uint64
	pu64  *uint64
	ppu64 **uint64
	f32   float32
	F32   float32
	pf32  *float32
	ppf32 **float32
	f64   float64
	F64   float64
	pf64  *float64
	ppf64 **float64
	str   string
	Str   string
	pstr  *string
	ppstr **string
	//sli1d  []*TestNest
	//psli1d []*TestNest
	//sli2d  [][]*TestNest
	//psli2d [][]*TestNest
	//arr1d  [3]*TestNest
	//parr1d [3]*TestNest
	//arr2d  [3][3]TestNest
	//parr2d [3][3]*TestNest
	m   map[int]*TestNest
	pm  map[int]*TestNest
	pmi map[int]interface{}
}

func (tn TestNest) String() {
	fmt.Println(tn)
}

func NewTestNest() *TestNest {
	tn := &TestNest{
		i:     1,
		I:     2,
		pi:    new(int),  //3
		ppi:   new(*int), //4
		i8:    5,
		I8:    6,
		pi8:   new(int8),  //7
		ppi8:  new(*int8), //8
		i16:   9,
		I16:   10,
		pi16:  new(int16),  //11
		ppi16: new(*int16), //12
		i32:   13,
		I32:   14,
		pi32:  new(int32),  //15
		ppi32: new(*int32), //16
		i64:   17,
		I64:   18,
		pi64:  new(int64),  //19
		ppi64: new(*int64), //20
		u:     21,
		U:     22,
		pu:    new(uint),  //23
		ppu:   new(*uint), //24
		u8:    25,
		U8:    26,
		pu8:   new(uint8),  //27
		ppu8:  new(*uint8), //28
		u16:   29,
		U16:   30,
		pu16:  new(uint16),  //32
		ppu16: new(*uint16), //33
		u32:   34,
		U32:   35,
		pu32:  new(uint32),  //36
		ppu32: new(*uint32), //37
		u64:   38,
		U64:   39,
		pu64:  new(uint64),  //40
		ppu64: new(*uint64), //41
		f32:   42,
		F32:   43,
		pf32:  new(float32),  //44
		ppf32: new(*float32), //45
		f64:   46,
		F64:   47,
		pf64:  new(float64),  //48
		ppf64: new(*float64), //49
		str:   "50",
		Str:   "51",
		pstr:  new(string),  //52
		ppstr: new(*string), //53
	}
	return tn
}

type Test struct {
	*TestNest
	T string
}

func (t Test) String() {
	fmt.Println(t)
}

func init() {
	RegisteType(&Test{})
}
