package components

import (
	"github.com/charmbracelet/lipgloss"
)

type ConfirmAction int

const (
	ConfirmDeleteCurrent ConfirmAction = iota
	ConfirmDeleteSelected
)

type Confirm struct {
	Active   bool
	Action   ConfirmAction
	Message  string
	YesLabel string
	NoLabel  string
	Width    int
	Height   int
}

func NewConfirm() *Confirm {
	return &Confirm{
		Active:   false,
		YesLabel: "Yes (y)",
		NoLabel:  "No (n/esc)",
		Width:    40,
		Height:   7,
	}
}

func (c *Confirm) Show(action ConfirmAction, message string) {
	c.Active = true
	c.Action = action
	c.Message = message
}

func (c *Confirm) Hide() {
	c.Active = false
}

func (c *Confirm) View() string {
	if !c.Active {
		return ""
	}

	warningColor := lipgloss.Color("214")
	secondaryColor := lipgloss.Color("240")

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(warningColor).
		Padding(1, 2).
		Align(lipgloss.Center).
		Width(c.Width)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(warningColor)

	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	buttonStyle := lipgloss.NewStyle().
		Foreground(secondaryColor)

	content := titleStyle.Render("CONFIRM") + "\n\n" +
		messageStyle.Render(c.Message) + "\n\n" +
		buttonStyle.Render("["+c.YesLabel+"]  ["+c.NoLabel+"]")

	return modalStyle.Render(content)
}
