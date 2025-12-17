package components

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/diogo/droid-config/internal/config"
)

const (
	FieldDisplayName = iota
	FieldModelID
	FieldBaseURL
	FieldAPIKey
	FieldProvider
	FieldMaxTokens
	FieldCount
)

type Form struct {
	inputs          []textinput.Model
	providerIndex   int
	focusIndex      int
	Width           int
	Height          int
	showAPIKey      bool
	validationError map[int]string
}

func NewForm() *Form {
	f := &Form{
		inputs:          make([]textinput.Model, FieldCount),
		providerIndex:   0,
		focusIndex:      0,
		Width:           50,
		Height:          20,
		showAPIKey:      false,
		validationError: make(map[int]string),
	}

	for i := range f.inputs {
		t := textinput.New()
		t.CharLimit = 256

		switch i {
		case FieldDisplayName:
			t.Placeholder = "Display Name (required)"
		case FieldModelID:
			t.Placeholder = "e.g., gpt-4, claude-3-opus"
		case FieldBaseURL:
			t.Placeholder = "https://api.example.com/v1"
		case FieldAPIKey:
			t.Placeholder = "sk-..."
			t.EchoMode = textinput.EchoPassword
		case FieldProvider:
			// Provider is handled separately as a selector
		case FieldMaxTokens:
			t.Placeholder = "e.g., 4096"
		}

		f.inputs[i] = t
	}

	f.inputs[0].Focus()
	return f
}

func (f *Form) inputTextWidth() int {
	// Inner width available for the textinput (excluding the input box border+padding=4).
	w := f.Width - 4
	if w < 1 {
		w = 1
	}
	return w
}

func padRight(s string, width int) string {
	if width <= 0 {
		return ""
	}
	current := lipgloss.Width(s)
	if current >= width {
		return ansi.Truncate(s, width, "")
	}
	return s + strings.Repeat(" ", width-current)
}

func (f *Form) LoadModel(m *config.CustomModel) {
	if m == nil {
		for i := range f.inputs {
			f.inputs[i].SetValue("")
		}
		f.providerIndex = 0
		return
	}

	f.inputs[FieldDisplayName].SetValue(m.DisplayName)
	f.inputs[FieldModelID].SetValue(m.Model)
	f.inputs[FieldBaseURL].SetValue(m.BaseURL)
	f.inputs[FieldAPIKey].SetValue(m.APIKey)
	f.inputs[FieldMaxTokens].SetValue("")
	if m.MaxTokens > 0 {
		f.inputs[FieldMaxTokens].SetValue(strconv.Itoa(m.MaxTokens))
	}

	f.providerIndex = 0
	for i, p := range config.Providers {
		if p == m.Provider {
			f.providerIndex = i
			break
		}
	}
}

func (f *Form) GetModel() config.CustomModel {
	maxTokens := 0
	if v, err := strconv.Atoi(f.inputs[FieldMaxTokens].Value()); err == nil {
		maxTokens = v
	}

	provider := ""
	if f.providerIndex >= 0 && f.providerIndex < len(config.Providers) {
		provider = config.Providers[f.providerIndex]
	}

	return config.CustomModel{
		DisplayName: f.inputs[FieldDisplayName].Value(),
		Model:       f.inputs[FieldModelID].Value(),
		BaseURL:     f.inputs[FieldBaseURL].Value(),
		APIKey:      f.inputs[FieldAPIKey].Value(),
		Provider:    provider,
		MaxTokens:   maxTokens,
	}
}

func (f *Form) Validate() (bool, string) {
	f.validationError = make(map[int]string)

	if f.inputs[FieldDisplayName].Value() == "" {
		f.validationError[FieldDisplayName] = "Required"
		return false, "Display Name is required"
	}

	maxTokensStr := f.inputs[FieldMaxTokens].Value()
	if maxTokensStr != "" {
		v, err := strconv.Atoi(maxTokensStr)
		if err != nil || v < 0 {
			f.validationError[FieldMaxTokens] = "Must be positive integer"
			return false, "Max Tokens must be a positive integer"
		}
	}

	return true, ""
}

func (f *Form) ToggleAPIKeyVisibility() {
	f.showAPIKey = !f.showAPIKey
	if f.showAPIKey {
		f.inputs[FieldAPIKey].EchoMode = textinput.EchoNormal
	} else {
		f.inputs[FieldAPIKey].EchoMode = textinput.EchoPassword
	}
}

func (f *Form) ClearValidationErrors() {
	f.validationError = make(map[int]string)
}

func (f *Form) Focus() {
	f.inputs[f.focusIndex].Focus()
}

func (f *Form) Blur() {
	for i := range f.inputs {
		f.inputs[i].Blur()
	}
}

func (f *Form) FocusNext() {
	f.inputs[f.focusIndex].Blur()
	f.focusIndex = (f.focusIndex + 1) % FieldCount
	if f.focusIndex != FieldProvider {
		f.inputs[f.focusIndex].Focus()
	}
}

func (f *Form) FocusPrev() {
	f.inputs[f.focusIndex].Blur()
	f.focusIndex = (f.focusIndex - 1 + FieldCount) % FieldCount
	if f.focusIndex != FieldProvider {
		f.inputs[f.focusIndex].Focus()
	}
}

func (f *Form) FocusIndex() int {
	return f.focusIndex
}

