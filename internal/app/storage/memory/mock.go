package memory

import "github.com/olkonon/shortener/internal/app/common"

var MockID1 = common.GenHashedString("http://test.com/test?v=3")
var MockID2 = common.GenHashedString("http://test.com/test")

// NewMockStorage - создает заполненный InMemory для тестов
func NewMockStorage() *InMemory {
	return &InMemory{
		storeByID: map[string]map[string]Record{
			common.TestUser: {
				MockID1: Record{
					OriginalURL: "http://test.com/test?v=3",
					IsDeleted:   false,
				},
				MockID2: Record{
					OriginalURL: "http://test.com/test",
					IsDeleted:   false,
				},
			},
		},
	}
}
