package tui

import (
	"fmt"
	open_meteo_api "go_education/weather/open_meteo_api"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DataFetchedMessage struct {
	Action  string
	Payload open_meteo_api.ForecastResponse
}

type model struct {
	Tabs       [3]string
	activeTab  int
	TabContent [3]string
	ClientData string
}

type TabContentWidget struct {
	Error   [3]string
	Loading [3]string
	Content [3]string
}

func buildTable(forecast *open_meteo_api.ForecastResponse) [3]string {

	columns := make([]table.Column, 0)
	temperatureRow := strings.Builder{}
	precipitationRow := strings.Builder{}
	windRow := strings.Builder{}
	separator := "!"

	for i, v := range forecast.Hourly.HourlyTimeTime {
		columns = append(columns, table.Column{Title: strings.Split(v, "T")[1], Width: 6})
		temperatureRow.WriteString(fmt.Sprintf("%v %s%s", math.Ceil(forecast.Hourly.Temperature[i]), forecast.Units.Temperature, separator))
		precipitationRow.WriteString(fmt.Sprintf("%v%s%s", math.Ceil(forecast.Hourly.Precipitation[i]), forecast.Units.Precipitation, separator))
		windRow.WriteString(fmt.Sprintf("%v%s%s", math.Ceil(forecast.Hourly.WindSpeed[i]), forecast.Units.WindSpeed, separator))
	}

	temperatureRows := []table.Row{
		strings.Split(temperatureRow.String(), separator)[0:24],
	}
	precipitationRows := []table.Row{
		strings.Split(precipitationRow.String(), separator)[0:24],
	}
	windRows := []table.Row{
		strings.Split(windRow.String(), separator)[0:24],
	}

	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	temperatureTable := table.New(
		table.WithColumns(columns),
		table.WithRows(temperatureRows),
		table.WithHeight(3),
	)
	precipitationTable := table.New(
		table.WithColumns(columns),
		table.WithRows(precipitationRows),
		table.WithHeight(3),
	)
	windTable := table.New(
		table.WithColumns(columns),
		table.WithRows(windRows),
		table.WithHeight(3),
	)

	temperatureTable.SetStyles(s)
	precipitationTable.SetStyles(s)
	windTable.SetStyles(s)

	tableVec := [3]string{temperatureTable.View(), precipitationTable.View(), windTable.View()}
	return tableVec
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
	case DataFetchedMessage:
		if msg.Action == "error" {
			m.TabContent = [3]string{"Temperature fetch error", "Precipitation fetch error", "Wind fetch error"}

		} else {
			t := time.Now()
			m.ClientData = fmt.Sprintf(
				"Getting forecast by IP: %s\nCurrent location: %s\nCurrent time: %02d:%02d",
				msg.Payload.IpAdd,
				msg.Payload.TimeZone,
				t.Hour(),
				t.Minute(),
			)
			m.TabContent = buildTable(&msg.Payload)
		}
		return m, nil

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
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 27)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

func (m model) View() string {
	doc := strings.Builder{}
	doc.WriteString(fmt.Sprintf("\n\n\n%s\n", m.ClientData))
	var renderedTabs []string

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
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
	doc.WriteString(
		windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).
			Render(m.TabContent[m.activeTab]))
	return docStyle.Render(doc.String())
}

func initialModel() model {
	return model{
		Tabs:       [3]string{"Temperature", "Precipitation", "Wind"},
		activeTab:  0,
		TabContent: [3]string{"Loading temperature...", "Loading precipitation...", "Loading wind..."},
		ClientData: "Getting forecast by IP:\nCurrent location:\nCurrent time:",
	}
}

func CreateApp() *tea.Program {

	app := tea.NewProgram(initialModel())

	return app
}
