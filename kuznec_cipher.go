package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"sync"
)

const BlockSize int = 16

var Text_lenth int

//Вспомогательная таблица для S преобразования
var Pi_table = [256]uint8{

	252, 238, 221, 17, 207, 110, 49, 22,
	251, 196, 250, 218, 35, 197, 4, 77,
	233, 119, 240, 219, 147, 46, 153, 186,
	23, 54, 241, 187, 20, 205, 95, 193,
	249, 24, 101, 90, 226, 92, 239, 33,
	129, 28, 60, 66, 139, 1, 142, 79,
	5, 132, 2, 174, 227, 106, 143, 160,
	6, 11, 237, 152, 127, 212, 211, 31,
	235, 52, 44, 81, 234, 200, 72, 171,
	242, 42, 104, 162, 253, 58, 206, 204,
	181, 112, 14, 86, 8, 12, 118, 18,
	191, 114, 19, 71, 156, 183, 93, 135,
	21, 161, 150, 41, 16, 123, 154, 199,
	243, 145, 120, 111, 157, 158, 178, 177,
	50, 117, 25, 61, 255, 53, 138, 126,
	109, 84, 198, 128, 195, 189, 13, 87,
	223, 245, 36, 169, 62, 168, 67, 201,
	215, 121, 214, 246, 124, 34, 185, 3,
	224, 15, 236, 222, 122, 148, 176, 188,
	220, 232, 40, 80, 78, 51, 10, 74,
	167, 151, 96, 115, 30, 0, 98, 68,
	26, 184, 56, 130, 100, 159, 38, 65,
	173, 69, 70, 146, 39, 94, 85, 47,
	140, 163, 165, 125, 105, 213, 149, 59,
	7, 88, 179, 64, 134, 172, 29, 247,
	48, 55, 107, 228, 136, 217, 231, 137,
	225, 27, 131, 73, 76, 63, 248, 254,
	141, 83, 170, 144, 202, 216, 133, 97,
	32, 113, 103, 164, 45, 43, 9, 91,
	203, 155, 37, 208, 190, 229, 108, 82,
	89, 166, 116, 210, 230, 244, 180, 192,
	209, 102, 175, 194, 57, 75, 99, 182,
}

