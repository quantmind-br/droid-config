package components

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

type StatusType int

const (
	StatusInfo StatusType = iota
	StatusSuccess
	StatusError
	StatusWarning
)

type Status struct {
	Message   string
	Type      StatusType
	ExpiresAt time.Time
	Width     int
}

func NewStatus() *Status {
	return &Status{
		Width: 80,
	}
}

func (s *Status) Set(msg string, t StatusType, duration time.Duration) {
	s.Message = msg
	s.Type = t
	s.ExpiresAt = time.Now().Add(duration)
}

func (s *Status) SetSuccess(msg string) {
	s.Set(msg, StatusSuccess, 3*time.Second)
}

func (s *Status) SetError(msg string) {
	s.Set(msg, StatusError, 5*time.Second)
}

func (s *Status) SetInfo(msg string) {
	s.Set(msg, StatusInfo, 3*time.Second)
}

func (s *Status) SetWarning(msg string) {
	s.Set(msg, StatusWarning, 4*time.Second)
}

func (s *Status) Clear() {
	s.Message = ""
}

func (s *Status) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *Status) View() string {
	if s.Message == "" || s.IsExpired() {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("Ready")
	}

	var style lipgloss.Style
	switch s.Type {
	case StatusSuccess:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Bold(true)
	case StatusError:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	case StatusWarning:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
	default:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	}

	return style.Render(s.Message)
}
