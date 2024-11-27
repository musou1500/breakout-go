package main

/*
#cgo LDFLAGS: -lncursesw
#include <ncurses.h>
#include <stdlib.h>
int go_mvprintw(int y, int x, char *name) {
	return mvprintw(y, x, name);
}
*/
import "C"
import (
	"math"
	"time"
	"unsafe"
)

var (
	hasBall = true
	px      = 40
	py      = 23
	playing = true
	bx      = 0.0
	by      = 0.0
	vx      = 0.0
	vy      = 0.0
)

func draw() {
	C.clear()
	ball := C.CString("*")
	paddle := C.CString("=====")
	defer C.free(unsafe.Pointer(ball))
	defer C.free(unsafe.Pointer(paddle))
	if hasBall {
		C.go_mvprintw(C.int(py-1), C.int(px), ball)
	}

	C.go_mvprintw(C.int(py), C.int(px-2), paddle)
	if !hasBall {
		C.go_mvprintw(C.int(by), C.int(bx), ball)
	}

	C.refresh()
}

func checkPaddleCollision() {
	if by < 23.0 || bx < float64(px-2) || bx > float64(px+3) {
		return
	}

	by = 23
	theta := math.Pi * ((float64(px)-bx+1.5)/8.0 + 0.25)
	vx = 0.5 * math.Cos(theta)
	vy = 0.5 * -math.Sin(theta)
}

func moveBall() {
	if hasBall {
		return
	}

	checkPaddleCollision()
	bx += vx
	by += vy

	if bx < 0 {
		bx = 0
		vx = math.Abs(vx)
	}

	if by < 0 {
		by = 0
		vy = math.Abs(vy)
	}

	if bx > 80 {
		bx = 80
		vx = -math.Abs(vx)
	}

	if by > 24 {
		by = 24
		hasBall = true
	}
}

func gameloop() {
	for playing {
		moveBall()
		draw()
		time.Sleep(15 * time.Millisecond)
	}
}

func main() {
	C.initscr()
	C.noecho()
	C.curs_set(0)
	C.keypad(C.stdscr, true)
	C.nodelay(C.stdscr, true)
	C.mousemask(C.REPORT_MOUSE_POSITION|C.ALL_MOUSE_EVENTS, nil)

	draw()
	go func() {
		gameloop()
	}()

	for {
		time.Sleep(15 * time.Millisecond)
		ch := int(C.getch())
		if ch == C.ERR {
			continue
		}

		if ch == 'q' {
			break
		}

		if hasBall && ch == ' ' {
			hasBall = false
			bx = float64(px)
			by = float64(py) - 1.0
			theta := 0.25 * math.Pi
			vx = 0.5 * math.Cos(theta)
			vy = 0.5 * -math.Sin(theta)
		}

		if ch == C.KEY_LEFT {
			px -= 3
			if px < 2 {
				px = 2
			}
		}

		if ch == C.KEY_RIGHT {
			px += 3
			if px > 77 {
				px = 77
			}
		}
	}

	playing = false
	C.endwin()
}
