package main

import "github.com/charmbracelet/lipgloss"

// Colors: https://htmlcolorcodes.com/color-names/

var ColorDeepPink = lipgloss.Color("#FF1493")
var ColorOrangeRed = lipgloss.Color("#FF4500")
var ColorLimeGreen = lipgloss.Color("#32CD32")
var ColorDeepSkyBlue = lipgloss.Color("#00BFFF")
var ColorGrey = lipgloss.Color("#808080")
var ColorWhite = lipgloss.Color("#FFFFFF")
var ColorGold = lipgloss.Color("#FFD700")

var StyledMain = lipgloss.NewStyle().PaddingTop(1).PaddingBottom(1).PaddingLeft(2).PaddingRight(2)
var StyledHeader = lipgloss.NewStyle().AlignHorizontal(lipgloss.Center).Bold(true).Border(lipgloss.RoundedBorder()).BorderForeground(ColorDeepPink)

var StyledChangesRequested = lipgloss.NewStyle().Foreground(ColorOrangeRed)
var StyledApproved = lipgloss.NewStyle().Foreground(ColorLimeGreen)
var StyledAwaiting = lipgloss.NewStyle().Foreground(ColorDeepSkyBlue)
var StyledDraft = lipgloss.NewStyle().Foreground(ColorGrey)
var StyledCommented = lipgloss.NewStyle().Foreground(ColorGold)

var StyledHelpShortcut = lipgloss.NewStyle().Foreground(ColorWhite)
var StyledHelpDescription = lipgloss.NewStyle().Foreground(ColorGrey)

var StyledUnderline = lipgloss.NewStyle().Underline(true)