// Inverse Pi (S) substitution lookup table.
var Pi_Reverse_table = [256]uint8{

	0xA5, 0x2D, 0x32, 0x8F, 0x0E, 0x30, 0x38, 0xC0,
	0x54, 0xE6, 0x9E, 0x39, 0x55, 0x7E, 0x52, 0x91,
	0x64, 0x03, 0x57, 0x5A, 0x1C, 0x60, 0x07, 0x18,
	0x21, 0x72, 0xA8, 0xD1, 0x29, 0xC6, 0xA4, 0x3F,
	0xE0, 0x27, 0x8D, 0x0C, 0x82, 0xEA, 0xAE, 0xB4,
	0x9A, 0x63, 0x49, 0xE5, 0x42, 0xE4, 0x15, 0xB7,
	0xC8, 0x06, 0x70, 0x9D, 0x41, 0x75, 0x19, 0xC9,
	0xAA, 0xFC, 0x4D, 0xBF, 0x2A, 0x73, 0x84, 0xD5,
	0xC3, 0xAF, 0x2B, 0x86, 0xA7, 0xB1, 0xB2, 0x5B,
	0x46, 0xD3, 0x9F, 0xFD, 0xD4, 0x0F, 0x9C, 0x2F,
	0x9B, 0x43, 0xEF, 0xD9, 0x79, 0xB6, 0x53, 0x7F,
	0xC1, 0xF0, 0x23, 0xE7, 0x25, 0x5E, 0xB5, 0x1E,
	0xA2, 0xDF, 0xA6, 0xFE, 0xAC, 0x22, 0xF9, 0xE2,
	0x4A, 0xBC, 0x35, 0xCA, 0xEE, 0x78, 0x05, 0x6B,
	0x51, 0xE1, 0x59, 0xA3, 0xF2, 0x71, 0x56, 0x11,
	0x6A, 0x89, 0x94, 0x65, 0x8C, 0xBB, 0x77, 0x3C,
	0x7B, 0x28, 0xAB, 0xD2, 0x31, 0xDE, 0xC4, 0x5F,
	0xCC, 0xCF, 0x76, 0x2C, 0xB8, 0xD8, 0x2E, 0x36,
	0xDB, 0x69, 0xB3, 0x14, 0x95, 0xBE, 0x62, 0xA1,
	0x3B, 0x16, 0x66, 0xE9, 0x5C, 0x6C, 0x6D, 0xAD,
	0x37, 0x61, 0x4B, 0xB9, 0xE3, 0xBA, 0xF1, 0xA0,
	0x85, 0x83, 0xDA, 0x47, 0xC5, 0xB0, 0x33, 0xFA,
	0x96, 0x6F, 0x6E, 0xC2, 0xF6, 0x50, 0xFF, 0x5D,
	0xA9, 0x8E, 0x17, 0x1B, 0x97, 0x7D, 0xEC, 0x58,
	0xF7, 0x1F, 0xFB, 0x7C, 0x09, 0x0D, 0x7A, 0x67,
	0x45, 0x87, 0xDC, 0xE8, 0x4F, 0x1D, 0x4E, 0x04,
	0xEB, 0xF8, 0xF3, 0x3E, 0x3D, 0xBD, 0x8A, 0x88,
	0xDD, 0xCD, 0x0B, 0x13, 0x98, 0x02, 0x93, 0x80,
	0x90, 0xD0, 0x24, 0x34, 0xCB, 0xED, 0xF4, 0xCE,
	0x99, 0x10, 0x44, 0x40, 0x92, 0x3A, 0x01, 0x26,
	0x12, 0x1A, 0x48, 0x68, 0xF5, 0x81, 0x8B, 0xC7,
	0xD6, 0x20, 0x0A, 0x08, 0x00, 0x4C, 0xD7, 0x74,
}

// L-function (transformation) vector.
var kB = [16]uint8{148, 32, 133, 16, 194, 192, 1, 251, 1, 192, 194, 16, 133, 32, 148, 1}

//X-func - XOR для 32 байт
//Вход: 2 строки по 32 байт для преобразования
//Выход: Результирующая строка XOR 32 байта
func X_func32(str1, str2 [BlockSize * 2]uint8) [BlockSize * 2]uint8 {
	var outdata [BlockSize * 2]uint8

	//Проверка на пустоту входных данных
	if len(str1) == 0 || len(str2) == 0 {
		fmt.Println("X_func_32: Internal Error!!")
		os.Exit(11)
	}

	// выполняем побитовое ложение по модулю 2
	for i := 0; i < BlockSize*2; i++ {
		//XOR для каждого байта
		outdata[i] = str1[i] ^ str2[i]
	}

	//Вывод результирующей строки
	return outdata
}

//X-func - XOR для 16 байт
//Вход: 2 строки по 16 байт для преобразования
//Выход: Результирующая строка XOR 16 байта
func X_func(str1, str2 [BlockSize]uint8) [BlockSize]uint8 {
	var outdata [BlockSize]uint8

	//Проверка на пустоту входных данных
	if len(str1) == 0 || len(str2) == 0 {
		fmt.Println("X_func: Internal Error!!")
		os.Exit(11)
	}

	// выполняем побитовое ложение по модулю 2
	for i := 0; i < BlockSize; i++ {
		//XOR для каждого байта
		outdata[i] = str1[i] ^ str2[i]
	}

	//Вывод результирующей строки
	return outdata
}

//S функция для нелинейного преобразования
//Вход:входной блок текста
//Выход:выходной блок текста
func S_func(indata [BlockSize]uint8) [BlockSize]uint8 {
	var outdata [BlockSize]uint8

	//Проверка на пустоту входных данных
	if len(indata) == 0 {
		fmt.Println("S_func: Internal Error!!")
		os.Exit(11)
	}

	for i := 0; i < BlockSize; i++ {
		//Нелинейное преобразования
		outdata[i] = Pi_table[indata[i]]
	}
	//Вывод результирующей строки
	return outdata
}

