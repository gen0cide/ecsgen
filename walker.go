package ecsgen

import (
	"errors"
)

// Walkable represents types that can be walked within ecsgen. This allows walking
// from arbitrary points within the graph, as well as from the root.
type Walkable interface {
	ListChildren() <-chan *Node
}

// ErrSkipChildren is a used as a return value from WalkFuncs to indicate that the
// callback should not be called for any children of the examined Node.
var ErrSkipChildren = errors.New("node walker: skip remaining children")

// WalkFunc is the type of the function called for each child of a node.
type WalkFunc func(n *Node) error

// Walk is a simple depth first walker for traversing the Schema from a starting Node.
func Walk(root Walkable, fn WalkFunc) error {
	// enumerate all children of the root
	for elm := range root.ListChildren() {

		// call the walk func on the element
		err := fn(elm)

		// if the returned error is an ErrSkipChildren, simply stop walking
		// this branch and continue
		if err == ErrSkipChildren {
			continue
		}

		// something else happened, stop immediately and return the error
		if err != nil {
			return err
		}

		// recursively call Walk on the element, bubbling up any errors
		// that arise in the recursive call
		err = Walk(elm, fn)
		if err != nil {
			return err
		}
	}

	return nil
}
