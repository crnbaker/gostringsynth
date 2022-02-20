# gostringsynth
## A finite-difference string synthesizer written in Go.

*gostringsynth* is a real-time, polyphonic physical-modelling synthesizer that simulates the vibration of strings using the finite-difference time-domain method. It's a standalone, GUI-less program controlled and played using a PC keyboard.

*gostringsynth* is very much a work in progress and has not yet been released, but is usable. This README describes its current state.

### Physics

The vibration of a string is modelled by discretising the 1D, second order wave equation using the finite-difference time-domain (FTDT) method. The explicit FTDT scheme used is described in [Numerical Sound Synthesis: Finite Difference Schemes and Simulation in Musical Acoustics by Stefan Bilbao](https://www.wiley.com/en-gb/Numerical+Sound+Synthesis%3A+Finite+Difference+Schemes+and+Simulation+in+Musical+Acoustics-p-9780470510469). Internally, *gostringsynth* holds a discretised representation of a string as an array of displacements (a Go slice). This string is clamped at both ends by permenantly setting its displacement to zero at those points. The state of the string at the previous two time steps is held, as required by the second-order scheme.

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

## Usage

You can play the synth using a QWERTY keyboard with mappings the same as for [Ableton Live](https://www.ableton.com/en/manual/live-keyboard-shortcuts/#36-13-key-midi-map-mode-and-the-computer-midi-keyboard):
* A C to C octave is mapped across keys A to K,
* The octave is lower and raised using Z and X,
* The Velocity is decreased and increased using C and V.

Press Q to quit.

The eventual plan is for the synth to be playable with a MIDI keyboard.

Each time a key is pressed, a new string is created and plucked. The pitch is controlled by setting the length of the newly created string. Currently the string wavespeed, pickup position, pluck position, pluck shape and decay time are set constant, but will be controllable in a future version. I also plan to enable sustained (i.e. bowed) notes in a future version.

Although *gostringsynth* has been designed for infinite polyphony, the maximum number if simultaneous voices (i.e. strings) is currently set to 8 to reduce the likelihood of overloadingyour CPU. This will also be controllable in a future version.

*gostringsynth* models discretised strings made up of the maximum number of sections possible while keeping the simulation stable (see [Bilbao's book](https://www.wiley.com/en-gb/Numerical+Sound+Synthesis%3A+Finite+Difference+Schemes+and+Simulation+in+Musical+Acoustics-p-9780470510469) for more info). Lower-pitched (i.e. longer) strings are discretised into more sections, and are therefore more computationally intensive to simulate. If your machine struggles, you may find its OK at higher pitches.

## Why?

Just an excuse to learn a new language. My MSc thesis (2009, supervised by Stefan Bilbao) was about finite difference simulations for musical applications and have thought since then that it would be cool to make a real-time finite difference synth. Go seemed like a good option, coming from Python and wanting to learn a compiled language.
