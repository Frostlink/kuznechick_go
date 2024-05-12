package main

import (
	"fmt"
	"os/exec"
)

func parser_json() {
	// Создаем объект типа Cmd
	cmd := exec.Command("C:/Users/Lenovo/GolandProjects/Diplom/venv/Scripts/python.exe", "parser.py")

	// Запускаем команду и ждем, пока она завершится
	err := cmd.Run()
	if err != nil {
		fmt.Println("Ошибка запуска скрипта:", err)
		return
	}

	fmt.Println("Скрипт json парсер успешно выполнен")
}
