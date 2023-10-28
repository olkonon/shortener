package memory

// NewMockStorage - создает заполненный InMemory для тестов
func NewMockStorage() *InMemory {
	return &InMemory{
		storeByID: map[string]string{
			"srewfrEW": "http://test.com/test?v=3",
			"rfdsgd":   "http://test.com/test",
		},
		storeByURL: map[string]string{
			"http://test.com/test?v=3": "srewfrEW",
			"http://test.com/test":     "rfdsgd",
		},
	}
}
