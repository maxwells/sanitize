package sanitize

import (
	"io"
	"bytes"
	"code.google.com/p/go.net/html"
	"errors"
	"encoding/json"
	"strings"
)

type Whitelist struct {
	StripWhitespace	bool				`json:"stripWhitespace"`
	StripComments 	bool				`json:"stripComments"`
	Elements		map[string][]string	`json:"elements"`
}

func (w *Whitelist) AddElement(elementTag string, attributes []string) {
	w.Elements[elementTag] = attributes
}

func (w *Whitelist) HasElement(elementTag string) (bool) {
	_, ok := w.Elements[elementTag]
	return ok
}

func (w *Whitelist) GetAttributesForElement(elementTag string) ([]string) {
	val, _ := w.Elements[elementTag]
	return val
}

func (w *Whitelist) HasAttributeForElement(elementTag string, attributeName string) (bool) {
	val, ok := w.Elements[elementTag]
	if !ok {
		return false
	}
	for _, attribute := range(val) {
		if attribute == attributeName {
			return true
		}
	}
	return false
}

func (w *Whitelist) ToJSON() (string, error) {
	v, err := json.Marshal(w)
	return string(v), err
}

// sanitizeRemove traverses pre-order over the nodes,
// removing any element nodes that are not whitelisted
// and and removing any attributes that are not whitelisted
// from a given element node
func (w *Whitelist) sanitizeRemove(n *html.Node) (error) {
	switch n.Type {
	case html.ErrorNode:
		return errors.New("Unable to parse HTML")
  	case html.TextNode:
  		if w.StripWhitespace {
  			n.Data = strings.TrimSpace(n.Data)
  		}
  	case html.DocumentNode:
  	case html.ElementNode:
  		if !w.HasElement(n.Data) {
  			if n.Parent != nil {
  				n.Parent.RemoveChild(n)
  			}
  			break
  		}

  		attributesToKeep := make([]html.Attribute, 0)

  		for _, attribute := range(n.Attr) {
			if w.HasAttributeForElement(n.Data, attribute.Key) {
				attributesToKeep = append(attributesToKeep, attribute)
			}
		}
		n.Attr = attributesToKeep
  	case html.CommentNode:
  		if w.StripComments {
  			if n.Parent != nil {
  				n.Parent.RemoveChild(n)
  			}
  			break
  		}
  	case html.DoctypeNode:
	}

	// loop through child nodes
	var nextChild *html.Node
	for c := n.FirstChild; c != nil; c = nextChild {

		// grab a reference to the next child before
		// processing the current one; it may be removed
		// in processing
		nextChild = c.NextSibling
		err := w.sanitizeRemove(c)
		if err != nil {
			return err
		}
	}

	return nil
}

// remove non whitelisted elements entirely
func (w *Whitelist) SanitizeRemove(reader io.Reader) (string, error) {
	var buffer bytes.Buffer

	doc, err := html.Parse(reader)
	if err != nil {
		return buffer.String(), err
	}
	
	err = w.sanitizeRemove(doc)
	if err != nil {
		return buffer.String(), err
	}

	err = html.Render(&buffer, doc)

	return buffer.String(), err
}

// unwrap non whitelisted elements
func (w *Whitelist) SanitizeUnwrap(reader io.Reader) (string, error) {
	return "", nil
}