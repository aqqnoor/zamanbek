package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"zmb-assistant/services"
)

type ragIn struct {
	Query string `json:"query"`
	K     int    `json:"k,omitempty"`
}
type ragOut struct {
	Answer   string                   `json:"answer"`
	Snippets []map[string]any         `json:"snippets"`
}

func RagChatHandler(w http.ResponseWriter, r *http.Request) {
	var in ragIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest); return
	}
	if in.K <= 0 { in.K = 5 }

	idx, err := services.LoadIndex(services.IndexPath())
	if err != nil { http.Error(w, err.Error(), 502); return }

	// эмбеддинг запроса
	ec := services.NewEmbedding(os.Getenv("OPENAI_BASE_URL"), os.Getenv("OPENAI_API_KEY"))
	vec, err := ec.Embed(in.Query)
	if err != nil { http.Error(w, "embed error: "+err.Error(), 502); return }

	// поиск
	hits := idx.Search(vec, in.Query, in.K)

	// inside RagChatHandler, после hits := idx.Search(...)

	langKZ := services.IsKazakh(in.Query)
	intent := services.DetectIntent(in.Query)

	// собираем контекст
	var b strings.Builder
	if langKZ {
		b.WriteString("Сен исламдық қаржыландыру бойынша ассистентсің. Жауапты қазақ тілінде бер.\n")
	} else {
		b.WriteString("Ты ассистент по исламским финансам. Отвечай на языке запроса.\n")
	}
	b.WriteString("Не выдумывай фактов; пользуйся только Контекстом ниже.\n")
	b.WriteString("Если данных мало, сначала кратко укажи, чего не хватает, затем задай 1 уточняющий вопрос.\n\n")

	b.WriteString("Контекст:\n")
	snippets := make([]map[string]any, 0, len(hits))
	for _, h := range hits {
		b.WriteString("---\n")
		b.WriteString(h.Doc.Title + "\n")
		b.WriteString(h.Doc.Content + "\n")
		snippets = append(snippets, map[string]any{"id": h.Doc.ID, "title": h.Doc.Title, "score": h.Score})
	}

	// Специализированный формат для halal_check
	if intent == services.IntentHalalCheck {
		if langKZ {
			b.WriteString("\nТапсырма: Пайдаланушының сұрағы өнімнің халал талаптарға сәйкестігін тексеру туралы.\n")
			b.WriteString("Егер Контекстте нақты өнім көрсетілмесе, қысқа түрде қандай ақпарат керек екенін айт та, 1 сұрақ қой.\n")
			b.WriteString("Егер жеткілікті дерек бар болса, келесі форматта жауап бер:\n")
			b.WriteString("1) Қысқа «Иә/Жоқ/Анықтау керек».\n")
			b.WriteString("2) 2–4 критерий бойынша дәлел (риба жоқ па, тыйым салынған салалар, активке негізделуі, ашық шарттар).\n")
			b.WriteString("3) 2–3 жол басылым — қандай шектеулер бар.\n")
			b.WriteString("4) «Дереккөздер» бөлімінде Контексттегі тақырыптарды тізімде.\n")
		} else {
			b.WriteString("\nЗадача: Проверить соответствие продукта халяль-требованиям.\n")
			b.WriteString("Если в Контексте нет явного указания на продукт, сначала укажи, каких данных не хватает, и задай 1 уточняющий вопрос.\n")
			b.WriteString("Если данных достаточно, ответь в формате:\n")
			b.WriteString("1) Короткое «Да/Нет/Нужно уточнить».\n")
			b.WriteString("2) 2–4 критерия-доказательства (нет риба, запретные категории, опора на актив, прозрачные условия).\n")
			b.WriteString("3) 2–3 строки ограничений.\n")
			b.WriteString("4) Раздел «Источники» со списком заголовков из Контекста.\n")
		}
	} else {
		if langKZ {
			b.WriteString("\nЖауап форматы: 3–6 сөйлем; соңында «Дереккөздер» бөлімінде тақырыптарды бер.\n")
		} else {
			b.WriteString("\nФормат ответа: 3–6 предложений; затем раздел «Источники» со списком заголовков.\n")
		}
	}

	b.WriteString("\nПайдаланушы сұрағы / Вопрос пользователя:\n")
	b.WriteString(in.Query)

	// контекст
	var b strings.Builder
	b.WriteString("Ты — ассистент по исламским финансам. Отвечай кратко и точно, цитируй заголовки источников.\n")
	b.WriteString("Если запрос противоречит шариату — объясни почему и предложи дозволенные альтернативы.\n\n")
	b.WriteString("Контекст:\n")
	snippets := make([]map[string]any, 0, len(hits))
	for _, h := range hits {
		b.WriteString("---\n")
		b.WriteString(h.Doc.Title+"\n")
		b.WriteString(h.Doc.Content+"\n")
		snippets = append(snippets, map[string]any{
			"id": h.Doc.ID, "title": h.Doc.Title, "score": h.Score,
		})
	}
	b.WriteString("\nВопрос пользователя:\n")
	b.WriteString(in.Query)
	b.WriteString("\n\nФормат ответа: 3–6 предложений + маркированный список заголовков источников.")

	// LLM
	llm := services.NewLLM(os.Getenv("OPENAI_BASE_URL"), os.Getenv("OPENAI_API_KEY"))
	ans, err := llm.Chat("gpt-4o-mini",
		[]services.ChatMessage{
			{Role:"system", Content:"Следуй инструкциям, не придумывай фактов, отвечай на языке вопроса."},
			{Role:"user",   Content:b.String()},
		},
		0.2,
	)
	if err != nil { http.Error(w, err.Error(), 502); return }

	_ = json.NewEncoder(w).Encode(ragOut{
		Answer: ans, Snippets: snippets,
	})
}