package main

import (
	"io"
	"log"

	"github.com/wcharczuk/go-chart" //exposes "chart"
	"github.com/wcharczuk/go-chart/drawing"
)

var (
	colors = map[EventType]drawing.Color{
		EventDoWork:         chart.ColorBlue.WithAlpha(255),
		EventFinished:       chart.ColorOrange.WithAlpha(255),
		EventPQOpen:         chart.ColorRed.WithAlpha(255),
		EventAppendFreeConn: chart.ColorGreen.WithAlpha(255),
		EventErrorInsert:    chart.ColorBlack.WithAlpha(255),
	}
)

func createChart(out io.Writer, min, max, step uint64, counts map[EventType][]uint64) {
	xValues := make([]float64, 0, (max-min)/step)
	for v := min + step; v <= max; v += step {
		xValues = append(xValues, float64(v))
	}
	xValues = append(xValues, float64(max))

	series := make([]chart.Series, 0, len(counts))
	for e, _ := range counts {
		yValues := make([]float64, len(xValues)) // all 0
		for i, v := range counts[e] {
			yValues[i] = float64(v)
		}
		series = append(series, chart.ContinuousSeries{
			Name: e.String(),
			Style: chart.Style{
				Show:        true,
				StrokeColor: colors[e],
			},
			XValues: xValues,
			YValues: yValues,
		})
	}

	graph := chart.Chart{
		Width:  *width,
		Height: *height,
		XAxis: chart.XAxis{
			Style: chart.Style{Show: true},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{Show: true},
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 200,
			},
		},
		Series: series,
	}

	graph.Elements = []chart.Renderable{
		chart.LegendLeft(&graph),
	}

	if err := graph.Render(chart.SVG, out); err != nil {
		log.Printf("cannot render; err= %v\n", err)
		return
	}
}
