package web

import "fmt"

type nodeType int8

const (
	static nodeType = iota
	root
	param
	catchAll
)

type partHandle func(part string)
type node struct {
	pattern  string
	part     string
	children []*node
	types    nodeType
	isWild   bool
	handle   Handle
}

func (n *node) matchChild(part string) *node {
	if part[0] == '*' && len(part) <= 1 {
		panic("* behind not key")
	}
	if n.types == catchAll {
		panic("* Only at the tail")
	}
	for _, child := range n.children {
		if child.part == part || child.isWild {
			if child.types == catchAll {
				//   /z/g/z/c
				//   /s/p/*z
				if part[0] == '*' {
					panic("Only one peer can exist *")
				}
				continue
			}
			if child.types == param {
				if part[0] == ':' {
					panic(fmt.Sprintf("param %s and %s the same level cannot exist", n.part+"/"+child.part, n.part+"/"+part))
				}
				continue
			}
			return child
		}
	}
	return nil
}
func (n *node) matchChildren(part string) []*node {
	var child []*node
	for _, node := range n.children {
		if node.part == part || node.isWild {
			child = append(child, node)
		}
	}
	return child
}
func (n *node) insert(pattern string, parts []string, height int, handle Handle) {
	if len(parts) == height {
		n.pattern = pattern
		n.handle = handle
		return
	}
	if pattern == "/" {
		n.handle = handle
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	var nodeType_ nodeType
	var whether bool
	if child == nil {
		switch part[0] {
		case ':':
			nodeType_ = param
			whether = true
		case '*':
			nodeType_ = catchAll
			whether = true
		default:
			nodeType_ = static
		}
		child = &node{part: part, types: nodeType_, isWild: whether}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1, handle)
}
func (n *node) search(parts []string, height int, handle partHandle) (*node, map[string]string) {
	siz := len(parts) - 1

	handle("/")
	if parts[0] == "/" {
		return n, nil
	}
	part := parts[height]
	nodeAll := n.matchChildren(part)

	params := make(map[string]string)
	goto find
find:
	{
		if height > siz {
			return nil, nil
		}
		part = parts[height]
		for _, v := range nodeAll {
			if v.part == part || v.isWild {
				if v.types == param || v.types == catchAll {
					params[v.part[1:]] = part

				} else {
					handle(v.part)
				}
				if len(v.children) == 0 || v.children == nil {
					if height >= siz {
						return v, params
					}
					continue
				}
				nodeAll = v.children
				height++
				goto find
			}
		}

	}
	return nil, nil
}
