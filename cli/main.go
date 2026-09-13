package main 


//importo librerias
import (
	"fmt" 
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

//Func main() se crea un nuevo programa BubbleTea usando NewModel()
// si el programa tiene un error lo guarda en err y sino ejecuta 
func main() {
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error de Programa: %v\n", err)
		os.Exit(1)
	}
}
