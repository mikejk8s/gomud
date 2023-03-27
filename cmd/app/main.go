package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/balibuild/tunnelssh"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/gliderlabs/ssh"
	"github.com/muesli/termenv"
	"gorm.io/gorm"

	"github.com/mikejk8s/gmud/pkg/backend"
	mn "github.com/mikejk8s/gmud/pkg/menus"
	"github.com/mikejk8s/gmud/pkg/models"
	postgrespkg "github.com/mikejk8s/gmud/pkg/postgrespkg"
)

const (
	Host = "127.0.0.1"
	Port = 2222
)

var RunningOnDocker = false

func passHandler(ctx ssh.Context, password string) bool {
	if ctx.User() == "" {
		return false
	}

	conn, err := postgrespkg.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close(*gorm.DB)

	if err := conn.GetSQLConn("users"); err != nil {
		log.Fatalln(err)
	}

	user := models.User{}
	query := fmt.Sprintf("SELECT password, email, name, username from users where username = '%s'", ctx.User())
	rows, err := conn.DB.Raw(query).Rows()
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&user.Password, &user.Email, &user.Name, &user.Username)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if err = rows.Err(); err != nil {
		log.Fatalln(err)
	}

	return user.CheckPassword(password) == nil
}

func main() {
	tempEnvVar, ok := os.LookupEnv("RUNNING_ON_DOCKER")
	if ok {
		RunningOnDocker = tempEnvVar == "true"
		RunningOnDocker = RunningOnDocker
		log.Println("Running on docker is set to", RunningOnDocker)
	} else {
		//  set default
		RunningOnDocker = true
	}

	if RunningOnDocker {
		postgrespkg.POSTGRES_USER = os.Getenv("POSTGRES_USER")
		postgrespkg.POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
		postgrespkg.POSTGRESS_HOST = os.Getenv("POSTGRES_HOST")
	} else {
		postgrespkg.POSTGRES_USER = "gmud"
		postgrespkg.POSTGRES_PASSWORD = "gmud"
		postgrespkg.POSTGRESS_HOST = "127.0.0.1:5432"
	}
	go backend.StartWSServer()
	go backend.StartWebPageBackend(RunningOnDocker, 6969)

	initialTableCreation, err := postgrespkg.Connect()
	if err != nil {
		log.Println(err)
	}
	defer initialTableCreation.Close()

	characterDB, err := postgrespkg.Connect()
	if err != nil {
		log.Println(err)
	}
	defer characterDB.Close()

	usersDB, err := postgrespkg.Connect()
	if err != nil {
		log.Println(err)
	}
	defer usersDB.Close()

	// Characters table creation.
	go func() {
		err := postgrespkg.CreateCharacterTable(characterDB)
		if err != nil {
			panic(err)
		}
	}()

	// Users table creation and migration.
	go func() {
		err := postgrespkg.CreateUsersTable(usersDB)
		if err != nil {
			log.Fatalln(err)
		}

		err = postgrespkg.UsersMigration(usersDB)
		if err != nil {
			log.Fatalln(err)
		}

		defer usersDB.DB.Close()
	}()


	// Users table creation and migration.
	go func() {
		err := postgrespkg.CreateUsersTable(usersDB)
		if err != nil {
			log.Fatalln(err)
		}

		err = postgrespkg.UsersMigration(usersDB)
		if err != nil {
			log.Fatalln(err)
		}

		defer usersDB.Close()
	}()




// Define the default shell
const DefaultShell = "bash"

// Initialize the SSH server
s, err := wish.NewServer(
    wish.WithIdleTimeout(30*time.Minute), // 30-minute idle timer, in case someone forgets to log out.
    wish.WithPasswordAuth(passHandler),
    wish.WithAddress(fmt.Sprintf("%s:%d", Host, Port)),
    wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
        return true
    }),
    wish.WithMiddleware(
        lm.Middleware(),
        loginBubbleteaMiddleware(),
    ),
)

if err != nil {
    log.Fatalln(err)
}

s.Handle(DefaultShell, func(sess wish.Session) {
    user, err := getUserFromDB(sess.Context(), usersDB)
    if err != nil {
        log.Fatalln(err)
    }
    // Do something with the user struct here.
})

	log.Printf("Starting SSH server on %s:%d...\n", Host, Port)
	defer s.Shutdown(context.Background())
	if err := s.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}

	done := make(chan os.Signal, 1)
	<-done
	log.Println("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}






func getUserFromDB(ctx ssh.Context) (*models.User, error) {
    usersConn, err := postgrespkg.Connect()
    if err != nil {
        return nil, err
    }
    defer usersConn.DB.Close()

    user := &models.User{}
    query := fmt.Sprintf("SELECT password, email, name, username from users where username = '%s'", ctx.User())
    rows, err := usersConn.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer func() {
        rows.Close()
        postgrespkg.Close()
    }()

    for rows.Next() {
        err = rows.Scan(&user.Password, &user.Email, &user.Name, &user.Username)
        if err != nil {
            return nil, err
        }
    }

    return user, nil
}








func loginBubbleteaMiddleware() wish.Middleware {
	login := func(m tea.Model, opts ...tea.ProgramOption) *tea.Program {
		p := tea.NewProgram(m, opts...)
		go func() {
			for {
				<-time.After(1 * time.Second)
				p.Send(timeMsg(time.Now()))
			}
		}()
		return p
	}

	sshHandler := func(s ssh.Session) *tea.Program {
		pty, _, _ := s.Pty()
		m := model{
			SSHSession: s,
			Width:      pty.Window.Width,
			Height:     pty.Window.Height,
			time:       time.Now(),
			accOwner:   s.User(),
		}
		return login(m, tea.WithInput(s), tea.WithOutput(s), tea.WithAltScreen(), tea.WithMouseCellMotion())
	}

	bmHandler := bm.MiddlewareWithProgramHandler(sshHandler, termenv.ANSI256)
	return bmHandler
}

type model struct {
	SSHSession ssh.Session
	time       time.Time
	Height     int
	Width      int
	accOwner   string // Account owner will be used for matching the characters created from this account.
}

type timeMsg time.Time

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timeMsg:
		m.time = time.Time(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "l", "ctrl+l":
			return mn.InitialModel(m.accOwner, m.SSHSession), nil // Go to the login page with passing account owner
		case "n", "ctrl+n":
			//mn.NewAccount()
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "Welcome to gmud!\n\n"
	s += "Date> " + m.time.Format(time.RFC1123) + "\n\n"
	s += "Press 'l' to go in.\n"
	s += m.SSHSession.LocalAddr().String() + "\n"
	s += m.SSHSession.RemoteAddr().String() + "\n"
	return fmt.Sprintln(s, m.Height, m.Width)
}
