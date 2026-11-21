package tifaweb

import (
	"testing"

	"github.com/liustack/tifa/tifatest"
)

// A Chinese file shows its body in place of the English doc, in the same place in the sidebar.
func TestChineseBodyReplacesTheEnglishOne(t *testing.T) {
	tifatest.LoadContent(t)
	en, err := Docs.Get("", "start/tutorial")
	if err != nil {
		t.Fatal(err)
	}
	zh := en
	zh.Title = "教程"
	tifatest.SetContent(t, DocsZH, zh)
	got, err := FindDoc(Chinese, "start/tutorial")
	if err != nil {
		t.Fatal(err)
	}
	if !got.Own || got.Title != "教程" || got.Group != en.Group || got.Order != en.Order {
		t.Errorf("Chinese doc = %+v", got)
	}
	if e, _ := FindDoc("en-US", "start/tutorial"); e.Title != en.Title || !e.Own {
		t.Errorf("English doc = %+v", e)
	}
	if err := CheckChinese(); err != nil {
		t.Error(err)
	}
}

// A doc with no Chinese file shows the English body in Chinese, and says it is not Own.
func TestEnglishBodyInChineseSiteIsMarked(t *testing.T) {
	tifatest.LoadContent(t)
	got, err := FindDoc(Chinese, "start/tutorial")
	if err != nil {
		t.Fatal(err)
	}
	if got.Own {
		t.Error("a doc with no Chinese file is not Own in Chinese")
	}
}

func TestChineseFileWithoutADocIsAnError(t *testing.T) {
	tifatest.LoadContent(t)
	z := Doc{Title: "x"}
	z.ID, z.Path, z.Slug = "nowhere/page", "nowhere/page", "page"
	tifatest.SetContent(t, DocsZH, z)
	if err := CheckChinese(); err == nil {
		t.Error("CheckChinese found nothing wrong")
	}
}

// With no file in content/docs-zh, which is how the site starts, the check passes.
func TestNoChineseFilesIsFine(t *testing.T) {
	tifatest.LoadContent(t)
	if err := CheckChinese(); err != nil {
		t.Error(err)
	}
}
