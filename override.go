package qb

import (
	"reflect"
	"runtime"
)

// OverrideMap allows a driver to override functions from qf and qc.
// This type is not intended to be used directly
type OverrideMap map[string]interface{}

// Add adds an override to the map
func (m OverrideMap) Add(target, new interface{}) {
	rt, rn := reflect.TypeOf(target), reflect.TypeOf(new)

	if rt.Kind() != reflect.Func || rn.Kind() != reflect.Func {
		panic(`Cannot use non-function arguments in OverrideMap.Add`)
	}

	if rt != rn {
		panic(`Arguments in OverrideMap.Add must be the same type`)
	}

	if !isQbType(rt) {
		panic(`Arguments must be qb types`)
	}

	m[runtime.FuncForPC(reflect.ValueOf(target).Pointer()).Name()] = new
}

func isQbType(rt reflect.Type) bool {
	field := reflect.TypeOf((*Field)(nil)).Elem()
	condition := reflect.TypeOf((*Condition)(nil)).Elem()

	return rt.NumOut() == 1 && (rt.Out(0).Implements(field) || rt.Out(0) == condition)
}

// Condition gets an override for qc, if there is no entry in the map fallback will be used
func (m OverrideMap) Condition(source string, fallback interface{}, in []interface{}) Condition {
	return m.execute(source, fallback, in).(Condition)
}

// Field gets an override for qf, if there is no entry in the map fallback will be used
func (m OverrideMap) Field(source string, fallback interface{}, in []interface{}) Field {
	return m.execute(source, fallback, in).(Field)
}

func (m OverrideMap) execute(source string, fallback interface{}, in []interface{}) interface{} {
	values := make([]reflect.Value, len(in))
	for k, v := range in {
		if v, ok := v.(reflect.Value); ok {
			values[k] = v
			continue
		}
		if v == nil {
			values[k] = reflect.ValueOf(&in[k]).Elem()
			continue
		}
		values[k] = reflect.ValueOf(v)
	}

	if v, ok := m[source]; ok {
		return reflect.ValueOf(v).Call(values)[0].Interface()
	}

	if fallback == nil {
		panic(`Function "` + source + `" not implemented by driver`)
	}

	return reflect.ValueOf(fallback).Call(values)[0].Interface()
}
