package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	textinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	modeIdle   = 0
	modeList   = 1
	modeAdd    = 2
	modeDelete = 3
)

var checkedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("202"))
var cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))

type Task struct {
	description string
	completed   bool
}

type model struct {
	tasks       []Task
	textinput   textinput.Model
	currentMode int
	cursor      int
	fileName    string
}

func LoadTasks(fileName string) ([]Task, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	var tasks []Task
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		completed := false
		if len(record) > 1 && record[1] == "true" {
			completed = true
		}
		tasks = append(tasks, Task{
			description: record[0],
			completed:   completed,
		})
	}
	return tasks, nil
}

func SaveTask(fileName string, tasks []Task) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, task := range tasks {
		record := []string{task.description, strconv.FormatBool(task.completed)}
		if err = writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func initialModel() model {
	filename := "tasks.csv"
	tasks, err := LoadTasks(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading tasks: %v\n", err)
		os.Exit(1)
	}

	ti := textinput.New()
	ti.Placeholder = "Enter the task"
	ti.Focus()
	return model{
		tasks:       tasks,
		textinput:   ti,
		currentMode: modeIdle,
		fileName:    filename,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.currentMode {
		case modeIdle:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit

			case "+":
				m.currentMode = modeAdd
				m.textinput.Reset()
				m.textinput.Placeholder = "Enter a new task"
				m.textinput.Focus()
				return m, nil

			case "l":
				m.currentMode = modeList
				return m, nil

			case "d":
				m.currentMode = modeDelete
				m.textinput.Reset()
				m.textinput.Placeholder = "Enter the task number"
				m.textinput.Focus()
				return m, nil
			}

		case modeAdd:
			m.textinput, cmd = m.textinput.Update(msg)
			switch msg.String() {
			case "enter":
				if m.textinput.Value() != "" {
					m.tasks = append(m.tasks, Task{
						description: m.textinput.Value(),
						completed:   false,
					})
					if err := SaveTask(m.fileName, m.tasks); err != nil {
						fmt.Fprintf(os.Stderr, "Error saving tasks: %v\n", err)
					}
				}
				m.textinput.Reset()
				m.currentMode = modeIdle
				return m, cmd

			case "esc":
				m.textinput.Reset()
				m.currentMode = modeIdle
				return m, cmd
			}

		case modeList:
			switch msg.String() {
			case "up":
				if m.cursor > 0 {
					m.cursor--
				}
				return m, nil

			case "down":
				if m.cursor < len(m.tasks)-1 {
					m.cursor++
				}
				return m, nil

			case "enter":
				if m.cursor < len(m.tasks) {
					m.tasks[m.cursor].completed = !m.tasks[m.cursor].completed
					if err := SaveTask(m.fileName, m.tasks); err != nil {
						fmt.Fprintf(os.Stderr, "Error saving tasks: %v\n", err)
					}
				}
				return m, nil

			case "d":
				for i := len(m.tasks) - 1; i >= 0; i-- {
					if m.tasks[i].completed {
						m.tasks = append(m.tasks[:i], m.tasks[i+1:]...)
					}
				}
				if err := SaveTask(m.fileName, m.tasks); err != nil {
					fmt.Fprintf(os.Stderr, "Error saving tasks: %v\n", err)
				}
				return m, nil

			case "esc":
				m.currentMode = modeIdle
				return m, nil
			}

		case modeDelete:
			m.textinput, cmd = m.textinput.Update(msg)
			switch msg.String() {
			case "enter":
				n, err := strconv.Atoi(m.textinput.Value())
				if err == nil && n > 0 && n <= len(m.tasks) {
					m.tasks = append(m.tasks[:n-1], m.tasks[n:]...)
					SaveTask(m.fileName, m.tasks)
					m.textinput.Reset()
					m.currentMode = modeIdle
				} else {
					m.textinput.SetValue("")
					m.textinput.Placeholder = "Invalid task number! Try again."
				}
				return m, cmd

			case "esc":
				m.textinput.Reset()
				m.currentMode = modeIdle
				return m, cmd
			}
		}
	}
	return m, cmd
}

func (m model) View() string {
	var s string
	switch m.currentMode {
	case modeIdle:
		s += "Enter a command:\n"
		s += "- Press '+' to add a task.\n"
		s += "- Press 'l' to view all tasks.\n"
		s += "- Press 'd' to delete a single task.\n"
		s += "- Press 'q' to quit.\n"

	case modeAdd:
		s += "Add the task:\n"
		s += m.textinput.View()
		s += "\nPress ENTER to save task or ESC to cancel.\n"

	case modeList:
		closeBracket := ")"
		if len(m.tasks) > 0 {
			s += "Your tasks:\n"
			for i, t := range m.tasks {
				stringI := strconv.Itoa(i + 1)
				cursor := " "
				if i == m.cursor {
					cursor = ">"
				}
				if t.completed {
					s += (fmt.Sprintf("%s %s%s %s\n", cursor, checkedStyle.Render(stringI), checkedStyle.Render(closeBracket), t.description))
				} else if i == m.cursor {
					s += (fmt.Sprintf("%s %s%s %s\n", (cursor), cursorStyle.Render(stringI), cursorStyle.Render(closeBracket), cursorStyle.Render(t.description)))
				} else {
					s += fmt.Sprintf("%s %s) %s\n", cursor, stringI, t.description)
				}
			}
			s += "\nPress ENTER to toggle completion.\nPress 'd' to delete completed tasks\nPress ESC to go back.\n"
		} else {
			s += "No tasks added.\nPress ESC to go back."
		}

	case modeDelete:
		if len(m.tasks) > 0 {
			s += "Enter the task number to delete:\n"
			s += m.textinput.View()
			s += "\nPress ENTER to delete the task or ESC to cancel.\n"
		} else {
			s += "No tasks to delete.\nPress ESC to go back."
		}
	}
	return s
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting program: %v\n", err)
		os.Exit(1)
	}
}