package render

import (
	"strings"
	"testing"

	"gnodoc/model"
)

func TestMarkdownRenderer_Overview(t *testing.T) {
	pkg := &model.DocPackage{
		Name:       "mypkg",
		ImportPath: "example.com/mypkg",
		Summary:    "Package mypkg provides utilities.",
		Doc:        "Package mypkg provides utilities.\n\nThis is a longer description.",
	}

	r := NewMarkdownRenderer(DefaultOptions())
	result := r.RenderOverview(pkg)

	// Should contain package name as title
	if !strings.Contains(result, "# mypkg") {
		t.Error("expected package name as title")
	}

	// Should contain import path
	if !strings.Contains(result, "example.com/mypkg") {
		t.Error("expected import path")
	}

	// Should contain doc content
	if !strings.Contains(result, "This is a longer description.") {
		t.Error("expected doc content")
	}
}

func TestMarkdownRenderer_Overview_NoDoc(t *testing.T) {
	pkg := &model.DocPackage{
		Name:       "mypkg",
		ImportPath: "example.com/mypkg",
	}

	r := NewMarkdownRenderer(DefaultOptions())
	result := r.RenderOverview(pkg)

	// Should still have title
	if !strings.Contains(result, "# mypkg") {
		t.Error("expected package name as title")
	}
}

func TestMarkdownRenderer_Index(t *testing.T) {
	pkg := &model.DocPackage{
		Name: "mypkg",
		Funcs: []model.DocFunc{
			{DocNode: model.DocNode{Name: "NewFoo", Kind: model.KindFunc, Exported: true}},
			{DocNode: model.DocNode{Name: "Process", Kind: model.KindFunc, Exported: true}},
		},
		Types: []model.DocType{
			{DocNode: model.DocNode{Name: "Foo", Kind: model.KindType, Exported: true}},
		},
	}

	r := NewMarkdownRenderer(DefaultOptions())
	result := r.RenderIndex(pkg)

	// Should contain Index header
	if !strings.Contains(result, "## Index") {
		t.Error("expected Index header")
	}

	// Should contain function links
	if !strings.Contains(result, "[NewFoo]") {
		t.Error("expected NewFoo in index")
	}
	if !strings.Contains(result, "[Process]") {
		t.Error("expected Process in index")
	}

	// Should contain type link
	if !strings.Contains(result, "[Foo]") {
		t.Error("expected Foo in index")
	}
}

func TestMarkdownRenderer_Constants(t *testing.T) {
	pkg := &model.DocPackage{
		Name: "mypkg",
		Consts: []model.DocValueGroup{
			{
				DocNode: model.DocNode{Doc: "Size constants."},
				Specs: []model.DocValueSpec{
					{
						DocNode: model.DocNode{Name: "KB", Kind: model.KindConst, Exported: true},
						Type:    "int",
						Value:   "1024",
					},
					{
						DocNode: model.DocNode{Name: "MB", Kind: model.KindConst, Exported: true},
						Type:    "int",
						Value:   "1024 * KB",
					},
				},
			},
		},
	}

	r := NewMarkdownRenderer(DefaultOptions())
	result := r.RenderConstants(pkg)

	// Should contain Constants header
	if !strings.Contains(result, "## Constants") {
		t.Error("expected Constants header")
	}

	// Should contain constant names
	if !strings.Contains(result, "KB") {
		t.Error("expected KB constant")
	}
	if !strings.Contains(result, "MB") {
		t.Error("expected MB constant")
	}
}

func TestMarkdownRenderer_Variables(t *testing.T) {
	pkg := &model.DocPackage{
		Name: "mypkg",
		Vars: []model.DocValueGroup{
			{
				DocNode: model.DocNode{Doc: "Default configuration."},
				Specs: []model.DocValueSpec{
					{
						DocNode: model.DocNode{Name: "DefaultConfig", Kind: model.KindVar, Exported: true},
						Type:    "*Config",
					},
				},
			},
		},
	}

	r := NewMarkdownRenderer(DefaultOptions())
	result := r.RenderVariables(pkg)

	// Should contain Variables header
	if !strings.Contains(result, "## Variables") {
		t.Error("expected Variables header")
	}

	// Should contain variable name
	if !strings.Contains(result, "DefaultConfig") {
		t.Error("expected DefaultConfig variable")
	}
}

func TestMarkdownRenderer_Functions(t *testing.T) {
	pkg := &model.DocPackage{
		Name: "mypkg",
		Funcs: []model.DocFunc{
			{
				DocNode: model.DocNode{
					Name:     "NewFoo",
					Kind:     model.KindFunc,
					Summary:  "NewFoo creates a new Foo.",
					Doc:      "NewFoo creates a new Foo.\n\nIt returns nil if x is negative.",
					Exported: true,
				},
				Params: []model.DocParam{
					{Name: "x", Type: "int"},
				},
				Results: []model.DocParam{
					{Type: "*Foo"},
				},
			},
		},
	}

	r := NewMarkdownRenderer(DefaultOptions())
	result := r.RenderFunctions(pkg)

	// Should contain signature
	if !strings.Contains(result, "func NewFoo") {
		t.Error("expected function signature")
	}

	// Should contain doc
	if !strings.Contains(result, "It returns nil if x is negative.") {
		t.Error("expected function doc")
	}
}

