package provider

import "reflect"

type ErrNotAFunction struct{}

func (e ErrNotAFunction) Error() string {
	return "not a function"
}

func analyzeConstructor(constructFunction any) ([]reflect.Type, error) {
	if reflect.TypeOf(constructFunction).Kind() != reflect.Func {
		return nil, ErrNotAFunction{}
	}

	constructor := reflect.ValueOf(constructFunction)
	var args []reflect.Type

	for i := 0; i < constructor.Type().NumIn(); i++ {
		args = append(args, constructor.Type().In(i))
	}

	return args, nil
}
