from pydub import AudioSegment
from pydub.playback import play

#Load an audio file
audio1 = "/Users/komatsu/Downloads/test1.mp3"
audio2 = "/Users/komatsu/Downloads/test2.mp3"

sound1 = AudioSegment.from_mp3(audio1)
sound2 = AudioSegment.from_mp3(audio2)

combined = sound1.overlay(sound2)
play(combined)