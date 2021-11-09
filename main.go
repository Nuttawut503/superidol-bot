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

var bot *discordgo.Session
var yt *youtube.Client
var prefix = "%"

var superidolYTIDs = []string{"HTGdfE2s4Hw", "chY9p-XLHHk", "DKpaKHUlyBY", "8ywlhKFWAzg"}

var song = map[string]string{
	"cn": "Super idol的笑容\n" +
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

func getYoutubeLink(youtubeID string) string {
	return "https://youtu.be/" + youtubeID
}

func getRandomSuperIdolYTLink() string {
	return getYoutubeLink(superidolYTIDs[rand.Intn(len(superidolYTIDs))])
}

func MessageResponseHandler(bot *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot || !strings.HasPrefix(m.Content, prefix) {
		return
	}
	message := strings.TrimPrefix(m.Content, prefix)
	if message == "link" {
		bot.ChannelMessageSend(m.ChannelID, getRandomSuperIdolYTLink())
	}
	if strings.HasPrefix(message, "lyrics") {
		split := strings.Split(message, " ")
		lang := "cn"
		if len(split) > 1 {
			lang = split[1]
		}
		if lyrics, ok := song[lang]; ok {
			bot.ChannelMessageSend(m.ChannelID, lyrics)
		}
	}
	if strings.HasPrefix(message, "gift") {
		split := strings.Split(message, " ")
		var username string
		if len(split) < 2 {
			return
		}
		username = split[1]
		bot.ChannelMessageSend(m.ChannelID, username+" "+getRandomSuperIdolYTLink())
	}
	if message == "play" {
		channel, err := bot.State.Channel(m.ChannelID)
		if err != nil {
			return
		}
		guild, err := bot.State.Guild(channel.GuildID)
		if err != nil {
			return
		}
		yt = &youtube.Client{}
		video, err := yt.GetVideo(getRandomSuperIdolYTLink())
		if err != nil {
			bot.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
		formats := video.Formats.WithAudioChannels()
		url, err := yt.GetStreamURL(video, &formats[0])
		if err != nil {
			bot.ChannelMessageSend(m.ChannelID, "Get Stream URL: "+err.Error())
			return
		}
		options := dca.StdEncodeOptions
		options.RawOutput = true
		options.Bitrate = 24
		options.Application = "lowdelay"
		encodingSession, err := dca.EncodeFile(url, options)
		if err != nil {
			bot.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
		defer encodingSession.Cleanup()
		for _, vs := range guild.VoiceStates {
			if vs.UserID != m.Author.ID {
				continue
			}
			vc, err := bot.ChannelVoiceJoin(guild.ID, vs.ChannelID, false, true)
			if err != nil {
				bot.ChannelMessageSend(m.ChannelID, err.Error())
				return
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
			return
		}
	}
}

func EmojiResponseHandler(bot *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Username == "Chaletlnwza007" {
		bot.MessageReactionAdd(m.ChannelID, m.ID, "😡")
		bot.MessageReactionAdd(m.ChannelID, m.ID, "🤢")
		return
	}
}

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

func main() {
	bot.AddHandler(func(bot *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})

	err := bot.Open()
	if err != nil {
		log.Fatalf("Can't start session: %v", err)
	}
	defer bot.Close()

	bot.AddHandler(MessageResponseHandler)
	bot.AddHandler(EmojiResponseHandler)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