//S функция для обратного нелинейного преобразования
//Вход:входной блок текста
//Выход:выходной блок текста
func S_Reverse_func(indata [BlockSize]uint8) [BlockSize]uint8 {
	var outdata [BlockSize]uint8

	//Проверка на пустоту входных данных
	if len(indata) == 0 {
		fmt.Println("S_Reverse_func: Internal Error!!")
		os.Exit(11)
	}

	for i := 0; i < BlockSize; i++ {
		//Обратное нелинейное преобразования
		outdata[i] = Pi_Reverse_table[indata[i]]
	}
	//Вывод результирующей строки
	return outdata
}

//R вспомогательная функция для линейного преобразования
//Вход: входной блок текста
//Выход: выходной блок текста
func R_func(indata [BlockSize]uint8) [BlockSize]uint8 {
	var outdata [BlockSize]uint8
	var sum byte = 0 //

	//Проверка на пустоту входных данных
	if len(indata) == 0 {
		fmt.Println("R_func: Internal Error!!")
		os.Exit(11)
	}

	for i := 0; i < BlockSize; i++ {
		sum ^= multTable[uint16(indata[i])*256+uint16(kB[i])]
	}

	outdata[0] = sum
	copy(outdata[1:], indata[:15])
	return outdata
}

//R вспомогательная функция для обратного линейного преобразования
//Вход: входной блок текста
//Выход: выходной блок текста
func R_Reverse_func(indata [BlockSize]uint8) [BlockSize]uint8 {
	var outdata [BlockSize]uint8
	var tmp [BlockSize]uint8
	var sum byte = 0

	//Проверка на пустоту входных данных
	if len(indata) == 0 {
		fmt.Println("R_Reverse_func: Internal Error!!")
		os.Exit(11)
	}

	copy(tmp[:15], indata[1:])
	tmp[15] = indata[0]

	for i := 0; i < BlockSize; i++ {
		sum ^= multTable[uint16(tmp[i])*256+uint16(kB[i])]
	}

	copy(outdata[:15], tmp[:15])
	outdata[15] = sum
	return outdata
}

//Функция нелинейного преобразования
//Вход: входной блок текста
//Выход: выходной блок текста
func L_func(indata [BlockSize]uint8) [BlockSize]uint8 {
	var outdata [BlockSize]uint8
	var tmp [BlockSize]uint8

	//Проверка на пустоту входных данных
	if len(indata) == 0 {
		fmt.Println("L_func: Internal Error!!")
		os.Exit(11)
	}

	//реализация сдвига
	copy(tmp[:], indata[:16])
	for i := 0; i < BlockSize; i++ {
		tmp = R_func(tmp)
		copy(outdata[:], tmp[:])
	}
	return outdata
}

//Функция обратного нелинейного преобразования
//Вход: входной блок текста
//Выход: выходной блок текста
func L_Reverse_func(indata [BlockSize]uint8) [BlockSize]uint8 {
	var outdata [BlockSize]uint8
	var tmp [BlockSize]uint8

	//Проверка на пустоту входных данных
	if len(indata) == 0 {
		fmt.Println("L_Reverse_func: Internal Error!!")
		os.Exit(11)
	}

	//реализация сдвига
	copy(tmp[:], indata[:16])
	for i := 0; i < BlockSize; i++ {
		tmp = R_Reverse_func(tmp)
		copy(outdata[:], tmp[:])
	}
	return outdata
}

