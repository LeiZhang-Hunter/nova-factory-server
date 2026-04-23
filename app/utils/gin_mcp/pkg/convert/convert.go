package convert

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"go/ast"
	"go/parser"
	"go/token"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"nova-factory-server/app/utils/gin_mcp/pkg/types"
)

// isDebugMode returns true if Gin is in debug mode
func isDebugMode() bool {
	return gin.Mode() == gin.DebugMode
}

// ConvertRoutesToTools converts Gin routes into a list of MCP Tools and an operation map.
func ConvertRoutesToTools(routes gin.RoutesInfo, registeredSchemas map[string]types.RegisteredSchemaInfo) ([]types.Tool, map[string]types.Operation) {
	ttools := make([]types.Tool, 0)
	operations := make(map[string]types.Operation)
	usedOperationIDs := make(map[string]bool)

	if isDebugMode() {
		log.Printf("Starting conversion of %d routes to MCP tools...", len(routes))
	}

	for _, route := range routes {
		// Default operation ID generation (e.g., GET_users_id)
		operationID := strings.ToUpper(route.Method) + strings.ReplaceAll(strings.ReplaceAll(route.Path, "/", "_"), ":", "")

		filePath, handlerName := getRouteHandlerInfo(route)
		// Parse handler function comments
		handlerDoc, _ := parseHandlerComments(filePath, handlerName)

		// Override with custom @operationId if present
		if handlerDoc != nil && handlerDoc.OperationID != "" {
			customOpID := handlerDoc.OperationID
			if usedOperationIDs[customOpID] {
				// Always log duplicate warnings, not just in debug mode
				log.Printf("Warning: Duplicate @operationId '%s' for route %s %s. Skipping this route to maintain consistency. First declaration wins.", customOpID, route.Method, route.Path)
				continue // Skip this route entirely
			}
			operationID = customOpID
		}

		// Check for duplicates even with default IDs (shouldn't happen but be defensive)
		if usedOperationIDs[operationID] {
			log.Printf("Warning: Duplicate operation ID '%s' for route %s %s. Skipping this route.", operationID, route.Method, route.Path)
			continue
		}

		// Track used operation IDs
		usedOperationIDs[operationID] = true

		if isDebugMode() {
			log.Printf("Processing route: %s %s -> OpID: %s", route.Method, route.Path, operationID)
		}

		// Generate description information
		description := fmt.Sprintf("Handler for %s %s", route.Method, route.Path)
		if handlerDoc != nil {
			if handlerDoc.Summary != "" {
				description = handlerDoc.Summary
			}
			if handlerDoc.Description != "" {
				description += "\n\n" + handlerDoc.Description
			}
		}

		// Generate schema for the tool's input
		inputSchema := generateInputSchema(route, registeredSchemas)

		// Add parameter descriptions to schema if available
		if handlerDoc != nil && len(handlerDoc.Params) > 0 {
			for paramName, paramDesc := range handlerDoc.Params {
				if prop, ok := inputSchema.Properties[paramName]; ok {
					prop.Description = paramDesc
				}
			}
		}

		// Extract tags from handler doc
		var tags []string
		if handlerDoc != nil && len(handlerDoc.Tags) > 0 {
			tags = handlerDoc.Tags
		}

		tool := types.Tool{
			Name:        operationID,
			Description: description,
			InputSchema: inputSchema,
			Tags:        tags,
		}

		ttools = append(ttools, tool)
		operations[operationID] = types.Operation{
			Method: route.Method,
			Path:   route.Path,
			Tags:   tags,
		}
	}

	if isDebugMode() {
		log.Printf("Finished route conversion. Generated %d tools.", len(ttools))
	}

	return ttools, operations
}

// PathParamRegex is used to find path parameters like :id or *action
var PathParamRegex = regexp.MustCompile(`[:\*]([a-zA-Z0-9_]+)`)

