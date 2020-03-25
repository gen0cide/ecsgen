package ecsgen

import (
	"sort"
	"strings"
)

// Node represents a single element within the ECS definition graph. A node represents
// either an "object" - something that has child fields - or a "field" - something that
// represents a simple scalar field. The goal is for all nodes within ECS to be mapped
// to a tree graph, where each Node represents a vertex. Any vertex with children thus
// is an "object", and any "leafs" are therefor "fields".
type Node struct {
	// Name is the single name of the Node. For example,
	// the Node "client.nat", the name should equal "nat".
	Name string

	// Path is the absolute path from the root of a Node. For example,
	// the Node "client.nat" has a Path equal to "client.nat".
	Path string

	// Parent is a pointer to this Node's parent. For top level objects,
	// this will be nil. For all children this will equal the parent.
	// Example: Node("client.nat").Parent => Node("client")
	Parent *Node

	// Root is a reference to the root namespace.
	Root *Root

	// Children are a map of all child nodes that belong to this node. The key
	// is the Name field of the child, and the value is a pointer to the Node itself.
	// Example: Node("client.nat") has a child ["ip"] => Node("client.nat.ip").
	Children map[string]*Node

	// Definition is used to link back to the source of truth YAML definition
	// that was parsed for this Node. This is garanteed to be non-nil for all Nodes
	// that are of type "field", but is generally nil for objects, as ECS
	// treats them mostly as implicit.
	Definition *Definition
}

// IsTopLevel returns true if the Node has no Parent, therefor indicating it
// belongs in top level namespace.
func (n *Node) IsTopLevel() bool {
	return n.Parent == nil
}

// Child is used to resolve a child Node from a specific Node. The name argument
// should be a relative "Name" and *not* an absolute path. For example, if you wanted
// to retrieve the "client.nat" child from the "client" Node, you would pass "nat".
// If the child does not exist, it is created. If it does exist, the existing child is returned.
// If you pass it an empty string, the Node will simply return itself.
func (n *Node) Child(name string) *Node {
	// short circuit check to see if we have this child already
	if child, found := n.Children[name]; found {
		return child
	}

	// if the caller passes us an empty child, that must mean they
	// think we're to be resolved.
	if name == "" {
		return n
	}

	// Create the new Node
	newChild := &Node{
		Name:     name,
		Parent:   n,
		Root:     n.Root,
		Path:     strings.Join([]string{n.Path, name}, "."), // dynamically create the path from the parent
		Children: map[string]*Node{},
	}

	// Add the new Node to the top level index
	n.Root.Index[newChild.Path] = newChild

	// Add the new Node to the current Node's children
	n.Children[name] = newChild

	return newChild
}

// TypeIdent creates an Identifier based on the Node's type. This is almost never called
// for fields, but is required for objects. *Node.GoType() uses this to create object types.
// The returned Identifier is equal to NewIdentifier(n.Path).
func (n *Node) TypeIdent() Identifier {
	return NewIdentifier(n.Path)
}

// FieldIdent creates an Identifier that can be used as a field reference. This is used
// to create Go struct fields for every Node. The returned Identifier is equal to NewIdentifier(n.Name).
func (n *Node) FieldIdent() Identifier {
	return NewIdentifier(n.Name)
}

// IsImplied IsImplied is used to determine if a node is implied via the schema. This is
// true for the majority of objects in the schema, while false for all fields.
func (n *Node) IsImplied() bool {
	return n.Definition == nil
}

// IsObject attempts to determine if the Node is an "object" or a "field" and returns
// accordingly. This function contains edge-case code within the schema to account
// for one offs, and should be trusted to do the right thing.
func (n *Node) IsObject() bool {
	// if it doesn't have a definition, then the node
	// is most certainly an implicit object
	if n.Definition == nil {
		return true
	}

	// edge case for ".labels"
	// that field is a top level map[string]string{}
	if n.Definition.Name == "labels" {
		return false
	}

	// there's a number of spots where explicit objects are specified
	// as if they're a scalar field
	if n.Definition.Type == "object" {
		return true
	}

	return false
}

// IsArray is used to determine if the given Node within the ECS schema is
// actually an array of it's noted data type.
func (n *Node) IsArray() bool {
	// Not an array if it's implied
	if n.IsImplied() {
		return false
	}

	// not an array if we don't have anything in Normalize
	if len(n.Definition.Normalize) == 0 {
		return false
	}

	// only if normalize definition key has an "array" element,
	// do we return true
	for _, elm := range n.Definition.Normalize {
		if elm == "array" {
			return true
		}
	}

	// No it's not!
	return false
}

// ListChildren implements the Walkable interface.
func (n *Node) ListChildren() <-chan *Node {
	ret := make(chan *Node, len(n.Children))

	// short circuit if there's no children
	if len(n.Children) == 0 {
		close(ret)
		return ret
	}

	// get all the child keys
	keys := []string{}
	for k := range n.Children {
		keys = append(keys, k)
	}

	// sort them alphabetically
	sort.Strings(keys)

	// add the keys in order to the channel
	for _, k := range keys {
		ret <- n.Children[k]
	}

	close(ret)

	// return the channel
	return ret
}
