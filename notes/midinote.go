package notes

import "math"

// MidiNote holds MIDI pitch and velocity. MidiNotes are generated from key presses by the PublishNotes func.
type MidiNote struct {
	Pitch    int
	Velocity int
}

func newNote(pitch int, velocity int) MidiNote {
	return MidiNote{pitch, velocity}
}

func changeNoteOctave(note MidiNote, octave int) MidiNote {
	note.Pitch += (octave + 2) * 12
	return note
}

// MidiPitchToFreq converts a MIDI pitch to a fundemental frequency
func MidiPitchToFreq(pitch int) float64 {
	return math.Pow(2, ((float64(pitch)-69)/12)) * 440
}

// MidiVelocityToAmplitude converts a MIDI velocity to a normalised amplitude
func MidiVelocityToAmplitude(velocity int) float64 {
	return float64(velocity) / float64(127)
}