// generateInputSchema creates the JSON schema for the tool's input parameters.
// This is a simplified version using basic reflection and not an external library.
func generateInputSchema(route gin.RouteInfo, registeredSchemas map[string]types.RegisteredSchemaInfo) *types.JSONSchema {
	// Base schema structure
	schema := &types.JSONSchema{
		Type:       "object",
		Properties: make(map[string]*types.JSONSchema),
		Required:   make([]string, 0),
	}
	properties := schema.Properties
	required := schema.Required

	// 1. Extract Path Parameters
	matches := PathParamRegex.FindAllStringSubmatch(route.Path, -1)
	for _, match := range matches {
		if len(match) > 1 {
			paramName := match[1]
			properties[paramName] = &types.JSONSchema{
				Type:        "string",
				Description: fmt.Sprintf("Path parameter: %s", paramName),
			}
			required = append(required, paramName) // Path params are always required
		}
	}

	// 2. Incorporate Registered Query and Body Types
	schemaKey := route.Method + " " + route.Path
	if schemaInfo, exists := registeredSchemas[schemaKey]; exists {
		if isDebugMode() {
			log.Printf("Using registered schema for %s", schemaKey)
		}

		// Reflect Query Parameters (if applicable for method and type exists)
		if (route.Method == "GET" || route.Method == "DELETE") && schemaInfo.QueryType != nil {
			reflectAndAddProperties(schemaInfo.QueryType, properties, &required, "query")
		}

		// Reflect Body Parameters (if applicable for method and type exists)
		if (route.Method == "POST" || route.Method == "PUT" || route.Method == "PATCH") && schemaInfo.BodyType != nil {
			reflectAndAddProperties(schemaInfo.BodyType, properties, &required, "body")
		}
	}

	// Update the required slice in the main schema
	schema.Required = required

	// If no properties were added (beyond path params), handle appropriately.
	// Depending on the spec, an empty properties object might be required.
	// if len(properties) == 0 { // Keep properties map even if empty
	// 	// Return schema with empty properties
	// 	return schema
	// }

	return schema
}

// reflectAndAddProperties uses basic reflection to add properties to the schema.
func reflectAndAddProperties(goType interface{}, properties map[string]*types.JSONSchema, required *[]string, source string) {
	if goType == nil {
		return // Handle nil input gracefully
	}
	t := types.ReflectType(reflect.TypeOf(goType)) // Use helper from types pkg
	if t == nil || t.Kind() != reflect.Struct {
		if isDebugMode() {
			log.Printf("Skipping schema generation for non-struct type: %v (%s)", reflect.TypeOf(goType), source)
		}
		return
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		formTag := field.Tag.Get("form")             // Used for query params often
		jsonschemaTag := field.Tag.Get("jsonschema") // Basic support

		fieldName := field.Name // Default to field name
		ignoreField := false

		// Determine field name from tags (prefer json, then form)
		if jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] == "-" {
				ignoreField = true
			} else {
				fieldName = parts[0]
			}
			if len(parts) > 1 && parts[1] == "omitempty" {
				// omitempty = true // Variable removed
			}
		} else if formTag != "" {
			parts := strings.Split(formTag, ",")
			if parts[0] == "-" {
				ignoreField = true
			} else {
				fieldName = parts[0]
			}
			// form tag doesn't typically have omitempty in the same way
		}

		if ignoreField || !field.IsExported() {
			continue
		}

		propSchema := schemaFromReflectType(field.Type, map[reflect.Type]bool{})

		// Basic 'required' and 'description' handling from jsonschema tag
		isRequired := false // Default to not required
		if jsonschemaTag != "" {
			parts := strings.Split(jsonschemaTag, ",")
			for _, part := range parts {
				trimmed := strings.TrimSpace(part)
				if trimmed == "required" {
					isRequired = true
				} else if strings.HasPrefix(trimmed, "description=") {
					propSchema.Description = strings.TrimPrefix(trimmed, "description=")
				} else if strings.HasPrefix(trimmed, "minimum=") {
					if minimum, ok := parseFloatTagValue(strings.TrimPrefix(trimmed, "minimum=")); ok {
						propSchema.Minimum = &minimum
					}
				} else if strings.HasPrefix(trimmed, "maximum=") {
					if maximum, ok := parseFloatTagValue(strings.TrimPrefix(trimmed, "maximum=")); ok {
						propSchema.Maximum = &maximum
					}
				} else if strings.HasPrefix(trimmed, "enum=") {
					enumValues := parseEnumTagValues(strings.TrimPrefix(trimmed, "enum="), field.Type)
					if len(enumValues) > 0 {
						propSchema.Enum = enumValues
					}
				}
			}
		}

		// Add to properties map
		properties[fieldName] = propSchema

		// Add to required list if necessary
		if isRequired {
			*required = append(*required, fieldName)
		}
	}
}

