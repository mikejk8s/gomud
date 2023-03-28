package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/muesli/termenv"

	"github.com/mikejk8s/gmud/pkg/backend"
	"github.com/mikejk8s/gmud/pkg/menus"
	"github.com/mikejk8s/gmud/pkg/models"
	"github.com/mikejk8s/gmud/pkg/postgrespkg"
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

	conn, err := postgrespkg.Connect(postgrespkg.POSTGRES_USERS_DB)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := conn.SqlDB.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	if err := conn.ConnectSQLToSchema("users"); err != nil {
		log.Fatalln(err)
	}

	user := models.User{}
	query := fmt.Sprintf("SELECT name, password_hash, remember_hash WHERE name = '%s'", ctx.User())
	rows, err := conn.DB.Raw(query).Rows()
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&user.Name, &user.PasswordHash, &user.RememberHash)
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
		log.Println("Running on docker is set to", RunningOnDocker)
	} else {
		//  set default
		RunningOnDocker = true
	}

	if RunningOnDocker {
		postgrespkg.POSTGRES_USER = os.Getenv("POSTGRES_USER")
		postgrespkg.POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
		postgrespkg.POSTGRESS_HOST = os.Getenv("POSTGRES_HOST")
		if postgrespkg.POSTGRES_USER == "" || postgrespkg.POSTGRES_PASSWORD == "" {
			postgrespkg.POSTGRES_USER = "gmud"
			postgrespkg.POSTGRES_PASSWORD = "gmud"
		}
	} else {
		postgrespkg.POSTGRES_USER = "gmud"
		postgrespkg.POSTGRES_PASSWORD = "gmud"
		postgrespkg.POSTGRESS_HOST = "127.0.0.1:5432"
	}
	go backend.StartWSServer()
	go backend.StartWebPageBackend(RunningOnDocker, 6969)
	initDB, err := postgrespkg.Connect("")
	if err != nil {
		log.Fatalln(err)
	}
	initDB.CreateDatabases()
	characterDB, err := postgrespkg.Connect(postgrespkg.POSTGRES_CHARACTERS_DB)
	if err != nil {
		panic(err)
	}
	usersDB, err := postgrespkg.Connect(postgrespkg.POSTGRES_USERS_DB)
	if err != nil {
		log.Println(err)
	}

	// Characters table creation.
	go func() {
		err := characterDB.CreateCharacterTable()
		if err != nil {
			panic(err)
		}
		err = characterDB.MigrateCharacters()
		if err != nil {
			panic(err)
		}
	}()

	// Users table creation and migration.
	go func() {
		err := usersDB.CreateUsersTable()
		if err != nil {
			log.Fatalln(err)
		}

		err = usersDB.MigrateUsers()
		if err != nil {
			log.Fatalln(err)
		}
	}()
	// defer close whole db's
	defer func() {
		characterDB.SqlDB.Close()
		usersDB.SqlDB.Close()
	}()

	// Define the default shell
	const DefaultShell = "bash"

	// Initialize the SSH server
	s, err := wish.NewServer(
		wish.WithIdleTimeout(30*time.Minute), // 30-minute idle timer, in case someone forgets to log out.
		wish.WithPasswordAuth(passHandler),
		wish.WithAddress(fmt.Sprintf("%s:%d", Host, Port)),
		wish.WithMiddleware(
			logging.Middleware(),
			loginBubbleteaMiddleware(),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}

	/* s.Handle(DefaultShell, func(sess wish.Session) {
		user, err := getUserFromDB(sess.Context(), usersDB)
		if err != nil {
			log.Fatalln(err)
		}
		// Do something with the user struct here.
	}) */

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
	usersConn, err := postgrespkg.Connect(postgrespkg.POSTGRES_USERS_DB)
	if err != nil {
		return nil, err
	}
	defer usersConn.SqlDB.Close()

	user := &models.User{}
	query := fmt.Sprintf("SELECT created_at, updated_at, name from users where name = '%s'", ctx.User())
	rows, err := usersConn.SqlDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		rows.Close()
	}()

	for rows.Next() {
		err = rows.Scan(&user.CreatedAt, &user.UpdatedAt, &user.Name)
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

	bmHandler := bubbletea.MiddlewareWithProgramHandler(sshHandler, termenv.ANSI256)
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
			return menus.InitialModel(m.accOwner, m.SSHSession), nil // Go to the login page with passing account owner
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
