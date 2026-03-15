package main

import (
	"fmt"
	"maestro/pkg/parser"
)

func main() {
	result1 := parser.ParseFile("/tmp/testyaml/format1.md")
	fmt.Printf("Format 1 [Bob]: %v\n", result1.Frontmatter.Assignee)

	result2 := parser.ParseFile("/tmp/testyaml/format2.md")
	fmt.Printf("Format 2 block: %v\n", result2.Frontmatter.Assignee)
}
