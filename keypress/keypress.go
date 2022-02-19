package keypress

import (
	"fmt"

	"github.com/crnbaker/gostringsynth/errors"
	tty "github.com/mattn/go-tty"
)

func KeyDispatcher(noteSendChannel chan float64) {
	tty, err := tty.Open()
	errors.Chk(err)
	defer tty.Close()

	for {
		r, err := tty.ReadRune()
		errors.Chk(err)
		fmt.Println(r, " key pressed")
		switch r {
		case 'a':
			noteSendChannel <- 100
		case 's':
			noteSendChannel <- 200
		case 'd':
			noteSendChannel <- 300
		case 'q':
			close(noteSendChannel)
		}
	}
}
