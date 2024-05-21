package main

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/0x3alex/gee"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	result            int
	allowedNums       [4]string
	allowedOperations []string
	history           []string
	historyIdx        = len(history) + 1
	infoStr           string
)

func newQuiz() {
	min, max := 10, 100
	result = rand.Intn(max-min+1) + min
	for i := range len(allowedNums) {
		allowedNums[i] = strconv.Itoa(rand.Intn(9-1+1) + 1)
	}
}

func validateAnswer(s string) bool {
	history = append(history, s)
	valT, res, err := gee.Eval(s)
	if err != nil {
		infoStr = "an error occured"
		return false
	}
	if valT == 0 {
		i := res.(float64)
		if i == float64(result) {
			infoStr = "Correct! Next Question."
			return true
		} else {
			infoStr = fmt.Sprintf("Wrong! Your input is: %.2f", i)
			return false
		}
	} else {
		return false
	}

}

func main() {
	allowedOperations = []string{"+", "-", "*", "/", "^", "(", ")"}
	newQuiz()
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type model struct {
	textInput textinput.Model
	err       error
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Your solution"
	ti.Focus()

	return model{
		textInput: ti,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			historyIdx = len(history) + 1
			if validateAnswer(m.textInput.Value()) {
				newQuiz()
				history = history[:0]
			}
			m.textInput.SetValue("")
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyUp:
			if historyIdx > 0 {
				historyIdx--
			}
			if historyIdx >= 0 && historyIdx < len(history) {
				m.textInput.SetValue(history[historyIdx])
			}

		case tea.KeyDown:
			if historyIdx <= len(history) {
				historyIdx++
			}
			if historyIdx < len(history) {
				m.textInput.SetValue(history[historyIdx])
			} else {
				m.textInput.SetValue("")
			}
		case tea.KeyBackspace:
		case tea.KeyRunes:
			str := msg.String()
			flag := false
			for _, v := range allowedNums {
				if str == v {
					flag = true
				}
			}
			for _, v := range allowedOperations {
				if str == v {
					flag = true
				}
			}
			if !flag {
				return m, cmd
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var allowedOpString, allowedNumsString string
	allowedOpString = strings.Join(allowedOperations, ",")
	allowedNumsString = strings.Join(allowedNums[:], ",")
	return fmt.Sprintf(
		"=Mathemann-cli=\n%s\nDesired result :%d\nAllowed Numbers: %s\nAllowed Operations: %s\n%s",
		infoStr, result, allowedNumsString, allowedOpString, m.textInput.View(),
	) + "\n"
}
