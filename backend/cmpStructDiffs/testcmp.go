package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type A struct {
	I     int
	S     string `diff:"ignore"`
	u     string
	Slc   []string
	Bslc  []B `diff:"non-recursive"`
	Bslc2 []B
}

type B struct {
	Ib int
	Sb string
	Cb C
}

func (b B) GetName() string {
	return b.Sb
}

type C struct {
	Ic int
	Sc string
}

type GetNamer interface {
	GetName() string
}

func TransformByTags() cmp.Option {
	return cmp.FilterPath(tagTransFilter, cmp.Transformer("tagNonRecursive", transNonRecursive))
}

func tagTransFilter(p cmp.Path) bool {
	//	fmt.Printf("tagTransFilter: %s\n", p.GoString())
	if p.Last().Type().Kind() != reflect.Struct {
		return false
	}

	sps := getParentStruct(p, len(p)-1)
	if sps < 0 {
		return false
	}

	superField, ok2 := p[sps+1].(cmp.StructField)
	if !ok2 {
		return false
	}

	superStruct := p[sps]
	superFieldRecord := superStruct.Type().Field(superField.Index())
	tag2 := superFieldRecord.Tag.Get("diff")
	tag2parsed := strings.Split(tag2, ",")
	if tag2parsed[0] == "non-recursive" {
		//		fmt.Printf("non-recursive: %v\n", p)
		return true
	}

	return false
}

func transNonRecursive(x interface{}) interface{} {
	//	fmt.Printf("Transforming: %v (%v)\n", x, reflect.TypeOf(x))
	if gmx, ok := x.(GetNamer); ok {
		//		fmt.Printf("To: %s\n", gmx.GetName())
		return gmx.GetName()
	} else {
		fmt.Printf("%v does not implement GetNamer interface\n", x)
	}

	return x
}

func IgnoreByTags() cmp.Option {
	return cmp.FilterPath(tagFilter, cmp.Ignore())
}

func getParentStruct(p cmp.Path, start int) int {
	for i := start - 1; i >= 0; i-- {
		if p[i].Type().Kind() == reflect.Struct {
			return i
		}
	}

	return -1
}

func tagFilter(p cmp.Path) bool {
	//	fmt.Printf("tagFilter: %s\n", p.GoString())
	field, ok := p.Index(-1).(cmp.StructField)
	if !ok {
		return false
	}

	// Has to exist, because field is a field
	struc := p.Index(-2)
	fieldRecord := struc.Type().Field(field.Index())
	tag1 := fieldRecord.Tag.Get("diff")
	if tag1 == "ignore" {
		//		fmt.Printf("ignore: %v\n", p)
		return true
	}

	sps := getParentStruct(p, len(p)-2)
	if sps < 0 {
		return false
	}

	superField, ok2 := p[sps+1].(cmp.StructField)
	if !ok2 {
		return false
	}

	superStruct := p[sps]
	superFieldRecord := superStruct.Type().Field(superField.Index())
	tag2 := superFieldRecord.Tag.Get("diff")
	tag2parsed := strings.Split(tag2, ",")
	if tag2parsed[0] == "non-recursive" {
		//		fmt.Printf("non-recursive: %v (%d)\n", p.GoString(), len(p))
		//		fmt.Printf("sps: %d\n", sps)
		//		fmt.Printf("field: %v, struct: %v, superStruct: %v\n", field.Type(), struc.Type(), superStruct.Type())
		return true
	}

	return false
}

func main() {
	b1 := []B{
		B{
			Ib: 211,
			Sb: "B11",
			Cb: C{
				Ic: 2111,
				Sc: "BC11",
			},
		},
		B{
			Ib: 212,
			Sb: "B12",
			Cb: C{
				Ic: 2211,
				Sc: "BC12",
			},
		},
	}

	b2 := []B{
		B{
			Ib: 221,
			Sb: "B11",
			Cb: C{
				Ic: 2221,
				Sc: "BC21",
			},
		},
		B{
			Ib: 212,
			Sb: "B12",
			Cb: C{
				Ic: 2211,
				Sc: "BC12",
			},
		},
		B{
			Ib: 223,
			Sb: "B23",
			Cb: C{
				Ic: 2223,
				Sc: "BC23",
			},
		},
	}

	a1 := A{
		I:     10,
		S:     "A1",
		u:     "u1",
		Slc:   []string{"A11", "A12", "A13"},
		Bslc:  b1,
		Bslc2: b1,
	}

	a2 := A{
		I:     11,
		S:     "A2",
		u:     "u2",
		Slc:   []string{"A11", "A21", "A13"},
		Bslc:  b2,
		Bslc2: b2,
	}

	fmt.Printf(cmp.Diff(&a1, &a2, cmpopts.IgnoreUnexported(a1), IgnoreByTags(), TransformByTags()))
}
