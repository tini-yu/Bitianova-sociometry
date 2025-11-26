package service

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

type ListService struct{}

func NewListService() *ListService {
	return &ListService{}
}


var current_names = make(map[string]bool)
var names = make([]string, 0) //not nil, but len==0
var SavedList []string
var NameIds map[string]int
var SessionId string

func CreateSessionId() {
	SessionId = uuid.NewString()
	log.Println("SessionId: " + SessionId)
}

func (s *ListService) ListAdd(name string) error {
	if !current_names[name] {
		names = append(names, name)
		fmt.Println(names)
		current_names[name] = true
	} else {
		return fmt.Errorf("тестируемый уже в списке")
	}
	CreateSessionId()
	return nil
}

func (s *ListService) ListGet() []string {
	return names
}

func (s *ListService) ListRemove(id int) error {
	current_names[names[id]] = false
	names = append(names[:id], names[id+1:]...)
	CreateSessionId()
	return nil
}

type ListFile struct {
    Labels []string `json:"labels"`
    UUID   string   `json:"uuid"`
}

func (s *ListService) ListExport() ListFile {
    SavedList = names
	log.Println("List exported successfully")

    return ListFile{
        Labels: SavedList,
        UUID:   SessionId,
    }
}