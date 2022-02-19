package keypress

import (
	"fmt"

	"github.com/crnbaker/gostringsynth/errors"
	tty "github.com/mattn/go-tty"
)

var letterPitchMap = map[rune]int{
	'a': 0,
	'w': 1,
	's': 2,
	'e': 3,
	'd': 4,
	'f': 5,
	't': 6,
	'g': 7,
	'y': 8,
	'h': 9,
	'u': 10,
	'j': 11,
	'k': 12,
}

func NoteDispatcher(noteSendChannel chan MidiNote) {

	octave := 2
	velocity := 64

	tty, err := tty.Open()
	errors.Chk(err)
	defer tty.Close()

	fmt.Println("Play notes with keyboard mapped across keys from A to K")
	fmt.Println("Increase octave with X")
	fmt.Println("Decrease octave with Z")
	fmt.Println("Increase velocity with >")
	fmt.Println("Increase velocity with <")
	fmt.Println("Q to quit")

	for {
		letter, err := tty.ReadRune()
		errors.Chk(err)
		pitch, ok := letterPitchMap[letter]
		if ok {
			noteSendChannel <- changeNoteOctave(newNote(pitch, velocity), octave)
		} else {
			switch letter {
			case 'q':
				close(noteSendChannel)
			case 'x':
				octave++
			case 'z':
				octave--
			case '.':
				velocity += 5
				if velocity > 127 {
					velocity = 127
				}
				fmt.Println("Velocity increased to", velocity)
			case ',':
				velocity -= 5
				if velocity < 0 {
					velocity = 0
				}
				fmt.Println("Velocity decreased to", velocity)
			}
		}
	}
}
