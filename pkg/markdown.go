package pkg

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type NamedCodeBlockAttributes struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type NamedCodeBlock struct {
	Attributes NamedCodeBlockAttributes `json:"attributes"`
	Content    string                   `json:"content"`
	Language   string                   `json:"language"`
}

func LoadNamedCodeBlocks(markdownContent []byte) (map[string]NamedCodeBlock, error) {
	// Create a Goldmark Markdown parser
	md := goldmark.New()

	// Parse the Markdown content into an AST (Abstract Syntax Tree)
	reader := text.NewReader(markdownContent)
	doc := md.Parser().Parse(reader)

	namedCodeBlocks := make(map[string]NamedCodeBlock)

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

			var namedCodeBlockAttributes NamedCodeBlockAttributes
			err := json.Unmarshal([]byte(attributes), &namedCodeBlockAttributes)
			if err != nil {
				fmt.Println("Error parsing JSON:", err)
				return 0, nil
			}

			// Append the NamedCodeBlock to the list
			namedCodeBlocks[namedCodeBlockAttributes.Name] = NamedCodeBlock{
				Attributes: namedCodeBlockAttributes,
				Content:    literal,
				Language:   string(fencedCodeBlock.Language(markdownContent)),
			}
		}
		return ast.WalkContinue, nil
	}

	if err := ast.Walk(doc, walker); err != nil {
		return nil, fmt.Errorf("failed to walk AST: %w", err)
	}

	return namedCodeBlocks, nil
}
