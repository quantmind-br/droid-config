package ui

import "github.com/charmbracelet/lipgloss"

var (
	primaryColor   = lipgloss.Color("39")
	secondaryColor = lipgloss.Color("240")
	accentColor    = lipgloss.Color("205")
	successColor   = lipgloss.Color("82")
	errorColor     = lipgloss.Color("196")
	warningColor   = lipgloss.Color("214")
	dimmedColor    = lipgloss.Color("242")

	DimmedStyle = lipgloss.NewStyle().Foreground(dimmedColor)

	TitleBackgroundStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("0")).
				Background(primaryColor).
				Padding(0, 1)

	ProviderBadges = map[string]string{
		"anthropic":                   lipgloss.NewStyle().Foreground(lipgloss.Color("141")).Render("[A]"),
		"openai":                      lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Render("[O]"),
		"generic-chat-completion-api": lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Render("[G]"),
	}

	ValidationErrorStyle = lipgloss.NewStyle().Foreground(errorColor).Bold(true)
	HintStyle            = lipgloss.NewStyle().Foreground(dimmedColor).Italic(true)

	SidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(0, 1)

	FormStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(0, 1)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor)

	SelectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("0")).
			Background(primaryColor)

	FocusedStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	BlurredStyle = lipgloss.NewStyle().
			Foreground(secondaryColor)

	LabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Width(14)

	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(0, 1)

	FocusedInputStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor).
				Padding(0, 1)

	StatusBarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(0, 1)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(secondaryColor)

	CheckboxChecked   = lipgloss.NewStyle().Foreground(successColor).Render("[x]")
	CheckboxUnchecked = lipgloss.NewStyle().Foreground(secondaryColor).Render("[ ]")

	ModalStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(warningColor).
			Padding(1, 2).
			Align(lipgloss.Center)

	ButtonStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor)

	FocusedButtonStyle = lipgloss.NewStyle().
				Padding(0, 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor).
				Foreground(primaryColor).
				Bold(true)

	// Scrollbar styles
	ScrollbarTrackChar = "░"
	ScrollbarThumbChar = "█"
	ScrollbarTrackStyle = lipgloss.NewStyle().Foreground(dimmedColor)
	ScrollbarThumbStyle = lipgloss.NewStyle().Foreground(primaryColor)

	// Edge indicator styles
	EdgeIndicatorUp   = "▲"
	EdgeIndicatorDown = "▼"
	EdgeIndicatorStyle = lipgloss.NewStyle().Foreground(secondaryColor)
)