func schemaFromReflectType(fieldType reflect.Type, visiting map[reflect.Type]bool) *types.JSONSchema {
	if fieldType == nil {
		return &types.JSONSchema{Type: "string"}
	}

	for fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}

	if fieldType == nil {
		return &types.JSONSchema{Type: "string"}
	}

	if fieldType == reflect.TypeOf(time.Time{}) {
		return &types.JSONSchema{Type: "string"}
	}

	switch fieldType.Kind() {
	case reflect.String:
		return &types.JSONSchema{Type: "string"}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &types.JSONSchema{Type: "integer"}
	case reflect.Float32, reflect.Float64:
		return &types.JSONSchema{Type: "number"}
	case reflect.Bool:
		return &types.JSONSchema{Type: "boolean"}
	case reflect.Slice, reflect.Array:
		return &types.JSONSchema{
			Type:  "array",
			Items: schemaFromReflectType(fieldType.Elem(), visiting),
		}
	case reflect.Map:
		return &types.JSONSchema{
			Type:                 "object",
			AdditionalProperties: schemaFromReflectType(fieldType.Elem(), visiting),
		}
	case reflect.Struct:
		schema := &types.JSONSchema{
			Type:       "object",
			Properties: make(map[string]*types.JSONSchema),
			Required:   make([]string, 0),
		}
		if visiting[fieldType] {
			return schema
		}
		visiting[fieldType] = true
		reflectTypeFields(fieldType, schema.Properties, &schema.Required, visiting)
		delete(visiting, fieldType)
		return schema
	case reflect.Interface:
		return &types.JSONSchema{Type: "object"}
	default:
		return &types.JSONSchema{Type: "string"}
	}
}

func reflectTypeFields(t reflect.Type, properties map[string]*types.JSONSchema, required *[]string, visiting map[reflect.Type]bool) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		formTag := field.Tag.Get("form")
		jsonschemaTag := field.Tag.Get("jsonschema")
		fieldName := field.Name
		ignoreField := false

		if jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] == "-" {
				ignoreField = true
			} else if parts[0] != "" {
				fieldName = parts[0]
			}
		} else if formTag != "" {
			parts := strings.Split(formTag, ",")
			if parts[0] == "-" {
				ignoreField = true
			} else if parts[0] != "" {
				fieldName = parts[0]
			}
		}

		if ignoreField {
			continue
		}

		if field.Anonymous && jsonTag == "" && formTag == "" {
			embeddedType := field.Type
			for embeddedType.Kind() == reflect.Ptr {
				embeddedType = embeddedType.Elem()
			}
			if embeddedType.Kind() == reflect.Struct {
				reflectTypeFields(embeddedType, properties, required, visiting)
				continue
			}
		}

		propSchema := schemaFromReflectType(field.Type, visiting)
		isRequired := false
		if jsonschemaTag != "" {
			parts := strings.Split(jsonschemaTag, ",")
			for _, part := range parts {
				trimmed := strings.TrimSpace(part)
				if trimmed == "required" {
					isRequired = true
				} else if strings.HasPrefix(trimmed, "description=") {
					propSchema.Description = strings.TrimPrefix(trimmed, "description=")
				} else if strings.HasPrefix(trimmed, "minimum=") {
					if minimum, ok := parseFloatTagValue(strings.TrimPrefix(trimmed, "minimum=")); ok {
						propSchema.Minimum = &minimum
					}
				} else if strings.HasPrefix(trimmed, "maximum=") {
					if maximum, ok := parseFloatTagValue(strings.TrimPrefix(trimmed, "maximum=")); ok {
						propSchema.Maximum = &maximum
					}
				} else if strings.HasPrefix(trimmed, "enum=") {
					enumValues := parseEnumTagValues(strings.TrimPrefix(trimmed, "enum="), field.Type)
					if len(enumValues) > 0 {
						propSchema.Enum = enumValues
					}
				}
			}
		}

		properties[fieldName] = propSchema
		if isRequired {
			*required = append(*required, fieldName)
		}
	}
}

func parseFloatTagValue(raw string) (float64, bool) {
	value, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		return 0, false
	}
	return value, true
}

func parseEnumTagValues(raw string, fieldType reflect.Type) []any {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	for fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}

	items := strings.Split(raw, "|")
	values := make([]any, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}

		switch fieldType.Kind() {
		case reflect.Bool:
			parsed, err := strconv.ParseBool(item)
			if err == nil {
				values = append(values, parsed)
				continue
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			parsed, err := strconv.ParseInt(item, 10, 64)
			if err == nil {
				values = append(values, parsed)
				continue
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			parsed, err := strconv.ParseUint(item, 10, 64)
			if err == nil {
				values = append(values, parsed)
				continue
			}
		case reflect.Float32, reflect.Float64:
			parsed, err := strconv.ParseFloat(item, 64)
			if err == nil {
				values = append(values, parsed)
				continue
			}
		}

		values = append(values, item)
	}
	return values
}

// HandlerDoc stores function documentation
type HandlerDoc struct {
	Summary     string
	Description string
	Params      map[string]string
	Returns     string
	Tags        []string
	OperationID string
}

