package gimple

import (
	"os"
	"reflect"
	"testing"

	"github.com/MarvinJWendt/testza"
)

// TestMain
//
//	@param m
func TestMain(m *testing.M) {
	exitVal := m.Run()

	os.Exit(exitVal)
}

const (
	STRING_VALUE_1   string = "val"
	INJECTION_NAME_1 string = "name1"
)

// Foo
type Foo struct {
	Name string
}

// TestWithClosure
//
//	@param t
func TestWithClosure(t *testing.T) {
	c := NewContainer()

	i := NewInjection(INJECTION_NAME_1, func(c ContainerInterface) *Foo {
		return &Foo{}
	})

	c.Add(i)

	f, err := c.Get(INJECTION_NAME_1)

	testza.AssertNoError(t, err)

	_, ok := f.(*Foo)

	testza.AssertTrue(t, ok)

}

// TestInjectionNotExists
//
//	@param t
func TestInjectionNotExists(t *testing.T) {
	c := NewContainer()
	_, err := c.Get(INJECTION_NAME_1)
	testza.AssertNotNil(t, err)
}

// TestWithString
//
//	@param t
func TestWithString(t *testing.T) {
	c := NewContainer()

	i := NewInjection(INJECTION_NAME_1, STRING_VALUE_1)

	result, err := c.Add(i).Get(INJECTION_NAME_1)
	testza.AssertNoError(t, err)
	testza.AssertEqual(t, result, STRING_VALUE_1)

}

// TestWithClosureSameInstance
//
//	@param t
func TestWithClosureSameInstance(t *testing.T) {
	c := NewContainer()

	i := NewInjection(INJECTION_NAME_1, func(c ContainerInterface) *Foo {
		return &Foo{Name: STRING_VALUE_1}
	})

	result1, err := c.Add(i).Get(INJECTION_NAME_1)
	testza.AssertNoError(t, err)
	result2, err := c.Add(i).Get(INJECTION_NAME_1)
	testza.AssertNoError(t, err)

	p1 := reflect.ValueOf(result1).Pointer()
	p2 := reflect.ValueOf(result2).Pointer()

	testza.AssertEqual(t, p1, p2)

}

// TestWithClosureFactory
//
//	@param t
func TestWithClosureFactory(t *testing.T) {
	c := NewContainer()

	i := NewInjection(INJECTION_NAME_1, func(c ContainerInterface) *Foo {
		return &Foo{Name: STRING_VALUE_1}
	}).SetFactory(true)

	result1, err := c.Add(i).Get(INJECTION_NAME_1)
	testza.AssertNoError(t, err)
	result2, err := c.Add(i).Get(INJECTION_NAME_1)
	testza.AssertNoError(t, err)

	p1 := reflect.ValueOf(result1).Pointer()
	p2 := reflect.ValueOf(result2).Pointer()

	testza.AssertNotEqual(t, p1, p2)

}

// TestWithClosureProtected
//
//	@param t
func TestWithClosureProtected(t *testing.T) {
	c := NewContainer()

	closure := func(c ContainerInterface) *Foo {
		return &Foo{Name: STRING_VALUE_1}
	}

	i := NewInjection(INJECTION_NAME_1, closure).SetProtected(true)

	result, err := c.Add(i).Get(INJECTION_NAME_1)
	testza.AssertNoError(t, err)

	isFunc := reflect.TypeOf(result).Kind() == reflect.Func
	testza.AssertTrue(t, isFunc)

	p1 := reflect.ValueOf(closure).Pointer()
	p2 := reflect.ValueOf(result).Pointer()

	testza.AssertEqual(t, p1, p2)

}
