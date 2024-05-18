package main

import (
	"club/internal/service"
	"fmt"
	"os"
)

func main() {
	// Проверяем, передан ли входной файл в аргументах командной строки
	if len(os.Args) < 2 {
		fmt.Println("Usage: club <input_file>")
		return
	}

	// Получаем имя входного файла
	inputFile := os.Args[1]

	// Обрабатываем входной файл и выводим ошибки, если они возникли
	err := service.Process(inputFile)
	if err != nil {
		fmt.Println("Error processing file:", err)
	}
}
