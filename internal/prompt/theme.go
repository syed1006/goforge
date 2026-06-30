package prompt

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// forgeTheme builds the goforge-branded huh theme — warm amber accents on a
// neutral palette, adaptive for light and dark terminals.
func forgeTheme() *huh.Theme {
	t := huh.ThemeBase()

	var (
		fg     = lipgloss.AdaptiveColor{Light: "235", Dark: "252"}
		dim    = lipgloss.AdaptiveColor{Light: "240", Dark: "243"}
		amber  = lipgloss.AdaptiveColor{Light: "#B45309", Dark: "#F59E0B"} // primary accent
		ember  = lipgloss.AdaptiveColor{Light: "#9A3412", Dark: "#F97316"} // selectors/cursors
		ok     = lipgloss.AdaptiveColor{Light: "#047857", Dark: "#10B981"} // confirm/selected
		err    = lipgloss.AdaptiveColor{Light: "#B91C1C", Dark: "#F87171"}
		cream  = lipgloss.Color("#FFFBEB")
		border = lipgloss.AdaptiveColor{Light: "250", Dark: "238"}
	)

	t.FieldSeparator = lipgloss.NewStyle().SetString("\n")

	t.Focused.Base = t.Focused.Base.
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(amber).
		PaddingLeft(1)
	t.Focused.Card = t.Focused.Base
	t.Focused.Title = t.Focused.Title.Foreground(amber).Bold(true)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(amber).Bold(true).MarginBottom(1)
	t.Focused.Description = t.Focused.Description.Foreground(dim)
	t.Focused.Directory = t.Focused.Directory.Foreground(amber)

	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(err)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(err)

	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(ember).SetString("❯ ")
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(ember)
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(ember)

	t.Focused.Option = t.Focused.Option.Foreground(fg)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(ember).SetString("❯ ")
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(ok)
	t.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(ok).SetString("✓ ")
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(dim).SetString("• ")
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(fg)

	t.Focused.FocusedButton = t.Focused.FocusedButton.
		Foreground(cream).
		Background(amber).
		Bold(true)
	t.Focused.Next = t.Focused.FocusedButton
	t.Focused.BlurredButton = t.Focused.BlurredButton.
		Foreground(dim).
		Background(lipgloss.AdaptiveColor{Light: "254", Dark: "236"})

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(ember)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(dim)
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(amber).SetString("❯ ")
	t.Focused.TextInput.Text = t.Focused.TextInput.Text.Foreground(fg)

	// Blurred state inherits but loses the side bar so unfocused fields read as quiet.
	t.Blurred = t.Focused
	t.Blurred.Base = t.Focused.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Card = t.Blurred.Base
	t.Blurred.Title = t.Focused.Title.Foreground(dim).Bold(false)
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	t.Group.Title = t.Focused.Title
	t.Group.Description = t.Focused.Description.MarginBottom(1)

	t.Help.Ellipsis = t.Help.Ellipsis.Foreground(dim)
	t.Help.ShortKey = t.Help.ShortKey.Foreground(amber)
	t.Help.ShortDesc = t.Help.ShortDesc.Foreground(dim)
	t.Help.ShortSeparator = t.Help.ShortSeparator.Foreground(border)
	t.Help.FullKey = t.Help.FullKey.Foreground(amber)
	t.Help.FullDesc = t.Help.FullDesc.Foreground(dim)
	t.Help.FullSeparator = t.Help.FullSeparator.Foreground(border)

	return t
}
