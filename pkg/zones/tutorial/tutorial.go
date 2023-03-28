package tutorial

// TUTORIAL ZONE
//         COMBAT TUTORIAL (combattutorial.go)
//					^
//					|
// 	???	   <- FIRST ROOM (YOU ARE HERE) -> ?????
//				    |
//
// 	         SECOND ROOM ??????????????
import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/mikejk8s/gmud/pkg/backend"
	"github.com/mikejk8s/gmud/pkg/models"
	"github.com/mikejk8s/gmud/pkg/zones/combattutorial"
)

const useHighPerformanceRenderer = false

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

type model struct {
	SSHSession ssh.Session
	content    string
	ready      bool
	viewport   viewport.Model
	Character  *models.Character
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func InitialModel(char *models.Character, SSHSess ssh.Session) model {
	// it is succesful, but I still dont know how the fuck I can implement whois function.
	// Load some text for our viewport
	content, err := os.ReadFile("./textfiles/tutorial.md")
	if err != nil {
		panic(err)
	}
	// dial to websocket server
	wsUtil, err := backend.NewWebsocketUtil()
	if err != nil {
		log.Println(err)
		return model{}
	}
	// then send a message to the websocket server
	if _, err = wsUtil.Conn.Write([]byte("Welcome, time has no meaning here!")); err != nil {
		log.Fatal(err)
		return model{}
	}
	// read the response from the websocket server
	var msg = make([]byte, 512)
	var n int
	if n, err = wsUtil.Conn.Read(msg); err != nil {
		log.Fatal(err)
		return model{}
	}
	// print the response
	content = append(content, msg[:n]...)
	// TODO: do this for whois, let the users know who is in the room
	return model{
		SSHSession: SSHSess,
		content:    string(content),
		ready:      false,
		Character:  char,
	}
}
func (m model) Init() tea.Cmd {
	// m.GetCharacterDB()
	return nil
}
func (m model) headerView() string {
	title := titleStyle.Render("Tutorial Zone - Center Area")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			return combattutorial.InitialModel(m.Character, m.SSHSession), nil
		case "down", "j":
			return combattutorial.InitialModel(m.Character, m.SSHSession), nil
		}
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if m.ready == false {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			width := msg.Width
			wrapped := lipgloss.NewStyle().Width(width).Bold(true).Render(m.content)
			m.viewport.SetContent(wrapped)
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	// If user didnt resize the window, we render it by whatever the size is right now.
	default:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if m.ready == false {
			// Get terminal size
			ptySize, _, _ := m.SSHSession.Pty()
			m.viewport = viewport.New(ptySize.Window.Width, ptySize.Window.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			width := ptySize.Window.Width
			wrapped := lipgloss.NewStyle().Width(width).Render(m.content)
			m.viewport.SetContent(wrapped)
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		}

		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)

}
func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}
