package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Variables
var (
	Token        string
	quoteCommand string
	quotes       []string
	err          error
	buffer       = make([][]byte, 0)
)

//Leer desde Archivo
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

//Escribir en Archivo
func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func main() {

	Token = "DiscordBotToken"
	quoteCommand = "exec"

	//Cargamos las citas.
	quotes, err = readLines("quotes.txt")
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	//Creamos una sesión de DiscordGo utilizando el Token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(ready)         //Agregamos un handler con la función Ready. (Status del Bot)
	dg.AddHandler(guildCreate)   //Agregamos la función guildCreate (Cuando el Bot entra al servidor)
	dg.AddHandler(messageCreate) //Agregamos la función messageCreate (Cuando se reciben mensajes)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

//Ready
func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Actualizamos el status del bot.
	s.UpdateStatus(0, "DOOM")
}

func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		return
	}
	//Añadimos el mensaje que mandará el bot en cada canal que se encuentre. (y tenga permiso)
	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, _ = s.ChannelMessageSend(channel.ID, "¡QuoteBot está de vuelta!")
			return
		}
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	//Generamos un random
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	//ignoramos los mensajes que envíe el bot.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, quoteCommand) { //Si el mensaje contiene el comando "exec" envia un mensaje random del quotes.txt
		number := r1.Intn(len(quotes))
		s.ChannelMessageSend(m.ChannelID, quotes[number])
	}
}

//https://discordapp.com/api/oauth2/authorize?client_id=433765018748846080&permissions=3660864&scope=bot
