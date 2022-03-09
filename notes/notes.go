/*
The notes package provides structs and functions for converting user input published MIDI notes.

The entrypoint is PublishNotes, a function that listens for keypresses, publishing a MIDI note to
a channel each time a key is pressed. The package also provides the MidiNote type - a struct that
holds pitch and velocity - and functions for converting from these values to frequency and amplitude.
*/
package notes

import (
	"sync"

	"github.com/crnbaker/gostringsynth/errors"
	"github.com/crnbaker/gostringsynth/gui"
	"github.com/crnbaker/gostringsynth/numeric"
	"github.com/crnbaker/gostringsynth/voicepub"

	tty "github.com/mattn/go-tty"
)

const maxOctave = 8
const minOctave = -2

// letterPitchMap maps QWERTY keyboard keys to MIDI notes (in MIDI octave -2)
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

// userSettings stores note and string properties for sending to other software module
type userSettings struct {
	midiNoteSettings
	stringSettings
}

// defaultUserSettings returns a UserSettings struct configured with default values
func defaultUserSettings() userSettings {
	return userSettings{defaultMidiNoteSettings(), defaultStringSettings()}
}

// PublishNotes listens for key presses and publishes MIDI notes to noteChannel until user quits
func PublishNotes(waitGroup *sync.WaitGroup, noteChannel chan voicepub.StringNote, synthParamsChan chan gui.SynthParameters) {

	defer waitGroup.Done()
	defer close(noteChannel)
	defer close(synthParamsChan)

	settings := defaultUserSettings()
	synthParamsChan <- &settings

	// key press listener
	tty, err := tty.Open()
	errors.Chk(err)
	defer tty.Close()

UserInputLoop:
	for {
		letter, err := tty.ReadRune()
		errors.Chk(err)
		pitch, ok := letterPitchMap[letter]
		if ok {
			noteChannel <- newStringMidiNote(pitch, settings.midiNoteSettings, settings.stringSettings)
		} else {
			switch letter {
			case 'q':
				// Quit the app
				synthParamsChan <- &settings
				break UserInputLoop
			case 'x':
				settings.octave++
				if settings.octave > maxOctave {
					settings.octave = maxOctave
				}
				synthParamsChan <- &settings
			case 'z':
				settings.octave--
				if settings.octave < minOctave {
					settings.octave = minOctave
				}
				synthParamsChan <- &settings
			case 'v':
				settings.velocity += 5
				if settings.velocity > 127 {
					settings.velocity = 127
				}
				synthParamsChan <- &settings
			case 'c':
				settings.velocity -= 5
				if settings.velocity < 0 {
					settings.velocity = 0
				}
				synthParamsChan <- &settings
			case '.':
				settings.pluckPos += 0.05
				settings.pluckPos = numeric.Clip(settings.pluckPos, 0.05, 0.95)
				synthParamsChan <- &settings
			case ',':
				settings.pluckPos -= 0.05
				settings.pluckPos = numeric.Clip(settings.pluckPos, 0.05, 0.95)
				synthParamsChan <- &settings
			case '>':
				settings.pluckWidth += 0.05
				settings.pluckWidth = numeric.Clip(settings.pluckWidth, 0, 0.9)
				synthParamsChan <- &settings
			case '<':
				settings.pluckWidth -= 0.05
				settings.pluckWidth = numeric.Clip(settings.pluckWidth, 0, 0.9)
				synthParamsChan <- &settings
			case ']':
				settings.pickupPos += 0.05
				settings.pickupPos = numeric.Clip(settings.pickupPos, 0.05, 0.95)
				synthParamsChan <- &settings
			case '[':
				settings.pickupPos -= 0.05
				settings.pickupPos = numeric.Clip(settings.pickupPos, 0.05, 0.95)
				synthParamsChan <- &settings
			case '=':
				settings.decayTimeS += 0.2
				settings.decayTimeS = numeric.Clip(settings.decayTimeS, 0.2, 10)
				synthParamsChan <- &settings
			case '-':
				settings.decayTimeS -= 0.2
				settings.decayTimeS = numeric.Clip(settings.decayTimeS, 0.2, 10)
				synthParamsChan <- &settings
			}
		}
	}
}
