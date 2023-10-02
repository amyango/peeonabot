package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"discordgo"

	"github.com/astralservices/go-dalle"
	"github.com/ayush6624/go-chatgpt"
)

func usage_exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

var (
	baseDir          string = "./"
	tokenFile        string = baseDir + "credentials/discord.token"
	openAPITokenFile string = baseDir + "credentials/openai.token"
	logFile          string = baseDir + "log.txt"
	openAPIToken     string
	dalbby           dalle.Client
	gptbby           *chatgpt.Client
	globalSession    *discordgo.Session
)

// Reads file discord.token and returns the discord bot token
func getToken() string {
	contents, err := os.ReadFile(tokenFile)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Trim(string(contents), "\n")
}

func getOpenAIToken() string {
	contents, err := os.ReadFile(openAPITokenFile)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Trim(string(contents), "\n")
}

func main() {
	var err error

	openAPIToken = getOpenAIToken()
	dalbby = dalle.NewClient(openAPIToken)
	gptbby, err = chatgpt.NewClient(openAPIToken)

	// Create bot
	bot, err := discordgo.New("Bot " + getToken())
	if err != nil {
		log.Fatal(err)
	}

	// Output Ready Status
	bot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)

		globalSession = s

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
	//log.Println("Peeona bot got a message: " + m.Content)
	//log.Println("The channelID is: " + m.ChannelID)

	// Ignore messages from bots
	if m.Author.Bot {
		return
	}

	//if m.ChannelID != "1066254778742616084" {
	//	return
	//}

	if strings.HasPrefix(m.Content, "/pp") {
		s.ChannelMessageSend(m.ChannelID, "woof woof")
		return
	}

	if strings.HasPrefix(m.Content, "/movie") {
		doMovies(s, m)
		return
	}

	if strings.HasPrefix(m.Content, "/dalle") {
		doDalle(s, m)
		return
	}

	if strings.HasPrefix(m.Content, "/ask") {
		doGPT(s, m)
	}

	if strings.HasPrefix(m.Content, "$n") {
		doNumbers(s, m)
		return
	}

	if strings.HasPrefix(m.Content, "/t") {
		changeTime(s, m)
		return
	}

	if strings.HasPrefix(m.Content, "/squat") {
		doSquats(s, m)
		return
	}

	if strings.HasPrefix(m.Content, "/squad") {
		printSquats(s, m)
		return
	}
}

func printSquats(s *discordgo.Session, m *discordgo.MessageCreate) {
	g := GetGuildByID(s, m)

	embed := discordgo.MessageEmbed{
		Title:       "Big Booty Squad :peach:",
		Description: "Add squats with the /squat command!\nIf you do even one squat, you're added to the BBC role for the day :)",
	}

	for _, squatter := range g.Squatters {
		msg := fmt.Sprintf("Today Squats: %v\nLifetime Squats: %v", squatter.TodaySquats, squatter.TotalSquats)
		name := squatter.UserName
		if squatter.MaxWeight > 0 {
			msg = fmt.Sprintf("Today Squats: %v/%vlbs\nLifetime Squats: %v/%vlbs",
				squatter.TodaySquats, squatter.TodayWeight, squatter.TotalSquats, squatter.TotalWeight)

			for i := 0; i < squatter.TotalSquats/1000; i++ {
				name = fmt.Sprintf("%v :peach:", name)
			}
			name = fmt.Sprintf("%v [%v lbs]", name, squatter.MaxWeight)
		} else {

		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   name,
			Value:  msg,
			Inline: true,
		})
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &embed)
}

