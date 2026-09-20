package service

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"soc_test/embedded" // ← поменяй на свой путь
)

// CopyEmbeddedExample копирует пример проекта, если папка ещё не существует
func CopyEmbeddedExample() error {
	targetDir := filepath.Join(SaveDir, "Проект для примера")

	// Если папка уже есть — ничего не делаем (чтобы не затирать данные пользователя)
	// if _, err := os.Stat(targetDir); err == nil {
	// 	return nil // уже существует
	// }

	log.Println("Создатся пример проекта...")

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("не удалось создать директорию примера: %w", err)
	}

	// Копируем всё содержимое embedded/example/*
	return copyEmbedDir(embedded.ExampleFS, embedded.ExampleDirName, targetDir)
}

// Вспомогательная функция копирования
func copyEmbedDir(efs embed.FS, sourceDir, targetDir string) error {
	return fs.WalkDir(efs, sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Пропускаем саму корневую папку
		if path == sourceDir {
			return nil
		}

		relPath, _ := filepath.Rel(sourceDir, path)
		targetPath := filepath.Join(targetDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// Это файл — копируем
		data, err := efs.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(targetPath, data, 0644)
	})
}