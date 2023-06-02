package results

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"go-airline-crew-rostering/input"
	"go-airline-crew-rostering/metrics"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/text"
	"gonum.org/v1/plot/vg"
)

type ticker struct {
	Ticker plot.Ticker
	Format string
	// Time   func(t float64) time.Time
}

func tickers(min float64, max float64) (int, int) {
	minorTicker := 1
	majorTicker := 2
	if max-min > 25 {
		minorTicker = 2
		majorTicker = 10
	}
	if max-min > 50 {
		minorTicker = 5
		majorTicker = 10
	}
	if max-min > 100 {
		minorTicker, majorTicker = tickers(min/10, max/10)
		minorTicker = minorTicker * 10
		majorTicker = majorTicker * 10
	}
	return minorTicker, majorTicker
}

func (t ticker) Ticks(min float64, max float64) []plot.Tick {
	minorTicker, majorTicker := tickers(min, max)
	ticks := []plot.Tick{}
	i := math.Floor(float64(min)/float64(minorTicker)) * float64(minorTicker)
	for {
		t := plot.Tick{Value: i}
		if int(i)%majorTicker == 0 {
			t.Label = fmt.Sprintf("%.0f\n", i)
		}
		ticks = append(ticks, t)
		i = i + float64(minorTicker)
		if i > max {
			break
		}
	}
	return ticks
}

func Min(min float64, max float64) float64 {
	_, majorTicker := tickers(min, max)
	return math.Floor(float64(min)/float64(majorTicker)) * float64(majorTicker)
}

func Max(min float64, max float64) float64 {
	_, majorTicker := tickers(min, max)
	return math.Ceil(float64(max)/float64(majorTicker)) * float64(majorTicker)
}

func drawIterationMetricsPlot(plotFile string, m *metrics.Metrics, args *input.ArgumentCollection) string {
	// draw a plot depicting the progression of best, worst and average cost with each iteration
	points := func(data []float64) plotter.XYs {
		points := make(plotter.XYs, len(data))
		for i := range points {
			points[i].X = float64(i)
			points[i].Y = data[i]
		}
		return points
	}

	p := plot.New()
	// p.BackgroundColor = color.RGBA{R: 204, G: 201, B: 239, A: 255}

	if args.Algorithm == "multiCSO" {
		p.Title.Text = fmt.Sprintf("Cost Comparison per Iteration\n(Multi-step CSO, FL=%02.1f)", *args.FL)
	} else if args.Algorithm == "AOA" {
		p.Title.Text = fmt.Sprintf("Cost Comparison per Iteration\n(AOA, C1=%02.1f, C2=%02.1f, C3=%02.1f, C4=%02.1f)", args.Constants[0], args.Constants[1], args.Constants[2], args.Constants[3])
	}

	p.Title.TextStyle.XAlign = text.XCenter
	p.Title.Padding = vg.Points(15)
	p.Title.TextStyle.Font.Typeface = "Arial"
	p.Title.TextStyle.Font.Size = 16
	p.Title.TextStyle.Font.Weight = 3

	p.X.Label.Text = "Generation"
	p.X.Label.TextStyle.Font.Typeface = "Arial"
	p.X.Label.TextStyle.Font.Size = 14
	p.X.Label.TextStyle.Font.Weight = 3
	p.X.Label.Padding = vg.Points(10)

	p.Y.Label.Text = "Solution Cost"
	p.Y.Label.TextStyle.Font.Typeface = "Arial"
	p.Y.Label.TextStyle.Font.Size = 14
	p.Y.Label.TextStyle.Font.Weight = 3
	p.Y.Label.Padding = vg.Points(10)

	p.Legend.TextStyle.Font.Typeface = "Arial"
	p.Legend.Top = true
	p.Legend.Left = true
	p.Legend.Padding = vg.Millimeter

	plottedData, err := plotter.NewLine(points(m.IterBestCost))
	if err != nil {
		log.Panic(err)
	}
	plottedData.Color = color.RGBA{R: 31, G: 179, B: 58, A: 255}
	p.Add(plottedData)
	p.Legend.Add("Iter Best Cost", plottedData)

	plottedData, err = plotter.NewLine(points(m.IterWorstCost))
	if err != nil {
		log.Panic(err)
	}
	plottedData.Color = color.RGBA{R: 153, G: 29, B: 77, A: 255}
	p.Add(plottedData)
	p.Legend.Add("Iter Worst Cost", plottedData)

	plottedData, err = plotter.NewLine(points(m.IterAverageCost))
	if err != nil {
		log.Panic(err)
	}
	plottedData.Color = color.RGBA{R: 106, G: 182, B: 224, A: 255}
	p.Add(plottedData)
	p.Legend.Add("Iter Average Cost", plottedData)

	p.X.Tick.Marker = ticker{}

	p.Y.Tick.Marker = ticker{}

	p.X.Max *= 1.01
	p.Y.Max *= 1.01

	err = p.Save(15*vg.Centimeter, 15*vg.Centimeter, plotFile)
	if err != nil {
		log.Panic(err)
	}

	return plotFile
}
