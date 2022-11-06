package main

import(
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"strconv"
	"log"
)

func usage_exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

var (
	baseDir string = "/Users/amandaliem/git/peeonabot/"
	tokenFile string = baseDir + "credentials/discord.token"
	logFile string = baseDir + "log.txt"
)

// Reads file discord.token and returns the discord bot token
func getToken() string {
	contents, err := os.ReadFile(tokenFile)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Trim(string(contents), "\n")
}

// Sets up logging file
func logInit() {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
}

// Sets up PID file
func pidInit() {
	// clear out the pid directory
	os.RemoveAll(baseDir + "pids/")
	err := os.Mkdir(baseDir + "pids", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	_, err = os.Create(baseDir + "pids/" + strconv.Itoa(os.Getpid()))
	if err != nil {
		log.Fatal(err)
	}
}

// Performs Cleanup
func cleanup() {
	// clear out the pid directory
	os.RemoveAll(baseDir + "pids/")
}

func main(){
	// Set up logging
	logInit()

	// Set up Pidfile
	pidInit()

	// Create bot
	bot, err := discordgo.New("Bot " + getToken())
	if err != nil {
		log.Fatal(err)
	}

	// Output Ready Status
	bot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)
	})
	bot.AddHandler(messageCreate)
	bot.Identify.Intents = discordgo.IntentsGuildMessages

	// Begin Listening
	err = bot.Open()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Peeona Bot is now listening %v", os.Getpid())

	// Wait here
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Kill the bot
	log.Println("Peeona Bot is now sleeping")
	cleanup()
	bot.Close()
}

// Monitor messages sent in the server
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Println("Peeona bot got a message: " + m.Content)

	// Ignore messages from bots
	if m.Author.Bot {
		return
	}

	if strings.HasPrefix(m.Content, "/pp") {
		s.ChannelMessageSend(m.ChannelID, "woof woof")
		return
	}
}

