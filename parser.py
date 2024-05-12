import json

with open('text.json', encoding='utf-8') as f:
    data = json.load(f)

# Открываем файл для записи
with open("text.txt", "w") as file:
    # Ищем значения 'text' во всех объектах 'alternatives'
    for chunk in data['response']['chunks']:
        for alternative in chunk['alternatives']:
            file.write(alternative['text']+ "\n")