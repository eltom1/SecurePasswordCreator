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


type msgResetCopiado struct{}

//
var (
	estiloBase = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Foreground(lipgloss.Color("250"))

	estiloAcento = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true)

	estiloTitulo = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFF00")).
		Bold(true).
		Align(lipgloss.Center).
		Width(50).
		Padding(0, 1)

	estiloClave = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Bold(true).
		Padding(0, 1).
		Background(lipgloss.Color("#2e3440"))
)


type modelo struct {
	table       table.Model
	inputLargo  textinput.Model
	opciones     Opciones
	clave       string
	res         ResultadoVal
	cursor      int
	err         error
	copiado      bool
}

//Crea y configura el estado inicial
func NuevoModelo() tea.Model {
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
			{"Longitud", "X"},
			{"Mayusculas", "X"},
			{"Minusculas", "X"},
			{"Numeros", "X"},
			{"Caracteres Especiales", "X"},
		}),
		table.WithFocused(false),
		table.WithHeight(5),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).Bold(true)
	s.Selected = s.Selected.Foreground(lipgloss.Color("#7D56F4")).Bold(true)
	t.SetStyles(s)

	return modelo{
		table:       t,
		inputLargo:  ti,
		opciones:     Opciones{UsaMayus: true, UsaMinus: true, UsaNum: true, UsaEsp: true},
		cursor:      0,
	}
}


//inicializa bubbleTea y activa el input del mouse
func (m modelo) Init() tea.Cmd {
	return textinput.Blink
}


//Se encarga de procesar las teclas y actualizar el estado
func (m modelo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			case 1: m.opciones.UsaMayus = !m.opciones.UsaMayus
			case 2: m.opciones.UsaMinus = !m.opciones.UsaMinus
			case 3: m.opciones.UsaNum = !m.opciones.UsaNum
			case 4: m.opciones.UsaEsp = !m.opciones.UsaEsp
			}
		case "c":
			if m.clave != "" {
				clipboard.WriteAll(m.clave)
				m.copiado = true
				return m, tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
					return msgResetCopiado{}
				})
			}
		case "enter":
			l, err := strconv.Atoi(m.inputLargo.Value())
			if err != nil || l <= 0 {
				m.err = fmt.Errorf("Longitud No Valida")
				return m, nil
			}
			m.err = nil
			psw, err := Generar(l, m.opciones)
			if err != nil {
				m.err = err
				return m, nil
			}
			m.clave = psw
			m.res = Validar(psw)
			m.actualizarTabla()
		}
	case msgResetCopiado:
		m.copiado = false
	}

	if m.cursor == 0 {
		m.inputLargo, cmd = m.inputLargo.Update(msg)
	}

	return m, cmd
}


//actualiza la tabla a la hora de crear la contrasena
func (m *modelo) actualizarTabla() {
	rows := []table.Row{
		{"Longitud", fmt.Sprintf("%s (%d)", formatoEstado(m.res.LargoOk), len(m.clave))},
		{"Mayusculas", fmt.Sprintf("%s (%d)", formatoEstado(m.res.MayusOk), m.res.CantMayus)},
		{"Minusculas", fmt.Sprintf("%s (%d)", formatoEstado(m.res.MinusOk), m.res.CantMinus)},
		{"Numeros", fmt.Sprintf("%s (%d)", formatoEstado(m.res.NumOk), m.res.CantNum)},
		{"Caracteres Especiales", fmt.Sprintf("%s (%d)", formatoEstado(m.res.EspOk), m.res.CantEsp)},
	}
	m.table.SetRows(rows)
}


//convierte el resultado de val en estado visible
func formatoEstado(ok bool) string {
	if ok { return "OK" }
	return "NO"
}

//genera la interfaz visual de la applicacion
func (m modelo) View() string {
	var s string
	s += estiloTitulo.Render("SECURE PASSWORD CREATOR") + "\n\n"

	// Length
	prefix := "  "
	if m.cursor == 0 { prefix = "> " }
	s += fmt.Sprintf("%sLongitud: %s\n\n", prefix, m.inputLargo.View())

	// Options
	s += "Opciones (Espacio para cambiar):\n"
	options := []struct {
		label string
		val   bool
		idx   int
	}{
		{"Mayusculas", m.opciones.UsaMayus, 1},
		{"Minusculas", m.opciones.UsaMinus, 2},
		{"Numeros", m.opciones.UsaNum, 3},
		{"Caracteres Especiales", m.opciones.UsaEsp, 4},
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

	if m.clave != "" {
		copyMsg := ""
		if m.copiado {
			copyMsg = "  Copiado al Portapapeles"
		}
		s += fmt.Sprintf("Contraseña: %s%s\n\n", estiloClave.Render(m.clave), estiloAcento.Render(copyMsg))
		s += "Validacion:\n"
		s += estiloBase.Render(m.table.View()) + "\n"
	}

	s += "\nControles: Arriba/Abajo: Mover | Espacio: Alternar | Enter: Generar | 'c': Copiar | Esc: Salir"

	return lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Render(s)
}
