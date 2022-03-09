package notes

type stringSettings struct {
	pluckPos     float64
	pluckWidth   float64
	decayTimeS   float64
	pickupPos    float64
	wavespeedMpS float64
}

func defaultStringSettings() stringSettings {
	return stringSettings{0.3, 0.0, 6, 0.15, 200}
}

// stringMidiNote stores midi note information and string physical properties for sending to synth module
type stringMidiNote struct {
	midiNote
	stringSettings
}

// newStringMidiNote creates a new StringMidiNote given pitch, midi and string settings
func newStringMidiNote(pitch int, midiSettings midiNoteSettings, settings stringSettings) *stringMidiNote {
	note := midiNote{midiSettings, pitch}
	return &stringMidiNote{note, settings}
}

// LengthM returnss the length of the string required to play the note given its pitch
func (note *stringMidiNote) LengthM() float64 {
	return FreqToStringLength(MidiPitchToFreq(note.Pitch()), note.wavespeedMpS)
}

// PluckPos returns the position of the pluck relative to the length of the string
func (str *stringSettings) PluckPos() float64 {
	return str.pluckPos
}

// PluckWidth returns the width of the pluck relative to the length of the string
func (str *stringSettings) PluckWidth() float64 {
	return str.pluckWidth
}

// PickupPos returns the position of the pickup relative to the length of the string
func (str *stringSettings) PickupPos() float64 {
	return str.pickupPos
}

// DecayTimeS return the decay time of the string in seconds
func (str *stringSettings) DecayTimeS() float64 {
	return str.decayTimeS
}

// WavespeedMpS return speed of the wave on the string in meters per second
func (note *stringSettings) WavespeedMpS() float64 {
	return note.wavespeedMpS
}

// FreqToStringLength converts a fundemental frequency to a string length, given a string wave speed in m/s
func FreqToStringLength(freqHz float64, waveSpeedMpS float64) float64 {
	return waveSpeedMpS / freqHz
}
