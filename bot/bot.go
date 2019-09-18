package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/Moekr/sword/common/constant"
)

const (
	badQueryMessage        = "不能识别的查询"
	queryErrorMessage      = "未知错误"
	serverForbiddenMessage = "认证失败"
	serverErrorMessage     = "服务端错误"
	noSuchTargetMessage    = "监测目标不存在"
)

func main() {
	api, server, token, debug := parseArgs()
	bot, err := tgbotapi.NewBotAPI(api)
	if err != nil {
		panic(err)
	}
	bot.Debug = debug
	fmt.Printf("[Bot] authorized on account %s", bot.Self.UserName)
	prefix := "@" + bot.Self.UserName + " "
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		var query string
		if update.Message != nil {
			text := update.Message.Text
			if update.Message.Chat.IsPrivate() {
				query = text
			} else if strings.HasPrefix(text, prefix) {
				query = strings.TrimPrefix(text, prefix)
			}
		}
		if len(query) == 0 {
			continue
		}
		response := doQuery(server, token, query)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ParseMode = "Markdown"
		bot.Send(msg)
	}
}

func doQuery(server, token, query string) string {
	qs := strings.Split(query, " ")
	if len(qs) == 1 {
		qs = append(qs, "")
	} else if len(qs) == 2 {
		query = qs[0]
	} else {
		return badQueryMessage
	}
	id, err := strconv.ParseInt(query, 10, 64)
	if err != nil {
		return badQueryMessage
	}
	url := fmt.Sprintf("%s/api/stat?t=%d&i=%s", server, id, qs[1])
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(constant.TokenHeader, token)
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return queryErrorMessage
	}
	switch rsp.StatusCode {
	case http.StatusForbidden:
		return serverForbiddenMessage
	case http.StatusInternalServerError:
		return serverErrorMessage
	case http.StatusNotFound:
		return noSuchTargetMessage
	case http.StatusOK:
		if bs, err := ioutil.ReadAll(rsp.Body); err != nil {
			return queryErrorMessage
		} else {
			return parseBody(bs)
		}
	default:
		return queryErrorMessage
	}
}

func parseBody(body []byte) (response string) {
	defer func() {
		if v := recover(); v != nil {
			response = queryErrorMessage
		}
	}()
	var obj map[string]interface{}
	if err := json.Unmarshal(body, &obj); err != nil {
		return queryErrorMessage
	}
	data := obj["data"].(map[string]interface{})
	target, stats := data["target"].(map[string]interface{}), data["stats"].([]interface{})
	b := strings.Builder{}
	b.WriteString("```\n")
	b.WriteString(target["name"].(string))
	b.WriteString("\n  Avg  /  Max  /  Min  /  Std  / Lost\n")
	for _, stat := range stats {
		m := stat.(map[string]interface{})
		b.WriteString(padStart(fmt.Sprint(m["avg"].(float64)), 4, ' '))
		b.WriteString("ms /")
		b.WriteString(padStart(fmt.Sprint(m["max"].(float64)), 4, ' '))
		b.WriteString("ms /")
		b.WriteString(padStart(fmt.Sprint(m["min"].(float64)), 4, ' '))
		b.WriteString("ms /")
		b.WriteString(padStart(fmt.Sprint(m["std"].(float64)), 4, ' '))
		b.WriteString("ms /")
		b.WriteString(padStart(fmt.Sprint(m["los"].(float64)), 3, ' '))
		b.WriteString("% from ")
		b.WriteString(m["client"].(map[string]interface{})["name"].(string))
		b.WriteString("\n")
	}
	b.WriteString("```\n")
	return b.String()
}

func parseArgs() (api, server, token string, debug bool) {
	flag.StringVar(&api, "b", "", "telegram bot api token")
	flag.StringVar(&server, "s", "http://localhost:7901", "sword server address")
	flag.StringVar(&token, "t", "", "token used in communication with sword server")
	flag.BoolVar(&debug, "v", false, "identify debug mode or not")
	flag.Parse()
	return api, server, token, debug
}

func padStart(s string, l int, r rune) string {
	if len(s) >= l {
		return s
	}
	b := strings.Builder{}
	for i := 0; i < l-len(s); i++ {
		b.WriteRune(r)
	}
	b.WriteString(s)
	return b.String()
}
