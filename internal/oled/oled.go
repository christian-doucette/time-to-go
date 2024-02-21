package oled

import (
	"image"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/ssd1306"
	"periph.io/x/devices/v3/ssd1306/image1bit"
	"periph.io/x/host/v3"
)

// clears the OLED display of input
func ClearDisplay(bus string) {
	// just calls DisplayTextLines() with an empty list of lines to display
	DisplayTextLines([]string{}, 0, 0, bus)
}

// prints the lines to the OLED display
func DisplayTextLines(lines []string, startingDepth int, recurringDepth int, bus string) {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		panic(err.Error())
	}

	// Use i2creg I²C bus registry to open the given bus
	b, err := i2creg.Open(bus)
	if err != nil {
		panic(err.Error())
	}
	defer b.Close()

	// creates SSD1306 device from bus
	dev, err := ssd1306.NewI2C(b, &ssd1306.DefaultOpts)
	if err != nil {
		panic("failed to initialize ssd1306: " + err.Error())
	}

	// Set up image from device
	img := image1bit.NewVerticalLSB(dev.Bounds())

	f := basicfont.Face7x13
	drawer := font.Drawer{
		Dst:  img,
		Src:  &image.Uniform{image1bit.On},
		Face: f,
		Dot:  fixed.P(0, startingDepth),
	}

	// Draw each line to drawer
	for i, line := range lines {
		drawer.Dot = fixed.P(0, startingDepth+i*recurringDepth)
		drawer.DrawString(line)
	}

	if err := dev.Draw(dev.Bounds(), img, image.Point{}); err != nil {
		panic(err.Error())
	}
}
