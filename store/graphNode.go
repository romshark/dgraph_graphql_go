package store

// AGraphNode represents an arbitrary graph node
type AGraphNode interface {
	// NodeID returns the internal graph node identifier
	NodeID() string
}

// GraphNode represents a graph node base
type GraphNode struct {
	UID string
}

// NodeID implements the AGraphNode interface
func (gn GraphNode) NodeID() string {
	return gn.UID
}
