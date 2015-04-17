// 2015-04-03 Adam Bryt
//
// Modyfikacja pliku: code.google.com/p/plotinum/plotter/grid.go
// Zostało dodane rysowanie pionowych linii dla minor ticks i rysowanie
// głównych (wyróżnionych) linii: jednej poniomej i jednej pionowej.

package main

import (
	"image/color"

	"biorytm/Godeps/_workspace/src/code.google.com/p/plotinum/plot"
	"biorytm/Godeps/_workspace/src/code.google.com/p/plotinum/vg"
)

var (
	// DefaultGridLineStyle is the default style for grid lines.
	DefaultGridLineStyle = plot.LineStyle{
		Color: color.Gray{128},
		Width: vg.Points(0.25),
	}

	// DefaultGridMainLineStyle is the default style for main lines.
	DefaultGridMainLineStyle = plot.LineStyle{
		Color: color.Gray{128},
		Width: vg.Points(1.0),
	}
)

// Grid implements the plot.Plotter interface, drawing a set of grid lines.
type Grid struct {
	vStyle     plot.LineStyle
	hStyle     plot.LineStyle
	vMain      float64 // wartość na osi x wyróżnionej pionowej lini
	hMain      float64 // wartość na osi y wyróżnionej poziomej lini
	vMainStyle plot.LineStyle
	hMainStyle plot.LineStyle
}

// NewGrid returns a new grid.
func NewGrid(v, h float64) *Grid {
	return &Grid{
		vStyle:     DefaultGridLineStyle,
		hStyle:     DefaultGridLineStyle,
		vMainStyle: DefaultGridMainLineStyle,
		hMainStyle: DefaultGridMainLineStyle,
		vMain:      v,
		hMain:      h,
	}
}

// Plot implements the plot.Plotter interface.
func (g *Grid) Plot(da plot.DrawArea, plt *plot.Plot) {
	trX, trY := plt.Transforms(&da)

	if g.vStyle.Color == nil {
		goto horiz
	}
	for _, tk := range plt.X.Tick.Marker(plt.X.Min, plt.X.Max) {
		//if tk.IsMinor() {
		//	continue
		//}
		x := trX(tk.Value)
		da.StrokeLine2(g.vStyle, x, da.Min.Y, x, da.Min.Y+da.Size.Y)
	}

horiz:
	if g.hStyle.Color == nil {
		return
	}
	for _, tk := range plt.Y.Tick.Marker(plt.Y.Min, plt.Y.Max) {
		if tk.IsMinor() {
			continue
		}
		y := trY(tk.Value)
		da.StrokeLine2(g.hStyle, da.Min.X, y, da.Min.X+da.Size.X, y)
	}

	// główna pionowa linia
	x := trX(g.vMain)
	da.StrokeLine2(g.vMainStyle, x, da.Min.Y, x, da.Min.Y+da.Size.Y)

	// główna pozioma linia
	y := trY(g.hMain)
	da.StrokeLine2(g.hMainStyle, da.Min.X, y, da.Min.X+da.Size.X, y)
}
