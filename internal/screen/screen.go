package screen

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/eiannone/keyboard"
	"github.com/snusEbjoer/projectManager/internal/config"
	"github.com/snusEbjoer/projectManager/internal/state"
)

const characters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var (
	focusedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA"))
	inputStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Blink(true)
	hintsStyle   = lipgloss.NewStyle().Faint(true)
)

type Screen struct {
	controls  *Controls
	eventator *Eventator
	state     *state.State[[]Project]
	curr      *state.State[int]
	input     *state.State[string]
	projects  []Project
	config    config.Config
}

func NewScreen() *Screen {
	e := NewEventator()
	s := &Screen{
		controls:  NewControls(e.GetChan()),
		eventator: e,
		state:     state.UseState([]Project{}),
		input:     state.UseState(""),
		curr:      state.UseState(0),
		projects:  NewFetcher().FetchProjects(),
		config:    config.NewConfig(),
	}
	s.init()
	return s
}

func (s *Screen) Render() string {
	var buffer bytes.Buffer
	for z, v := range s.state.Value() {
		if z == s.curr.Value() {
			v.name = focusedStyle.Render(v.name)
		}
		buffer.WriteString(v.name + "\n")
	}

	buffer.WriteString("\n" + inputStyle.Render("search: "+s.input.Value()))
	buffer.WriteString(hintsStyle.Render("\n\nType project name to search"))
	buffer.WriteString(hintsStyle.Render("\nPress enter to open project"))
	buffer.WriteString(hintsStyle.Render("\nPress ESC to exit\n"))
	return buffer.String()
}

func (s *Screen) initInput() {
	for _, char := range characters {
		s.controls.AddHandler(RuneKey(string(char)), func(a ...any) {
			s.input.SetState(func(prev string) string {
				return prev + string(char)
			})
		})
	}
	s.controls.AddHandler(OtherKey(keyboard.KeyBackspace2), func(a ...any) {
		s.input.SetState(func(prev string) string {
			if len(prev) < 1 {
				return prev
			}
			return prev[0 : len(prev)-1]
		})
	})
	s.controls.AddHandler(OtherKey(keyboard.KeySpace), func(a ...any) {
		s.input.SetState(func(prev string) string {
			return prev + " "
		})
	})
}

func (s *Screen) init() {
	s.state.SetState(func(prev []Project) []Project {
		return s.projects
	})

	s.controls.AddHandler(OtherKey(keyboard.KeyArrowUp), func(a ...any) {
		if s.curr.Value() == 0 {
			s.curr.SetState(func(prev int) int {
				if len(s.state.Value()) > 0 {
					return len(s.state.Value()) - 1
				}
				return 0
			})
		} else {
			s.curr.SetState(func(prev int) int {
				return prev - 1
			})
		}
	})
	s.controls.AddHandler(OtherKey(keyboard.KeyArrowDown), func(a ...any) {
		if s.curr.Value() == len(s.state.Value())-1 {
			s.curr.SetState(func(prev int) int {
				return 0
			})
		} else {
			s.curr.SetState(func(prev int) int {
				return prev + 1
			})
		}
	})
	s.controls.AddHandler(OtherKey(keyboard.KeyEsc), func(a ...any) {
		log.Fatal("BYE BYE")
	})
	s.initInput()
	s.controls.AddHandler(OtherKey(keyboard.KeyEnter), func(a ...any) {
		Shellout(fmt.Sprintf("code %s", s.projects[s.curr.Value()].workDir+s.projects[s.curr.Value()].name))
	})
}

func hasPrefix(s string, prefix string) bool { // fuck case sensitivity
	return strings.HasPrefix(strings.ToLower(s), strings.ToLower(prefix))
}

func (s *Screen) Search() {
	s.state.SetState(func(prev []Project) []Project {
		if s.input.Value() == "" {
			return s.projects
		}
		res := []Project{}
		for _, el := range s.projects {
			if hasPrefix(el.name, s.input.Value()) {
				res = append(res, el)
			}
		}
		s.curr.SetState(func(prev int) int {
			return 0
		})
		return res
	})
}

func (s *Screen) rerender() {
	fmt.Print("\033[H\033[J")
	fmt.Println(s.Render())
}

func (s *Screen) Run() {
	s.rerender()

	go state.UseEffect(func() {
		s.Search()
	}, []state.StateAny{s.input})

	state.UseEffect(func() {
		s.rerender()
	}, []state.StateAny{s.curr, s.state})
}
