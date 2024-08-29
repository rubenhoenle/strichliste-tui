package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const apiBaseUrl = "https://demo.strichliste.org/api"

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

type statusMsg int

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

type User struct {
	ID      int    `json:"id"`
	NAME    string `json:"name"`
	BALANCE int    `json:"balance"`
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func getUsers() []User {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := c.Get(apiBaseUrl + "/user")
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}
	defer res.Body.Close() // nolint:errcheck

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	type Response struct {
		Users []User `json:"users"`
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil
	}

	var ids []int
	for _, user := range response.Users {
		ids = append(ids, user.ID)
	}

	fmt.Println("IDs:", ids)
	return response.Users
}

func main() {
	users := getUsers()
	items := []list.Item{}
	for _, element := range users {
		items = append(items, item{title: element.NAME, desc: "abc"})
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Strichliste"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
