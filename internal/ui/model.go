package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/diogo/droid-config/internal/config"
	"github.com/diogo/droid-config/internal/ui/components"
)

type FocusArea int

const (
	FocusSidebar FocusArea = iota
	FocusForm
)

type Model struct {
	config        *config.ConfigData
	list          *components.List
	form          *components.Form
	status        *components.Status
	confirm       *components.Confirm
	help          help.Model
	focusArea     FocusArea
	width         int
	height        int
	sidebarWidth  int
	formWidth     int
	sidebarHeight int
	formHeight    int
	contentHeight int
	stackedLayout bool
	ready         bool
	quitting      bool
	dirty         bool
}

func NewModel() Model {
	cfg, _ := config.Load()
	if cfg == nil {
		cfg = &config.ConfigData{CustomModels: []config.CustomModel{}}
	}

	list := components.NewList()
	list.SetItems(cfg.CustomModels)

	form := components.NewForm()
	if len(cfg.CustomModels) > 0 {
		form.LoadModel(&cfg.CustomModels[0])
	}

	return Model{
		config:    cfg,
		list:      list,
		form:      form,
		status:    components.NewStatus(),
		confirm:   components.NewConfirm(),
		help:      help.New(),
		focusArea: FocusSidebar,
		ready:     false,
	}
}

type statusClearMsg struct{}

func statusClearCmd() tea.Cmd {
	return tea.Tick(3*time.Second, func(time.Time) tea.Msg {
		return statusClearMsg{}
	})
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Layout: content area (panels) + status bar (3 lines) + help line (1 line).
		statusHeight := 3
		helpHeight := 1
		m.contentHeight = msg.Height - statusHeight - helpHeight
		if m.contentHeight < 4 {
			m.contentHeight = 4
		}

		minSidebar := 28
		minForm := 44
		if msg.Width < minSidebar+minForm {
			m.stackedLayout = true
			m.sidebarWidth = msg.Width
			m.formWidth = msg.Width

			// Split vertical space between list and form.
			minPanelHeight := 4
			sidebarHeight := m.contentHeight / 3
			if sidebarHeight < minPanelHeight {
				sidebarHeight = minPanelHeight
			}
			formHeight := m.contentHeight - sidebarHeight
			if formHeight < minPanelHeight {
				formHeight = minPanelHeight
				sidebarHeight = m.contentHeight - formHeight
				if sidebarHeight < minPanelHeight {
					// Too small to satisfy minimums; split evenly.
					sidebarHeight = m.contentHeight / 2
					formHeight = m.contentHeight - sidebarHeight
				}
			}
			m.sidebarHeight = sidebarHeight
			m.formHeight = formHeight
		} else {
			m.stackedLayout = false
			sidebarWidth := msg.Width / 3
			if sidebarWidth < minSidebar {
				sidebarWidth = minSidebar
			}
			if sidebarWidth > msg.Width-minForm {
				sidebarWidth = msg.Width - minForm
			}
			m.sidebarWidth = sidebarWidth
			m.formWidth = msg.Width - sidebarWidth
			m.sidebarHeight = m.contentHeight
			m.formHeight = m.contentHeight
		}

		// Panels have a border (2 cols) and inner horizontal padding (2 cols).
		// Components receive the inner content width (outer - 4).
		m.list.Width = max(1, m.sidebarWidth-4)
		m.list.Height = max(1, m.sidebarHeight-2)
		m.form.Width = max(1, m.formWidth-4)
		m.form.Height = max(1, m.formHeight-2)

		m.status.Width = msg.Width
		m.confirm.Width = min(40, max(20, msg.Width-10))
		m.ready = true
		return m, nil

	case statusClearMsg:
		if m.status.IsExpired() {
			m.status.Clear()
		}
		return m, nil

	case tea.KeyMsg:
		if m.confirm.Active {
			return m.handleConfirmKeys(msg)
		}

		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "tab":
			if m.focusArea == FocusSidebar {
				m.focusArea = FocusForm
				m.form.Focus()
			} else {
				m.form.FocusNext()
			}
			return m, nil

		case "shift+tab":
			if m.focusArea == FocusForm {
				if m.form.FocusIndex() == 0 {
					m.focusArea = FocusSidebar
					m.form.Blur()
				} else {
					m.form.FocusPrev()
				}
			}
			return m, nil

		case "esc":
			if m.focusArea == FocusForm {
				m.focusArea = FocusSidebar
				m.form.Blur()
				if currentModel := m.list.CurrentModel(); currentModel != nil {
					m.form.LoadModel(currentModel)
				}
			}
			return m, nil

		case "ctrl+s":
			return m.saveCurrentModel()

		case "ctrl+v":
			if m.focusArea == FocusForm {
				m.form.ToggleAPIKeyVisibility()
			}
			return m, nil

		case "ctrl+up":
			if m.focusArea == FocusSidebar {
				if m.list.MoveItemUp() {
					m.dirty = true
					m.status.SetSuccess("Model moved up")
					return m.saveConfig()
				}
			}
			return m, nil

		case "ctrl+down":
			if m.focusArea == FocusSidebar {
				if m.list.MoveItemDown() {
					m.dirty = true
					m.status.SetSuccess("Model moved down")
					return m.saveConfig()
				}
			}
			return m, nil
		}

		if m.focusArea == FocusSidebar {
			return m.handleSidebarKeys(msg)
		} else {
			return m.handleFormKeys(msg)
		}
	}

	if m.focusArea == FocusForm && m.form.FocusIndex() != components.FieldProvider {
		if input := m.form.CurrentInput(); input != nil {
			newInput, cmd := input.Update(msg)
			m.form.UpdateInput(newInput)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleSidebarKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.list.MoveUp()
		if currentModel := m.list.CurrentModel(); currentModel != nil {
			m.form.LoadModel(currentModel)
		}
		return m, nil

	case "down", "j":
		m.list.MoveDown()
		if currentModel := m.list.CurrentModel(); currentModel != nil {
			m.form.LoadModel(currentModel)
		}
		return m, nil

	case " ":
		m.list.ToggleSelected()
		return m, nil

	case "a":
		m.list.SelectAll()
		return m, nil

	case "n":
		return m.addNewModel()

	case "d":
		return m.handleDelete()

	case "enter":
		m.focusArea = FocusForm
		m.form.Focus()
		return m, nil
	}

	return m, nil
}

