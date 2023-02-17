# TUI Code Review

TUI Code Review is a terminal-based application that allows users to view a list of pull requests on GitHub for which
they have been marked for code review. The application fetches data from the GitHub API to display a list of pull
requests awaiting review and another list of pull requests that the user has already reviewed.

Users can configure the repositories that they want to display in the application.

The application's user interface is simple and easy to navigate. Users can view the pull requests in a tabular format
and sort them based on different criteria. The user can select a pull request to open it in the default browser.

TUI Code Review simplifies the code review process by providing an easy-to-use interface to manage pull requests. The
application is especially useful for users who need to review pull requests across multiple repositories.

### Technologies

- [bubbletea](https://pkg.go.dev/github.com/charmbracelet/bubbletea) is a framework for building tui applications in go.
- [bubbles](https://github.com/charmbracelet/bubbles) provides a set of components that go along with bubbletea.
- [lipgloss](https://github.com/charmbracelet/lipgloss) lets you easily style tui components.
- [termenv](https://github.com/muesli/termenv) allows to add colors and styles to text.
