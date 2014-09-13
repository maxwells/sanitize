package sanitize

import (
	"io"
	"bytes"
	"code.google.com/p/go.net/html"
	"errors"
)

type Attribute struct {
	Name  		string
}

type Element struct {
	Tag			string
	Attributes 	[]*Attribute // if no attributes are specified, all are allowed
}

type Whitelist struct {
	Elements	[]*Element
}

func (w *Whitelist) AddElement(e *Element) {
	w.Elements = append(w.Elements, e)
}

func (w *Whitelist) GetElement(elementName string) (*Element, bool) {
	for _, element := range(w.Elements) {
		if element.Tag == elementName {
			return element, true
		}
	}
	return nil, false
}

func (e *Element) AddAttribute(a *Attribute) {
	e.Attributes = append(e.Attributes, a)
}

func (e *Element) GetAttribute(attributeName string) (*Attribute, bool) {
	for _, attribute := range(e.Attributes) {
		if attribute.Name == attributeName {
			return attribute, true
		}
	}
	return nil, false
}

// traverseTree takes a node and a function to perform on that node.
// It processess the tree withNode in pre-order.
//
// if withNode returns false, the subtree below this node will be
// skipped entirely.
func traverseTree(node *html.Node, withNode func(*html.Node) (bool, error)) (error) {
	processSubtree, err := withNode(node)
	if err != nil || !processSubtree {
		return err
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		err := traverseTree(c, withNode)
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
	switch n.Type {
	case html.ErrorNode:
		return errors.New("Unable to parse HTML")
  	case html.TextNode:
  	case html.DocumentNode:
  	case html.ElementNode:
  		element, hasElement := w.GetElement(n.Data)
  		if !hasElement {
  			if n.Parent != nil {
  				n.Parent.RemoveChild(n)
  			}
  			break
  		}

  		attributesToKeep := make([]html.Attribute, 0)

  		for _, attribute := range(n.Attr) {
			_, hasAttribute := element.GetAttribute(attribute.Key)
			if hasAttribute {
				attributesToKeep = append(attributesToKeep, attribute)
			}
		}
		n.Attr = attributesToKeep
  	case html.CommentNode:
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