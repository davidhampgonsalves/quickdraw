package main

import (
  "fmt"
  "encoding/json"
  "bufio"
  "os"
  "log"
  "math/rand"
  "math"
  "time"
  "flag"

  "github.com/tdewolff/canvas"
)

type Drawing struct {
  Word string
  CountryCode string
  Recognized bool
  Drawing [][][]float64
}

func main() {
  var w float64
  var h float64
  var m float64
  var countPerRow int

  flag.Float64Var(&w, "w", 1000, "width of output")
  flag.Float64Var(&h, "h", 1000, "height of output")
  flag.Float64Var(&m, "m", 25, "margin around drawings")
  flag.IntVar(&countPerRow, "c", 20, "drawing count per row")
  help := flag.Bool("help", false, "print this help message")
  flag.Parse()

  if(*help) {
    flag.PrintDefaults()
    return
  }

  fileName := flag.Arg(0)

  if(len(fileName) == 0) {
    fmt.Println("error: file must be provided.")
    return
  }

  rand.Seed(time.Now().UTC().UnixNano())

  boxSize := (w - (m * float64((countPerRow + 1)))) / float64(countPerRow)
  mBoxSize := boxSize + m
  numberOfCols := int(math.Floor((h - m) / mBoxSize))

  boxScale := boxSize / 255

  file, err := os.Open(fileName)
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  c := canvas.New(w, h)
  c.SetStrokeColor(canvas.Black)
  c.SetFillColor(canvas.Transparent)

  scanner := bufio.NewScanner(file)
  randoms := []float64{}
  drawings := []Drawing{}

  fmt.Printf("picking %v random drawings.", numberOfCols * countPerRow)
  line := 0
  for scanner.Scan() {
    line += 1
    if(line % 1000 == 0) { fmt.Print(".") }
    drawingJSON := scanner.Text()
    drawing := Drawing{}
    json.Unmarshal([]byte(drawingJSON), &drawing)

    if(!drawing.Recognized) { continue }

    random := rand.Float64()
    if(len(randoms) >= (countPerRow * numberOfCols)) {
      min := 2.0
      i := -1

      for j, val := range randoms {
        if(val > min) { continue }

        i = j
        min = val
      }

      if(random > min) {
        randoms[i] = random
        drawings[i] = drawing
      }
    } else {
      randoms = append(randoms, random)
      drawings = append(drawings, drawing)
    }
  }
  fmt.Println("\ngenerating svg.")

  offsetX := 0
  offsetY := 0
  for _, drawing := range drawings {
    for _, path := range drawing.Drawing {
      p := &canvas.Path{}
      for i, x := range path[0] {
        y := boxSize - (path[1][i] * boxScale)
        x = boxScale - (x * boxScale) - m

        if(i == 0) {
          p.MoveTo(x, y)
        } else {
          p.LineTo(x, y)
        }
      }

      c.DrawPath((float64(offsetX) * mBoxSize) + m, (float64(offsetY) * mBoxSize) + m, p)
    }

    offsetX += 1
    if(offsetX >= countPerRow) {
      offsetY += 1
      offsetX = 0
    }

    if(float64(offsetY) * mBoxSize >= h) { break }
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }

  c.Fit(m)
  c.SaveSVG("./out.svg")

  fmt.Println("done: out.svg")
}
