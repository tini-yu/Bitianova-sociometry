package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type ResultService struct{}

func NewResultService() *ResultService {
	return &ResultService{}
}

type Results struct {
	Labels []string            `json:"labels"`
	UUID   string              `json:"uuid"`
	Maps   []map[string]string `json:"maps"` // всего 8 мап
}

func (s *ResultService) GetResults() Results {
	choiceCalc()
	statusCalc()
	mutualChoicesCalc()
	autosociometryCalc()
	referentChoiceCalc()

	mat1 := convertMap(positiveChoices)
	mat2 := convertMap(negativeChoices)
	mat4 := convertMap(mutualPositiveChoices)
	mat5 := convertMap(mutualNegativeChoices)
	mat6 := convertMap(contradictoryChoices)
	mat8 := convertMap(referentChoices)

	maps := []map[string]string{mat1, mat2, status, mat4, mat5, mat6, autosociometryChoices, mat8}

	// fmt.Println("SessionId: ", SessionId)
	return Results{
		Labels: names,
		UUID:   SessionId,
		Maps:   maps,
	}
}

// Конвертируем карты т.к. таблица в .ts работает только с map[string]string
func convertMap(originalMap map[string]int) map[string]string {
	newMap := make(map[string]string)
	for key, value := range originalMap {
		newMap[key] = strconv.Itoa(value)
	}
	return newMap
}

type ResultFile struct {
	UUID string     `json:"uuid"`
	Data [][]string `json:"data"`
}

func (s *ResultService) SaveResult(filename string, data [][]string, uuid string) error {

	fullPath := filepath.Join(SaveDir, filename)
	fmt.Println(fullPath)

	file, err := os.Create(fullPath)
	if err != nil {
		log.Printf("Error creating savefile at %s: %v", fullPath, err)
		return err
	}
	defer file.Close()

	payload := ResultFile{
		UUID: uuid,
		Data: data,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(payload)
}

func (s *ResultService) LoadResult(filename string) (string, error) {
	fullPath := filepath.Join(SaveDir, filename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", errors.New("save file not found")
	}
	file, err := os.Open(fullPath)
	if err != nil {
		log.Println("Error opening saved result file")
		return "", err
	}
	defer file.Close()

	var resultFile ResultFile
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&resultFile)
	if err != nil {
		log.Println("Error decoding saved result")
		return "", err
	}

	jsonData, err := json.Marshal(resultFile)
	if err != nil {
		log.Println("Error marshalling saved result")
		return "", err
	}
	return string(jsonData), nil
}

// Проверяем, что все матрицы описывают один и тот же список, чтобы не считать результаты в ином случае
func (s *ResultService) CheckMatrices() bool {
	matrix1, err := getMatrixUUID(M1)
	if err != nil {
		log.Println("Error getting matrix1 for calcs")
	}
	matrix2, err := getMatrixUUID(M2)
	if err != nil {
		log.Println("Error getting matrix2 for calcs")
	}
	matrix3, err := getMatrixUUID(M3)
	if err != nil {
		log.Println("Error getting matrix3 for calcs")
	}
	if matrix1 == matrix2 && matrix1 == matrix3 && matrix1 == SessionId {
		return true
	}

	return false
}
// Проверяем что текущее сохранение актуально.
func (s *ResultService) CheckResults(currentUUID string) bool {
	fullPath := filepath.Join(SaveDir, "results.json")

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		fmt.Println(err)
		return false
	}
	file, err := os.Open(fullPath)
	if err != nil {
		log.Println("Error opening saved matrix file")
		return false
	}
	defer file.Close()

	var resultsFile ResultFile
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&resultsFile)
	if err != nil {
		log.Println("Error decoding saved matrix")
		return false
	}

	if resultsFile.UUID == currentUUID {
		return true
	}

	return false
	
}

func getResultsAsStringSlices() ([][]string, error) {
	fullPath := filepath.Join(SaveDir, "results.json")

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, errors.New("save file not found")
	}
	file, err := os.Open(fullPath)
	if err != nil {
		log.Println("Error opening saved matrix file")
		return nil, err
	}
	defer file.Close()

	var originalMatrix originalMatrix
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&originalMatrix)
	if err != nil {
		log.Println("Error decoding saved matrix")
		return nil, err
	}

	return originalMatrix.Data, nil
}