package provider

import (
	"context"
	"reflect"
	"strings"
)

type Provider struct {
	constructors map[reflect.Type]map[reflect.Value]struct{}
	container    map[reflect.Type]any
}

func New() *Provider {
	return &Provider{
		constructors: make(map[reflect.Type]map[reflect.Value]struct{}),
		container:    make(map[reflect.Type]any),
	}
}

func (p *Provider) Register(constructFunction any) error {
	args, err := analyzeConstructor(constructFunction)
	if err != nil {
		return err
	}

	for _, arg := range args {
		if _, ok := p.constructors[arg]; !ok {
			p.constructors[arg] = make(map[reflect.Value]struct{})
		}
		p.constructors[arg][reflect.ValueOf(constructFunction)] = struct{}{}
	}

	return nil
}

func Get[T any](provider *Provider) (T, bool) {
	v, ok := provider.container[reflect.TypeOf(*new(T))].(T)
	return v, ok
}

type ErrNotProvided struct {
	Type reflect.Type
}

func (e ErrNotProvided) Error() string {
	return "not provided: " + e.Type.String()
}

type ErrInvalidConstructorReturn struct{}

func (e ErrInvalidConstructorReturn) Error() string {
	return "invalid constructor return"
}

type ErrMaybeCyclicDependency struct {
	cons []reflect.Value
}

func (e ErrMaybeCyclicDependency) Error() string {
	sb := strings.Builder{}
	sb.WriteString("maybe cyclic dependency: ")
	for i, con := range e.cons {
		sb.WriteString(con.String())
		if i != len(e.cons)-1 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

func getContextType() reflect.Type {
	return reflect.TypeOf((*context.Context)(nil)).Elem()
}

func (p *Provider) Construct(ctx context.Context) error {
	p.container[getContextType()] = ctx
	count := 0
	for len(p.constructors) > 0 {
		for arg, constructors := range p.constructors {
		ConsLoop:
			for con := range constructors {
				args := make([]reflect.Value, con.Type().NumIn())
				for i := 0; i < con.Type().NumIn(); i++ {
					at := con.Type().In(i)
					v, ok := p.container[at]
					if !ok {
						continue ConsLoop
					}
					args[i] = reflect.ValueOf(v)
				}

				returns := con.Call(args)

				for _, ret := range returns {
					if ret.Type().Kind().String() == "error" {
						if !ret.IsNil() {
							return ret.Interface().(error)
						}
					}

					p.container[ret.Type()] = ret.Interface()
				}

				count++

				delete(p.constructors[arg], con)
			}

			if len(p.constructors[arg]) == 0 {
				delete(p.constructors, arg)
			}
		}

		if count == 0 {
			return ErrMaybeCyclicDependency{}
		}
	}

	return nil
}
