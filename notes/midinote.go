package notes

import "math"

type Note interface {
	Pitch() int
}
type MidiNoteSettings struct {
	Velocity int
	Octave   int
}

func DefaultMidiNoteSettings() MidiNoteSettings {
	return MidiNoteSettings{64, 3}
}

// MidiNote holds MIDI pitch and velocity. MidiNotes are generated from key presses by the PublishNotes func.
type MidiNote struct {
	MidiNoteSettings
	rawPitch int
}
type StringSettings struct {
	PluckPos     float64
	PluckWidth   float64
	WaveSpeedMpS float64
	DecayTimeS   float64
	PickupPos    float64
}

func DefaultStringSettings() StringSettings {
	return StringSettings{0.2, 0.05, 200, 3, 0.1}
}

type StringMidiNote struct {
	MidiNote
	StringSettings
}

func NewStringMidiNote(pitch int, midiSettings MidiNoteSettings, stringSettings StringSettings) StringMidiNote {
	note := MidiNote{midiSettings, pitch}
	return StringMidiNote{note, stringSettings}
}

func (note *MidiNote) Pitch() int {
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
