package notes

import "math"

type midiNoteSettings struct {
	Velocity int
	Octave   int
}

func defaultMidiNoteSettings() midiNoteSettings {
	return midiNoteSettings{64, 5}
}

// midiNote holds MIDI pitch and velocity. MidiNotes are generated from key presses by the PublishNotes func.
type midiNote struct {
	midiNoteSettings
	rawPitch int
}
type stringSettings struct {
	PluckPos   float64
	PluckWidth float64
	DecayTimeS float64
	PickupPos  float64
}

func defaultStringSettings() stringSettings {
	return stringSettings{0.3, 0.0, 6, 0.15}
}

// StringMidiNote stores midi note information and string physical properties for sending to synth module
type StringMidiNote struct {
	midiNote
	stringSettings
}

// NewStringMidiNote creates a new StringMidiNote given pitch, midi and string settings
func NewStringMidiNote(pitch int, midiSettings midiNoteSettings, settings stringSettings) StringMidiNote {
	note := midiNote{midiSettings, pitch}
	return StringMidiNote{note, settings}
}

// Pitch returns the pitch of a MIDI note taking into account its octave
func (note *midiNote) Pitch() int {
	return note.rawPitch + (note.Octave+2)*12
}

// MidiPitchToFreq converts a MIDI pitch to a fundemental frequency
func MidiPitchToFreq(pitch int) float64 {
	return math.Pow(2, ((float64(pitch)-69)/12)) * 440
}

// MidiVelocityToAmplitude converts a MIDI velocity to a normalised amplitude
func MidiVelocityToAmplitude(velocity int) float64 {
	return float64(velocity) / float64(127)
}
