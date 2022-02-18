package keypress

import (
	"fmt"

	"github.com/crnbaker/gostringsynth/errors"
	tty "github.com/mattn/go-tty"
)

func KeyDispatcher(noteChannel chan float64) {
	tty, err := tty.Open()
	errors.Chk(err)
	defer tty.Close()

	for {
		r, err := tty.ReadRune()
		errors.Chk(err)
		fmt.Println(r, " key pressed")
		switch r {
		case 'a':
			noteChannel <- 100
		case 's':
			noteChannel <- 200
		case 'd':
			noteChannel <- 300
		case 'q':
			close(noteChannel)
		}
	}
}
