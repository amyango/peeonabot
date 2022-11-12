package main

import(
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"log"
)

func usage_exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

var (
	baseDir string = "./"
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

func main(){

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

	if strings.HasPrefix(m.Content, "/movie") {
		doMovies(s, m)
		return
	}

}

func doMovies(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Split(m.Content, " ")
	movies := Find_movies(args[1:])
	if (len(movies) <= 0) {
		s.ChannelMessageSend(m.ChannelID, "No results found")
		return
	}
	movie := movies[0]
	response := fmt.Sprintf("[%s] %s (%v)", movie.Release_date, movie.Original_title, movie.Vote_average)
	s.ChannelMessageSend(m.ChannelID, response)
	if movie.Poster_path != "" {
		s.ChannelMessageSend(m.ChannelID, "https://image.tmdb.org/t/p/original/" + movie.Poster_path)
	}
}
