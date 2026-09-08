package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type resetCopiedMsg struct{}

var (
    baseStyle = lipgloss.NewStyle().
        BorderStyle(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color("240")).
        Foreground(lipgloss.Color("250"))

    accentStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#7D56F4")).
        Bold(true)

	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFF00")).
		Bold(true).
		Align(lipgloss.Center).
		Width(50).
		Padding(0, 1)

    passwordStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#00FF00")).
        Bold(true).
        Padding(0, 1).
        Background(lipgloss.Color("#2e3440"))
)

type model struct {
	table       table.Model
	lengthInput textinput.Model
	options     Options
	password    string
	result      ValidationResult
	cursor      int
	err         error
	copied      bool
}

func NewModel() tea.Model {
	ti := textinput.New()
	ti.Placeholder = "8"
	ti.Focus()
	ti.CharLimit = 3
	ti.Width = 5

	columns := []table.Column{
		{Title: "Requirement", Width: 20},
		{Title: "Status", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{
			{"Length", "X"},
			{"Uppercase", "X"},
			{"Lowercase", "X"},
			{"Numbers", "X"},
			{"Special", "X"},
		}),
		table.WithFocused(false),
		table.WithHeight(5),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).Bold(true)
	s.Selected = s.Selected.Foreground(lipgloss.Color("#7D56F4")).Bold(true)
	t.SetStyles(s)

	return model{
		table:       t,
		lengthInput: ti,
		options:     Options{UseUpper: true, UseLower: true, UseNumbers: true, UseSpecial: true},
		cursor:      0,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 { m.cursor-- }
		case "down", "j":
			if m.cursor < 4 { m.cursor++ }
		case " ":
			switch m.cursor {
			case 1: m.options.UseUpper = !m.options.UseUpper
			case 2: m.options.UseLower = !m.options.UseLower
			case 3: m.options.UseNumbers = !m.options.UseNumbers
			case 4: m.options.UseSpecial = !m.options.UseSpecial
			}
		case "c":
			if m.password != "" {
				clipboard.WriteAll(m.password)
				m.copied = true
				return m, tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
					return resetCopiedMsg{}
				})
			}
		case "enter":
			l, err := strconv.Atoi(m.lengthInput.Value())
			if err != nil || l <= 0 {
				m.err = fmt.Errorf("invalid length")
				return m, nil
			}
			m.err = nil
			psw, err := Generate(l, m.options)
			if err != nil {
				m.err = err
				return m, nil
			}
			m.password = psw
			m.result = Validate(psw)
			m.updateTable()
		}
	case resetCopiedMsg:
		m.copied = false
	}

	if m.cursor == 0 {
		m.lengthInput, cmd = m.lengthInput.Update(msg)
	}

	return m, cmd
}

func (m *model) updateTable() {
	rows := []table.Row{
		{"Length", fmt.Sprintf("%s (%d)", formatStatus(m.result.LengthOk), len(m.password))},
		{"Uppercase", fmt.Sprintf("%s (%d)", formatStatus(m.result.UpperOk), m.result.UpperCount)},
		{"Lowercase", fmt.Sprintf("%s (%d)", formatStatus(m.result.LowerOk), m.result.LowerCount)},
		{"Numbers", fmt.Sprintf("%s (%d)", formatStatus(m.result.NumOk), m.result.NumCount)},
		{"Special", fmt.Sprintf("%s (%d)", formatStatus(m.result.SpecOk), m.result.SpecCount)},
	}
	m.table.SetRows(rows)
}

func formatStatus(ok bool) string {
	if ok { return "OK" }
	return "NO"
}

func (m model) View() string {
	var s string
	s += titleStyle.Render("SECURE PASSWORD CREATOR") + "\n\n"

	// Length
	prefix := "  "
	if m.cursor == 0 { prefix = "> " }
	s += fmt.Sprintf("%sLength: %s\n\n", prefix, m.lengthInput.View())

	// Options
	s += "Options (Space to toggle):\n"
	options := []struct {
		label string
		val   bool
		idx   int
	}{
		{"Uppercase", m.options.UseUpper, 1},
		{"Lowercase", m.options.UseLower, 2},
		{"Numbers", m.options.UseNumbers, 3},
		{"Special", m.options.UseSpecial, 4},
	}
	
	for _, opt := range options {
		p := "  "
		if m.cursor == opt.idx { p = "> " }
		check := "[ ]"
		if opt.val { check = "[x]" }
		s += fmt.Sprintf("%s%s %s\n", p, check, opt.label)
	}

	if m.err != nil {
		s += fmt.Sprintf("\nError: %v\n", m.err)
	}

	if m.password != "" {
							copyMsg := ""
			if m.copied {
				copyMsg = "  Copied to clipboard! ✅"
			}
			s += fmt.Sprintf("\nPassword: %s%s\n\n", passwordStyle.Render(m.password), accentStyle.Render(copyMsg))
		s += "Validation:\n"
		s += baseStyle.Render(m.table.View()) + "\n"
	}

	s += "\nControls: Up/Down: Move | Space: Toggle | Enter: Generate | 'c': Copy | Esc: Quit"

	return lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Render(s)
}
