package views

import (
	"bytes"
	"context"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestLayout(t *testing.T) {
	component := Layout("Test Title", "light")
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

	htmlNode := doc.Find("html")
	if attr, exists := htmlNode.Attr("data-bs-theme"); !exists || attr != "light" {
		t.Errorf("expected data-bs-theme='light', got %s", attr)
	}
	if attr, exists := htmlNode.Attr("data-theme"); !exists || attr != "light" {
		t.Errorf("expected data-theme='light', got %s", attr)
	}
	if doc.Find("#theme-toggle").Length() == 0 {
		t.Errorf("expected theme-toggle button to be present")
	}
	if doc.Find(".theme-icon-moon").Length() == 0 {
		t.Errorf("expected theme-icon-moon element to be present")
	}
	if doc.Find(".theme-icon-sun").Length() == 0 {
		t.Errorf("expected theme-icon-sun element to be present")
	}
}

func TestLayout_DarkMode(t *testing.T) {
	component := Layout("Dark Title", "dark")
	var buf bytes.Buffer
	err := component.Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("failed to render component: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(&buf)
	if err != nil {
		t.Fatalf("failed to parse html: %v", err)
	}

	htmlNode := doc.Find("html")
	if attr, exists := htmlNode.Attr("data-bs-theme"); !exists || attr != "dark" {
		t.Errorf("expected data-bs-theme='dark', got %s", attr)
	}
	if attr, exists := htmlNode.Attr("data-theme"); !exists || attr != "dark" {
		t.Errorf("expected data-theme='dark', got %s", attr)
	}
}

func TestLayout_FailingWriter(_ *testing.T) {
	for i := 0; i < 15000; i++ {
		_ = Layout("Test", "light").Render(context.Background(), &failWriter{target: i})
	}
}
