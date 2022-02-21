/*
The notes package provides structs and functions for converting user input published MIDI notes.

The entrypoint is NotePublisher, a function that listens for keypresses, publishing a MIDI note to
a channel each time a key is pressed. The package also provides the MidiNote type - a struct that
holds pitch and velocity - and functions for converting from these values to frequency and amplitude.
*/
package notes

import (
	"fmt"
	"sync"

	"github.com/crnbaker/gostringsynth/errors"
	tty "github.com/mattn/go-tty"
)

// letterPitchMap Maps QWERTY keyboard keys to MIDI notes (in MIDI octave -2)
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

// NotePublisher listens for key presses and publishes MIDI notes to noteChannel until user quits
func NotePublisher(waitGroup *sync.WaitGroup, noteChannel chan MidiNote) {

	defer waitGroup.Done()

	octave := 3
	velocity := 64

	// key press listener
	tty, err := tty.Open()
	errors.Chk(err)
	defer tty.Close()

	fmt.Println("Play notes with keyboard mapped across keys from A to K")
	fmt.Println("Increase octave with X")
	fmt.Println("Decrease octave with Z")
	fmt.Println("Increase velocity with V")
	fmt.Println("Decrease velocity with C")
	fmt.Println("Q to quit")

UserInputLoop:
	for {
		letter, err := tty.ReadRune()
		errors.Chk(err)
		pitch, ok := letterPitchMap[letter]
		if ok {
			noteChannel <- changeNoteOctave(newNote(pitch, velocity), octave)
		} else {
			switch letter {
			case 'q':
				// Quit the app
				close(noteChannel)
				break UserInputLoop
			case 'x':
				octave++
			case 'z':
				octave--
			case 'v':
				velocity += 5
				if velocity > 127 {
					velocity = 127
				}
				fmt.Println("Velocity increased to", velocity)
			case 'c':
				velocity -= 5
				if velocity < 0 {
					velocity = 0
				}
				fmt.Println("Velocity decreased to", velocity)
			}
		}
	}
}
