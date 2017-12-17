package valid

import (
	"flag"
	"reflect"
	"testing"
)

func testParse(t *testing.T, fs *flag.FlagSet, act interface{}, ok bool, exp interface{}, args ...string) bool {
	// reset *act by zero.
	el := reflect.ValueOf(act).Elem()
	el.Set(reflect.Zero(el.Type()))
	err := Parse(fs, args)
	// when parse/validation failed
	if !ok {
		if err == nil {
			t.Errorf("should be failed: %v", args)
			return false
		}
		return true
	}
	// when parse/validation succeeded
	if err != nil {
		t.Errorf("should be succeeded: %v %v", exp, args)
		return false
	}
	v := el.Interface()
	if !reflect.DeepEqual(v, exp) {
		t.Errorf("the option is not match: act=%v exp=%v %v", v, exp, args)
		return false
	}
	return true
}

func TestIntMustSet(t *testing.T) {
	var opt int
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Var(Int(&opt, 0).MustSet(), "opt", "")

	testParse(t, fs, &opt, false, nil)
	testParse(t, fs, &opt, true, 123, "-opt", "123")
	testParse(t, fs, &opt, true, 999, "-opt", "999")
}

func TestIntMin(t *testing.T) {
	var opt int
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Var(Int(&opt, 0).Min(10), "opt", "")

	testParse(t, fs, &opt, false, nil)
	testParse(t, fs, &opt, false, nil, "-opt", "9")
	testParse(t, fs, &opt, true, 10, "-opt", "10")
	testParse(t, fs, &opt, true, 123, "-opt", "123")
}

func TestIntMax(t *testing.T) {
	var opt int
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Var(Int(&opt, 0).Max(10), "opt", "")

	testParse(t, fs, &opt, false, nil, "-opt", "11")
	testParse(t, fs, &opt, true, 10, "-opt", "10")
	testParse(t, fs, &opt, true, 0)
	testParse(t, fs, &opt, true, -123, "-opt", "-123")
}

func TestIntOneOf(t *testing.T) {
	var opt int
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Var(Int(&opt, 0).OneOf(100, 200, 300), "opt", "")

	testParse(t, fs, &opt, false, nil)
	testParse(t, fs, &opt, true, 100, "-opt", "100")
	testParse(t, fs, &opt, true, 200, "-opt", "200")
	testParse(t, fs, &opt, true, 300, "-opt", "300")
	testParse(t, fs, &opt, false, nil, "-opt", "1")
	testParse(t, fs, &opt, false, nil, "-opt", "99")
	testParse(t, fs, &opt, false, nil, "-opt", "101")
	testParse(t, fs, &opt, false, nil, "-opt", "900")
}
