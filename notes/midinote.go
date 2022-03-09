package notes

import "math"

type midiNoteSettings struct {
	velocity int
	octave   int
}

func defaultMidiNoteSettings() midiNoteSettings {
	return midiNoteSettings{64, 5}
}

// midiNote holds MIDI pitch and velocity. MidiNotes are generated from key presses by the PublishNotes func.
type midiNote struct {
	midiNoteSettings
	rawPitch int
}

// Velocity returns the velocity value of a midi note
func (note *midiNoteSettings) Velocity() int {
	return note.velocity
}

// Octave returns the midi octave value of a midi note
func (note *midiNoteSettings) Octave() int {
	return note.octave
}

// Pitch returns the pitch of a MIDI note taking into account its octave
func (note *midiNote) Pitch() int {
	return note.rawPitch + (note.octave+2)*12
}

// Amplitude scales a note's velocity between 0 and 1 for use as a waveform amplitude
func (note *midiNote) Amplitude() float64 {
	return float64(note.velocity) / float64(127)
}

// MidiPitchToFreq converts a MIDI pitch to a fundemental frequency
func MidiPitchToFreq(pitch int) float64 {
	return math.Pow(2, ((float64(pitch)-69)/12)) * 440
}
