package views

import (
	"bytes"
	"context"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestLayout(t *testing.T) {
	component := Layout("Test Title")
	var buf bytes.Buffer
	err := component.Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("failed to render component: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(&buf)
	if err != nil {
		t.Fatalf("failed to parse html: %v", err)
	}

	if doc.Find("title").Text() != "Test Title - Spikeball League Tracker" {
		t.Errorf("expected title not found, got %s", doc.Find("title").Text())
	}
}

func TestLayout_FailingWriter(_ *testing.T) {
	for i := 0; i < 15000; i++ {
		_ = Layout("Test").Render(context.Background(), &failWriter{target: i})
	}
}
