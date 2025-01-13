package pkg

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type NamedCodeBlockAttribute struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type NamedCodeBlock struct {
	Attributes NamedCodeBlockAttribute `json:"attributes"`
	Content    string                  `json:"content"`
	Language   string                  `json:"language"`
}

func LoadNamedCodeBlocks(markdownContent []byte) ([]NamedCodeBlock, error) {
	// Create a Goldmark Markdown parser
	md := goldmark.New()

	// Parse the Markdown content into an AST (Abstract Syntax Tree)
	reader := text.NewReader(markdownContent)
	doc := md.Parser().Parse(reader)

	var namedCodeBlocks []NamedCodeBlock

	// Walk through the AST to find FencedCodeBlock nodes
	walker := func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		// Check if the current node is a FencedCodeBlock
		if entering && node.Kind() == ast.KindFencedCodeBlock {
			fencedCodeBlock := node.(*ast.FencedCodeBlock)

			// Extract the content
			var literal string
			for i := 0; i < fencedCodeBlock.Lines().Len(); i++ {
				line := fencedCodeBlock.Lines().At(i)
				literal += string(line.Value(markdownContent))
			}

			// Skip if no language, no content, or no attributes
			if fencedCodeBlock.Language(markdownContent) == nil || literal == "" || len(fencedCodeBlock.Info.Value(markdownContent)) < 4 {
				return ast.WalkContinue, nil
			}

			// Extract the attributes
			firstLine := strings.Fields(string(fencedCodeBlock.Info.Value(markdownContent)))
			attributes := strings.Join(firstLine[1:], " ")

			var namedCodeBlockAttribute NamedCodeBlockAttribute
			err := json.Unmarshal([]byte(attributes), &namedCodeBlockAttribute)
			if err != nil {
				fmt.Println("Error parsing JSON:", err)
				return 0, nil
			}

			// Append the NamedCodeBlock to the list
			namedCodeBlocks = append(namedCodeBlocks, NamedCodeBlock{
				Attributes: namedCodeBlockAttribute,
				Content:    literal,
				Language:   string(fencedCodeBlock.Language(markdownContent)),
			})
		}
		return ast.WalkContinue, nil
	}

	if err := ast.Walk(doc, walker); err != nil {
		return nil, fmt.Errorf("failed to walk AST: %w", err)
	}

	// Sort the list based on the Name attribute
	sort.Slice(namedCodeBlocks, func(i, j int) bool {
		return namedCodeBlocks[i].Attributes.Name < namedCodeBlocks[j].Attributes.Name
	})

	return namedCodeBlocks, nil
}

func GetNamedCodeBlock(namedCodeBlocks []NamedCodeBlock, name string) (NamedCodeBlock, error) {
	for _, block := range namedCodeBlocks {
		if block.Attributes.Name == name && block.Content != "" {
			return block, nil
		}
	}

	return NamedCodeBlock{}, fmt.Errorf("code block '%s' not found", name)
}