func doSquats(s *discordgo.Session, m *discordgo.MessageCreate) {

	w := 0

	resp := strings.Split(m.Content, " ")
	if len(resp) < 2 {
		s.ChannelMessageSend(m.ChannelID, "do /squat <# of squats> <weight (optional)>")
		return
	}

	n, err := strconv.Atoi(resp[1])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "do /squat <# of squats> <weight (optional)>")
		return
	}

	if len(resp) > 2 {
		weight, err := strconv.Atoi(resp[2])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "do /squat <# of squats> <weight (optional)>")
			return
		}
		w = weight
	}

	g := GetGuildByID(s, m)
	g.Squat(*m.Author, n, w, s)
	s.ChannelMessageSend(m.ChannelID, "nice")
}

func doGPT(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Split(m.Content, " ")
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "/dalle <stuff>")
		return
	}

	desc := args[1:]
	resp, err := askPeeona(chatgpt.GPT35Turbo, strings.Join(desc, " "))

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Donowalled by ChatGPT")
		fmt.Println(err)
	} else {
		s.ChannelMessageSend(m.ChannelID, resp)
	}
}

func askPeeona(model chatgpt.ChatGPTModel, question string) (string, error) {
	req := chatgpt.ChatCompletionRequest{
		Model: model,
		Messages: []chatgpt.ChatMessage{
			{
				Role:    chatgpt.ChatGPTModelRoleSystem,
				Content: question,
			},
		},
	}

	ctx := context.Background()
	res, err := gptbby.Send(ctx, &req)
	resp := ""

	if err != nil {
		return "", err
	} else {
		for _, choice := range res.Choices {
			resp = fmt.Sprintf("%v %v", resp, choice.Message.Content)
		}
	}

	return resp, nil
}

func doDalle(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Split(m.Content, " ")
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "/dalle <stuff>")
		return
	}

	desc := args[1:]
	msg := fmt.Sprintf("Thinking about %v...", strings.Join(desc, " "))
	s.ChannelMessageSend(m.ChannelID, msg)
	data, err := dalbby.Generate(strings.Join(desc, " "), nil, nil, nil, nil)

	if err != nil {
		fmt.Println(err)
	} else {
		s.ChannelMessageSend(m.ChannelID, data[0].URL)
	}
}

var num string = ""
var ms int = 750

func changeTime(s *discordgo.Session, m *discordgo.MessageCreate) {

	args := strings.Split(m.Content, " ")

	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "/t <time in millisecs>")
		return
	}

	i, err := strconv.Atoi(args[1])

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "/t <time in millisecs>")
		return
	}

	ms = i
	s.ChannelMessageSend(m.ChannelID, "$n time is now "+args[1])
}

func doNumbers(s *discordgo.Session, m *discordgo.MessageCreate) {

	if num == "" {

		numstr := ""
		ch := m.ChannelID

		for i := 0; i < 8; i++ {
			numstr = numstr + strconv.Itoa(rand.Intn(10))
		}

		msg, _ := s.ChannelMessageSend(m.ChannelID, numstr)
		num = numstr

		time.Sleep(time.Duration(ms) * time.Millisecond)

		r := "XXXXXXXX"
		s.ChannelMessageEdit(ch, msg.ID, r)
	} else {
		// compare
		args := strings.Split(m.Content, " ")

		if len(args) < 2 {
			s.ChannelMessageSend(m.ChannelID, "please type the number")
			return
		}

		if num == args[1] {
			s.ChannelMessageSend(m.ChannelID, "CORRECT!!!")
		} else {
			s.ChannelMessageSend(m.ChannelID, "Wrong it was "+num)
		}

		num = ""
	}
}

func doMovies(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Split(m.Content, " ")
	movies := Find_movies(args[1:])
	if len(movies) <= 0 {
		s.ChannelMessageSend(m.ChannelID, "No results found")
		return
	}
	movie := movies[0]
	response := fmt.Sprintf("[%s] %s (%v)", movie.Release_date, movie.Original_title, movie.Vote_average)
	s.ChannelMessageSend(m.ChannelID, response)
	if movie.Poster_path != "" {
		s.ChannelMessageSend(m.ChannelID, "https://image.tmdb.org/t/p/original/"+movie.Poster_path)
	}
}