//Функция раунда шифрования
//Вход: блок текста для шифрования, раундовый ключ
//Выход: блок текста шифрования
func LSX_func(a, b [BlockSize]uint8) [BlockSize]uint8 {
	var temp1 [BlockSize]uint8
	var temp2 [BlockSize]uint8
	var outdata [BlockSize]uint8

	//Проверка на пустоту входных данных
	if len(b) == 0 || len(a) == 0 {
		fmt.Println("LSX_func: Internal Error!!")
		os.Exit(11)
	}

	//XOR входного блока и ключа
	temp1 = X_func(a, b)
	//Функция линейного преобразования
	temp2 = S_func(temp1)
	//Функция нелинейного преобразования
	outdata = L_func(temp2)

	return outdata
}

//Функция раунда дешифрации
//Вход: блок текста для дешифрации, раундовый ключ
//Выход:блок текста дешифрации
func LSX_Reverse_func(a, b [BlockSize]uint8) [BlockSize]uint8 {
	var temp1 [BlockSize]uint8
	var temp2 [BlockSize]uint8
	var outdata [BlockSize]uint8

	//Проверка на пустоту входных данных
	if len(b) == 0 || len(a) == 0 {
		fmt.Println("LSX_Reverse_func: Internal Error!!")
		os.Exit(11)
	}

	//XOR входного блока и ключа
	temp1 = X_func(a, b)
	//Обратная функция нелинейного преобразования
	temp2 = L_Reverse_func(temp1)
	//Обратная функция линейного преобразования
	outdata = S_Reverse_func(temp2)

	return outdata
}

//Функция для для реализации ячейки фейстеля
//Вход: 1 раундовый ключ, 2 раундовый ключ, номер раунда
//Выход: 1 выходной ключ, 2 выходной ключ
func F_func(inputKey, inputKeySecond, iterationConst [BlockSize]uint8) ([BlockSize]uint8, [BlockSize]uint8) {
	//вспомогательная переменная
	var temp1 [BlockSize]uint8
	//вспомогательная переменная
	var temp2 [BlockSize]uint8
	//Первый выходной ключ
	var outputKey [BlockSize]uint8
	//Второй выходной ключ
	var outputKeySecond [BlockSize]uint8

	//Проверка на пустоту входных данных
	if len(inputKey) == 0 || len(inputKeySecond) == 0 || len(iterationConst) == 0 {
		fmt.Println(" F_func: Internal Error!!")
		os.Exit(11)
	}

	temp1 = LSX_func(inputKey, iterationConst)
	temp2 = X_func(temp1, inputKeySecond)

	copy(outputKeySecond[:], inputKey[:])
	copy(outputKey[:], temp2[:])

	//Вывод ключей
	return outputKey, outputKeySecond
}

//Функция расчета констант
//Вход: номер раунда
//Выход: сгенерированное значение из номера раунда
func C_func(number uint8) [BlockSize]uint8 {
	var temp1 [BlockSize]uint8
	var output [BlockSize]uint8

	temp1[15] = number
	output = L_func(temp1)

	return output
}

//Функция для генерации раундовых ключей
//Вход: мастер-ключ 32 байта
//Выход: Раундовые ключи 10 по 16 байт
func ExpandKey(masterKey [BlockSize * 2]uint8) [BlockSize * 10]uint8 {
	//Переменная для хранения раундовых ключей
	var keys [BlockSize * 10]uint8
	var C [BlockSize]uint8
	var temp1 [BlockSize]uint8
	var temp2 [BlockSize]uint8
	//Вспомогательная переменная для итераций
	var i, j uint8

	//Проверка на пустоту входных данных
	if len(masterKey) == 0 {
		fmt.Println("ExpandKey: Internal Error!!")
		os.Exit(11)
	}

	//Запись первой половины мастер ключа как 1 раундового ключа
	copy(keys[:16], masterKey[:16])
	//Запись второй половины мастер ключа как 2 раундового ключа
	copy(keys[16:], masterKey[16:])

	//генерация раундовых ключей
	for j = 0; j < 4; j++ {
		copy(temp1[:], keys[j*2*16:(j*2+1)*16])
		copy(temp2[:], keys[(j*2+1)*16:(j*2+2)*16])

		for i = 1; i < 8; i++ {
			C = C_func(j*8 + i)
			temp1, temp2 = F_func(temp1, temp2, C)
		}

		C = C_func(j*8 + 8)
		//Использование сети фейстеля для генерации
		temp1, temp2 = F_func(temp1, temp2, C)

		copy(keys[(j*2+2)*16:(j*2+3)*16], temp1[:])
		copy(keys[(j*2+3)*16:(j*2+4)*16], temp2[:])

	}
	//Вывод раундовых ключей
	return keys
}