func (m Model) handleFormKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.form.FocusIndex() == components.FieldProvider {
		switch msg.String() {
		case "left", "h":
			m.form.PrevProvider()
			return m, nil
		case "right", "l":
			m.form.NextProvider()
			return m, nil
		case "up", "k":
			m.form.FocusPrev()
			return m, nil
		case "down", "j":
			m.form.FocusNext()
			return m, nil
		}
	} else {
		switch msg.String() {
		case "up":
			m.form.FocusPrev()
			return m, nil
		case "down":
			m.form.FocusNext()
			return m, nil
		case "enter":
			m.form.FocusNext()
			return m, nil
		}

		if input := m.form.CurrentInput(); input != nil {
			newInput, cmd := input.Update(msg)
			m.form.UpdateInput(newInput)
			return m, cmd
		}
	}

	return m, nil
}

func (m Model) handleConfirmKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		action := m.confirm.Action
		m.confirm.Hide()

		switch action {
		case components.ConfirmDeleteCurrent:
			if m.list.DeleteCurrent() {
				m.dirty = true
				m.status.SetSuccess("Model deleted")
				if currentModel := m.list.CurrentModel(); currentModel != nil {
					m.form.LoadModel(currentModel)
				} else {
					m.form.LoadModel(nil)
				}
				return m.saveConfig()
			}
		case components.ConfirmDeleteSelected:
			count := m.list.DeleteSelected()
			if count > 0 {
				m.dirty = true
				m.status.SetSuccess("Deleted " + string(rune('0'+count)) + " model(s)")
				if currentModel := m.list.CurrentModel(); currentModel != nil {
					m.form.LoadModel(currentModel)
				} else {
					m.form.LoadModel(nil)
				}
				return m.saveConfig()
			}
		}
		return m, nil

	case "n", "N", "esc":
		m.confirm.Hide()
		return m, nil
	}

	return m, nil
}

func (m Model) addNewModel() (tea.Model, tea.Cmd) {
	newModel := config.CustomModel{
		DisplayName: "New Model",
		Provider:    "openai",
	}
	m.list.AddModel(newModel)
	m.form.LoadModel(&newModel)
	m.focusArea = FocusForm
	m.form.SetFocusIndex(0)
	m.form.Focus()
	m.dirty = true
	m.status.SetInfo("New model created - edit and save")

	return m.saveConfig()
}

func (m Model) handleDelete() (tea.Model, tea.Cmd) {
	selected := m.list.GetSelectedIndices()
	if len(selected) > 0 {
		m.confirm.Show(components.ConfirmDeleteSelected,
			"Delete "+string(rune('0'+len(selected)))+" selected model(s)?")
	} else if len(m.list.Items) > 0 {
		modelName := m.list.CurrentModel().DisplayName
		if modelName == "" {
			modelName = "this model"
		}
		m.confirm.Show(components.ConfirmDeleteCurrent,
			"Delete \""+modelName+"\"?")
	}
	return m, nil
}

func (m Model) saveCurrentModel() (tea.Model, tea.Cmd) {
	valid, errMsg := m.form.Validate()
	if !valid {
		m.status.SetError(errMsg)
		return m, statusClearCmd()
	}

	updatedModel := m.form.GetModel()
	m.list.UpdateCurrentModel(updatedModel)
	m.dirty = true
	m.status.SetSuccess("Changes saved!")

	return m.saveConfig()
}

func (m Model) saveConfig() (tea.Model, tea.Cmd) {
	m.config.CustomModels = m.list.GetModels()
	if err := config.Save(m.config); err != nil {
		m.status.SetError("Failed to save: " + err.Error())
		return m, statusClearCmd()
	}
	m.dirty = false
	return m, statusClearCmd()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