// parseHandlerComments parses function documentation from source code
func parseHandlerComments(filePath string, handlerName string) (*HandlerDoc, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		log.Printf("Failed to parse file %s: %v", filePath, err)
		return nil, err
	}

	var doc *HandlerDoc

	// Iterate through top-level declarations
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if fn.Name.String() == handlerName {
				doc = &HandlerDoc{
					Params: make(map[string]string),
				}
				if fn.Doc != nil {
					// Parse comments
					lines := strings.Split(fn.Doc.Text(), "\n")
					for _, line := range lines {
						line = strings.ToLower(strings.TrimSpace(line))
						switch {
						case strings.HasPrefix(line, "@summary"):
							doc.Summary = strings.TrimSpace(strings.TrimPrefix(line, "@summary"))
						case strings.HasPrefix(line, "@description"):
							doc.Description = strings.TrimSpace(strings.TrimPrefix(line, "@description"))
						case strings.HasPrefix(line, "@param"):
							paramText := strings.TrimSpace(strings.TrimPrefix(line, "@param"))
							parts := strings.SplitN(paramText, " ", 2)
							if len(parts) == 2 {
								paramName := strings.TrimSpace(parts[0])
								paramDesc := strings.TrimSpace(parts[1])
								doc.Params[paramName] = paramDesc
							}
						case strings.HasPrefix(line, "@return"):
							doc.Returns = strings.TrimSpace(strings.TrimPrefix(line, "@return"))
						case strings.HasPrefix(line, "@tags"):
							tagsText := strings.TrimSpace(strings.TrimPrefix(line, "@tags"))
							// Split on spaces and commas, trim whitespace, ignore empty entries
							var tags []string
							for _, sep := range []string{",", " "} {
								parts := strings.Split(tagsText, sep)
								for _, part := range parts {
									trimmed := strings.TrimSpace(part)
									if trimmed != "" && !contains(tags, trimmed) {
										tags = append(tags, trimmed)
									}
								}
								// After first pass with commas, rejoin and split by spaces
								if sep == "," {
									tagsText = strings.Join(tags, " ")
									tags = []string{}
								}
							}
							doc.Tags = tags
						case strings.HasPrefix(line, "@operationId"):
							opID := strings.TrimSpace(strings.TrimPrefix(line, "@operationId"))
							if opID != "" && doc.OperationID == "" {
								doc.OperationID = opID
							}
						}
					}
				}
			}
		}
	}

	return doc, nil
}

func getHandlerInfo(handler gin.HandlerFunc) (string, string) {
	// Get function reflection value
	v := reflect.ValueOf(handler)
	ptr := v.Pointer()

	// Get function name and info
	funcInfo := runtime.FuncForPC(ptr)
	if funcInfo == nil {
		return "", ""
	}

	fullName := funcInfo.Name()
	filePath, _ := funcInfo.FileLine(ptr)
	shortName := normalizeHandlerName(fullName)

	return filePath, shortName
}

func getRouteHandlerInfo(route gin.RouteInfo) (string, string) {
	filePath, handlerName := getHandlerInfo(route.HandlerFunc)
	fallbackPath, fallbackName := resolveHandlerInfo(route.Handler)

	if fallbackName != "" {
		handlerName = fallbackName
	}
	if !isSourceFilePath(filePath) && fallbackPath != "" {
		filePath = fallbackPath
	}

	return filePath, handlerName
}

func isSourceFilePath(filePath string) bool {
	return filePath != "" && filePath != "<autogenerated>" && strings.HasSuffix(filePath, ".go")
}

type handlerIdentity struct {
	PackagePath  string
	ReceiverName string
	HandlerName  string
}

func resolveHandlerInfo(fullName string) (string, string) {
	identity := parseHandlerIdentity(fullName)
	if identity.HandlerName == "" {
		return "", ""
	}

	packageDir, err := resolvePackageDir(identity.PackagePath)
	if err != nil {
		return "", identity.HandlerName
	}

	filePath, err := findHandlerSourceFile(packageDir, identity.HandlerName, identity.ReceiverName)
	if err != nil {
		return "", identity.HandlerName
	}

	return filePath, identity.HandlerName
}

