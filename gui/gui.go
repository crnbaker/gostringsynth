/* The GUI package provides a termui terminal GUI for gostring synth.

It only displayes information received on its input channels. User input is handled by the notes package.
*/

package gui

import (
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/crnbaker/gostringsynth/notes"
	"github.com/crnbaker/gostringsynth/numeric"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// StartUILoop creates a GUI, and then listens for incoming pluck and settings data and displays them
func StartUILoop(waitGroup *sync.WaitGroup, pluckPlotChan chan []float64, userSettingsChan chan notes.UserSettings) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer waitGroup.Done()
	defer ui.Close()

	const guiWidth = 70
	const topBoxesRatio = 0.55

	plot := makePluckPlot(guiWidth)
	horLineBetweenBoxes := int(math.Floor(guiWidth * topBoxesRatio))
	instructions := makeInstructionsBox(0, horLineBetweenBoxes)
	settingsBox := newSettingsBox(horLineBetweenBoxes, guiWidth)

	for pluckPlotChan != nil && userSettingsChan != nil {
		select {
		case pluckPlot, ok := <-pluckPlotChan:
			if !ok {
				pluckPlotChan = nil
			} else if len(pluckPlot) > 0 {
				plot.Title = fmt.Sprintf("pluck shape (ampl.: %.3f)", numeric.Max(pluckPlot))
				plot.Data = makePlotData(pluckPlot, guiWidth)
			}
		case userSettings, ok := <-userSettingsChan:
			if !ok {
				userSettingsChan = nil
			} else {
				settingsBox.update(userSettings)
			}
		}
		ui.Render(plot, instructions, settingsBox)
	}
}

func makeInstructionsBox(hStart int, hStop int) *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Title = "gostringsynth"
	p.Text = `A finite-difference time-domain string synthesiser.
	
	Christian Baker 2022.
	
	Musical keyboard mapped across keys "a" to "k"
	Press q to quit.`
	p.SetRect(hStart, 0, hStop, 10)
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan
	return p
}

type settingsBox struct {
	*widgets.Paragraph
	settings notes.UserSettings
}

func newSettingsBox(hStart int, hStop int) settingsBox {
	p := widgets.NewParagraph()
	p.Title = "parameters"
	p.SetRect(hStart, 0, hStop, 10)
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan
	return settingsBox{p, notes.UserSettings{}}
}

func (s settingsBox) update(u notes.UserSettings) {
	s.Text = fmt.Sprintf(`Param.       Control   Value

	Octave       z x       %d
	Velocity     c v       %d
	Pluck pos    , .       %.3f 
	Pluck width  < >       %.3f
	Decay (s)    - =       %.3f
	Pickup pos   [ ]       %.3f`,
		u.Octave, u.Velocity, u.PluckPos, u.PluckWidth, u.DecayTimeS, u.PickupPos)
}
