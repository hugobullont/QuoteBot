package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

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

	dg.AddHandler(ready)       //Agregamos un handler con la función Ready.
	dg.AddHandler(guildCreate) //Agregamos la función guildCreate (Cuando el Bot entra al servidor)
}

//Ready
func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Actualizamos el status del bot. En este caso pondremos !help
	s.UpdateStatus(0, "!help")
}

func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		return
	}

	//Añadimos el mensaje que mandará el bot en cada canal que se encuentre. (y tenga permiso)
	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, _ = s.ChannelMessageSend(channel.ID, "¡QuoteBot está de vuelta! ¡Escribe !help para descubrir las nuevas funciones!")
			return
		}
	}
}
