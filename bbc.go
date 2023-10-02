package main

import (
	"discordgo"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/robfig/cron/v3"
)

// bbc creates the bbc role with a peach icon
// then allows people to record the number of squats they've done
// you're kicked from the bbc role every day at midnight
var servers guildList

type guildList map[string]*guild
type guild struct {
	GuildID   string
	Squatters map[string]*squatter
	RoleID    string
}
type squatter struct {
	TotalSquats int    `json:"totalSquats,omitempty"`
	TotalWeight int    `json:"totalWeight,omitempty"`
	Days        []int  `json:"days,omitempty"`
	TodaySquats int    `json:"todaySquats,omitempty"`
	TodayWeight int    `json:"todayWeight,omitempty"`
	MaxWeight   int    `json:"maxWeight,omitempty"`
	UserID      string `json:"userID"`
	UserName    string `json:"userName"`
}

const filepath string = "./db/bbc.json"

// initialize the guild maps
func init() {

	servers = make(guildList)

	file, err := ioutil.ReadFile(filepath)
	if err == nil {
		err = json.Unmarshal(file, &servers)
		if err != nil {
			fmt.Println("Failed to unmarshal bbc.json into servers: ", err)
		}
	}

	//	loc, err := time.LoadLocation("America/Chicago")
	//	if err != nil {
	//		fmt.Println("Failed to set location to America/Chicago")
	//	}
	cronHandler := cron.New()
	id, err := cronHandler.AddFunc("0 5 * * *", dailyCron)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(id, time.Now())
	cronHandler.Start()
}

func Save() {

	fmt.Println("Saving File")

	file, err := json.MarshalIndent(servers, "", " ")

	if err != nil {
		fmt.Println("Failed to marshal servers: ", err)
	}

	err = ioutil.WriteFile(filepath, file, 0644)

	if err != nil {
		fmt.Println("Failed to write servers to bbc.json: ", err)
	}
	fmt.Println("Done saving file... ")

}

func GetGuildByID(s *discordgo.Session, m *discordgo.MessageCreate) *guild {

	// Guild has been previously registered
	if g, ok := servers[m.GuildID]; ok {
		return g
	}

	// Set up the new role and save server information
	perms := int64(0)
	mentionable := true
	roleparams := discordgo.RoleParams{
		Name:         "Big Booty Club",
		Permissions:  &perms,
		Mentionable:  &mentionable,
		UnicodeEmoji: "\U0001F351",
	}
	role, err := s.GuildRoleCreate(m.GuildID, &roleparams)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Get the new guild
	newGuild := guild{
		GuildID:   m.GuildID,
		Squatters: make(map[string]*squatter),
		RoleID:    role.ID,
	}

	servers[m.GuildID] = &newGuild
	return &newGuild
}

func (g guild) GetSquatter(u discordgo.User) *squatter {
	s, ok := g.Squatters[u.ID]
	if !ok {
		newSquatter := squatter{
			UserID:   u.ID,
			UserName: u.Username,
		}
		g.Squatters[u.ID] = &newSquatter
		s = &newSquatter
	}
	return s
}

func (g guild) Squat(u discordgo.User, nsquats int, weight int, s *discordgo.Session) {

	fmt.Println(weight)
	sq := g.GetSquatter(u)
	if sq.TodaySquats == 0 {
		// add to the BBC role
		err := s.GuildMemberRoleAdd(g.GuildID, u.ID, g.RoleID)
		if err != nil {
			fmt.Printf("failed to add user %v to %v: %v", u.ID, g.RoleID, err)
		}
	}
	if weight > 0 {
		sq.TodayWeight += nsquats * weight
		sq.TotalWeight += nsquats * weight
		if weight > sq.MaxWeight {
			sq.MaxWeight = weight
		}
	}
	sq.TodaySquats += nsquats
	sq.TotalSquats += nsquats
	Save()
}

func dailyCron() {
	for guild, squatters := range servers {
		fmt.Printf("BBC for guild %v", guild)
		for _, squatter := range squatters.Squatters {
			days := squatter.Days
			days = append(days, squatter.TodaySquats)
			if len(days) > 7 {
				days = days[:len(days)-1]
			}
			squatter.Days = days
			squatter.TodaySquats = 0
			squatter.TodayWeight = 0
			err := globalSession.GuildMemberRoleRemove(squatters.GuildID, squatter.UserID, squatters.RoleID)
			if err != nil {
				fmt.Println("Failed to remove user from role: ", err)
			}
		}
	}
	Save()
}
