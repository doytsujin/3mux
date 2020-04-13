/*
Package pane manages one goncurses window displaying a shell.
*/
package pane

import (
	"os"
	"os/exec"
	"os/signal"
	"sync/atomic"
	"syscall"

	"github.com/aaronjanse/3mux/vterm"
	"github.com/kr/pty"
	gc "github.com/rthornton128/goncurses"
)

type Pane struct {
	renderer    *renderer
	Dead        bool
	onDeath     func()
	stdin       chan<- rune
	stdout      chan rune
	ptmx        *os.File
	cmd         *exec.Cmd
	byteCounter uint64
	vt          *vterm.VTerm
}

// NewPane creates a new Pane using an initialized goncurses Window, launching $SHELL
func NewPane(x, y, w, h int, onDeath func()) (*Pane, error) {
	gcWin, err := gc.NewWindow(h, w, y, x)
	if err != nil {
		return nil, err
	}

	p := &Pane{
		renderer: &renderer{
			gcWin,
		},
		stdin:       make(chan rune, 3200000),
		stdout:      make(chan rune, 3200000),
		cmd:         exec.Command(os.Getenv("SHELL")),
		byteCounter: 0,
	}

	ptmx, err := pty.Start(p.cmd)
	if err != nil {
		p.freezeWithError(err)
	}
	p.ptmx = ptmx
	p.resizeShell(w, h)

	// feed ptmx output to stdout channel
	go func() {
		defer func() {
			if r := recover(); r != nil {
				if r.(error).Error() != "send on closed channel" {
					p.freezeWithError(r.(error))
				}
			}
		}()

		for {
			bs := make([]byte, 1000)
			_, err := ptmx.Read(bs)
			if err != nil {
				if err.Error() == "read /dev/ptmx: input/output error" {
					break // ^D
				} else if err.Error() == "EOF" {
					break
				} else {
					panic(err)
				}
			}
			for _, b := range bs {
				atomic.AddUint64(&p.byteCounter, 1)
				p.stdout <- rune(b)
			}
		}
	}()

	p.renderer.gcWin.Print("HEY HEY")
	p.renderer.gcWin.Refresh()

	go func() {
		p.cmd.Wait()
		p.Kill()
		p.onDeath()
	}()

	// FIXME: implement parentSetCursor
	p.vt = vterm.NewVTerm(&p.byteCounter, p.renderer, func(x, y int) {}, p.stdout, p.stdin)
	p.vt.Reshape(x, y, w, h)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				p.freezeWithError(r.(error))
			}
		}()
		p.vt.ProcessStream()
	}()

	return p, nil
}

func (p *Pane) Move(x, y int) {
	// p.renderer.Move(x, y)
}

func (p *Pane) Resize(w, h int) {
	// p.renderer.Resize(w, h)
}

func (p *Pane) resizeShell(w, h int) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			err := pty.Setsize(p.ptmx, &pty.Winsize{
				Rows: uint16(h), Cols: uint16(w),
				X: 16 * uint16(w), Y: 16 * uint16(h),
			})
			if err != nil {
				p.freezeWithError(err)
			}
		}
	}()
	ch <- syscall.SIGWINCH
}

func (p *Pane) Kill() {
	p.Dead = true

	close(p.stdin)
	close(p.stdout)

	// FIXME: handle error
	p.ptmx.Close()
	// FIXME: handle error
	p.cmd.Process.Kill()
}

func (p *Pane) freezeWithError(err error) {
	p.Kill()
	p.renderer.FreezeWithError(err)
}