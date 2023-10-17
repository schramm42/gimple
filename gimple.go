package gimple

import (
	"fmt"
	"reflect"
)

// InjectionInterface
type InjectionInterface interface {
	GetName() string
	SetFactory(factory bool) InjectionInterface
	SetProtected(protected bool) InjectionInterface
}

// ContainerInterface
type ContainerInterface interface {
	Get(name string) (interface{}, error)
	Add(injection InjectionInterface) ContainerInterface
}

// injection
type injection struct {
	name      string
	value     interface{}
	factory   bool
	protected bool
}

// NewInjection
//
//	@param name
//	@param value
//	@return InjectionInterface
func NewInjection(name string, value interface{}) InjectionInterface {
	return &injection{name: name, value: value}
}

// GetName
//
//	@receiver i injection
//	@return string
func (i *injection) GetName() string {
	return i.name
}

// SetFactory
//
//	@receiver i injection
//	@param factory
//	@return InjectionInterface
func (i *injection) SetFactory(factory bool) InjectionInterface {
	i.factory = factory

	return i
}

// SetProtected
//
//	@receiver i injection
//	@param protected
//	@return InjectionInterface
func (i *injection) SetProtected(protected bool) InjectionInterface {
	i.protected = protected

	return i
}

// container
type container struct {
	injections map[string]*injection
	cache      map[string]*interface{}
}

// NewContainer
//
//	@return ContainerInterface
func NewContainer() ContainerInterface {
	c := new(container)
	c.injections = make(map[string]*injection)
	c.cache = make(map[string]*interface{})

	return c
}

// Get
//
//	@receiver c container
//	@param name
//	@return interface{}
//	@return error
func (c *container) Get(name string) (interface{}, error) {
	if value, ok := c.cache[name]; ok {
		return *value, nil
	}

	injection, exists := c.injections[name]
	if !exists {
		return nil, fmt.Errorf("injection %s does not exists", name)
	}

	value := injection.value
	functionArgs := []reflect.Value{reflect.ValueOf(c)}

	if isFunction(value) && !injection.protected {

		ref := reflect.ValueOf(value)
		value = ref.Call(functionArgs)[0].Interface()
		if !injection.factory {
			c.cache[name] = &value
		}
	}

	return value, nil
}

// Add
//
//	@receiver c container
//	@param newInjection
//	@return ContainerInterface
func (c *container) Add(newInjection InjectionInterface) ContainerInterface {
	c.injections[newInjection.GetName()] = newInjection.(*injection)

	return c
}

// isFunction
//
//	@param value
//	@return bool
func isFunction(value interface{}) bool {
	return reflect.TypeOf(value).Kind() == reflect.Func
}
