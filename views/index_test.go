package views

import (
	"bytes"
	"context"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestIndex(t *testing.T) {
	component := Index()
	var buf bytes.Buffer
	err := component.Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("failed to render component: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(&buf)
	if err != nil {
		t.Fatalf("failed to parse html: %v", err)
	}

	if doc.Find("h1").Text() != "Spikeball Leagues" {
		t.Errorf("expected title not found, got %s", doc.Find("h1").Text())
	}
	if doc.Find("form[hx-post='/leagues']").Length() == 0 {
		t.Errorf("expected form not found")
	}
}

func TestIndex_FailingWriter(_ *testing.T) {
	for i := 0; i < 15000; i++ {
		_ = Index().Render(context.Background(), &failWriter{target: i})
	}
}
