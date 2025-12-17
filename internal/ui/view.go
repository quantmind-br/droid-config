package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func (m Model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	if !m.ready {
		return "Loading...\n"
	}

	primaryColor := lipgloss.Color("39")
	secondaryColor := lipgloss.Color("240")

	modelName := ""
	if currentModel := m.list.CurrentModel(); currentModel != nil {
		modelName = currentModel.DisplayName
	}

	sidebarContent := m.list.View(m.focusArea == FocusSidebar, m.dirty)
	formContent := m.form.View(m.focusArea == FocusForm, modelName)

	if m.focusArea != FocusSidebar {
		sidebarContent = DimmedStyle.Render(sidebarContent)
	}
	if m.focusArea != FocusForm {
		formContent = DimmedStyle.Render(formContent)
	}

	sidebarStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Width(max(0, m.sidebarWidth-2)).
		Height(max(0, m.sidebarHeight-2)).
		Padding(0, 1)
	if m.focusArea == FocusSidebar {
		sidebarStyle = sidebarStyle.BorderForeground(primaryColor)
	}

	formStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Width(max(0, m.formWidth-2)).
		Height(max(0, m.formHeight-2)).
		Padding(0, 1)
	if m.focusArea == FocusForm {
		formStyle = formStyle.BorderForeground(primaryColor)
	}

	sidebar := sidebarStyle.Render(sidebarContent)
	form := formStyle.Render(formContent)

	var content string
	if m.stackedLayout {
		content = lipgloss.JoinVertical(lipgloss.Left, sidebar, form)
	} else {
		content = lipgloss.JoinHorizontal(lipgloss.Top, sidebar, form)
	}

	statusStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Width(max(0, m.width-2)).
		Padding(0, 1)

	statusLineWidth := max(0, m.width-4)
	statusContent := padOrTruncate("Status: "+m.status.View(), statusLineWidth)
	statusBar := statusStyle.Render(statusContent)

	helpStyle := lipgloss.NewStyle().
		Foreground(secondaryColor).
		Padding(0, 1)

	helpText := "tab: form | ↑↓/jk: nav | space: select | a: all | n: new | d: del | ctrl+↑↓: move | ctrl+s: save | ctrl+c: quit"
	if m.focusArea == FocusForm {
		helpText = "tab/↑↓: fields | ←→: provider | ctrl+v: show key | ctrl+s: save | esc: back | ctrl+c: quit"
	}
	helpBar := helpStyle.Render(padOrTruncate(helpText, max(0, m.width-2)))

	full := lipgloss.JoinVertical(lipgloss.Left, content, statusBar, helpBar)

	if m.confirm.Active {
		return m.renderWithModal(full)
	}

	return full
}

func (m Model) renderWithModal(background string) string {
	modalView := m.confirm.View()

	// Use lipgloss.Place to center the modal over the background
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modalView,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("240")),
	)
}

func padOrTruncate(s string, width int) string {
	if width <= 0 {
		return ""
	}
	out := ansi.Truncate(s, width, "")
	if w := lipgloss.Width(out); w < width {
		return out + strings.Repeat(" ", width-w)
	}
	return out
}
