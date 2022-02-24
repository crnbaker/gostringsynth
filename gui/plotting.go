package gui

import (
	"math"

	"github.com/crnbaker/gostringsynth/numeric"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"gonum.org/v1/gonum/interp"
)

type pluckMarker struct {
	pos   float64
	width float64
}

func MakePluckPlot(width int) *widgets.Plot {
	plot := widgets.NewPlot()
	plot.Marker = widgets.MarkerBraille
	plot.SetRect(0, 10, width, 20)
	plot.ShowAxes = false
	plot.AxesColor = ui.ColorWhite
	plot.LineColors[0] = ui.ColorCyan
	plot.LineColors[1] = ui.ColorRed
	plot.LineColors[2] = ui.ColorGreen
	plot.PlotType = widgets.LineChart
	return plot
}

func makePlotData(data []float64, figureWidth int, pluckMarker pluckMarker, pickupMarker float64) [][]float64 {
	dataWidth := len(data)
	plotData := make([][]float64, 3)
	plotData[0] = make([]float64, figureWidth)
	plotData[1] = make([]float64, figureWidth)
	pickupMarkerIndex := int(math.Round(pickupMarker / float64(figureWidth)))
	plotData[1][pickupMarkerIndex] = numeric.Max(data)
	pluckMarkerIndex := int(math.Round(pluckMarker.pos / float64(figureWidth)))
	// pluckMarkerHalfWidth := 0.5 * pluckMarker.width / float64(figureWidth)
	// pluckMarkerStart := int(math.Floor(float64(pluckMarkerIndex) - pluckMarkerHalfWidth))
	// pluckMarkerStop := int(math.Floor(float64(pluckMarkerIndex) + pluckMarkerHalfWidth))

	plotData[2] = make([]float64, figureWidth)
	plotData[2][pluckMarkerIndex] = numeric.Max(data)
	// for i := pluckMarkerStart; i <= pluckMarkerStop; i++ {
	// 	plotData[2][i] = numeric.Max(data)
	// }

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
