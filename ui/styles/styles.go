package styles

import "github.com/charmbracelet/lipgloss"


var WarningStyle = lipgloss.NewStyle().
		   Bold(true).
		   Foreground(lipgloss.Color("160"))

var SuccessStyle = lipgloss.NewStyle().
		   Foreground(lipgloss.Color("#00ff00"))

var FailStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("9"))

var IPStyle = lipgloss.NewStyle().
              Foreground(lipgloss.Color("129")).
	      Padding(0, 1)


var OSStyle = lipgloss.NewStyle().
	      Italic(true).
	      Foreground(lipgloss.Color("23")).
	      Padding(0, 1)

var HeaderStyle = lipgloss.NewStyle().
		  Bold(true).
		  Foreground(lipgloss.Color("12")).
		  Background(lipgloss.Color("0")).
		  Padding(0, 1)

var EvenRowStyle = lipgloss.NewStyle().
		   Background(lipgloss.Color("8")).
		   Padding(0,1)

var OddRowStyle = lipgloss.NewStyle().
		  Background(lipgloss.Color("0")).
		  Padding(0,1)

var TitleStyle = lipgloss.NewStyle().
		 Bold(true).
		 Foreground(lipgloss.Color("3")).
		 Background(lipgloss.Color("black")).
		 Align(lipgloss.Center)
		 


