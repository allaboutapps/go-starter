//go:build scripts

package scaffold

import (
	"fmt"
	"go/ast"
	"go/parser"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/go-openapi/swag"
	"golang.org/x/mod/modfile"
)

// Scaffolding tool to auto-generate basic CRUD handlers for a given database model.

type FieldType struct {
	Name string
}

type Field struct {
	Name string
	Type FieldType
}

type StorageResource struct {
	Name   string
	Fields []Field
}

func ParseModel(path string) (*StorageResource, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	regex, err := regexp.Compile(`type ([^\s]+) (struct {[^}]+})`)
	if err != nil {
		return nil, err
	}

	matches := regex.FindStringSubmatch(string(content))
	resourceName := matches[1]
	expression := matches[2]

	node, err := parser.ParseExpr(expression)
	if err != nil {
		return nil, err
	}

	structType, ok := node.(*ast.StructType)
	if !ok {
		return nil, fmt.Errorf("%s is not a struct definition", expression)
	}

	fields := []Field{}
	for _, field := range structType.Fields.List {
		name := field.Names[0].Name
		fieldType := expression[field.Type.Pos()-1 : field.Type.End()-1]
		field := Field{
			Name: name,
			Type: FieldType{
				Name: fieldType,
			},
		}

		if field.Name == "R" || field.Name == "L" {
			break
		}

		fields = append(fields, field)
	}

	resource := StorageResource{
		Name:   resourceName,
		Fields: fields,
	}

	return &resource, nil
}

type Property struct {
	Name     string
	Type     string
	Required bool
	Format   *string
}

type SwaggerResource struct {
	Name              string
	URLName           string
	Properties        []Property
	PayloadProperties []Property
}

func GenerateSwagger(resource *StorageResource, outputPath string, force bool) error {
	definitionsPath := filepath.Join(outputPath, "definitions")
	pathsPath := filepath.Join(outputPath, "paths")
	if err := createDirIfAbsent(definitionsPath); err != nil {
		return err
	}
	if err := createDirIfAbsent(pathsPath); err != nil {
		return err
	}

	swaggerResource := toSwaggerResource(resource)
	definitionsSpecPath := filepath.Join(definitionsPath, swaggerResource.URLName+".yml")
	pathsSpecPath := filepath.Join(pathsPath, swaggerResource.URLName+".yml")

	if err := executeTemplate(swaggerDefinitionsTemplate, definitionsSpecPath, swaggerResource, force); err != nil {
		return err
	}

	return executeTemplate(swaggerPathsTemplate, pathsSpecPath, swaggerResource, force)
}

var payloadExcluded = []string{"ID", "CreatedAt", "UpdatedAt"}

func toSwaggerResource(resource *StorageResource) *SwaggerResource {
	properties := make([]Property, 0, len(resource.Fields))
	payloadProperties := make([]Property, 0, len(resource.Fields))
	for _, field := range resource.Fields {
		property := fieldToProperty(field)
		properties = append(properties, property)

		if !util.ContainsString(payloadExcluded, field.Name) {
			payloadProperties = append(payloadProperties, property)
		}
	}

	swaggerResource := SwaggerResource{
		Name:              resource.Name,
		URLName:           strings.ToLower(resource.Name), // TODO: Use dash separator
		Properties:        properties,
		PayloadProperties: payloadProperties,
	}

	return &swaggerResource
}

func fieldToProperty(field Field) Property {
	propertyType := "string" // Fallback type
	required := true
	var format *string

	switch field.Type.Name {
	case "int":
		propertyType = "integer"
	case "null.Int":
		propertyType = "integer"
		required = false
	case "bool":
		propertyType = "boolean"
	case "null.Bool":
		propertyType = "boolean"
		required = false
	case "null.String":
		propertyType = "string"
		required = false
	case "types.Decimal":
		propertyType = "number"
	case "types.NullDecimal":
		propertyType = "number"
		required = false
	case "time.Time":
		format = swag.String("date-time")
	case "null.Time":
		format = swag.String("date-time")
		required = false
	}

	if strings.Contains(field.Name, "ID") {
		format = swag.String("uuid4")
	}

	return Property{
		Name:     goToJSNaming(field.Name),
		Type:     propertyType,
		Required: required,
		Format:   format,
	}
}

