package main

//Librerias necesarias
import (
	"fmt"  //para trabajar con msj en consola
	"strconv"  //convierte strings a enteros
	"time" // manejo de tiempo

	"github.com/atotto/clipboard"  // permite que se puede pegar en el portapapeles
	"github.com/charmbracelet/bubbles/table" // msj en la terminal
	"github.com/charmbracelet/bubbles/textinput" // entrada de texto
	tea "github.com/charmbracelet/bubbletea" // para crear la interfaz Cli 
	"github.com/charmbracelet/lipgloss" // + efectos a Cli 
)


type resetCopiedMsg struct{}

//
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

//Crea y configura el estado inicial
func NewModel() tea.Model {
	ti := textinput.New()
	ti.Placeholder = "8"
	ti.Focus()
	ti.CharLimit = 3
	ti.Width = 5

	columns := []table.Column{
		{Title: "Requisito", Width: 20},
		{Title: "Estado", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{
			{"Longiud", "X"},
			{"Mayuscula", "X"},
			{"Minusculas", "X"},
			{"Numeros", "X"},
			{"Caracter Especial", "X"},
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


//inicializa bubbleTea y activa el input del mouse 
func (m model) Init() tea.Cmd {
	return textinput.Blink
}


//Se encarga de procesar las teclas y actualizar el estado 
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
				m.err = fmt.Errorf("Longitud No Valida")
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


//actualiza la tabla a la hora de crear la contrasena
func (m *model) updateTable() {
	rows := []table.Row{
		{"Longitud", fmt.Sprintf("%s (%d)", formatStatus(m.result.LengthOk), len(m.password))},
		{"Mayusculas", fmt.Sprintf("%s (%d)", formatStatus(m.result.UpperOk), m.result.UpperCount)},
		{"Minuscuas", fmt.Sprintf("%s (%d)", formatStatus(m.result.LowerOk), m.result.LowerCount)},
		{"Numeros", fmt.Sprintf("%s (%d)", formatStatus(m.result.NumOk), m.result.NumCount)},
		{"Caracter Especial", fmt.Sprintf("%s (%d)", formatStatus(m.result.SpecOk), m.result.SpecCount)},
	}
	m.table.SetRows(rows)
}


//convierte el resultado de val en estado visible
func formatStatus(ok bool) string {
	if ok { return "OK" }
	return "NO"
}

//genera la interfaz visual de la applicacion
func (m model) View() string {
	var s string
	s += titleStyle.Render("\t\t\t\t\t\tSECURE PASSWORD CREATOR") + "\n\n"

	// Length
	prefix := "  "
	if m.cursor == 0 { prefix = "> " }
	s += fmt.Sprintf("%sLongitud: %s\n\n", prefix, m.lengthInput.View())

	// Options
	s += "Opciones (Espacio para cambiar):\n"
	options := []struct {
		label string
		val   bool
		idx   int
	}{
		{"Mayusculas", m.options.UseUpper, 1},
		{"Minusculas", m.options.UseLower, 2},
		{"Numeros", m.options.UseNumbers, 3},
		{"Caracteres Especiales", m.options.UseSpecial, 4},
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
				copyMsg = "  Copiado al Portapapeles"
			}
			s += fmt.Sprintf("Contraseña: %s%s\n\n", passwordStyle.Render(m.password), accentStyle.Render(copyMsg))
		s += "Validacion:\n"
		s += baseStyle.Render(m.table.View()) + "\n"
	}

	s += "\nControles: Arriba/Abajo: Mover | Espacio: Alternar | Enter: Generar | 'c': Copiar | Esc: Salir"

	return lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Render(s)
}
