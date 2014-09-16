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

// Remove all attributes on the provided node
// that are not contained within this whitelist
func (w *Whitelist) sanitizeAttributes(n *html.Node) {
	attributes := make([]html.Attribute, len(n.Attr))

	i := 0
	for _, attribute := range(n.Attr) {
		if w.HasAttributeForElement(n.Data, attribute.Key) {
			attributes[i] = attribute
			i += 1
		}
	}
	n.Attr = attributes[0:i]

}

// Remove the comment if this whitelist is configured
// with the StripComments configuration
func (w *Whitelist) handleComment(n *html.Node) {
	if w.StripComments {
			if n.Parent != nil {
			n.Parent.RemoveChild(n)
		}
	}
}

// Strip whitespace if this whitelist is configured
// with the StripWhitespace configuration
func (w *Whitelist) handleText(n *html.Node) {
	if w.StripWhitespace {
		n.Data = strings.TrimSpace(n.Data)
	}
}

// Helper function to process a specific node.
// Defers logic around how to handle ElementNodes to
// the caller.
//
// Returns the return value of handleElement:
// a boolean describing whether the children
// of the node should be further sanitized (ie. not skipped).
func (w *Whitelist) sanitizeNode(n *html.Node, handleElement func(*html.Node) (bool)) (error) {
	switch n.Type {
	case html.ErrorNode:
		return errors.New("Unable to parse HTML")
  	case html.TextNode:
  		w.handleText(n)
  	case html.DocumentNode:
  	case html.ElementNode:
  		if (!handleElement(n)) {
  			return nil
  		}
  		w.sanitizeAttributes(n)
  	case html.CommentNode:
  		w.handleComment(n)
  	case html.DoctypeNode:
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		err := w.sanitizeNode(c, handleElement)
		if err != nil {
			return err
		}
	}

	return nil
}

// sanitizeRemove traverses pre-order over the nodes,
// removing any element nodes that are not whitelisted
// and and removing any attributes that are not whitelisted
// from a given element node
func (w *Whitelist) sanitizeRemove(n *html.Node) (error) {
	return w.sanitizeNode(n, func(n *html.Node) (bool) {
  		if !w.HasElement(n.Data) {
  			if n.Parent != nil {
	  			nextSibling := n.NextSibling
  				n.Parent.RemoveChild(n)

  				// reset next sibling to support continuation
  				// of linked-list style traversal of parent node's children
  				n.NextSibling = nextSibling
  			}
  			return false
  		}
  		return true
	})
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

func (w *Whitelist) sanitizeUnwrap(n *html.Node) (error) {
	return w.sanitizeNode(n, func(n *html.Node) (bool) {
		if w.HasElement(n.Data) || n.Parent == nil {
			return true
		}

		insertBefore := n.NextSibling
		firstChild := n.FirstChild
		for c := n.FirstChild; c != nil; {
			nodeToUnwrap := c
			c = c.NextSibling
			
			n.RemoveChild(nodeToUnwrap)
			n.Parent.InsertBefore(nodeToUnwrap, insertBefore)
		}
		n.Parent.RemoveChild(n)

		// reset next sibling to support continuation
  		// of linked-list style traversal of parent node's children
		n.NextSibling = firstChild
		return false
	})
}

// unwrap non whitelisted elements
func (w *Whitelist) SanitizeUnwrap(reader io.Reader) (string, error) {
	var buffer bytes.Buffer

	doc, err := html.Parse(reader)
	if err != nil {
		return buffer.String(), err
	}
	
	err = w.sanitizeUnwrap(doc)
	if err != nil {
		return buffer.String(), err
	}

	err = html.Render(&buffer, doc)

	return buffer.String(), err
}