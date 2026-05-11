package service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
)

// Here will be code to create, delete and open projects (saves) using different directories in appdata

type ProjectService struct{}

func NewProjectService() *ProjectService {
	return &ProjectService{}
}

var ProjectDir string

func (s *ProjectService) CreateProject(folderName string) error {
	dirlist := s.GetProjectList()
	if slices.Contains(dirlist, folderName) {
		return fmt.Errorf("Проект с названием %s уже существует.", folderName)
	}
	if len(folderName) > 255 {
		return fmt.Errorf("Длина названия превышает допустимое значение.")
	} else if len(folderName) == 0 {
		return fmt.Errorf("Название не может быть пустым.")
	}
	ProjectDir = filepath.Join(SaveDir, folderName)
	if err := os.MkdirAll(ProjectDir, 0755); err != nil {
		return fmt.Errorf("Не удалось создать директорию проекта: %s", err)
	}
	return nil
}

func (s *ProjectService) GetProjectList() []string {
	dirList := []string{}

	entries, err := os.ReadDir(SaveDir)
	if err != nil {
		log.Fatal("Не удалось получить список проектов", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			dirList = append(dirList, entry.Name())
		}
	}
	return dirList
}

func (s *ProjectService) DeleteProject(folderName string) error {

	fullDir := filepath.Join(SaveDir, folderName)

	err := os.RemoveAll(fullDir)
	if err != nil {
		log.Printf("Не удалось удалить проект", err)
	}

	return err
}

func (s *ProjectService) ChangeProjectDir(folderName string) string {
	fmt.Println(folderName, ProjectDir)

	ProjectDir = filepath.Join(SaveDir, folderName)

	return ProjectDir
}
