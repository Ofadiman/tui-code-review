package main

import "github.com/charmbracelet/lipgloss"

// Colors: https://htmlcolorcodes.com/color-names/

var ColorDeepPink = lipgloss.Color("#FF1493")

var StyledMain = lipgloss.NewStyle().PaddingTop(1).PaddingBottom(1).PaddingLeft(2).PaddingRight(2)
var StyledHeader = lipgloss.NewStyle().AlignHorizontal(lipgloss.Center).Bold(true).Border(lipgloss.RoundedBorder()).BorderForeground(ColorDeepPink)
