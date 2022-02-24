# gostringsynth
## A finite-difference string synthesizer written in Go.

*gostringsynth* is a real-time, polyphonic physical-modelling synthesizer that simulates the vibration of strings using the finite-difference time-domain method. It's a standalone terminal program with a [termui](https://github.com/gizak/termui) GUI controlled and played using a PC keyboard.

<img width="603" alt="image" src="https://user-images.githubusercontent.com/31904251/155533740-ba7f8dc4-5953-4eb7-91cd-3ba251447b20.png">

### Physics

The vibration of a string is modelled by discretising the 1D, second order wave equation using the finite-difference time-domain (FTDT) method. An "explicit" FTDT scheme is used, taken from [Numerical Sound Synthesis: Finite Difference Schemes and Simulation in Musical Acoustics by Stefan Bilbao](https://www.wiley.com/en-gb/Numerical+Sound+Synthesis%3A+Finite+Difference+Schemes+and+Simulation+in+Musical+Acoustics-p-9780470510469). Internally, *gostringsynth* holds a discretised representation of a string as an array of displacements (a Go slice). This string is clamped at both ends by permenantly setting its displacement to zero at those points. The string can be "plucked" at different positions along the string and with different "finger" widths, the position of the pickup can also be changed, and damping is controlled by setting a 60 dB decay time in seconds.

## Installation

The audio stream is handled using the [Go interface](https://github.com/gordonklaus/portaudio) for [portaudio](http://www.portaudio.com), which needs to be installed separately.

With homebrew:
```
brew install portaudio
```

Once portaudio is installed, to build and run *gostringsynth* from source you'll need to install Go:
```
brew install go
```
Then you can clone the repository, navigate to it and build and run *gostringsynth* in one command with
```
go run main.go
```

Go dependencies:
* [termui](https://github.com/gizak/termui) for the terminal GUI
* [portaudio](https://github.com/gordonklaus/portaudio) for audio
* [go-tty](https://github.com/mattn/go-tty) for keypress handling
* [gonum](https://gonum.org/v1/gonum) for interpolation


## Usage


https://user-images.githubusercontent.com/31904251/155550348-af1b0d86-df55-46bf-b5ad-701501388722.mov


You can play the synth using a QWERTY keyboard with mappings the same as for [Ableton Live](https://www.ableton.com/en/manual/live-keyboard-shortcuts/#36-13-key-midi-map-mode-and-the-computer-midi-keyboard):
* A C to C chormatic scale is mapped across keys A to K,
* The octave is lower and raised using Z and X,
* The Velocity is decreased and increased using C and V.

The eventual plan is for the synth to be playable with a MIDI keyboard.

Each time a key is pressed, a new string is created and plucked. The pitch is controlled by setting the length of the newly created string.

The physical properties of the string model can be changed, taking effect from the next plucked string:
* Pluck position (relative to string length) can be moved left and right with "," and "."
* Pluck width (relative to string length) can be decreased and increased with "<" and ">" (requires holding shift)
* The pickup position (relative to string length) can be moved left and right with "[" and "]"
* The decay time in seconds can be decreased and increased with "-" and "="

The pluck shape is created by the synthesis module and therefore will only be displayed in the GUI after a note has been played, and only relates to the last plucked string.

Although *gostringsynth* has been designed for infinite polyphony, the maximum number if simultaneous voices (i.e. strings) is currently set to 6 to reduce the likelihood of overloading your CPU.

*gostringsynth* models discretised strings made up of the maximum number of sections possible while keeping the simulation stable (see [Bilbao's book](https://www.wiley.com/en-gb/Numerical+Sound+Synthesis%3A+Finite+Difference+Schemes+and+Simulation+in+Musical+Acoustics-p-9780470510469) for more info). Lower-pitched (i.e. longer) strings are discretised into more sections, and are therefore more computationally intensive to simulate. If your machine struggles, you may find it's OK at higher pitches.

## Why?

Just an excuse to learn a new language. My MSc thesis (2009, supervised by Stefan Bilbao) was about finite difference simulations for musical applications and I have thought since then that it would be cool to make a real-time finite difference synth. Go seemed like a good option, coming from Python and wanting to learn a compiled language.
