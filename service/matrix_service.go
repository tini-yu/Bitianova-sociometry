package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var SaveDir string

type MatrixService struct{}

func NewMatrixService() *MatrixService {
	return &MatrixService{}
}

type MatrixFile struct {
	UUID string  `json:"uuid"`
	Data [][]int `json:"data"`
}

func (s *MatrixService) SaveMatrix(filename string, data [][]int, uuid string) error {
	fullPath := filepath.Join(SaveDir, filename)
	fmt.Println(fullPath)

	file, err := os.Create(fullPath)
	if err != nil {
		log.Printf("Error creating savefile at %s: %v", fullPath, err)
		return err
	}
	defer file.Close()

	payload := MatrixFile{
		UUID: uuid,
		Data: data,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(payload)
}

func (s *MatrixService) LoadMatrix(filename string) (string, error) {
	fullPath := filepath.Join(SaveDir, filename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "{}", errors.New("save file not found")
	}
	file, err := os.Open(fullPath)
	if err != nil {
		log.Println("Error opening saved matrix file")
		return "", err
	}
	defer file.Close()

	var matrixFile MatrixFile
	decoder := json.NewDecoder(file)
	//декодер читает JSON из файла и распаковывает его по полям структуры matrixFile.
	err = decoder.Decode(&matrixFile) //проверяем совпадение структуры json и MatrixFile
	if err != nil {
		log.Println("Error decoding saved matrix")
		return "", err
	}

	jsonData, err := json.Marshal(matrixFile)
	if err != nil {
		log.Println("Error marshalling saved matrix")
		return "", err
	}
	return string(jsonData), nil
} //стоит переделать декодинк и маршалинг обратно - в простую проверку структуры и передачу json

func getMatrix(filename string) ([][]int, error) {
	fullPath := filepath.Join(SaveDir, filename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, errors.New("save file not found")
	}
	file, err := os.Open(fullPath)
	if err != nil {
		log.Println("Error opening saved matrix file")
		return nil, err
	}
	defer file.Close()

	var matrixFile MatrixFile
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&matrixFile)
	if err != nil {
		log.Println("Error decoding saved matrix")
		return nil, err
	}

	return matrixFile.Data, nil
}

func getMatrixUUID(filename string) (string, error) {
	fullPath := filepath.Join(SaveDir, filename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", errors.New("save file not found")
	}
	file, err := os.Open(fullPath)
	if err != nil {
		log.Println("Error opening saved matrix file")
		return "", err
	}
	defer file.Close()

	var matrixFile MatrixFile
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&matrixFile)
	if err != nil {
		log.Println("Error decoding saved matrix")
		return "", err
	}

	return matrixFile.UUID, nil
}
