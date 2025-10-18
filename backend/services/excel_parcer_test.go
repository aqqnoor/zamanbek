package services

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseExcelProducts_SkipIfMissing(t *testing.T) {
	path := "/mnt/data/Справочник по продуктам Zamanbank.xlsx"
	if st, err := os.Stat(path); err != nil || st.IsDir() {
		t.Skipf("excel not found at %s, skipping", path)
	}

	xs, err := ParseExcelProducts(path)
	if err != nil {
		t.Fatalf("ParseExcelProducts error: %v", err)
	}
	if len(xs) == 0 {
		t.Fatalf("expected some products from excel, got 0")
	}

	ps := ExcelToProducts(xs)
	if len(ps) == 0 {
		t.Fatalf("ExcelToProducts empty")
	}

	// базовые инварианты
	for _, p := range ps {
		if p.ID == "" {
			t.Errorf("product ID is empty for %+v", p)
		}
		if p.Name == "" {
			// допускаем, но отмечаем
			t.Logf("warning: product with empty name: %+v", p)
		}
		if filepath.Base(p.RawSource) == "" {
			t.Errorf("raw source (sheet) not set for %+v", p)
		}
	}
}