type HandlerField struct {
	Name             string
	Value            string
	PlaceholderValue string
}

type HandlerResource struct {
	Name   string
	Fields []HandlerField
}

type Handler struct {
	Module   string
	Package  string
	Resource *HandlerResource
}

func toHandlerResource(storageResource *StorageResource, swaggerResource *SwaggerResource) *HandlerResource {
	fields := make([]HandlerField, len(swaggerResource.Properties))
	for i, property := range swaggerResource.Properties {
		fields[i] = propertyToHandlerField(property)

		// Hack to get the proper field name.
		fields[i].Name = storageResource.Fields[i].Name
	}

	return &HandlerResource{
		Name:   swaggerResource.Name,
		Fields: fields,
	}
}

func propertyToHandlerField(property Property) HandlerField {
	placeholderValue := `swag.String("` + property.Name + `")` // Fallback placeholder

	switch property.Type {
	case "integer":
		placeholderValue = `swag.Int64(100)`
	case "boolean":
		placeholderValue = `swag.Bool(true)`
	case "number":
		placeholderValue = `swag.Float64(10.0)`
	}

	if property.Format != nil {
		switch *property.Format {
		case "date-time":
			placeholderValue = `conv.DateTime(strfmt.DateTime(time.Now()))`

		case "uuid4":
			placeholderValue = `conv.UUID4(strfmt.UUID4("1606388b-1f88-4f56-bd97-c27fbc3fe080"))`
		}
	}

	return HandlerField{
		Name:             property.Name,
		PlaceholderValue: placeholderValue,
	}
}

type handlerConfig struct {
	filePrefix string
	fileSuffix string
	template   string
}

var configuredHandlers = map[string]handlerConfig{
	"get-all": {"get_", "_list.go", getListHandlerTemplate},
	"get":     {"get_", ".go", getHandlerTemplate},
	"post":    {"post_", ".go", postHandlerTemplate},
	"put":     {"put_", ".go", putHandlerTemplate},
	"delete":  {"delete_", ".go", deleteHandlerTemplate},
}

func GenerateHandlers(resource *StorageResource, handlerBaseDir, modulePath string, methods []string, force bool) error {
	packageName := strings.ToLower(resource.Name)
	resourceBaseDir := filepath.Join(handlerBaseDir, packageName)

	if _, err := os.Stat(resourceBaseDir); os.IsNotExist(err) {
		if err := os.Mkdir(resourceBaseDir, 0755); err != nil {
			return err
		}
	}

	module, err := getModuleName(modulePath)
	if err != nil {
		return err
	}

	handler := Handler{
		Module:   module,
		Package:  packageName,
		Resource: toHandlerResource(resource, toSwaggerResource(resource)),
	}

	for _, method := range methods {
		handlerConfig, ok := configuredHandlers[method]
		if !ok {
			return fmt.Errorf("unsupported method: %s", method)
		}

		outputPath := filepath.Join(resourceBaseDir, handlerConfig.filePrefix+packageName+handlerConfig.fileSuffix)
		if err := executeTemplate(handlerConfig.template, outputPath, handler, force); err != nil {
			return err
		}
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func createDirIfAbsent(path string) error {
	if !fileExists(path) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory '%s': %w", path, err)
		}
	}

	return nil
}

func executeTemplate(templateStr, outputPath string, data interface{}, force bool) error {
	templ := template.Template{}
	if _, err := templ.Parse(templateStr); err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	if !force && fileExists(outputPath) {
		return fmt.Errorf("file '%s' already exists; call with --force to overwrite", outputPath)
	}

	file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", outputPath, err)
	}
	defer file.Close()

	if err := templ.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

func goToJSNaming(name string) string {
	// Hack to properly handle ID.
	switch name {
	case "ID":
		return "id"
	default:
		first := strings.ToLower(name[0:1])
		return first + name[1:]
	}
}

func getModuleName(absolutePathToGoMod string) (string, error) {
	dat, err := os.ReadFile(absolutePathToGoMod)

	if err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	mod := modfile.ModulePath(dat)

	return mod, nil
}
