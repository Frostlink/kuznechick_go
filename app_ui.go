package main

import (
	"bufio"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"os"
	"time"
)

var ticker *time.Ticker
var timerActive bool

func UI() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Шифрование речи")

	timer := widget.NewLabel("00:00:00")
	stopFlagChan := make(chan bool)
	progressBar := widget.NewProgressBar()
	progressBar.Min = 0
	progressBar.Max = 100
	recognizedText := widget.NewMultiLineEntry()
	EncryptedText := widget.NewMultiLineEntry()
	DecryptedText := widget.NewMultiLineEntry()
	Info_was_Recognized := widget.NewLabel("Содержимое аудофайла:")
	Info_was_Encrypted := widget.NewLabel("Содержимое зашифрованного текста:")
	Info_was_Decrypted := widget.NewLabel("Содержимое расшифрованного текста:")
	////////////////////////////////////////////////////////////////////////////////////////
	// кнопки

	DecryptButton := widget.NewButton("Расшифровать текст", func() {
		Text_decrypt()
		dialog.ShowInformation(
			"Уведомление",
			"Расшифровка текста выполнена успешно",
			myWindow,
		)
		DecryptedText.Show()
		decodedStr := reader("Plain.txt")
		DecryptedText.SetText(decodedStr)
	})

	EncryptButton := widget.NewButton("Зашифровать текст", func() {
		Text_encrypt()
		dialog.ShowInformation(
			"Уведомление",
			"Шифрование выполнено успешно",
			myWindow,
		)
		EncryptedText.Show()
		decodedStr := reader("Cipher.txt")
		EncryptedText.SetText(decodedStr)
		DecryptButton.Show()
	})

	recognizeButton := widget.NewButton("Отправить текст на распознавание", func() {
		converter()
		upload_to_storage()
		recognize()
		dialog.ShowInformation(
			"Уведомление",
			"Распознавание завершено",
			myWindow,
		)
		recognizedText.Show()
		decodedStr := reader("text.txt")
		recognizedText.SetText(decodedStr)
		EncryptButton.Show()
	})
	//Кнопка для остановки записи
	stopButton := widget.NewButton("Остановить запись", func() {
		stopTimer()
		stopFlagChan <- true
		recognizeButton.Show()
	})

	// Кнопка для начала записи
	startButton := widget.NewButton("Начать запись", func() {
		stopButton.Hide()
		EncryptButton.Hide()
		DecryptButton.Hide()
		recognizeButton.Hide()
		recognizedText.Hide()
		EncryptedText.Hide()
		Info_was_Decrypted.Hide()
		DecryptedText.Hide()
		Info_was_Encrypted.Hide()
		Info_was_Recognized.Hide()
		go runTimer(timer)
		flag := 1
		go record(stopFlagChan, flag)
		stopButton.Show()
	})

	//Кнопка выхода из программы
	exitButton := widget.NewButton("Выход", func() {
		os.Exit(0)
	})

	fastButton := widget.NewButton("Начать работу", func() {
		// Запуск таймера
		go runTimer(timer)
		//record()
		converter()
		progressBar.SetValue(10)
		fmt.Println("Progress 10%")
		fmt.Println("Converted done")
		upload_to_storage()
		progressBar.SetValue(30)
		fmt.Println("Progress 20%")
		fmt.Println("Upload to storage done")
		recognize()
		progressBar.SetValue(70)
		fmt.Println("Progress 80%")
		fmt.Println("Recognize audio done")
		Text_encrypt()
		progressBar.SetValue(100)
		fmt.Println("Progress 100%")
		fmt.Println("Text Encrypted")
		Text_decrypt()
		fmt.Println("Text Decrypted")
		stopTimer()
		dialog.ShowInformation(
			"Уведомление",
			"Шифрование выполнено успешно",
			myWindow,
		)
	})
	//////////////////////////////////////////////////////////////////////////////////////

	stopButton.Hide()
	EncryptButton.Hide()
	DecryptButton.Hide()
	recognizeButton.Hide()
	recognizedText.Hide()
	EncryptedText.Hide()
	Info_was_Decrypted.Hide()
	DecryptedText.Hide()
	Info_was_Encrypted.Hide()
	Info_was_Recognized.Hide()

	// Кнопка для отображения нового контента
	button_fast := widget.NewButton("Быстрый режим", func() {
		// Создаем новое содержимое для окна
		newContent := container.NewVBox(
			widget.NewLabel("Быстрый режим"),
			widget.NewLabel("Перед началом работы, поместите нужный файл формата .wav в папку проекта"),
			fastButton,
			timer,
			progressBar,
			exitButton,
		)

		// Обновляем содержимое окна
		myWindow.SetContent(newContent)
	})

	// Кнопка для отображения нового контента
	button_steps := widget.NewButton("Демонстрационный режим", func() {
		// Создаем новое содержимое для окна
		newContent := container.NewVBox(
			widget.NewLabel("Демонстрационный режим"),
			startButton,
			timer,
			stopButton,
			recognizeButton,
			Info_was_Recognized,
			recognizedText,
			EncryptButton,
			Info_was_Encrypted,
			EncryptedText,
			DecryptButton,
			Info_was_Decrypted,
			DecryptedText,
			exitButton,
		)

		// Обновляем содержимое окна
		myWindow.SetContent(newContent)
	})

	// Создаем содержимое для окна
	contentContainer := container.NewVBox(
		button_steps,
		button_fast,
		exitButton,
	)

	myWindow.Resize(fyne.NewSize(600, 700))
	myWindow.CenterOnScreen()
	myWindow.SetContent(contentContainer)
	myWindow.ShowAndRun()
}

func runTimer(timer *widget.Label) {
	if timerActive {
		return // Если таймер уже активен, выходим из функции
	}

	timerActive = true
	ticker = time.NewTicker(time.Second)
	duration := time.Duration(0)

	for {
		select {
		case <-ticker.C:
			duration += time.Second
			timeString := fmt.Sprintf("%02d:%02d:%02d", int(duration.Hours()), int(duration.Minutes())%60, int(duration.Seconds())%60)
			timer.SetText(timeString)
		}
	}
}

func stopTimer() {
	if ticker != nil {
		ticker.Stop()
		ticker = nil
		timerActive = false
	}
}

func reader(filename string) string {
	// Открываем файл
	file, err := os.Open(filename)
	if err != nil {
		// Обработка ошибок
	}
	defer file.Close()

	// Читаем содержимое файла
	scanner := bufio.NewScanner(file)
	content := ""
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}

	decoder := charmap.Windows1251.NewDecoder()
	decodedStr, _, _ := transform.String(decoder, content)
	return decodedStr
}
