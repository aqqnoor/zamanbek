package services

// NewLLM — фабрика: выбирает реализацию по флагу backend ("bank" или "mock").
func NewLLM(baseURL, apiKey, backend string) LLMClient {
    if backend == "bank" {
        return NewLLMBank(baseURL, apiKey)
    }
    return NewLLMMock()
}
