package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const outputFile = "app/utils/gin_mcp/pkg/convert/handler_docs_gen.go"

type handlerIdentity struct {
	PackagePath  string
	ReceiverName string
	HandlerName  string
}

type handlerDoc struct {
	Summary     string
	Description string
	Params      map[string]string
	Returns     string
	Tags        []string
	OperationID string
}

func main() {
	root, err := findModuleRoot()
	if err != nil {
		fatal(err)
	}

	moduleName, err := readModuleName(filepath.Join(root, "go.mod"))
	if err != nil {
		fatal(err)
	}

	entries := make(map[handlerIdentity]handlerDoc)
	appDir := filepath.Join(root, "app")
	err = walkGoFiles(appDir, func(path string) error {
		if err := extractFromFile(root, moduleName, path, entries); err != nil {
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", path, err)
		}
		return nil
	})
	if err != nil {
		fatal(err)
	}

	outPath := filepath.Join(root, outputFile)
	if err := writeGenerated(outPath, entries); err != nil {
		fatal(err)
	}

	fmt.Printf("Generated %s with %d handlers\n", outPath, len(entries))
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func walkGoFiles(root string, visit func(path string) error) error {
	return walkGoFilesInDir(root, make(map[string]struct{}), visit)
}

func walkGoFilesInDir(dir string, visiting map[string]struct{}, visit func(path string) error) error {
	realDir, err := filepath.EvalSymlinks(dir)
	if err == nil {
		realDir, err = filepath.Abs(realDir)
		if err == nil {
			if _, ok := visiting[realDir]; ok {
				return nil
			}
			visiting[realDir] = struct{}{}
			defer delete(visiting, realDir)
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	for _, entry := range entries {
		name := entry.Name()
		path := filepath.Join(dir, name)
		if entry.IsDir() {
			if shouldSkipDir(name) {
				continue
			}
			if err := walkGoFilesInDir(path, visiting, visit); err != nil {
				return err
			}
			continue
		}
		if entry.Type()&os.ModeSymlink != 0 {
			info, err := os.Stat(path)
			if err == nil && info.IsDir() {
				if shouldSkipDir(name) {
					continue
				}
				if err := walkGoFilesInDir(path, visiting, visit); err != nil {
					return err
				}
				continue
			}
		}
		if strings.HasSuffix(path, ".go") {
			if err := visit(path); err != nil {
				return err
			}
		}
	}
	return nil
}

func shouldSkipDir(name string) bool {
	return name == ".git" || name == "vendor"
}

func findModuleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}

func readModuleName(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			moduleName := strings.TrimSpace(strings.TrimPrefix(line, "module "))
			if moduleName != "" {
				return moduleName, nil
			}
		}
	}
	return "", fmt.Errorf("module name not found in %s", path)
}

func extractFromFile(root, moduleName, filePath string, out map[handlerIdentity]handlerDoc) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	relDir, err := filepath.Rel(root, filepath.Dir(filePath))
	if err != nil {
		return err
	}
	packagePath := moduleName + "/" + filepath.ToSlash(relDir)

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Doc == nil {
			continue
		}

		doc := parseDoc(fn.Doc.Text())
		if doc.isEmpty() {
			continue
		}

		identity := handlerIdentity{
			PackagePath:  packagePath,
			ReceiverName: receiverName(fn),
			HandlerName:  fn.Name.Name,
		}
		out[identity] = doc
	}
	return nil
}

func receiverName(fn *ast.FuncDecl) string {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return ""
	}
	return receiverTypeName(fn.Recv.List[0].Type)
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

func parseDoc(text string) handlerDoc {
	doc := handlerDoc{Params: make(map[string]string)}
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		lower := strings.ToLower(line)
		switch {
		case strings.HasPrefix(lower, "@summary"):
			doc.Summary = strings.TrimSpace(line[len("@summary"):])
		case strings.HasPrefix(lower, "@description"):
			doc.Description = strings.TrimSpace(line[len("@description"):])
		case strings.HasPrefix(lower, "@param"):
			paramText := strings.TrimSpace(line[len("@param"):])
			parts := strings.SplitN(paramText, " ", 2)
			if len(parts) == 2 {
				doc.Params[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		case strings.HasPrefix(lower, "@return"):
			doc.Returns = strings.TrimSpace(line[len("@return"):])
		case strings.HasPrefix(lower, "@tags"):
			doc.Tags = parseTags(strings.TrimSpace(line[len("@tags"):]))
		case strings.HasPrefix(lower, "@operationid"):
			operationID := strings.TrimSpace(line[len("@operationid"):])
			if operationID != "" && doc.OperationID == "" {
				doc.OperationID = operationID
			}
		}
	}
	return doc
}

func parseTags(raw string) []string {
	if raw == "" {
		return nil
	}
	fields := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t'
	})
	seen := make(map[string]bool, len(fields))
	tags := make([]string, 0, len(fields))
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" || seen[field] {
			continue
		}
		seen[field] = true
		tags = append(tags, field)
	}
	return tags
}

func (d handlerDoc) isEmpty() bool {
	return d.Summary == "" && d.Description == "" && len(d.Params) == 0 &&
		d.Returns == "" && len(d.Tags) == 0 && d.OperationID == ""
}

func writeGenerated(path string, entries map[handlerIdentity]handlerDoc) error {
	keys := make([]handlerIdentity, 0, len(entries))
	for key := range entries {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].PackagePath != keys[j].PackagePath {
			return keys[i].PackagePath < keys[j].PackagePath
		}
		if keys[i].ReceiverName != keys[j].ReceiverName {
			return keys[i].ReceiverName < keys[j].ReceiverName
		}
		return keys[i].HandlerName < keys[j].HandlerName
	})

	var buf bytes.Buffer
	buf.WriteString("// Code generated by handler-doc-gen; DO NOT EDIT.\n\n")
	buf.WriteString("package convert\n\n")
	buf.WriteString("var handlerDocByIdentity = map[handlerIdentity]HandlerDoc{\n")
	for _, key := range keys {
		doc := entries[key]
		fmt.Fprintf(&buf, "\t{PackagePath: %q, ReceiverName: %q, HandlerName: %q}: {\n", key.PackagePath, key.ReceiverName, key.HandlerName)
		fmt.Fprintf(&buf, "\t\tSummary: %q,\n", doc.Summary)
		fmt.Fprintf(&buf, "\t\tDescription: %q,\n", doc.Description)
		fmt.Fprintf(&buf, "\t\tReturns: %q,\n", doc.Returns)
		if len(doc.Params) > 0 {
			paramKeys := make([]string, 0, len(doc.Params))
			for name := range doc.Params {
				paramKeys = append(paramKeys, name)
			}
			sort.Strings(paramKeys)
			buf.WriteString("\t\tParams: map[string]string{\n")
			for _, name := range paramKeys {
				fmt.Fprintf(&buf, "\t\t\t%q: %q,\n", name, doc.Params[name])
			}
			buf.WriteString("\t\t},\n")
		}
		if len(doc.Tags) > 0 {
			buf.WriteString("\t\tTags: []string{")
			for i, tag := range doc.Tags {
				if i > 0 {
					buf.WriteString(", ")
				}
				fmt.Fprintf(&buf, "%q", tag)
			}
			buf.WriteString("},\n")
		}
		fmt.Fprintf(&buf, "\t\tOperationID: %q,\n", doc.OperationID)
		buf.WriteString("\t},\n")
	}
	buf.WriteString("}\n")

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	return os.WriteFile(path, formatted, 0644)
}