func TestMarkdownRenderer_MultilineDocLists(t *testing.T) {
	pkg := &model.DocPackage{
		Funcs: []model.DocFunc{{
			DocNode: model.DocNode{
				Name:     "CreatePool",
				Kind:     model.KindFunc,
				Exported: true,
				Doc: "CreatePool creates a pool.\n\nParameters:\n" +
					"  - cur: Current realm context; callers use cross(cur) when crossing into this\n" +
					"    realm.\n" +
					"  - token0Path: First token path, which is normalized into the\n" +
					"    canonical pool order.\n\n" +
					"Returns:\n" +
					"  - poolPath: Canonical pool path as a\n" +
					"    string.\n\n" +
					"Halt check: reverts while the Pool halt scope is active.",
			},
			Params:  []model.DocParam{{Name: "cur", Type: "realm"}, {Name: "token0Path", Type: "string"}},
			Results: []model.DocParam{{Name: "poolPath", Type: "string"}},
		}},
	}

	result := NewMarkdownRenderer(DefaultOptions()).RenderFunctions(pkg)
	for _, want := range []string{
		"| `cur` | realm | Current realm context; callers use cross(cur) when crossing into this realm. |",
		"| `token0Path` | string | First token path, which is normalized into the canonical pool order. |",
		"| `poolPath` | string | Canonical pool path as a string. |",
		"> **Halt check:** reverts while the Pool halt scope is active.",
	} {
		if !strings.Contains(result, want) {
			t.Errorf("missing rendered description %q:\n%s", want, result)
		}
	}
	if strings.Contains(result, "\n    realm.") || strings.Contains(result, "\n  - token0Path:") {
		t.Errorf("parameter continuation leaked into prose:\n%s", result)
	}
	if !strings.Contains(result, "```go\nfunc CreatePool(") || !strings.Contains(result, ") (poolPath string)\n```") {
		t.Errorf("signature code fence has padding lines:\n%s", result)
	}
	if !strings.Contains(result, "CreatePool creates a pool.\n\n> **Halt check:**") {
		t.Errorf("halt note is not separated from prose:\n%s", result)
	}
}

func TestMarkdownRenderer_Types(t *testing.T) {
	pkg := &model.DocPackage{
		Name: "mypkg",
		Types: []model.DocType{
			{
				DocNode: model.DocNode{
					Name:      "Foo",
					Kind:      model.KindType,
					Summary:   "Foo represents a foo.",
					Doc:       "Foo represents a foo.\n\nIt is used for bar.",
					Signature: "type Foo struct",
					Exported:  true,
				},
				TypeKind: model.TypeKindStruct,
				Fields: []model.DocField{
					{
						DocNode: model.DocNode{Name: "ID", Kind: model.KindField, Exported: true},
						Type:    "int",
					},
				},
				Methods: []model.DocFunc{
					{
						DocNode: model.DocNode{
							Name:     "String",
							Kind:     model.KindMethod,
							Summary:  "String returns a string representation.",
							Exported: true,
						},
						Receiver: &model.DocReceiver{Name: "f", Type: "Foo"},
						Results:  []model.DocParam{{Type: "string"}},
					},
				},
			},
		},
	}

	r := NewMarkdownRenderer(DefaultOptions())
	result := r.RenderTypes(pkg)

	// Should contain Types header
	if !strings.Contains(result, "## Types") {
		t.Error("expected Types header")
	}

	// Should contain type name as subheader
	if !strings.Contains(result, "### Foo") {
		t.Error("expected Foo as subheader")
	}

	// Should contain method
	if !strings.Contains(result, "String") {
		t.Error("expected String method")
	}

	if !strings.Contains(result, "```go\ntype Foo struct\n```") {
		t.Errorf("type declaration should not have blank lines inside its code fence:\n%s", result)
	}
}

func TestMarkdownRenderer_Examples(t *testing.T) {
	pkg := &model.DocPackage{
		Name: "mypkg",
		Examples: []model.DocExample{
			{
				Name:   "Example",
				Doc:    "Basic usage example.",
				Code:   `fmt.Println("hello")`,
				Output: "hello",
			},
		},
	}

	r := NewMarkdownRenderer(DefaultOptions())
	result := r.RenderExamples(pkg)

	// Should contain Examples header
	if !strings.Contains(result, "## Examples") {
		t.Error("expected Examples header")
	}

	// Should contain code block
	if !strings.Contains(result, "```go") {
		t.Error("expected go code block")
	}

	// Should contain code
	if !strings.Contains(result, `fmt.Println("hello")`) {
		t.Error("expected example code")
	}

	// Should contain output
	if !strings.Contains(result, "Output:") || !strings.Contains(result, "hello") {
		t.Error("expected example output")
	}

	if !strings.Contains(result, "```go\nfmt.Println(\"hello\")\n```") || !strings.Contains(result, "```\nhello\n```") {
		t.Errorf("example code and output should not have blank lines inside their fences:\n%s", result)
	}
}