//Функция для шифрования блока текста
//Вход: основной блок текста 16 байт, раундовые ключи
//Выход: Зашифрованный текст
func Encrypt(plainText [BlockSize]uint8, keys [BlockSize * 10]uint8) [BlockSize]uint8 {
	//Вспомогательная переменная для текста
	var tempX [16]uint8
	//Вспомогательная переменная для текста
	var tempY [16]uint8
	//Переменная для раундового ключа
	var key [16]uint8
	//Переменная для зашифрованного блока
	var cipherText [16]uint8
	//Вспомогательная переменная счетчика раундов
	var i uint8

	//Проверка на пустоту входных данных
	if len(plainText) == 0 || len(keys) == 0 {
		fmt.Println("Encrypt: Internal Error!!")
		os.Exit(11)
	}

	copy(tempX[:], plainText[:])

	//9 полных раундов дешифрации блока текста
	for i = 0; i < 9; i++ {
		//Забрать текущий раундовый ключ
		copy(key[:], keys[i*16:(i+1)*16])
		//Провести LSX раунд с обратными функциями
		tempY = LSX_func(tempX, key)
		copy(tempX[:], tempY[:])
	}
	copy(key[:], keys[9*16:10*16])
	cipherText = X_func(tempY, key)

	//Зашифрованный блок текста
	return cipherText
}

//Функция для дешифрации блока текста
//Вход:Зашифрованный блок текста 16 байт, раундовые ключи
//Выход:Расшифрованный текст
func Decrypt(cipherText [BlockSize]uint8, keys [BlockSize * 10]uint8) [BlockSize]uint8 {
	//Вспомогательная переменная для текста
	var tempX [16]uint8
	//Вспомогательная переменная для текста
	var tempY [16]uint8
	//Переменная для раундового ключа
	var key [16]uint8
	//Переменная для расшифрованного блока
	var plainText [16]uint8
	//Вспомогательная переменная счетчика раундов
	var i uint8

	//Проверка на пустоту входных данных
	if len(cipherText) == 0 || len(keys) == 0 {
		fmt.Println("Decrypt: Internal Error!!")
		os.Exit(11)
	}

	copy(tempX[:], cipherText[:])

	//9 полных раундов дешифрации блока текста
	for i = 0; i < 9; i++ {
		//Забрать текущий раундовый ключ
		copy(key[:], keys[(9-i)*16:(10-i)*16])
		//Провести LSX раунд с обратными функциями
		tempY = LSX_Reverse_func(tempX, key)
		copy(tempX[:], tempY[:])
	}

	copy(key[:], keys[:16])
	plainText = X_func(tempY, key)

	//Расшифрованный блок текста
	return plainText
}

//Функция для создания нового мастер ключа по номеру блока
//Вход: Мастер-Ключ 32 байта, номер блока
//Выход: новый мастер-ключ
func Counter_key(MasterKey [BlockSize * 2]uint8, counter [BlockSize * 2]uint8) [BlockSize * 2]uint8 {
	//Переменная для хранения вектор-строки
	var nonce [BlockSize * 2]uint8
	//Переменная для хранения вычисленной строки
	var key [BlockSize * 2]uint8
	//Переменная для хранения нового мастер ключа
	var outdata [BlockSize * 2]uint8
	//Вектор-строка для создания нового ключа
	copy(nonce[:], "001122334455AABBCCDDFF1234567890")

	//Проверка на пустоту входных данных
	if len(MasterKey) == 0 || len(counter) == 0 {
		fmt.Println("Counter_key: Internal Error!!")
		os.Exit(11)
	}

	//XOR вектор-строки и номера блока
	key = X_func32(nonce, counter)
	//XOR вычисленной строки и мастер ключа
	outdata = X_func32(key, MasterKey)

	//Вывод нового мастер-ключа
	return outdata
}

