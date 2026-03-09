package ingestion

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/sirupsen/logrus"
)

// GoFileParser handles Go source file parsing
type GoFileParser struct {
	logger *logrus.Logger
}

// NewGoFileParser creates a new Go file parser
func NewGoFileParser(logger *logrus.Logger) *GoFileParser {
	return &GoFileParser{
		logger: logger,
	}
}

// ParseGoFile extracts structured content from a Go source file
func (g *GoFileParser) ParseGoFile(ctx context.Context, content []byte) ([]GoCodeElement, error) {
	fset := token.NewFileSet()
	
	// Parse the Go source file
	file, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Go file: %w", err)
	}

	var elements []GoCodeElement

	// Extract package information
	if file.Name != nil {
		elements = append(elements, GoCodeElement{
			Type:        "package",
			Name:        file.Name.Name,
			Content:     fmt.Sprintf("package %s", file.Name.Name),
			Line:        fset.Position(file.Pos()).Line,
			DocComment:  "",
		})
	}

	// Extract imports
	for _, imp := range file.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		content := fmt.Sprintf("import %s", imp.Path.Value)
		if imp.Name != nil {
			content = fmt.Sprintf("import %s %s", imp.Name.Name, imp.Path.Value)
		}

		elements = append(elements, GoCodeElement{
			Type:        "import",
			Name:        importPath,
			Content:     content,
			Line:        fset.Position(imp.Pos()).Line,
			DocComment:  "",
		})
	}

	// Extract declarations
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			// Handle const, var, type declarations
			if d.Tok == token.CONST {
				elements = append(elements, g.extractConstDecl(fset, d)...)
			} else if d.Tok == token.VAR {
				elements = append(elements, g.extractVarDecl(fset, d)...)
			} else if d.Tok == token.TYPE {
				elements = append(elements, g.extractTypeDecl(fset, d)...)
			}

		case *ast.FuncDecl:
			// Handle function declarations
			elements = append(elements, g.extractFuncDecl(fset, d))
		}
	}

	// Extract comments
	for _, commentGroup := range file.Comments {
		commentText := strings.TrimSpace(commentGroup.Text())
		if commentText != "" {
			elements = append(elements, GoCodeElement{
				Type:        "comment",
				Name:        "",
				Content:     commentText,
				Line:        fset.Position(commentGroup.Pos()).Line,
				DocComment:  "",
			})
		}
	}

	g.logger.Infof("Extracted %d Go code elements", len(elements))
	return elements, nil
}

// GetGoFileMetadata extracts metadata from Go source file
func (g *GoFileParser) GetGoFileMetadata(ctx context.Context, content []byte) (map[string]interface{}, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Go file: %w", err)
	}

	metadata := make(map[string]interface{})

	// Basic file info
	metadata["package"] = file.Name.Name
	metadata["content_length"] = len(content)
	metadata["line_count"] = fset.Position(file.End()).Line

	// Count various elements
	var functions, types, constants, variables, imports int

	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok == token.CONST {
				constants++
			} else if d.Tok == token.VAR {
				variables++
			} else if d.Tok == token.TYPE {
				types++
			}
		case *ast.FuncDecl:
			functions++
		}
	}

	imports = len(file.Imports)

	metadata["functions"] = functions
	metadata["types"] = types
	metadata["constants"] = constants
	metadata["variables"] = variables
	metadata["imports"] = imports

	// Check for main package
	if file.Name.Name == "main" {
		metadata["is_main_package"] = true
		
		// Check for main function
		for _, decl := range file.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Name.Name == "main" {
				metadata["has_main_function"] = true
				break
			}
		}
	}

	return metadata, nil
}

// ValidateGoFile checks if the content is valid Go code
func (g *GoFileParser) ValidateGoFile(content []byte) error {
	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("invalid Go code: %w", err)
	}
	return nil
}

// Helper functions for extracting specific declaration types

func (g *GoFileParser) extractConstDecl(fset *token.FileSet, decl *ast.GenDecl) []GoCodeElement {
	var elements []GoCodeElement
	
	docComment := ""
	if decl.Doc != nil {
		docComment = strings.TrimSpace(decl.Doc.Text())
	}

	for _, spec := range decl.Specs {
		if valueSpec, ok := spec.(*ast.ValueSpec); ok {
			for i, name := range valueSpec.Names {
				content := fmt.Sprintf("const %s", name.Name)
				if i < len(valueSpec.Values) {
					content += fmt.Sprintf(" = %s", g.formatNode(valueSpec.Values[i]))
				}

				elements = append(elements, GoCodeElement{
					Type:       "constant",
					Name:       name.Name,
					Content:    content,
					Line:       fset.Position(name.Pos()).Line,
					DocComment: docComment,
				})
			}
		}
	}

	return elements
}

