/*
Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.
*/

package plotting

import (
	"fmt"
	"log"
	"sync"

	"github.com/crnbaker/gostringsynth/numeric"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"gonum.org/v1/gonum/interp"
)

func StartUI(waitGroup *sync.WaitGroup, pluckPlotChan chan []float64) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer waitGroup.Done()
	defer ui.Close()

	const plotWidth = 100

	plot := widgets.NewPlot()
	plot.Marker = widgets.MarkerBraille
	plot.SetRect(0, 10, plotWidth, 20)
	plot.ShowAxes = false
	plot.AxesColor = ui.ColorWhite
	plot.LineColors[0] = ui.ColorCyan
	plot.PlotType = widgets.LineChart

	instructions := makeInstructionsBox()

	for pluckPlot := range pluckPlotChan {
		plot.Title = fmt.Sprintf("Pluck shape (amp: %.3f)", numeric.Max(pluckPlot))
		plot.Data = makePlotData(pluckPlot, plotWidth)

		ui.Render(plot, instructions)
	}

}

func makeInstructionsBox() *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Title = "gostringsynth"
	p.Text = `keyboard mapped across keys from A to K
	Octave down/up:        Z, X
	Velocity down/up:      C, V
	Quit:                  Q`
	p.SetRect(0, 0, 100, 10)
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan
	return p
}

func makePlotData(data []float64, figureWidth int) [][]float64 {
	dataWidth := len(data)
	plotData := make([][]float64, 1)
	plotData[0] = make([]float64, figureWidth)

	fracOfDataLength := make([]float64, dataWidth)
	for i := 0; i < dataWidth; i++ {
		fracOfDataLength[i] = float64(i) / float64(dataWidth-1)
	}
	var interpolator interp.PiecewiseLinear
	interpolator.Fit(fracOfDataLength, data)

	var fracOfPlotLength float64
	for i := 0; i < figureWidth; i++ {
		fracOfPlotLength = float64(i) / float64(figureWidth-1)
		plotData[0][i] = interpolator.Predict(fracOfPlotLength)
	}
	return plotData
}
