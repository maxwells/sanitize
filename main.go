package main

import (
	"log"
	"os"
	"fmt"
	"github.com/maxwells/sanitizer/sanitize"
)

// init is called prior to main.
func init() {
	// Change the device for logging to stdout.
	log.SetOutput(os.Stdout)
}

// main is the entry point for the program.
func main() {
	// Perform the search for the specified term.
	// str := `
	// <html>
	// 	<head>
	// 		<title>My title</title>
	// 	</head>
	// 	<body>
	// 		<h1>Heading</h1>
	// 		<p>Paragraph</p>
	// 		<span>Span</span>
	// </html>
	// `

	// sanitize.Santize(str, )

	// PROGRAMMATIC
	var attributes []*sanitize.Attribute
	attributes = append(attributes, &sanitize.Attribute{
		Name:   "href",
	})
	attributes = append(attributes, &sanitize.Attribute{
		Name:   "src",
	})

	var elements []*sanitize.Element

	elements = append(elements, &sanitize.Element{
		Tag: "Hello",
		Attributes: attributes,
	})

	whitelist := &sanitize.Whitelist{
		Elements: elements,
	}


	for _, element := range whitelist.Elements {
		fmt.Printf("element name %s\n", element.Tag)
		for _, attribute := range element.Attributes {
			fmt.Printf("attribute %s\n", attribute.Name)
		}
	}

	// FROM FILE
	whitelist, _ = sanitize.NewWhitelist("./whitelist.json")
	for _, element := range whitelist.Elements {
		fmt.Printf("element name %s\n", element.Tag)
		for _, attribute := range element.Attributes {
			fmt.Printf("attribute %s\n", attribute.Name)
		}
	}

	f, _ := os.Open("./example.html")
	sanitized, _ := whitelist.SanitizeRemove(f)
	fmt.Printf("\n%s",sanitized)

}