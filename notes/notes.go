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
	"github.com/crnbaker/gostringsynth/numeric"
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

// UserSettings stores note and string properties for sending to other software module
type UserSettings struct {
	midiNoteSettings
	stringSettings
}

// DefaultUserSettings returns a UserSettings struct configured with default values
func DefaultUserSettings() UserSettings {
	return UserSettings{defaultMidiNoteSettings(), defaultStringSettings()}
}

// PublishNotes listens for key presses and publishes MIDI notes to noteChannel until user quits
func PublishNotes(waitGroup *sync.WaitGroup, noteChannel chan StringMidiNote, userSettingsChannel chan UserSettings) {

	defer waitGroup.Done()
	defer close(noteChannel)
	defer close(userSettingsChannel)

	settings := DefaultUserSettings()
	userSettingsChannel <- settings

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
			noteChannel <- NewStringMidiNote(pitch, settings.midiNoteSettings, settings.stringSettings)
		} else {
			switch letter {
			case 'q':
				// Quit the app
				userSettingsChannel <- settings
				break UserInputLoop
			case 'x':
				settings.Octave++
				if settings.Octave > maxOctave {
					settings.Octave = maxOctave
				}
				userSettingsChannel <- settings
			case 'z':
				settings.Octave--
				if settings.Octave < minOctave {
					settings.Octave = minOctave
				}
				userSettingsChannel <- settings
			case 'v':
				settings.Velocity += 5
				if settings.Velocity > 127 {
					settings.Velocity = 127
				}
				userSettingsChannel <- settings
			case 'c':
				settings.Velocity -= 5
				if settings.Velocity < 0 {
					settings.Velocity = 0
				}
				userSettingsChannel <- settings
			case '.':
				settings.PluckPos += 0.05
				settings.PluckPos = numeric.Clip(settings.PluckPos, 0.05, 0.95)
				userSettingsChannel <- settings
			case ',':
				settings.PluckPos -= 0.05
				settings.PluckPos = numeric.Clip(settings.PluckPos, 0.05, 0.95)
				userSettingsChannel <- settings
			case '>':
				settings.PluckWidth += 0.05
				settings.PluckWidth = numeric.Clip(settings.PluckWidth, 0, 0.9)
				userSettingsChannel <- settings
			case '<':
				settings.PluckWidth -= 0.05
				settings.PluckWidth = numeric.Clip(settings.PluckWidth, 0, 0.9)
				userSettingsChannel <- settings
			case ']':
				settings.PickupPos += 0.05
				settings.PickupPos = numeric.Clip(settings.PickupPos, 0.05, 0.95)
				userSettingsChannel <- settings
			case '[':
				settings.PickupPos -= 0.05
				settings.PickupPos = numeric.Clip(settings.PickupPos, 0.05, 0.95)
				userSettingsChannel <- settings
			case '=':
				settings.DecayTimeS += 0.2
				settings.DecayTimeS = numeric.Clip(settings.DecayTimeS, 0.2, 10)
				userSettingsChannel <- settings
			case '-':
				settings.DecayTimeS -= 0.2
				settings.DecayTimeS = numeric.Clip(settings.DecayTimeS, 0.2, 10)
				userSettingsChannel <- settings
			}
		}
	}
}
