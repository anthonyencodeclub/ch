package outfmt_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/anthonyencodeclub/ch/internal/outfmt"
)

func TestFromFlags_Defaults(t *testing.T) {
	mode, err := outfmt.FromFlags(false, false)
	if err != nil {
		t.Fatalf("FromFlags(false, false) error: %v", err)
	}
	if mode.JSON || mode.Plain {
		t.Errorf("default mode should have JSON=false, Plain=false, got JSON=%v, Plain=%v", mode.JSON, mode.Plain)
	}
}

func TestFromFlags_JSONOnly(t *testing.T) {
	mode, err := outfmt.FromFlags(true, false)
	if err != nil {
		t.Fatalf("FromFlags(true, false) error: %v", err)
	}
	if !mode.JSON {
		t.Error("JSON mode should be true")
	}
	if mode.Plain {
		t.Error("Plain mode should be false")
	}
}

func TestFromFlags_PlainOnly(t *testing.T) {
	mode, err := outfmt.FromFlags(false, true)
	if err != nil {
		t.Fatalf("FromFlags(false, true) error: %v", err)
	}
	if mode.JSON {
		t.Error("JSON mode should be false")
	}
	if !mode.Plain {
		t.Error("Plain mode should be true")
	}
}

func TestFromFlags_BothConflict(t *testing.T) {
	_, err := outfmt.FromFlags(true, true)
	if err == nil {
		t.Fatal("FromFlags(true, true) should return error")
	}
}

func TestContextRoundtrip(t *testing.T) {
	mode := outfmt.Mode{JSON: true, Plain: false}
	ctx := outfmt.WithMode(context.Background(), mode)

	got := outfmt.FromContext(ctx)
	if got.JSON != true || got.Plain != false {
		t.Errorf("FromContext() = %+v, want %+v", got, mode)
	}
}

func TestIsJSON(t *testing.T) {
	ctx := outfmt.WithMode(context.Background(), outfmt.Mode{JSON: true})
	if !outfmt.IsJSON(ctx) {
		t.Error("IsJSON should return true")
	}
	if outfmt.IsPlain(ctx) {
		t.Error("IsPlain should return false")
	}
}

func TestIsPlain(t *testing.T) {
	ctx := outfmt.WithMode(context.Background(), outfmt.Mode{Plain: true})
	if outfmt.IsPlain(ctx) != true {
		t.Error("IsPlain should return true")
	}
	if outfmt.IsJSON(ctx) {
		t.Error("IsJSON should return false")
	}
}

func TestFromContext_EmptyContext(t *testing.T) {
	mode := outfmt.FromContext(context.Background())
	if mode.JSON || mode.Plain {
		t.Errorf("empty context should return zero Mode, got %+v", mode)
	}
}

func TestWriteJSON(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]string{"name": "Test Corp", "number": "12345678"}

	err := outfmt.WriteJSON(&buf, data)
	if err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}

	out := buf.String()
	if len(out) == 0 {
		t.Fatal("WriteJSON produced empty output")
	}
	// Should contain the key-value pairs
	if !bytes.Contains(buf.Bytes(), []byte(`"name"`)) {
		t.Error("JSON output should contain 'name' key")
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"Test Corp"`)) {
		t.Error("JSON output should contain 'Test Corp' value")
	}
}

func TestWriteJSON_NoHTMLEscape(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]string{"url": "https://example.com?a=1&b=2"}

	err := outfmt.WriteJSON(&buf, data)
	if err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	// Should NOT escape & to \u0026
	if bytes.Contains(buf.Bytes(), []byte(`\u0026`)) {
		t.Error("WriteJSON should not HTML-escape ampersands")
	}
}

func TestFromEnv(t *testing.T) {
	// Default (no env vars set) should be all false
	mode := outfmt.FromEnv()
	// We can't guarantee env vars aren't set, but we test the function doesn't panic
	_ = mode
}
