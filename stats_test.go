package errorstats

import (
	"fmt"
	"testing"
)

type ErrOne struct {
	String   string
	ErrTwo   ErrTwo
	ErrThree *ErrThree
}

type ErrTwo struct {
	String string
}

type ErrThree struct {
	String string
}

func (e ErrOne) Error() string {
	return fmt.Sprint(e)
}

func (e ErrTwo) Error() string {
	return fmt.Sprint(e)
}

func (e *ErrThree) Error() string {
	return fmt.Sprint(e)
}

func newErr() error {
	return ErrOne{
		String:   "err-one",
		ErrTwo:   ErrTwo{String: "err-two"},
		ErrThree: &ErrThree{String: "err-three"},
	}
}

func ExampleStats_JSON() {
	s := New()

	s.SetEncoder(ErrOne{}, func(v interface{}) string {
		x := v.(ErrOne)
		return s.Visit(x.String, x.ErrTwo, x.ErrThree)
	})

	s.SetEncoder(ErrTwo{}, func(v interface{}) string {
		x := v.(ErrTwo)
		return x.String
	})

	s.SetEncoder(ErrThree{}, func(v interface{}) string {
		x := v.(ErrThree)
		return x.String
	})

	s.Log(newErr())

	fmt.Println(s.JSON())
	// Output: {"err-one/ err-two/ err-three":1}
}

type myErr struct{}

func (e *myErr) Error() string {
	return "myErr"
}

func TestEncoderFuncWithPointer(t *testing.T) {
	// panics if there is any conversion errors
	{
		s := New()
		s.SetEncoder(myErr{}, func(v interface{}) string {
			x := v.(myErr)
			return x.Error()
		})
		s.Log(&myErr{})
	}
	{
		s := New()
		s.SetEncoder(&myErr{}, func(v interface{}) string {
			x := v.(myErr)
			return x.Error()
		})
		s.Log(&myErr{})
	}
}

func TestEncoderFuncWithValue(t *testing.T) {
	// panics if there is any conversion errors
	{
		s := New()
		s.SetEncoder(myErr{}, func(v interface{}) string {
			x := v.(myErr)
			return x.Error()
		})
		s.Log(myErr{})
	}
	{
		s := New()
		s.SetEncoder(&myErr{}, func(v interface{}) string {
			x := v.(myErr)
			return x.Error()
		})
		s.Log(myErr{})
	}
}

func BenchmarkLog(b *testing.B) {
	s := New()

	s.SetEncoder(ErrOne{}, func(v interface{}) string {
		x := v.(ErrOne)
		return s.Visit(x.String, x.ErrTwo, x.ErrThree)
	})

	s.SetEncoder(ErrTwo{}, func(v interface{}) string {
		x := v.(ErrTwo)
		return x.String
	})

	s.SetEncoder(ErrThree{}, func(v interface{}) string {
		x := v.(ErrThree)
		return x.String
	})

	e := newErr()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Log(e)
	}
}
