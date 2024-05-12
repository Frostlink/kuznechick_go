from pydub import AudioSegment

# Загрузить файл .wav
audio = AudioSegment.from_wav("mic.wav")

# Конвертировать в формат .ogg с кодеком Opus
audio.export("mic.ogg", format="ogg", codec="libopus")