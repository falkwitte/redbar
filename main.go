package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
	"image"
	"image/color"
	_ "image/png"
	"math/rand"
	"os"
)

func Loadimage(path string) (pixel.Picture, error) {
	// open picture
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	// close file
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	// decode picture
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func run() {
	// window config
	cfg := pixelgl.WindowConfig{
		Title: "redbar",

		// create rectangle
		Bounds: pixel.R(0, 0, 1600, 900),

		// window updates with refreshrate of monitor
		VSync: true,

		AlwaysOnTop: true,
		Monitor:     pixelgl.PrimaryMonitor(),
	}

	// create window
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	//creating an atlas
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	//creating text
	redbar := text.New(pixel.ZV, atlas)
	countertext := text.New(pixel.V(200, 200), atlas)

	//text color
	redbar.Color = color.RGBA{R: 0xff, A: 0xff}
	countertext.Color = color.Black

	// printing the name of the game
	fmt.Fprintln(redbar, "redbar")

	// load redbarrel.png and scope.png
	redbarrel, err := Loadimage("red_barrel.png")
	if err != nil {
		panic(err)
	}
	scope, err := Loadimage("scope.png")
	if err != nil {
		panic(err)
	}

	// sprites
	redbarrelsprite := pixel.NewSprite(redbarrel, redbarrel.Bounds())
	scopesprite := pixel.NewSprite(scope, scope.Bounds())

	// count how many times scope has been moved
	counter := 0

	// scale sprites
	scopemat := pixel.IM
	scopemat = scopemat.Scaled(pixel.ZV, 5)

	// initial position scope
	scopemat = scopemat.Moved(win.Bounds().Center())

	// create rectangle with the same size as scopemat
	scoperect := scope.Bounds()
	// good enough
	scoperect = scoperect.Resized(scoperect.Center(), pixel.V(100, 100))

	//initial position scoperect
	scoperect = scoperect.Moved(win.Bounds().Center())

	//variable vectorval for movevec
	vectorvala := 1400
	vectorvalb := 700

	//variable counterpos for countertextmat
	// to always draw the counter at the top of the window
	var counterpos float64 = 200

	// main loop
	for !win.Closed() {
		// set a background color and clear screen
		win.Clear(color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xff})

		// scale, move and draw redbarrelsprite
		redbarrelmat := pixel.IM
		redbarrelmat = redbarrelmat.Scaled(pixel.ZV, 4)
		redbarrelmat = redbarrelmat.Moved(win.MousePosition())

		redbarrelsprite.Draw(win, redbarrelmat)

		// after redbarrel to be drawn on top
		scopesprite.Draw(win, scopemat)

		//draw redbar text ontop of everything
		redbartextmat := pixel.IM
		redbartextmat = redbartextmat.Scaled(redbar.Orig, 10)
		redbartextmat = redbartextmat.Moved(win.Bounds().Center().Sub(pixel.V(200, 200)))

		//draw the counter at the top of the window
		countertextmat := pixel.IM
		countertextmat = countertextmat.Scaled(countertext.Orig, 2)
		countertextmat = countertextmat.Moved(win.Bounds().Center().Add(pixel.V(-300, counterpos)))

		countertext.Draw(win, countertextmat)

		// random position with min and max
		movevec := pixel.V(float64(rand.Intn(vectorvala-100)+100), float64(rand.Intn(vectorvalb-100)+100))

		if counter == 0 {
			redbar.Draw(win, redbartextmat)
		}

		// if scoperect contains the mouse position apply a random delta vector to move
		if scoperect.Contains(win.MousePosition()) == true {
			counter++

			//clear countertext
			countertext.Clear()

			// printing the counter
			fmt.Fprintln(countertext, "counter: ", counter)

			// logging
			//fmt.Println("Count: ", counter)

			// needs to move scope and rect back to the beginning of the window to then apply further transformations
			// otherwise the transformations will just add up and if I move it back to the initpos scope will not be able to
			// spawn on the whole screen
			// this is the best solution that I have found for this problem
			if counter >= 1 {
				// move scope and rect back to the beginning of the window
				// move by delta vector: beggining of window(zero vector) - currentpos
				scopemat = scopemat.Moved(pixel.ZV.Sub(scopemat.Project(pixel.ZV)))
				scoperect = scoperect.Moved(pixel.ZV.Sub(scoperect.Center()))
			}
			// apply random delta vector
			scopemat = scopemat.Moved(movevec)
			scoperect = scoperect.Moved(movevec)
		}

		// disable the cursor if on window
		if win.MouseInsideWindow() == true {
			win.SetCursorVisible(false)
		}

		// enable/disable fullscreen
		if win.Monitor() == nil {
			if win.JustPressed(pixelgl.KeyF) {
				win.SetMonitor(pixelgl.PrimaryMonitor())
			}
			//change max rand value
			vectorvala = 1400
			vectorvalb = 700

			//change pos of countertext
			counterpos = 200

		} else {
			if win.JustPressed(pixelgl.KeyF) {
				win.SetMonitor(nil)
			}
			//change max rand value
			vectorvala = 1720
			vectorvalb = 880

			//change pos of countertext
			counterpos = 300
		}

		// close the window on pressing q
		if win.JustPressed(pixelgl.KeyQ) {
			win.SetClosed(true)
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