func (g *GoFileParser) extractVarDecl(fset *token.FileSet, decl *ast.GenDecl) []GoCodeElement {
	var elements []GoCodeElement
	
	docComment := ""
	if decl.Doc != nil {
		docComment = strings.TrimSpace(decl.Doc.Text())
	}

	for _, spec := range decl.Specs {
		if valueSpec, ok := spec.(*ast.ValueSpec); ok {
			for i, name := range valueSpec.Names {
				content := fmt.Sprintf("var %s", name.Name)
				if valueSpec.Type != nil {
					content += fmt.Sprintf(" %s", g.formatNode(valueSpec.Type))
				}
				if i < len(valueSpec.Values) {
					content += fmt.Sprintf(" = %s", g.formatNode(valueSpec.Values[i]))
				}

				elements = append(elements, GoCodeElement{
					Type:       "variable",
					Name:       name.Name,
					Content:    content,
					Line:       fset.Position(name.Pos()).Line,
					DocComment: docComment,
				})
			}
		}
	}

	return elements
}

func (g *GoFileParser) extractTypeDecl(fset *token.FileSet, decl *ast.GenDecl) []GoCodeElement {
	var elements []GoCodeElement
	
	docComment := ""
	if decl.Doc != nil {
		docComment = strings.TrimSpace(decl.Doc.Text())
	}

	for _, spec := range decl.Specs {
		if typeSpec, ok := spec.(*ast.TypeSpec); ok {
			content := fmt.Sprintf("type %s %s", typeSpec.Name.Name, g.formatNode(typeSpec.Type))

			elements = append(elements, GoCodeElement{
				Type:       "type",
				Name:       typeSpec.Name.Name,
				Content:    content,
				Line:       fset.Position(typeSpec.Pos()).Line,
				DocComment: docComment,
			})
		}
	}

	return elements
}

func (g *GoFileParser) extractFuncDecl(fset *token.FileSet, decl *ast.FuncDecl) GoCodeElement {
	docComment := ""
	if decl.Doc != nil {
		docComment = strings.TrimSpace(decl.Doc.Text())
	}

	content := fmt.Sprintf("func %s", decl.Name.Name)
	
	// Add receiver if it's a method
	if decl.Recv != nil {
		content += fmt.Sprintf(" %s", g.formatNode(decl.Recv))
	}
	
	// Add parameters
	content += fmt.Sprintf("(%s)", g.formatFieldList(decl.Type.Params))
	
	// Add return values
	if decl.Type.Results != nil {
		content += fmt.Sprintf(" %s", g.formatFieldList(decl.Type.Results))
	}

	return GoCodeElement{
		Type:       "function",
		Name:       decl.Name.Name,
		Content:    content,
		Line:       fset.Position(decl.Pos()).Line,
		DocComment: docComment,
	}
}

// Helper function to format AST nodes as strings
func (g *GoFileParser) formatNode(node ast.Node) string {
	if node == nil {
		return ""
	}
	// This is a simplified implementation
	// In a production system, you might use go/printer for proper formatting
	return fmt.Sprintf("%T", node)
}

func (g *GoFileParser) formatFieldList(fields *ast.FieldList) string {
	if fields == nil || len(fields.List) == 0 {
		return ""
	}
	
	var parts []string
	for _, field := range fields.List {
		for _, name := range field.Names {
			parts = append(parts, fmt.Sprintf("%s %s", name.Name, g.formatNode(field.Type)))
		}
	}
	
	return strings.Join(parts, ", ")
}

// GoCodeElement represents a structured element from Go source code
type GoCodeElement struct {
	Type       string `json:"type"`        // package, import, function, type, variable, constant, comment
	Name       string `json:"name"`        // name of the element (if applicable)
	Content    string `json:"content"`     // the actual code/content
	Line       int    `json:"line"`        // line number where it appears
	DocComment string `json:"doc_comment"` // associated documentation comment
}
