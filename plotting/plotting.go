/*
Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.
*/

package plotting

import (
	"log"
	"sync"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"gonum.org/v1/gonum/interp"
)

func TestPlot(waitGroup *sync.WaitGroup, pluckPlotChan chan []float64) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer waitGroup.Done()
	defer ui.Close()

	const plotWidth = 100

	plot := widgets.NewPlot()
	plot.Title = "Pluck shape"
	plot.Marker = widgets.MarkerDot
	plot.SetRect(0, 0, plotWidth, 20)
	plot.AxesColor = ui.ColorWhite
	plot.LineColors[0] = ui.ColorYellow
	plot.PlotType = widgets.ScatterPlot

	for pluckPlot := range pluckPlotChan {
		plot.Data = makePlotData(pluckPlot, plotWidth)
		ui.Render(plot)
	}

}

func makePlotData(data []float64, plotWidth int) [][]float64 {
	dataWidth := len(data)
	plotData := make([][]float64, 1)
	plotData[0] = make([]float64, plotWidth)

	dataHorAxs := make([]float64, dataWidth)
	for i := 0; i < dataWidth; i++ {
		dataHorAxs[i] = float64(i)
	}
	var interpolator interp.PiecewiseLinear
	interpolator.Fit(dataHorAxs, data)

	for i := 0; i < plotWidth; i++ {
		plotData[0][i] = interpolator.Predict(float64(dataWidth) * float64(i) / float64(plotWidth))
	}
	return plotData
}
