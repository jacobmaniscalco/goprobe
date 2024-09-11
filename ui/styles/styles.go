package styles

import "github.com/charmbracelet/lipgloss"


var WarningStyle = lipgloss.NewStyle().
		   Bold(true).
		   Foreground(lipgloss.Color("#FF0000"))


var IPStyle = lipgloss.NewStyle().
	      Bold(true).
              Foreground(lipgloss.Color("#00FF00"))


var OSStyle = lipgloss.NewStyle().
	      Italic(true).
	      Foreground(lipgloss.Color("4")).
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
		 Padding(0,1)