func (f *Form) SetFocusIndex(idx int) {
	if idx >= 0 && idx < FieldCount {
		f.inputs[f.focusIndex].Blur()
		f.focusIndex = idx
		if f.focusIndex != FieldProvider {
			f.inputs[f.focusIndex].Focus()
		}
	}
}

func (f *Form) NextProvider() {
	f.providerIndex = (f.providerIndex + 1) % len(config.Providers)
}

func (f *Form) PrevProvider() {
	f.providerIndex = (f.providerIndex - 1 + len(config.Providers)) % len(config.Providers)
}

func (f *Form) UpdateInput(msg textinput.Model) {
	if f.focusIndex >= 0 && f.focusIndex < len(f.inputs) && f.focusIndex != FieldProvider {
		f.inputs[f.focusIndex] = msg
	}
}

func (f *Form) CurrentInput() *textinput.Model {
	if f.focusIndex >= 0 && f.focusIndex < len(f.inputs) && f.focusIndex != FieldProvider {
		return &f.inputs[f.focusIndex]
	}
	return nil
}

func (f *Form) View(focused bool, modelName string) string {
	primaryColor := lipgloss.Color("39")
	secondaryColor := lipgloss.Color("240")
	errorColor := lipgloss.Color("196")
	dimmedColor := lipgloss.Color("242")

	titleBackgroundStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("0")).
		Background(primaryColor).
		Padding(0, 1)
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Bold(true)
	inputBorder := lipgloss.RoundedBorder()
	errorStyle := lipgloss.NewStyle().Foreground(errorColor).Bold(true)
	hintStyle := lipgloss.NewStyle().Foreground(dimmedColor).Italic(true)

	panelWidth := f.Width
	if panelWidth < 1 {
		panelWidth = 1
	}

	inputTextWidth := f.inputTextWidth()
	for i := range f.inputs {
		promptWidth := lipgloss.Width(f.inputs[i].Prompt)
		f.inputs[i].Width = max(1, inputTextWidth-promptWidth)
	}

	title := "NEW MODEL"
	if modelName != "" {
		title = "EDITING: " + modelName
	}
	titleContentWidth := max(0, panelWidth-2) // titleBackgroundStyle has horizontal padding=2
	titleLine := titleBackgroundStyle.Render(padRight(title, titleContentWidth))

	fieldLabels := []string{
		"Display Name:",
		"Model ID:",
		"Base URL:",
		"API Key:",
		"Provider:",
		"Max Tokens:",
	}

	fieldHints := []string{
		"Name displayed in the UI",
		"Identifier (e.g., gpt-4-turbo)",
		"API endpoint, usually ends in /v1",
		"Provider API key (Ctrl+V to toggle)",
		"← → to switch providers",
		"Maximum tokens per request",
	}

	inputStyleFor := func(active bool, err bool) lipgloss.Style {
		borderColor := secondaryColor
		if active {
			borderColor = primaryColor
		}
		if err {
			borderColor = errorColor
		}
		return lipgloss.NewStyle().
			Border(inputBorder).
			BorderForeground(borderColor).
			Padding(0, 1).
			Width(max(0, panelWidth-2))
	}

	blocks := make([][]string, 0, FieldCount)
	for i := 0; i < FieldCount; i++ {
		active := focused && f.focusIndex == i

		label := fieldLabels[i]
		if errMsg, hasErr := f.validationError[i]; hasErr && i != FieldProvider {
			label += " " + errorStyle.Render("! "+errMsg)
		}
		labelLine := ansi.Truncate(labelStyle.Render(label), panelWidth, "")

		var box string
		if i == FieldProvider {
			provider := ""
			if f.providerIndex >= 0 && f.providerIndex < len(config.Providers) {
				provider = config.Providers[f.providerIndex]
			}
			providerText := padRight("< "+provider+" >", inputTextWidth)
			box = inputStyleFor(active, false).Render(providerText)
		} else {
			box = inputStyleFor(active, false).Render(f.inputs[i].View())
		}

		block := []string{labelLine}
		block = append(block, strings.Split(box, "\n")...)
		if active {
			block = append(block, hintStyle.Render(ansi.Truncate("  "+fieldHints[i], panelWidth, "")))
		}
		blocks = append(blocks, block)
	}

	header := []string{titleLine}
	if f.Height >= 12 {
		header = append(header, "")
	}

	if f.Height <= len(header) {
		return strings.Join(header[:max(0, f.Height)], "\n")
	}

	available := f.Height - len(header)
	if available < 0 {
		available = 0
	}

	// Add extra spacing between fields only when we can show the whole form comfortably.
	fieldGap := 0
	if f.Height >= 31 {
		fieldGap = 1
	}

	start := 0
	end := 0
	for start <= f.focusIndex {
		used := 0
		end = start
		for end < len(blocks) {
			cost := len(blocks[end])
			if end > start {
				cost += fieldGap
			}
			if used+cost > available {
				break
			}
			used += cost
			end++
		}
		if f.focusIndex < end {
			break
		}
		start++
	}
	if start > f.focusIndex {
		start = f.focusIndex
	}

	var body []string
	if end > start {
		for i := start; i < end; i++ {
			if i > start && fieldGap > 0 {
				body = append(body, "")
			}
			body = append(body, blocks[i]...)
		}
	} else if available > 0 && f.focusIndex >= 0 && f.focusIndex < len(blocks) {
		// Fallback: show as much as possible of the focused block.
		block := blocks[f.focusIndex]
		if len(block) > available {
			block = block[:available]
		}
		body = append(body, block...)
	}

	return strings.Join(append(header, body...), "\n")
}
