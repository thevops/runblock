package pkg

import (
	"encoding/json"
	"fmt"
	"runblock/pkg/logger"
	"sort"
	"strings"

	"regexp"

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

// LoadNamedCodeBlocks parses the provided Markdown content to extract named code blocks.
// It returns a slice of NamedCodeBlock structs and an error if any occurred during parsing.
//
// Parameters:
// - markdownContent: A byte slice containing the Markdown content to be parsed.
//
// Returns:
// - A slice of NamedCodeBlock structs containing the extracted code blocks.
// - An error if any occurred during parsing.
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

			// Skip if no language or no content
			if fencedCodeBlock.Language(markdownContent) == nil || literal == "" {
				return ast.WalkContinue, nil
			}

			// Extract the attributes
			firstLine := strings.Fields(string(fencedCodeBlock.Info.Value(markdownContent)))
			// Use a regular expression to extract the JSON part of the attributes
			re := regexp.MustCompile(`\{.*\}`)
			attributes := re.FindString(strings.Join(firstLine[1:], " "))

			// Skip if no attributes
			if attributes == "" {
				return ast.WalkContinue, nil
			}

			var namedCodeBlockAttribute NamedCodeBlockAttribute
			err := json.Unmarshal([]byte(attributes), &namedCodeBlockAttribute)
			if err != nil {
				logger.Log.Fatalf("Error parsing JSON for attributes '%s': %v\n", attributes, err)
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

func CreateNamedCodeBlockMap(namedCodeBlocks []NamedCodeBlock) map[string]NamedCodeBlock {
	namedCodeBlockMap := make(map[string]NamedCodeBlock, len(namedCodeBlocks))
	for _, block := range namedCodeBlocks {
		if block.Content != "" {
			namedCodeBlockMap[block.Attributes.Name] = block
		}
	}
	return namedCodeBlockMap
}

// GetNamedCodeBlock retrieves a NamedCodeBlock from the map by its name.
// Parameters:
// - namedCodeBlockMap: a map where the key is the name of the code block and the value is the NamedCodeBlock.
// - name: the name of the code block to retrieve.
// Returns:
// - NamedCodeBlock: the code block with the specified name.
// - error: an error if the code block is not found.
func GetNamedCodeBlock(namedCodeBlockMap map[string]NamedCodeBlock, name string) (NamedCodeBlock, error) {
	if block, exists := namedCodeBlockMap[name]; exists {
		return block, nil
	}
	return NamedCodeBlock{}, fmt.Errorf("code block '%s' not found", name)
}