func parseHandlerIdentity(fullName string) handlerIdentity {
	fullName = strings.TrimSpace(strings.TrimSuffix(fullName, "-fm"))
	if fullName == "" {
		return handlerIdentity{}
	}

	if idx := strings.Index(fullName, ".func"); idx >= 0 {
		fullName = fullName[:idx]
	}

	if idx := strings.LastIndex(fullName, ")."); idx >= 0 {
		handlerName := normalizeHandlerName(fullName[idx+2:])
		prefix := fullName[:idx+1]
		packageIdx := strings.LastIndex(prefix, ".(")
		if packageIdx < 0 {
			return handlerIdentity{HandlerName: handlerName}
		}

		receiverName := normalizeReceiverName(prefix[packageIdx+2 : len(prefix)-1])
		return handlerIdentity{
			PackagePath:  prefix[:packageIdx],
			ReceiverName: receiverName,
			HandlerName:  handlerName,
		}
	}

	lastDot := strings.LastIndex(fullName, ".")
	if lastDot < 0 {
		return handlerIdentity{HandlerName: normalizeHandlerName(fullName)}
	}

	handlerName := normalizeHandlerName(fullName[lastDot+1:])
	prefix := fullName[:lastDot]
	lastSlash := strings.LastIndex(prefix, "/")
	suffix := prefix
	if lastSlash >= 0 {
		suffix = prefix[lastSlash+1:]
	}
	if receiverIdx := strings.LastIndex(suffix, "."); receiverIdx >= 0 {
		packagePath := prefix[:len(prefix)-len(suffix)] + suffix[:receiverIdx]
		return handlerIdentity{
			PackagePath:  packagePath,
			ReceiverName: normalizeReceiverName(suffix[receiverIdx+1:]),
			HandlerName:  handlerName,
		}
	}

	return handlerIdentity{
		PackagePath: prefix,
		HandlerName: handlerName,
	}
}

func normalizeHandlerName(fullName string) string {
	if fullName == "" {
		return ""
	}

	name := strings.TrimSuffix(strings.TrimSpace(fullName), "-fm")
	if idx := strings.LastIndex(name, "."); idx >= 0 {
		name = name[idx+1:]
	}
	if idx := strings.Index(name, ".func"); idx >= 0 {
		name = name[:idx]
	}

	return name
}

func normalizeReceiverName(receiver string) string {
	receiver = strings.TrimSpace(receiver)
	receiver = strings.TrimPrefix(receiver, "*")
	if idx := strings.Index(receiver, "["); idx >= 0 {
		receiver = receiver[:idx]
	}
	return receiver
}

func resolvePackageDir(packagePath string) (string, error) {
	moduleName, moduleRoot, err := getModuleInfo()
	if err != nil {
		return "", err
	}
	if packagePath == "" || !strings.HasPrefix(packagePath, moduleName) {
		return "", fmt.Errorf("package %s is outside current module %s", packagePath, moduleName)
	}

	relativePath := strings.TrimPrefix(strings.TrimPrefix(packagePath, moduleName), "/")
	return filepath.Join(moduleRoot, filepath.FromSlash(relativePath)), nil
}

func getModuleInfo() (string, string, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", "", fmt.Errorf("failed to determine current file path")
	}

	moduleRoot, err := findModuleRoot(filepath.Dir(currentFile))
	if err != nil {
		return "", "", err
	}

	goModContent, err := os.ReadFile(filepath.Join(moduleRoot, "go.mod"))
	if err != nil {
		return "", "", err
	}

	for _, line := range strings.Split(string(goModContent), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			moduleName := strings.TrimSpace(strings.TrimPrefix(line, "module "))
			if moduleName != "" {
				return moduleName, moduleRoot, nil
			}
		}
	}

	return "", "", fmt.Errorf("module name not found in go.mod")
}

func findModuleRoot(startDir string) (string, error) {
	dir := startDir
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s", startDir)
		}
		dir = parent
	}
}

func findHandlerSourceFile(dir string, handlerName string, receiverName string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	fset := token.NewFileSet()
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())
		file, err := parser.ParseFile(fset, filePath, nil, 0)
		if err != nil {
			continue
		}

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Name.Name != handlerName {
				continue
			}
			if receiverName != "" && !receiverMatches(fn, receiverName) {
				continue
			}
			if receiverName == "" && fn.Recv != nil {
				continue
			}
			return filePath, nil
		}
	}

	return "", fmt.Errorf("handler %s not found in %s", handlerName, dir)
}

func receiverMatches(fn *ast.FuncDecl, receiverName string) bool {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return false
	}

	return receiverTypeName(fn.Recv.List[0].Type) == receiverName
}

func receiverTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return receiverTypeName(t.X)
	case *ast.IndexExpr:
		return receiverTypeName(t.X)
	case *ast.IndexListExpr:
		return receiverTypeName(t.X)
	default:
		return ""
	}
}

// contains checks if a string slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
