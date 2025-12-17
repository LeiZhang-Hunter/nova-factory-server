package convert

// Config presents Gen configurations.
type Config struct {
	Debugger Debugger

	// SearchDir the swag would parse,comma separated if multiple
	SearchDir string

	// excludes dirs and files in SearchDir,comma separated
	Excludes string

	// outputs only specific extension
	ParseExtension string

	// OutputDir represents the output directory for all the generated files
	OutputDir string

	// OutputTypes define types of files which should be generated
	OutputTypes []string

	// MainAPIFile the Go file path in which 'swagger general API Info' is written
	MainAPIFile string

	// PropNamingStrategy represents property naming strategy like snake case,camel case,pascal case
	PropNamingStrategy string

	// MarkdownFilesDir used to find markdown files, which can be used for tag descriptions
	MarkdownFilesDir string

	// CodeExampleFilesDir used to find code example files, which can be used for x-codeSamples
	CodeExampleFilesDir string

	// InstanceName is used to get distinct names for different swagger documents in the
	// same project. The default value is "swagger".
	InstanceName string

	// ParseDepth dependency parse depth
	ParseDepth int

	// ParseVendor whether swag should be parse vendor folder
	ParseVendor bool

	// ParseDependencies whether swag should be parse outside dependency folder: 0 none, 1 models, 2 operations, 3 all
	ParseDependency int

	// UseStructNames stick to the struct name instead of those ugly full-path names
	UseStructNames bool

	// ParseInternal whether swag should parse internal packages
	ParseInternal bool

	// Strict whether swag should error or warn when it detects cases which are most likely user errors
	Strict bool

	// GeneratedTime whether swag should generate the timestamp at the top of docs.go
	GeneratedTime bool

	// RequiredByDefault set validation required for all fields by default
	RequiredByDefault bool

	// OverridesFile defines global type overrides.
	OverridesFile string

	// ParseGoList whether swag use go list to parse dependency
	ParseGoList bool

	// include only tags mentioned when searching, comma separated
	Tags string

	// LeftTemplateDelim defines the left delimiter for the template generation
	LeftTemplateDelim string

	// RightTemplateDelim defines the right delimiter for the template generation
	RightTemplateDelim string

	// PackageName defines package name of generated `docs.go`
	PackageName string

	// CollectionFormat set default collection format
	CollectionFormat string

	// Parse only packages whose import path match the given prefix, comma separated
	PackagePrefix string

	// State set host state
	State string

	// ParseFuncBody whether swag should parse api info inside of funcs
	ParseFuncBody bool

	// ParseGoPackages whether swag use golang.org/x/tools/go/packages to parse source.
	ParseGoPackages bool
}
