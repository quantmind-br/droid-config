package components

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/diogo/droid-config/internal/config"
)

// Constants for list layout
const (
	listHeaderLines = 4 // Title bar + action hints + separator + section title
	listFooterLines = 0
	scrollPadding   = 1 // Keep cursor this many items from edge during scrolling
	scrollbarWidth  = 2 // Width reserved for scrollbar (1 char + 1 space)
)

// Scrollbar and edge indicator styling (local to avoid import cycles)
var (
	scrollbarDimColor       = lipgloss.Color("242")
	scrollbarPrimaryColor   = lipgloss.Color("39")
	scrollbarSecondaryColor = lipgloss.Color("240")

	scrollbarTrackChar  = "░"
	scrollbarThumbChar  = "█"
	scrollbarTrackStyle = lipgloss.NewStyle().Foreground(scrollbarDimColor)
	scrollbarThumbStyle = lipgloss.NewStyle().Foreground(scrollbarPrimaryColor)

	edgeIndicatorUp    = "▲"
	edgeIndicatorDown  = "▼"
	edgeIndicatorStyle = lipgloss.NewStyle().Foreground(scrollbarSecondaryColor)
)

type ListItem struct {
	Model    config.CustomModel
	Selected bool
}

type List struct {
	Items  []ListItem
	Cursor int
	Height int
	Width  int
	offset int
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

func NewList() *List {
	return &List{
		Items:  []ListItem{},
		Cursor: 0,
		Height: 10,
		Width:  25,
	}
}

// getVisibleHeight returns the number of items that can be displayed
func (l *List) getVisibleHeight() int {
	visible := l.Height - l.headerLinesToShow() - listFooterLines
	if len(l.Items) > 0 && visible < 1 && l.Height > 1 {
		visible = 1
	}
	if visible < 0 {
		visible = 0
	}
	return visible
}

func (l *List) headerLinesToShow() int {
	if l.Height <= 0 {
		return 0
	}
	// Never exceed available height. Leave at least one line for list content when possible.
	header := listHeaderLines
	if header > l.Height {
		header = l.Height
	}
	if l.Height > 1 && header == l.Height {
		header = l.Height - 1
	}
	if header < 1 {
		header = 1
	}
	return header
}

func (l *List) SetItems(models []config.CustomModel) {
	l.Items = make([]ListItem, len(models))
	for i, m := range models {
		l.Items[i] = ListItem{Model: m, Selected: false}
	}
	if l.Cursor >= len(l.Items) {
		l.Cursor = max(0, len(l.Items)-1)
	}
}

func (l *List) GetModels() []config.CustomModel {
	models := make([]config.CustomModel, len(l.Items))
	for i, item := range l.Items {
		models[i] = item.Model
	}
	return models
}

func (l *List) MoveUp() {
	if l.Cursor > 0 {
		l.Cursor--
		l.updateOffset()
	}
}

func (l *List) MoveDown() {
	if l.Cursor < len(l.Items)-1 {
		l.Cursor++
		l.updateOffset()
	}
}

func (l *List) updateOffset() {
	visibleHeight := l.getVisibleHeight()

	// Padded scrolling: keep cursor away from edges when possible
	padding := scrollPadding
	if visibleHeight <= 3 {
		padding = 0 // Disable padding for very small lists
	}

	// Scroll up: maintain padding from top
	if l.Cursor < l.offset+padding {
		l.offset = max(0, l.Cursor-padding)
	}

	// Scroll down: maintain padding from bottom
	if l.Cursor >= l.offset+visibleHeight-padding {
		l.offset = l.Cursor - visibleHeight + padding + 1
	}

	// Clamp offset to valid range
	maxOffset := max(0, len(l.Items)-visibleHeight)
	if l.offset > maxOffset {
		l.offset = maxOffset
	}
	if l.offset < 0 {
		l.offset = 0
	}
}

func (l *List) ToggleSelected() {
	if l.Cursor < len(l.Items) {
		l.Items[l.Cursor].Selected = !l.Items[l.Cursor].Selected
	}
}

func (l *List) SelectAll() {
	allSelected := l.AllSelected()
	for i := range l.Items {
		l.Items[i].Selected = !allSelected
	}
}

func (l *List) AllSelected() bool {
	if len(l.Items) == 0 {
		return false
	}
	for _, item := range l.Items {
		if !item.Selected {
			return false
		}
	}
	return true
}

func (l *List) GetSelectedIndices() []int {
	var indices []int
	for i, item := range l.Items {
		if item.Selected {
			indices = append(indices, i)
		}
	}
	return indices
}

func (l *List) ClearSelections() {
	for i := range l.Items {
		l.Items[i].Selected = false
	}
}

func (l *List) CurrentModel() *config.CustomModel {
	if l.Cursor >= 0 && l.Cursor < len(l.Items) {
		return &l.Items[l.Cursor].Model
	}
	return nil
}

func (l *List) UpdateCurrentModel(m config.CustomModel) {
	if l.Cursor >= 0 && l.Cursor < len(l.Items) {
		l.Items[l.Cursor].Model = m
	}
}

func (l *List) AddModel(m config.CustomModel) {
	l.Items = append(l.Items, ListItem{Model: m, Selected: false})
	l.Cursor = len(l.Items) - 1
	l.updateOffset()
}

func (l *List) DeleteSelected() int {
	indices := l.GetSelectedIndices()
	if len(indices) == 0 {
		return 0
	}

	newItems := make([]ListItem, 0, len(l.Items)-len(indices))
	selectedMap := make(map[int]bool)
	for _, idx := range indices {
		selectedMap[idx] = true
	}

	for i, item := range l.Items {
		if !selectedMap[i] {
			newItems = append(newItems, item)
		}
	}

	l.Items = newItems
	if l.Cursor >= len(l.Items) {
		l.Cursor = max(0, len(l.Items)-1)
	}
	l.updateOffset()

	return len(indices)
}

func (l *List) DeleteCurrent() bool {
	if len(l.Items) == 0 || l.Cursor < 0 || l.Cursor >= len(l.Items) {
		return false
	}

	l.Items = append(l.Items[:l.Cursor], l.Items[l.Cursor+1:]...)
	if l.Cursor >= len(l.Items) {
		l.Cursor = max(0, len(l.Items)-1)
	}
	l.updateOffset()
	return true
}

func (l *List) MoveItemUp() bool {
	if l.Cursor <= 0 || l.Cursor >= len(l.Items) {
		return false
	}
	l.Items[l.Cursor], l.Items[l.Cursor-1] = l.Items[l.Cursor-1], l.Items[l.Cursor]
	l.Cursor--
	l.updateOffset()
	return true
}

func (l *List) MoveItemDown() bool {
	if l.Cursor < 0 || l.Cursor >= len(l.Items)-1 {
		return false
	}
	l.Items[l.Cursor], l.Items[l.Cursor+1] = l.Items[l.Cursor+1], l.Items[l.Cursor]
	l.Cursor++
	l.updateOffset()
	return true
}

func (l *List) View(focused bool, dirty bool) string {
	var itemStyle, selectedItemStyle lipgloss.Style
	primaryColor := lipgloss.Color("39")
	secondaryColor := lipgloss.Color("240")
	successColor := lipgloss.Color("82")

	titleBackgroundStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("0")).
		Background(primaryColor).
		Padding(0, 1)
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(primaryColor)
	itemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	selectedItemStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("0")).Background(primaryColor)

	checkboxChecked := lipgloss.NewStyle().Foreground(successColor).Render("[x]")
	checkboxUnchecked := lipgloss.NewStyle().Foreground(secondaryColor).Render("[ ]")

	providerBadges := map[string]string{
		"anthropic":                   lipgloss.NewStyle().Foreground(lipgloss.Color("141")).Render("[A]"),
		"openai":                      lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Render("[O]"),
		"generic-chat-completion-api": lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Render("[G]"),
	}

	headerWidth := l.Width
	if headerWidth < 1 {
		headerWidth = 1
	}

	allSelectText := "[A] Select All"
	if l.AllSelected() && len(l.Items) > 0 {
		allSelectText = "[A] Deselect All"
	}
	var lines []string
	headerLines := l.headerLinesToShow()
	if headerLines >= 1 {
		lines = append(lines, titleStyle.Render(padOrTruncate("[N] New  [D] Delete", headerWidth)))
	}
	if headerLines >= 2 {
		lines = append(lines, lipgloss.NewStyle().Foreground(secondaryColor).Render(padOrTruncate(allSelectText, headerWidth)))
	}
	if headerLines >= 3 {
		lines = append(lines, lipgloss.NewStyle().Foreground(secondaryColor).Render(strings.Repeat("─", headerWidth)))
	}

	titleText := "YOUR MODELS"
	if dirty {
		titleText += " *"
	}
	if headerLines >= 4 {
		titleContentWidth := max(0, headerWidth-2) // titleBackgroundStyle has horizontal padding=2
		lines = append(lines, titleBackgroundStyle.Render(padOrTruncate(titleText, titleContentWidth)))
	}

	visibleHeight := l.getVisibleHeight()
	if visibleHeight == 0 {
		return strings.Join(lines, "\n")
	}

	if len(l.Items) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(secondaryColor).Italic(true)
		emptyLines := []string{
			"  No models configured",
			"  Press 'N' to create your first model",
		}
		for i := 0; i < visibleHeight && i < len(emptyLines); i++ {
			lines = append(lines, emptyStyle.Render(padOrTruncate(emptyLines[i], headerWidth)))
		}
	} else {
		end := l.offset + visibleHeight
		if end > len(l.Items) {
			end = len(l.Items)
		}

		// Check if we need edge indicators
		hasItemsAbove := l.offset > 0
		hasItemsBelow := end < len(l.Items)

		// Generate scrollbar if list is scrollable
		scrollbar := l.renderScrollbar(visibleHeight)

		// Build list items with scrollbar
		contentWidth := headerWidth - scrollbarWidth
		if contentWidth < 10 {
			contentWidth = 10
		}

		var listLines []string
		for i := l.offset; i < end; i++ {
			item := l.Items[i]
			checkbox := checkboxUnchecked
			if item.Selected {
				checkbox = checkboxChecked
			}

			badge := providerBadges[item.Model.Provider]
			if badge == "" {
				badge = lipgloss.NewStyle().Foreground(secondaryColor).Render("[?]")
			}

			name := item.Model.DisplayName
			if name == "" {
				name = item.Model.Model
			}
			if name == "" {
				name = "(unnamed)"
			}

			// Truncate to keep each rendered line within the panel width.
			digits := len(strconv.Itoa(i + 1))
			maxNameLen := contentWidth - (digits + 10)
			if maxNameLen < 5 {
				maxNameLen = 5
			}
			if len(name) > maxNameLen {
				name = name[:maxNameLen-3] + "..."
			}

			rawLine := fmt.Sprintf("%s %s %d. %s", checkbox, badge, i+1, name)
			rawLine = padOrTruncate(rawLine, contentWidth)

			if i == l.Cursor && focused {
				listLines = append(listLines, selectedItemStyle.Render(rawLine))
			} else if i == l.Cursor {
				listLines = append(listLines, itemStyle.Bold(true).Render(rawLine))
			} else {
				listLines = append(listLines, itemStyle.Render(rawLine))
			}
		}

		// Combine list lines with scrollbar
		for idx, line := range listLines {
			scrollbarChar := " "
			if len(scrollbar) > idx {
				scrollbarChar = scrollbar[idx]
			}

			// Add edge indicator on first/last visible line
			if idx == 0 && hasItemsAbove {
				scrollbarChar = edgeIndicatorStyle.Render(edgeIndicatorUp)
			} else if idx == len(listLines)-1 && hasItemsBelow {
				scrollbarChar = edgeIndicatorStyle.Render(edgeIndicatorDown)
			}

			lines = append(lines, line+scrollbarChar)
		}
	}

	return strings.Join(lines, "\n")
}

// renderScrollbar generates a vertical scrollbar as a slice of characters
func (l *List) renderScrollbar(visibleHeight int) []string {
	totalItems := len(l.Items)
	if totalItems <= visibleHeight {
		// No scrollbar needed - return empty strings
		result := make([]string, visibleHeight)
		for i := range result {
			result[i] = " "
		}
		return result
	}

	result := make([]string, visibleHeight)

	// Calculate thumb size and position
	thumbSize := max(1, visibleHeight*visibleHeight/totalItems)
	if thumbSize > visibleHeight {
		thumbSize = visibleHeight
	}

	// Calculate thumb position based on offset
	maxOffset := totalItems - visibleHeight
	thumbPos := 0
	if maxOffset > 0 {
		thumbPos = l.offset * (visibleHeight - thumbSize) / maxOffset
	}

	// Build scrollbar
	for i := 0; i < visibleHeight; i++ {
		if i >= thumbPos && i < thumbPos+thumbSize {
			result[i] = scrollbarThumbStyle.Render(scrollbarThumbChar)
		} else {
			result[i] = scrollbarTrackStyle.Render(scrollbarTrackChar)
		}
	}

	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
