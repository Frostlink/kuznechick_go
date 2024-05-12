package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

//Функция для чтения файла
//Вход: название файла
//Выход: данные из файла
func read_file(filename string) []uint8 {
	var outdata []uint8

	//Проверка на пустоту входных данных
	if len(filename) == 0 {
		fmt.Println("read_file: Internal Error!!")
		os.Exit(11)
	}

	//Открытие файла
	file, err := os.Open(filename)

	//Проверка открытия файла
	if err != nil {
		fmt.Println("read_file: open file Error!!")
		os.Exit(12)
	}
	defer file.Close() //Закрытие файла

	// Читаем содержимое файла в []byte.
	outdata, err = ioutil.ReadAll(file)
	//Проверка чтения файла
	if err != nil {
		fmt.Println("read_file: read file Error!!")
		os.Exit(13)
	}

	//Вывод данных из файла
	return outdata
}

//Функция для записи в файл
//Вход: название файла, данные для записи
func write_file(filename string, indata []uint8) {

	//Проверка на пустоту входных данных
	if len(filename) == 0 {
		fmt.Println("write_file: Internal Error!!")
		os.Exit(11)
	}

	//Открытие файла
	file, err := os.Create(filename)

	//Проверка открытия файла
	if err != nil {
		fmt.Println("write_file: open file Error!!")
		os.Exit(12)
	}
	defer file.Close() //Закрытие файла

	_, err = file.Write(indata[:])
	if err != nil {
		fmt.Println("write_file: write file Error!!")
		os.Exit(14)
	}

}
