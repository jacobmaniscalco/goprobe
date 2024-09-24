package components

import (
	"fmt"

	"github.com/jacobmaniscalco/blue-caterpillar-cli/internal/scan"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type scanResultMsg struct {
	result string
	err    error
}

type ScannerModel struct {
	spinner  spinner.Model
	quit     bool
	scanOptions scan.ScanOptions
	Result   string
	Err      error
}

func NewModel(scanOptions scan.ScanOptions) ScannerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return ScannerModel{
		spinner: s,
		quit: false,
		scanOptions: scanOptions,
		Result: "",
		Err: nil,
	}
}

func (m ScannerModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, startScan(m.scanOptions))
}

func (m ScannerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case scanResultMsg: 
		if msg.err != nil {
			m.Err = msg.err
			return m, tea.Quit
		}
		m.Result = msg.result
		return m, tea.Quit

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quit = true
			return m, tea.Quit
		default:
			return m, nil
		}

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m ScannerModel) View() string {

	str := fmt.Sprintf("\n%s Scanning Target(s)...\n\n", m.spinner.View())
	
	if m.quit {
		return str + "\n"
	}
	return str
}

func startScan(scanOptions scan.ScanOptions) tea.Cmd {

	return func() tea.Msg {
		results, err := scan.PerformScan(scanOptions)
		if err != nil {
			return scanResultMsg{result: "", err: err}
		}

		return scanResultMsg{result: results, err: nil}
	}
}
