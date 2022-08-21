package router

import "reflect"

type endpoint struct {
	Name        string
	HandlerFunc reflect.Value
	In          reflect.Type
	Out         reflect.Type
}
