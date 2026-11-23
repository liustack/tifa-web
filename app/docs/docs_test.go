package docs

import (
	"testing"

	"github.com/liustack/tifa/tifatest"

	tifaweb "github.com/liustack/tifa-web"
)

func TestTreeFollowsTheLanguage(t *testing.T) {
	tifatest.LoadContent(t)
	en, zh := BuildTree("en-US"), BuildTree(tifaweb.Chinese)
	if len(en.Children) == 0 || len(en.Children) != len(zh.Children) {
		t.Fatalf("the trees have %d and %d sections", len(en.Children), len(zh.Children))
	}
	if en.Children[0].Name != "Start" || zh.Children[0].Name != "开始" {
		t.Errorf("first sections: %q and %q", en.Children[0].Name, zh.Children[0].Name)
	}
	if en.Children[0].Children[0].URL != zh.Children[0].Children[0].URL {
		t.Errorf("the pages are at different URLs: %q and %q", en.Children[0].Children[0].URL, zh.Children[0].Children[0].URL)
	}
}
