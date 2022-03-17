package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle()
	itemStyle         = lipgloss.NewStyle().PaddingLeft(1)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("5"))
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(1).PaddingBottom(1)
)

type CdfModel struct {
	list  list.Model
	state DirectoryWalkState
}

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}
	str := fmt.Sprint(i)
	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}
	fmt.Fprint(w, fn(str))
}

func getItems(items []string) []list.Item {
	newItems := make([]list.Item, len(items))
	for i, v := range items {
		newItems[i] = item(v)
	}
	return newItems
}

func initializeList(items []string) list.Model {
	newList := list.New([]list.Item(getItems(items)), itemDelegate{}, 0, 0)
	newList.SetShowStatusBar(false)
	newList.SetFilteringEnabled(false)
	newList.SetShowHelp(true)
	newList.SetShowPagination(false)
	newList.Styles.Title = titleStyle
	newList.Styles.HelpStyle = helpStyle
	return newList
}

type NewCandidateList struct {
	state DirectoryWalkState
}

type Walked struct {
	state DirectoryWalkState
}

type OutOfOptions bool

func prepareWalk(state DirectoryWalkState) tea.Cmd {
	return func() tea.Msg {
		newState := state.PrepareWalk()
		switch len(newState.pathOptions) {
		case 0:
			return OutOfOptions(true)
		default:
			return NewCandidateList{newState}
		}
	}
}

func acceptChoice(state DirectoryWalkState, choice string) tea.Cmd {
	return func() tea.Msg {
		newState := state.DoWalk(choice)
		return Walked{newState}
	}
}

func (m CdfModel) Init() tea.Cmd {
	return prepareWalk(m.state)
}

func (m CdfModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.list.SelectedItem() != nil {
				selectedFolder := string(m.list.SelectedItem().(item))
				return m, tea.Batch(acceptChoice(m.state, selectedFolder))
			}
		}
	case NewCandidateList:
		m.state = NewCandidateList(msg).state
		if len(m.state.pathOptions) == 1 {
			return m, acceptChoice(m.state, m.state.pathOptions[0])
		}
		cmd := m.list.SetItems(getItems(m.state.pathOptions))
		m.list.SetHeight(len(m.list.Items()) + 5)
		m.list.Select(0)
		return m, tea.Batch(cmd)
	case Walked:
		m.state = NewCandidateList(msg).state
		if len(m.state.remainingTargetPath) == 0 {
			return m, tea.Quit
		}
		return m, prepareWalk(m.state)
	case OutOfOptions:
		return m, tea.Quit
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m CdfModel) View() string {
	if len(m.state.pathOptions) == 0 {
		return ""
	}
	m.list.Title = "Select folder below " + m.state.pathWalked
	return m.list.View()
}

func Run(initialState DirectoryWalkState) string {
	l := initializeList(make([]string, 0))

	var ui = tea.NewProgram(CdfModel{
		list:  l,
		state: initialState,
	})

	result, err := ui.StartReturningModel()
	if err != nil {
        return ""
	}

    return result.(CdfModel).state.pathWalked
}

