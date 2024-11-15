package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	Tabs       []string
	TabContent []string
	activeTab  int
}

func NewMainModel() model {
    return model{
        Tabs:       []string{"Home", "Devices", "Logs", "Settings"},
        TabContent: []string{"Welcome to Home!", "List of Devices", "Log History", "Settings Page"}, // Populate TabContent
        activeTab:  0, // Start with the first tab active
    }
}


func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	}

	return m, nil
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

func (m model) View() string {
	doc := strings.Builder{}

    var renderedTabs []string

    for i, t := range m.Tabs {
        var style lipgloss.Style
        isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab
        if isActive {
            style = activeTabStyle
        } else {
            style = inactiveTabStyle
        }
        border, _, _, _, _ := style.GetBorder()
        if isFirst && isActive {
            border.BottomLeft = "│"
        } else if isFirst && !isActive {
            border.BottomLeft = "├"
        } else if isLast && isActive {
            border.BottomRight = "│"
        } else if isLast && !isActive {
            border.BottomRight = "┤"
        }
        style = style.Border(border)
        renderedTabs = append(renderedTabs, style.Render(t))
    }

    row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
    doc.WriteString(row)
    doc.WriteString("\n")

    // Safeguard against empty TabContent
    if m.activeTab < len(m.TabContent) {
        doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.TabContent[m.activeTab]))
    } else {
        doc.WriteString(windowStyle.Render("No content available for this tab."))
    }

    return docStyle.Render(doc.String())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
