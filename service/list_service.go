package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type ListService struct{}

func NewListService() *ListService {
	return &ListService{}
}

// Убедиться, что совпадает с константой в .ts коде
const LIST_NAME = "List.json"

var (
	current_names  = make(map[string]bool)
	names          = make([]string, 0) // not nil, but len==0
	SavedList      []string
	NameIds        map[string]int
	SessionId      string
	testDate       time.Time
	testDateString string
	prettyTestDate string
)

func CreateSessionId() {
	SessionId = uuid.NewString()
}

func (s *ListService) SetSessionId(uuid string) {
	SessionId = uuid
}

func (s *ListService) SetList(savedNames []string) {
	names = savedNames
}

func (s *ListService) ListAdd(name string) error {
	if !current_names[name] {
		names = append(names, name)
		// fmt.Println(names)
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
	SavedList, err := s.LoadListNames(LIST_NAME)
	if err != nil {
		fmt.Printf("Ошибка экспорта листа имен")
	}

	uuid := ""
	uuid, err = LoadSessionId(ProjectDir)
	if err != nil {
		fmt.Printf("не удалось загрузить uuid листа")
	}

	return ListFile{
		Labels: SavedList,
		UUID:   uuid, //SessionId (Актуальный из list.json)
	}
}

func getListInterface() *[]interface{} {
	listInterface := make([]interface{}, len(names))
	for id, val := range names {
		listInterface[id] = val
	}

	return &listInterface
}

func (s *ListService) SaveTestDate(date string) {
	inputLayout := "2006-01-02" // yyyy-mm-dd

	t, err := time.Parse(inputLayout, date)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}
	prettyTestDate = t.Format("02.01.2006")
	testDateString = date
	testDate = t
}

type ListSaveFile struct {
	Labels []string `json:"labels"`
	UUID   string   `json:"uuid"`
	Date   string   `json:"date"`
}

func (s *ListService) SaveList(filename string) error {
	fullPath := filepath.Join(ProjectDir, filename)

	file, err := os.Create(fullPath)
	if err != nil {
		log.Printf("Error creating savefile at %s: %v", fullPath, err)
		return err
	}
	defer file.Close()

	payload := ListSaveFile{
		Labels: names,
		UUID:   SessionId,
		Date:   testDateString,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(payload)
}

func createEmptyList(fullPath string) error {
	file, err := os.Create(fullPath)
	if err != nil {
		log.Printf("Error creating savefile at %s: %v", fullPath, err)
		return err
	}
	defer file.Close()

	payload := ListSaveFile{
		Labels: []string{},
		UUID:   "",
		Date:   "",
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(payload)
}

func (s *ListService) LoadList(filename string) (string, error) {
	fullPath := filepath.Join(ProjectDir, filename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		createEmptyList(fullPath)
		return "{}", errors.New("save file not found")
	}
	file, err := os.Open(fullPath)
	if err != nil {
		log.Println("Error opening saved namelist file")
		return "", err
	}
	defer file.Close()

	var listSaveFile ListSaveFile
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&listSaveFile) // проверяем совпадение структуры
	if err != nil {
		log.Println("Error decoding saved namelist")
		return "", err
	}

	jsonData, err := json.Marshal(listSaveFile)
	if err != nil {
		log.Println("Error marshalling saved namelist")
		return "", err
	}
	return string(jsonData), nil
}

func (s *ListService) LoadListNames(filename string) ([]string, error) {
	fullPath := filepath.Join(ProjectDir, filename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		createEmptyList(fullPath)
		return []string{}, errors.New("save file not found")
	}
	file, err := os.Open(fullPath)
	if err != nil {
		log.Println("Error opening saved namelist file")
		return []string{}, err
	}
	defer file.Close()

	var listSaveFile ListSaveFile
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&listSaveFile) // проверяем совпадение структуры
	if err != nil {
		log.Println("Error decoding saved namelist")
		return []string{}, err
	}

	return listSaveFile.Labels, nil
}

func LoadSessionId(listDir string) (string, error) {
	fullPath := filepath.Join(listDir, LIST_NAME)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		createEmptyList(fullPath)
		return "", errors.New("save file not found")
	}
	file, err := os.Open(fullPath)
	if err != nil {
		log.Println("Error opening saved namelist file")
		return "", err
	}
	defer file.Close()

	var listSaveFile ListSaveFile
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&listSaveFile) // проверяем совпадение структуры
	if err != nil {
		log.Println("Error decoding saved namelist")
		return "", err
	}

	return listSaveFile.UUID, nil
}
