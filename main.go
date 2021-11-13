package main

import (
	"io"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/kkdai/youtube/v2"
)

var (
	bot       *discordgo.Session
	prefix    = "%"
	serverIDs = []string{""} // add your server ID here
)

func init() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("No token provided...")
	}
	var err error
	bot, err = discordgo.New("Bot " + args[1])
	if err != nil {
		log.Fatalf("Can't create session from the token: %v", err)
	}
}

var (
	yt            *youtube.Client
	encodeOptions *dca.EncodeOptions
)

func init() {
	yt = &youtube.Client{}
	encodeOptions = dca.StdEncodeOptions
	encodeOptions.RawOutput = true
	encodeOptions.Bitrate = 24
	encodeOptions.Application = "lowdelay"
}

func getStreamURL(videoUrl string) (string, error) {
	video, err := yt.GetVideo(getRandomSuperIdolYTLink())
	if err != nil {
		return "", err
	}
	formats := video.Formats.WithAudioChannels()
	url, err := yt.GetStreamURL(video, &formats[0])
	if err != nil {
		return "", err
	}
	return url, nil
}

var (
	infamousPeople = map[string]string{
		"champ": "https://images-ext-1.discordapp.net/external/Lp_FlEMlN1S7iDm6h4BCI0Nu0jl0hQZKrfdKA_mWKTU/https/media.discordapp.net/attachments/842035790363885608/887279334258266182/sunglassesChampThink.gif?format=png",
		"fluke": "https://media.discordapp.net/attachments/842035790363885608/907606047403950090/fluk.png",
	}
	song = map[string]string{
		"chinese": "Super idol的笑容\n" +
			"都没你的甜\n" +
			"八月正午的阳光\n" +
			"都没你耀眼\n" +
			"热爱105°C度的你\n" +
			"滴滴清纯的蒸馏水\n" +
			"你不知道你有多可爱\n" +
			"跌倒后会傻笑着再站起来\n" +
			"你从来都不轻言失败\n" +
			"对梦想的执著一直不曾更改\n" +
			"很安心 当你对我说\n" +
			"不怕有我在\n" +
			"放著让我来\n" +
			"勇敢追自己的梦想\n" +
			"那坚定的模样\n",
		"pinyin": "Super Idol de xiaorong\n" +
			"dou mei ni de tian\n" +
			"ba yue zhengwu de yangguang\n" +
			"dou mei ni yaoyan\n" +
			"re’ai 105 °C de ni\n" +
			"di di qingchun de zhengliushui\n" +
			"ni bu zhidao ni you duo ke’ai\n" +
			"diedao hou hui shaxiaozhe zai zhan qilai\n" +
			"ni conglai dou bu qing yan shibai\n" +
			"dui mengxiang de zhizhuo yizhi buceng genggai\n" +
			"hen anxin dang ni dui wo shuo\n" +
			"bupa you wo zai\n" +
			"fangzhe rang wo lai\n" +
			"yonggan zhui ziji de mengxiang\n" +
			"na jianding de muyang\n",
		"romaji": "super idol no egao yori mo\n" +
			"ano hachigatsu no gogo yori mo\n" +
			"hyakkugosen shuu tou yori\n" +
			"hikaru kimi e\n" +
			"kawaii tto ierunara\n" +
			"koronde mo sugu warau kimi wa\n" +
			"yume wa tooi hazunanoni\n" +
			"yubi sashita hoshi ga chikazuita\n" +
			"yasashii kaze fuite\n" +
			"tonari ijou motto chikaku\n" +
			"futari nara daijoubu sou ittara\n",
		"thai": "Super Idol ก็ยิ้มไม่หวานได้เท่ากับเธอ\n" +
			"ดวงอาทิตย์ที่ว่าสดใสก็ยังไม่เท่าเธอ\n" +
			"องศารักที่ 105 นี้ได้กลั่นเป็นน้ำสะอาดใสไหลริน\n" +
			"เคยรู้ไหมว่าเธอน่ารักแค่ไหน\n" +
			"แม้ล้มลงไปกี่ครั้งก็จะลุกขึ้นใหม่\n" +
			"เรื่องไหนเธอก็ไม่เคยคิดถอดใจ\n" +
			"มุ่งมั่นวิ่งตามความฝันและไม่เคยผันแปรไป\n" +
			"เธอบอกฉันว่า เธอไม่ต้องกลัว ไม่ว่าเจอเรื่องใด เธอก็ยังมีฉัน\n" +
			"จงตั้งใจไล่ตามความฝันและจงไม่ยอมเลิกราไปง่ายๆ\n",
	}
	superidolYTIDs = []string{"https://youtu.be/HTGdfE2s4Hw", "https://youtu.be/chY9p-XLHHk", "https://youtu.be/DKpaKHUlyBY", "https://youtu.be/8ywlhKFWAzg"}
)

