package services

import "testing"

func TestSearchProductsOnMockData(t *testing.T) {
	// моковая база
	prodDB = ProductDB{
		Products: []Product{
			{ID: "card-halal", Category: "Retail", Name: "Халал карта", Type: "карта", Halal: true, Tags: []string{"карта","halal","visa"}},
			{ID: "murabaha", Category: "Retail", Name: "Мурабаха", Type: "финансирование", Halal: true, Tags: []string{"мурабаха","финансирование"}},
			{ID: "deposit-halal", Category: "Retail", Name: "Халал депозит", Type: "депозит", Halal: true, Tags: []string{"депозит","halal"}},
		},
	}
	// поиск по названию
	got := SearchProducts("мурабаха", "ru", 3)
	if len(got) == 0 || got[0].ID != "murabaha" {
		t.Fatalf("search by name failed, got=%+v", got)
	}
	// поиск по типу
	got = SearchProducts("депозит", "ru", 3)
	if len(got) == 0 || got[0].ID != "deposit-halal" {
		t.Fatalf("search by type failed, got=%+v", got)
	}
	// поиск по тэгу
	got = SearchProducts("visa", "ru", 3)
	if len(got) == 0 || got[0].ID != "card-halal" {
		t.Fatalf("search by tag failed, got=%+v", got)
	}
}

func TestBuildContextBullets(t *testing.T) {
	ps := []Product{
		{ID: "murabaha", Name: "Мурабаха", Type: "финансирование", Halal: true, Markup: "от 300", MinSum: "10000", MaxSum: "300000"},
	}
	kz := BuildContextBullets(ps, "kz")
	ru := BuildContextBullets(ps, "ru")
	en := BuildContextBullets(ps, "en")

	if len(kz) == 0 || len(ru) == 0 || len(en) == 0 {
		t.Fatalf("bullets must not be empty")
	}
}
