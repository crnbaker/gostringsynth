/*
Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.
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

func StartUI(waitGroup *sync.WaitGroup, pluckPlotChan chan []float64, userSettingsChan chan notes.UserSettings) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer waitGroup.Done()
	defer ui.Close()

	const guiWidth = 100

	pluckMarker := pluckMarker{}
	pickupPos := 0.0

	plot := MakePluckPlot(guiWidth)
	instructions := makeInstructionsBox(int(math.Floor(guiWidth / 2)))
	settingsBox := NewSettingsBox(int(math.Floor(guiWidth / 2)))

	for pluckPlotChan != nil && userSettingsChan != nil {
		select {
		case pluckPlot, ok := <-pluckPlotChan:
			if !ok {
				pluckPlotChan = nil
			} else if len(pluckPlot) > 0 {
				plot.Title = fmt.Sprintf("Pluck shape (amp: %.3f)", numeric.Max(pluckPlot))
				plot.Data = makePlotData(pluckPlot, guiWidth, pluckMarker, pickupPos)
			}
		case userSettings, ok := <-userSettingsChan:
			if !ok {
				userSettingsChan = nil
			} else {
				settingsBox.Update(userSettings)
				pluckMarker.pos = userSettings.PluckPos
				pluckMarker.width = userSettings.PluckWidth
				pickupPos = userSettings.PickupPos
			}
		}
		ui.Render(plot, instructions, settingsBox)
	}
}

func makeInstructionsBox(width int) *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Title = "gostringsynth"
	p.Text = `keyboard mapped across keys from a to k
	Quit:                       q`
	p.SetRect(0, 0, width, 10)
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan
	return p
}

type SettingsBox struct {
	*widgets.Paragraph
	settings notes.UserSettings
}

func NewSettingsBox(width int) SettingsBox {
	p := widgets.NewParagraph()
	p.Title = "parameters"
	p.SetRect(width, 0, 2*width, 10)
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan
	return SettingsBox{p, notes.UserSettings{}}
}

func (s SettingsBox) Update(u notes.UserSettings) {
	s.Text = fmt.Sprintf(`
	Octave : %d               down/up z/x
	Velocity: %d             down/up c/v
	Pluck pos: %.3f         left/right ,/.
	Pluck width: %.3f       down/up </>
	Decay time: %.3f        down/up -/=
	Pickup pos: %.3f        left/right [/]`,
		u.Octave, u.Velocity, u.PluckPos, u.PluckWidth, u.DecayTimeS, u.PickupPos)
}
