package schema

import (
	"go/ast"
	"gz/db/inter"
	"reflect"
)

type Field struct {
	Name string
	Type string
	Tag  string
}
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldsName []string
	fieldMap   map[string]*Field

	destReflect reflect.Value
}
type Table interface {
	TableName() string
}

func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}
func (s *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))

	var fieldValue []interface{}
	if destValue.Kind() == reflect.Slice {
		for i := 0; i < destValue.Len(); i++ {
			var tmpField []interface{}
			for _, field := range s.Fields {
				tmpField = append(tmpField, destValue.Index(i).FieldByName(field.Name).Interface())
			}
			fieldValue = append(fieldValue, tmpField)
		}
		return fieldValue
	}

	for _, field := range s.Fields {
		fieldValue = append(fieldValue, destValue.FieldByName(field.Name).Interface())
	}
	return []interface{}{fieldValue}
}
func Parse(structs interface{}, drive inter.Drive) *Schema {
	tmpRef := reflect.Indirect(reflect.ValueOf(structs))
	types := tmpRef.Type()
	if types.Kind() == reflect.Slice {
		//
		types = types.Elem()
	}
	var modelName string
	tableref := reflect.New(types)
	if table, ok := tableref.Interface().(Table); ok {
		modelName = table.TableName()
	} else {
		modelName = types.Name()
	}
	scheams := &Schema{
		Model:       structs,
		Name:        modelName,
		Fields:      nil,
		FieldsName:  nil,
		fieldMap:    make(map[string]*Field),
		destReflect: tmpRef,
	}

	for i := 0; i < types.NumField(); i++ {
		f := types.Field(i)
		if !f.Anonymous && ast.IsExported(f.Name) {
			field := &Field{
				Name: f.Name,
				Type: drive.FieldType(reflect.Indirect(reflect.New(f.Type))),
				Tag:  "",
			}
			if v, ok := f.Tag.Lookup("db"); ok {
				field.Tag = v
			}
			scheams.Fields = append(scheams.Fields, field)
			scheams.FieldsName = append(scheams.FieldsName, f.Name)
			scheams.fieldMap[f.Name] = field
		}
	}
	return scheams
}
