package memory

import "github.com/olkonon/shortener/internal/app/common"

var MockID1 = common.GenHashedString("http://test.com/test?v=3")
var MockID2 = common.GenHashedString("http://test.com/test")

// NewMockStorage - создает заполненный InMemory для тестов
func NewMockStorage() *InMemory {
	return &InMemory{
		storeByID: map[string]map[string]string{
			common.AnonymousUser: {
				MockID1: "http://test.com/test?v=3",
				MockID2: "http://test.com/test",
			},
		},
	}
}
