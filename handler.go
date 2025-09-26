package dorky

import (
	"context"
	"errors"
	"reflect"
)

type handlerArgKind int

const (
	handlerArgKindUnknown handlerArgKind = iota
	handlerArgKindCommand
	handlerArgKindEvent
)

// validateHandler ensures that a provided handler satisfies the handler
// interface requirements. If valid, it returns the reflected type of the
// handler Message argument.
//
// A valid handler:
// - is a function
// - must receive 2 arguments
// - the first argument must be a context.Context
// - the second argument must implement the Command or Event interface
// - must return 2 values
// - the first return value must be a slice of Events
// - the second return value must implement the error interface
func validateHandler(
	handler reflect.Value,
) (
	reflect.Type,
	handlerArgKind,
	error,
) {
	argKind := handlerArgKindUnknown

	if !handler.IsValid() {
		return nil, argKind, errors.New("handler value is not valid")
	}

	handlerType := handler.Type()

	// the handler must be a func
	if handlerType.Kind() != reflect.Func {
		return nil, argKind, errors.New("handler is not a Func")
	}

	// the handler must accept two arguments
	if handlerType.NumIn() != 2 {
		return nil, argKind, errors.New(
			"handler does not take two arguments",
		)
	}

	arg1Type := handlerType.In(0)

	// the handler's first argument must implement the Context interface
	if !arg1Type.Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		return nil, argKind, errors.New(
			"handler 1st argument type does not implement a Context interface",
		)
	}

	arg2Type := handlerType.In(1)

	// the handler's second argument must either implement the Command or
	// Event interface
	if arg2Type.Implements(reflect.TypeOf((*Command)(nil)).Elem()) {
		argKind = handlerArgKindCommand
	} else if arg2Type.Implements(reflect.TypeOf((*Event)(nil)).Elem()) {
		argKind = handlerArgKindEvent
	} else {
		return nil, argKind, errors.New(
			"handler 2nd argument type must implement Command or Event interface",
		)
	}

	// the handler must return two values
	if handlerType.NumOut() != 2 {
		return nil, argKind, errors.New(
			"handler must return ([]Message, error)",
		)
	}

	returnType1 := handlerType.Out(0)

	// the handler's first return value must be convertible to a slices of Events
	if !returnType1.ConvertibleTo(reflect.TypeOf(([]Event)(nil))) {
		return nil, argKind, errors.New(
			"handler must return ([]Message, error)",
		)
	}

	returnType2 := handlerType.Out(1)

	// the handler must return an error as its second return value
	if !returnType2.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return nil, argKind, errors.New(
			"handler must return ([]Message, error)",
		)
	}

	return arg2Type, argKind, nil
}

// invokeHandler calls a handler with the message that it needs to process
// The handler provided must be a valid handler that was validated using
// validateHandler()
func invokeHandler(
	ctx context.Context,
	handler reflect.Value,
	message any,
) (
	[]Event,
	error,
) {
	inputs := make([]reflect.Value, 2)
	inputs[0] = reflect.ValueOf(ctx)
	inputs[1] = reflect.ValueOf(message)

	ret := handler.Call(inputs)

	if len(ret) != 2 {
		panic("handler did not return two values")
	}

	// Convert the first return value to a slice of Events
	var events []Event

	ret1 := ret[0].Interface()

	if eventsRet, ok := ret1.([]Event); ok {
		events = eventsRet
	}

	// convert the second return value to an error
	ret2 := ret[1].Interface()

	if ret2 == nil {
		return events, nil
	}

	err, ok := ret2.(error)

	if !ok {
		panic("handler did not return an error")
	}

	return events, err
}