func getRandomSuperIdolYTLink() string {
	return superidolYTIDs[rand.Intn(len(superidolYTIDs))]
}

func MessageResponseHandler(bot *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot || !strings.HasPrefix(m.Content, prefix) {
		return
	}
	message := strings.TrimPrefix(m.Content, prefix)
	switch {
	case message == "play":
		channel, err := bot.State.Channel(m.ChannelID)
		if err != nil {
			return
		}
		guild, err := bot.State.Guild(channel.GuildID)
		if err != nil {
			return
		}
		channelID := ""
		for _, vs := range guild.VoiceStates {
			if vs.UserID == m.Author.ID {
				channelID = vs.ChannelID
				break
			}
		}
		if channelID == "" {
			bot.ChannelMessageSend(m.ChannelID, "You aren't in a voice channel")
			return
		}
		url, err := getStreamURL(getRandomSuperIdolYTLink())
		if err != nil {
			bot.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
		encodingSession, err := dca.EncodeFile(url, encodeOptions)
		if err != nil {
			bot.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
		defer encodingSession.Cleanup()
		vc, err := bot.ChannelVoiceJoin(guild.ID, channelID, false, true)
		if err != nil {
			if _, ok := bot.VoiceConnections[guild.ID]; ok {
				vc = bot.VoiceConnections[guild.ID]
			} else {
				bot.ChannelMessageSend(m.ChannelID, err.Error())
				return
			}
		}
		vc.Speaking(true)
		done := make(chan error)
		dca.NewStream(encodingSession, vc, done)
		err = <-done
		if err != nil && err != io.EOF {
			bot.ChannelMessageSend(m.ChannelID, err.Error())
		}
		vc.Speaking(false)
		vc.Disconnect()
	}
}

func EmojiResponseHandler(bot *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Username == "Chaletlnwza007" {
		bot.MessageReactionAdd(m.ChannelID, m.ID, "😡")
		bot.MessageReactionAdd(m.ChannelID, m.ID, "🤢")
		return
	}
}

func SlashCommandsHandler(bot *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd := i.ApplicationCommandData().Name
	if link, exist := infamousPeople[cmd]; exist {
		bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: link,
			},
		})
	}
	if cmd == "gift" {
		bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: getRandomSuperIdolYTLink(),
			},
		})
	}
	if cmd == "lyrics" {
		if len(i.ApplicationCommandData().Options) == 0 {
			return
		}
		optionName := i.ApplicationCommandData().Options[0].Name
		if optionName != "language" {
			return
		}
		if len(i.ApplicationCommandData().Options[0].Options) == 0 {
			return
		}
		selectedOption := i.ApplicationCommandData().Options[0].Options[0].Name
		if lyrics, ok := song[selectedOption]; ok {
			bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: lyrics,
				},
			})
		}
	}
}

func init() {
	bot.AddHandler(MessageResponseHandler)
	bot.AddHandler(EmojiResponseHandler)
	bot.AddHandler(SlashCommandsHandler)
}

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "champ",
		Description: "cool pictures from Champ",
	},
	{
		Name:        "fluke",
		Description: "cool pictures from Fluke",
	},
	{
		Name:        "gift",
		Description: "receive a random superidol link",
	},
	{
		Name:        "lyrics",
		Description: "show superidol lyrics",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Name:        "language",
				Description: "choose the lyrics' language",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "chinese",
						Description: "Chinese ver.",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
					},
					{
						Name:        "pinyin",
						Description: "Chinese ver. but pinyin",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
					},
					{
						Name:        "romaji",
						Description: "Japanese ver. with romaji",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
					},
					{
						Name:        "thai",
						Description: "Thai ver.",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
					},
				},
			},
		},
	},
}

func main() {
	bot.AddHandler(func(bot *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})

	err := bot.Open()
	if err != nil {
		log.Fatalf("Can't start session: %v", err)
	}
	defer bot.Close()

	for _, cmd := range commands {
		if len(serverIDs) == 0 {
			_, err = bot.ApplicationCommandCreate(bot.State.User.ID, "", cmd)
			if err != nil {
				log.Printf("Can't create a command: %v\n", err)
			}
		} else {
			for _, serverID := range serverIDs {
				_, err = bot.ApplicationCommandCreate(bot.State.User.ID, serverID, cmd)
				if err != nil {
					log.Printf("GuildID(%v) - Can't create a command: %v\n", serverID, err)
				}
			}
		}
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Bot is down...")
}
