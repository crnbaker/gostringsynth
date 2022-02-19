package keypress

import "math"

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

func MidiPitchToFreq(pitch int) float64 {
	return math.Pow(2, ((float64(pitch)-69)/12)) * 440
}

func MidiVelocityToAmplitude(velocity int) float64 {
	return float64(velocity) / float64(127)
}
