package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	mets "github.com/Diogenesoftoronto/ardi/internal/mets"
	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

// model for the file picker
type model struct {
	picker   filepicker.Model
	selected string
	quit     bool
	err      error
	tea.Model
}

// the leader key will be configurable with viper in the future.
var leader = ":"

var (
	invalidPath = fmt.Sprint("The path is not valid")
)

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return m.picker.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", leader + "q":
			m.quit = true
			return m, tea.Quit
		}
	case clearErrorMsg:
		m.err = nil
	}
	var cmd tea.Cmd
	m.picker, cmd = m.picker.Update(msg)

	if selected, path := m.picker.DidSelectFile(msg); selected {
		// we can turn selected into an array
		// but lets check if this works first
		m.selected = path
	}
	if selected, path := m.picker.DidSelectDisabledFile(msg); selected {
		m.err = errors.New(path + invalidPath)
		m.selected = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}
	return m, cmd

}

func (m model) View() string {
	if m.quit {
		// a log on quit would be nice. especially using the same logger
		// passed into the program.
		return ""
	}
	var s strings.Builder
	s.WriteString("\n ")
	if m.err != nil {
		s.WriteString(m.picker.Styles.DisabledCursor.Render(m.err.Error()))
	} else if m.selected == "" {
		s.WriteString("Choose a file:")
	} else {
		s.WriteString("Currently selected: " + m.picker.Styles.Selected.Render(m.selected))
	}
	s.WriteString("\n\n" + m.picker.View() + "\n")
	return s.String()
}

func main() {
	fp := filepicker.New()
	// allowed extensions
	fp.AllowedTypes = []string{mets.ZIP, mets.TAR, mets.Z7, mets.XML}
	fp.CurrentDirectory, _ = os.UserHomeDir()

	m := model{
		picker: fp,
	}
	tm, _ := tea.NewProgram(&m, tea.WithOutput(os.Stderr)).Run()
	mm := tm.(model)
	fmt.Println("\n You selected: " + m.picker.Styles.Selected.Render(mm.selected) + "\n")
}