//Функция для Шифрования блока текста при параллельной обработки
//Вход: Основной блок текста 16 байт для шифрования, Мастер-ключ 32 байта, номер блока
//Выход: Зашифрованный блок текста 16 байт
func CTR_Encrypt(PlainBlock [BlockSize]uint8, MasterKey [BlockSize * 2]uint8, counter int) [BlockSize]uint8 {
	//Переменная для хранения раундовых ключей
	var keys [BlockSize * 10]uint8
	//Переменная для нового мастер-ключа
	var Key [BlockSize * 2]uint8
	//Переменная для зашифрованного блока
	var CipherBlock [BlockSize]uint8
	var buf [BlockSize * 2]uint8

	//Проверка на пустоту входных данных
	if len(PlainBlock) == 0 || len(MasterKey) == 0 {
		fmt.Println("CTR_Encrypt: Internal Error!!")
		os.Exit(11)
	}

	//Перевод номера блока в 32 разрядный формат
	binary.LittleEndian.PutUint32(buf[:], uint32(counter))

	//Создание нового мастер-ключа ключа для блока
	Key = Counter_key(MasterKey, buf)
	//Генерация раундовых ключей из нового ключа
	keys = ExpandKey(Key)
	//Шифрование блока
	CipherBlock = Encrypt(PlainBlock, keys)

	//Вывод Зашифрованного блока
	return CipherBlock
}

//Функция для Расшифровки блока текста при параллельной обработки
//Вход: Зашифрованный блок 16 байт, Мастер-ключ 32 байта, номер блока
//Выход: Расшифрованный блок текста 16 байт
func CTR_Decrypt(PlainBlock [BlockSize]uint8, MasterKey [BlockSize * 2]uint8, counter int) [BlockSize]uint8 {
	//Переменная для хранения раундовых ключей
	var keys [BlockSize * 10]uint8
	//Переменная для нового мастер-ключа
	var Key [BlockSize * 2]uint8
	//Переменная для расшифрованного блока
	var CipherBlock [BlockSize]uint8
	var buf [BlockSize * 2]uint8

	//Проверка на пустоту входных данных
	if len(PlainBlock) == 0 || len(MasterKey) == 0 {
		fmt.Println("CTR_Decrypt: Internal Error!!")
		os.Exit(11)
	}

	//Перевод номера блока в 32 разрядный формат
	binary.LittleEndian.PutUint32(buf[:], uint32(counter))

	//Создание нового мастер-ключа ключа для блока
	Key = Counter_key(MasterKey, buf)
	//Генерация раундовых ключей из нового ключа
	keys = ExpandKey(Key)
	//Расшифрование блока
	CipherBlock = Decrypt(PlainBlock, keys)

	//Вывод расшифрованного блока
	return CipherBlock
}

