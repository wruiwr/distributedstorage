package drawer

import (
	"fmt"
	"github.com/ajstarks/svgo"
	"os"
	"strconv"
	"time"
)

// ref: SVGo: A Go library for SVG generation
// github.com/ajstarks/svgo
// See figure by Sguiggle tool or a browser

// Start Drawer
func StartDrawer() (*svg.SVG, int64) {

	var width = 1200
	var height = 400
	var startTime = time.Now().UnixNano()

	outputSVG := svg.New(os.Stdout)
	outputSVG.Start(width, height)
	outputSVG.Rect(10, 10, 1100, 100, "fill:#eeeeee;stroke:none")
	outputSVG.Text(20, 30, "Process 1 Timeline", "text-anchor:start;font-size:12px;fill:#333333")
	outputSVG.Rect(10, 130, 1100, 100, "fill:#eeeeee;stroke:none")
	outputSVG.Text(20, 150, "Process 2 Timeline", "text-anchor:start;font-size:12px;fill:#333333")
	for i := 0; i < 1201; i++ {
		timeText := strconv.FormatInt(int64(i), 10)
		if i%100 == 0 {
			outputSVG.Text(i, 380, timeText, "text-anchor:middle;font-size:10px;fill:#000000")
		} else if i%4 == 0 {
			outputSVG.Circle(i, 377, 1, "fill:#cccccc;stroke:none")
		}
		if i%10 == 0 {
			outputSVG.Rect(i, 0, 1, 400, "fill:#dddddd")
		}
		if i%50 == 0 {
			outputSVG.Rect(i, 0, 1, 400, "fill:#cccccc")
		}
	}

	outputSVG.Text(650, 360, "Run with goroutines", "text-anchor:start;font-size:12px;fill:#333333")
	return outputSVG, startTime
}

// Draw Point
func DrawPoint(osvg *svg.SVG, pnt int, process int, rw string, startTime int64) {
	sec := time.Now().UnixNano()
	diff := (int64(sec) - int64(startTime)) / 100000

	pointLocation := 0
	if int(diff) <= 100 {
		pointLocation = int(diff) / 10
	} else {
		pointLocation = int(diff) / 100
	}

	fmt.Println(pointLocation)

	pointLocationV := 0
	color := "#000000"

	textTime := strconv.Itoa(pointLocation)
	switch {
	case process == 1:
		pointLocationV = 60
		color = "#cc6666"
		osvg.Text(pointLocation, 50, textTime, "text-anchor:start;font-size:10px;fill:#333333")
		osvg.Text(pointLocation, 40, rw, "text-anchor:start;font-size:10px;fill:#333333")
	default:
		pointLocationV = 180
		color = "#66cc66"
		osvg.Text(pointLocation, 170, textTime, "text-anchor:start;font-size:10px;fill:#333333")
		osvg.Text(pointLocation, 160, rw, "text-anchor:start;font-size:10px;fill:#333333")
	}

	osvg.Rect(pointLocation, pointLocationV, 3, 5, "fill:"+color+";stroke:none;")

	time.Sleep(200 * time.Millisecond)
}

func EndDrawer(outputSVG *svg.SVG) {
	outputSVG.End()
}
