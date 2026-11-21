package docs

import "github.com/liustack/tifa"

// Layout gives every docs page the page tree for the sidebar.
func Layout(c *tifa.Ctx) (Tree, error) {
	return BuildTree(), nil
}