//Функция для параллельного шифрования
//Вход: Текст для шифрования, Мастер-ключ 32 байта
//Выход: Зашифрованный текст
func EncryptParallel(Text []uint8, MasterKey [BlockSize * 2]uint8) []uint8 {
	var numBlocks int
	var wg sync.WaitGroup //Переменная для  wait group

	//Проверка на пустоту входных данных
	if len(Text) == 0 || len(MasterKey) == 0 {
		fmt.Println("CTR_Decrypt: Internal Error!!")
		os.Exit(11)
	}

	//Вычисление размера текста в блоках по 16 байт
	numBlocks = len(Text) / BlockSize

	number := 0

	//Если текст не кратен 16, то добавить символы, чтобы было кратно
	if len(Text)%BlockSize != 0 {
		numBlocks += 1
		number = BlockSize - len(Text)%BlockSize
	}

	CipherText := make([]uint8, len(Text)+number)
	TextBlocks := make([]uint8, len(Text)+number)
	copy(TextBlocks[0:], Text[:])

	for i := 0; i < numBlocks; i++ {
		wg.Add(1)
		go func(numBlock int) {
			//Буферный блок 16 байт
			var Block [BlockSize]uint8
			var ciphertextBlock [BlockSize]uint8
			//Забрать блок 16 байт из расшифрованного текста в буфер
			copy(Block[:], TextBlocks[(numBlock)*BlockSize:(numBlock+1)*BlockSize])
			//Зашифровать блок текста из буфера
			ciphertextBlock = CTR_Encrypt(Block, MasterKey, numBlock)
			//Записать зашифрованный текст в основной буфер
			copy(CipherText[(numBlock)*BlockSize:(numBlock+1)*BlockSize], ciphertextBlock[:])
			//Убрать из очереди горутину
			defer wg.Done()
		}(i)
	}

	//Очередь для ожидания горутин
	wg.Wait()
	//Вывод зашифрованного текста
	return CipherText
}

//Функция для параллельной расшифровки
//Вход: Зашифрованный текст, Мастер-ключ 32 байта
//Выход: Расшифрованный текст
func DecryptParallel(Text []uint8, MasterKey [BlockSize * 2]uint8) []uint8 {
	var numBlocks int
	var wg sync.WaitGroup                   //Переменная для  wait group
	PlainText := make([]uint8, len(Text))   //Переменная для расшифрованного текста
	ReturnText := make([]uint8, Text_lenth) //Переменная для расшифрованного текста

	//Проверка на пустоту входных данных
	if len(Text) == 0 || len(MasterKey) == 0 {
		fmt.Println("CTR_Decrypt: Internal Error!!")
		os.Exit(11)
	}

	//Вычисление размера текста в блоках по 16 байт
	numBlocks = len(Text) / BlockSize

	TextBlocks := make([]uint8, len(Text))
	copy(TextBlocks[0:], Text[:])

	for i := 0; i < numBlocks; i++ {
		//Увеличить очередь горутины на 1
		wg.Add(1)
		go func(numBlock int) {
			//Буферный блок 16 байт
			var Block [BlockSize]uint8
			var ciphertextBlock [BlockSize]uint8
			//Забрать блок 16 байт из расшифрованного текста в буфер
			copy(Block[:], TextBlocks[(numBlock)*BlockSize:(numBlock+1)*BlockSize])
			//Расшифровать блок текста из буфера
			ciphertextBlock = CTR_Decrypt(Block, MasterKey, numBlock)
			//Записать зашифрованный текст в основной буфер
			copy(PlainText[(numBlock)*BlockSize:(numBlock+1)*BlockSize], ciphertextBlock[:])
			//Убрать из очереди горутину
			defer wg.Done()
		}(i)
	}

	//Очередь для ожидания горутин
	wg.Wait()
	copy(ReturnText[:], PlainText[:Text_lenth])
	//Вывод расшифрованного текста
	return ReturnText
}

//Функция шифрования
func Text_encrypt() {
	var Buf []uint8
	var Text []uint8
	var CipherText []uint8
	var MasterKey [BlockSize * 2]uint8

	Text = read_file("text.txt")
	Buf = read_file("masterkey.txt")

	Text_lenth = len(Text)
	copy(MasterKey[:], Buf[:32])

	CipherText = EncryptParallel(Text, MasterKey)

	write_file("Cipher.txt", CipherText)
}

//Функция расшифровки
func Text_decrypt() {
	var Buf []uint8
	var Text []uint8
	var PlainText []uint8
	var MasterKey [BlockSize * 2]uint8

	Text = read_file("Cipher.txt")
	Buf = read_file("masterkey.txt")

	copy(MasterKey[:], Buf[:32])
	PlainText = DecryptParallel(Text, MasterKey)

	write_file("Plain.txt", PlainText)
}