func TestMarkdownRenderer_Notes(t *testing.T) {
	pkg := &model.DocPackage{
		Name: "mypkg",
		Notes: []model.DocNote{
			{Kind: "BUG", Body: "This function has a known issue."},
			{Kind: "TODO", Body: "Add more tests."},
		},
		Deprecated: []model.DocDeprecated{
			{Body: "Use NewFoo instead."},
		},
	}

	r := NewMarkdownRenderer(DefaultOptions())
	result := r.RenderNotes(pkg)

	// Should contain Notes header
	if !strings.Contains(result, "## Notes") {
		t.Error("expected Notes header")
	}

	// Should contain BUG note
	if !strings.Contains(result, "BUG") {
		t.Error("expected BUG note")
	}

	// Should contain deprecated
	if !strings.Contains(result, "Deprecated") {
		t.Error("expected Deprecated section")
	}
}

func TestMarkdownRenderer_ExportedOnly(t *testing.T) {
	pkg := &model.DocPackage{
		Name: "mypkg",
		Funcs: []model.DocFunc{
			{DocNode: model.DocNode{Name: "PublicFunc", Kind: model.KindFunc, Exported: true}},
			{DocNode: model.DocNode{Name: "privateFunc", Kind: model.KindFunc, Exported: false}},
		},
	}

	opts := DefaultOptions()
	opts.ExportedOnly = true

	r := NewMarkdownRenderer(opts)
	result := r.RenderFunctions(pkg)

	// Should contain exported function
	if !strings.Contains(result, "PublicFunc") {
		t.Error("expected PublicFunc")
	}

	// Should NOT contain unexported function
	if strings.Contains(result, "privateFunc") {
		t.Error("unexpected privateFunc in output")
	}
}

func TestMarkdownRenderer_FullRender(t *testing.T) {
	pkg := &model.DocPackage{
		Name:       "mypkg",
		ImportPath: "example.com/mypkg",
		Summary:    "Package mypkg provides utilities.",
		Doc:        "Package mypkg provides utilities.\n\nThis is a detailed description.",
		Consts: []model.DocValueGroup{
			{
				Specs: []model.DocValueSpec{
					{
						DocNode: model.DocNode{Name: "MaxSize", Kind: model.KindConst, Exported: true},
						Value:   "1024",
					},
				},
			},
		},
		Funcs: []model.DocFunc{
			{
				DocNode: model.DocNode{
					Name:     "NewFoo",
					Kind:     model.KindFunc,
					Doc:      "NewFoo creates a new Foo.",
					Exported: true,
				},
				Results: []model.DocParam{{Type: "*Foo"}},
			},
		},
		Types: []model.DocType{
			{
				DocNode: model.DocNode{
					Name:     "Foo",
					Kind:     model.KindType,
					Doc:      "Foo is a type.",
					Exported: true,
				},
				TypeKind: model.TypeKindStruct,
			},
		},
		Examples: []model.DocExample{
			{
				Name: "Example",
				Code: `fmt.Println("hello")`,
			},
		},
	}

	r := NewMarkdownRenderer(DefaultOptions())
	result := r.Render(pkg)

	// Should contain all sections
	sections := []string{
		"# mypkg",
		"## Index",
		"## Constants",
		"## Types",
		"## Examples",
	}

	for _, section := range sections {
		if !strings.Contains(result, section) {
			t.Errorf("expected section %q in output", section)
		}
	}

	// Check order: Overview should come before Index
	overviewIdx := strings.Index(result, "# mypkg")
	indexIdx := strings.Index(result, "## Index")
	if overviewIdx > indexIdx {
		t.Error("Overview should come before Index")
	}

	// Index should come before Constants
	constsIdx := strings.Index(result, "## Constants")
	if indexIdx > constsIdx {
		t.Error("Index should come before Constants")
	}
}

func TestMarkdownRenderer_EmptyPackage(t *testing.T) {
	pkg := &model.DocPackage{
		Name:       "empty",
		ImportPath: "example.com/empty",
	}

	r := NewMarkdownRenderer(DefaultOptions())
	result := r.Render(pkg)

	// Should still have title
	if !strings.Contains(result, "# empty") {
		t.Error("expected package name")
	}

	// Should not have empty sections
	if strings.Contains(result, "## Functions") {
		t.Error("should not have Functions section for empty package")
	}
	if strings.Contains(result, "## Types") {
		t.Error("should not have Types section for empty package")
	}
}
