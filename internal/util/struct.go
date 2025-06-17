package util

import (
	"errors"
	"reflect"
)

// GetFieldsImplementing returns all fields of a struct implementing a certain interface.
// Returned fields are pointers to a type or interface objects.
//
// Parameter structPtr must be a pointer to a struct.
// Parameter interfaceObject must be given as a pointer to an interface,
// for example (*Insertable)(nil), where Insertable is an interface name.
func GetFieldsImplementing[T any](structPtr interface{}, interfaceObject *T) ([]T, error) {

	// Verify if structPtr is a pointer to a struct
	inputParamStructType := reflect.TypeOf(structPtr)
	if inputParamStructType == nil ||
		inputParamStructType.Kind() != reflect.Ptr ||
		inputParamStructType.Elem().Kind() != reflect.Struct {
		return nil, errors.New("invalid input structPtr param: should be a pointer to a struct")
	}

	inputParamIfcType := reflect.TypeOf(interfaceObject)
	// Verify if interfaceObject is a pointer to an interface
	if inputParamIfcType == nil ||
		inputParamIfcType.Kind() != reflect.Ptr ||
		inputParamIfcType.Elem().Kind() != reflect.Interface {

		return nil, errors.New("invalid input interfaceObject param: should be a pointer to an interface")
	}

	// We need the type, not the pointer to it.
	// By using Elem() we can get the value pointed by the pointer.
	interfaceType := inputParamIfcType.Elem()
	structType := inputParamStructType.Elem()

	structValue := reflect.ValueOf(structPtr).Elem()

	retFields := make([]T, 0)

	// Getting the VisibleFields returns all public fields in the struct
	for i, field := range reflect.VisibleFields(structType) {

		// Check if the field can be exported.
		// Interface() can be called only on exportable fields.
		if !field.IsExported() {
			continue
		}

		fieldValue := structValue.Field(i)

		// Depending on the field type, different checks apply.
		switch field.Type.Kind() {

		case reflect.Pointer:

			// Let's check if it implements the interface.
			if field.Type.Implements(interfaceType) {
				// Great, we can add it to the return slice
				retFields = append(retFields, fieldValue.Interface().(T))
			}

		case reflect.Interface:
			// If it's an interface, make sure it's not nil.
			if fieldValue.IsNil() {
				continue
			}

			// Now we can check if it's the same interface.
			if field.Type.Implements(interfaceType) {
				// Great, we can add it to the return slice
				retFields = append(retFields, fieldValue.Interface().(T))
			}

		default:
			// We can skip any other cases.
			continue
		}
	}

	return retFields, nil
}

// IsStructInitialized checks if all the struct fields are initalized (not zero).
// Members of the struct such as empty strings or numbers set to zero are interpreted as a zero value!
// Parameter structPtr needs to be a pointer to a struct.
func IsStructInitialized(structPtr interface{}) (err error) {
	inputType := reflect.TypeOf(structPtr)
	if inputType == nil ||
		inputType.Kind() != reflect.Pointer ||
		inputType.Elem().Kind() != reflect.Struct {

		return errors.New("invalid input structPtr param: should be a pointer to a struct")
	}

	// we want to access values of the struct, not value of the pointer, therefore we use Elem()
	structVal := reflect.ValueOf(structPtr).Elem()
	structType := inputType.Elem()

	for i := 0; i < structVal.NumField(); i++ {
		field := structVal.Field(i)
		if field.IsValid() && field.IsZero() {
			err = errors.Join(err, errors.New(structType.Field(i).Name+" is not initialized"))
		}
	}

	return err
}
