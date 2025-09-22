package handler

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go-telegram-expense-bot/service"

	"github.com/dustin/go-humanize"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramHandler struct {
	tg *service.Telegram
	gs *service.GoogleSheets
}

var categories = []string{"Makanan/Minuman", "Transportasi", "Belanja", "Tagihan", "Lainnya"}

func NewTelegramHandler(tg *service.Telegram, gs *service.GoogleSheets) *TelegramHandler {
	return &TelegramHandler{tg: tg, gs: gs}
}

func (h *TelegramHandler) Start() {
	for update := range h.tg.Updates() {
		if update.Message != nil {
			h.handleMessage(update.Message)
		}
		if update.CallbackQuery != nil {
			h.handleCallback(update.CallbackQuery)
		}
	}
}

func (h *TelegramHandler) handleMessage(msg *tgbotapi.Message) {
	if msg.Chat.Type != "private" {
		h.tg.Bot.Send(tgbotapi.NewMessage(msg.Chat.ID,
			"⚠️ Expenzo hanya bisa digunakan di *private chat*. Silakan chat langsung dengan Expenzo.",
		))
		return
	}

	if msg.Text == "/start" {
		name := msg.From.FirstName
		if msg.From.LastName != "" {
			name += " " + msg.From.LastName
		}

		greeting := fmt.Sprintf(
			"Halo %s, selamat datang di *Expenzo*!\n\n"+
				"Dengan Expenzo kamu bisa:\n"+
				"• Catat pengeluaran harian\n"+
				"• Pilih kategori\n"+
				"• Lihat total bulanan dengan /total\n"+
				"• Export pengeluaran bulan ini ke CSV dengan /export\n\n"+
				"Langsung coba ketik pengeluaran (contoh: `Beli kopi 15000`).",
			name,
		)

		msgReply := tgbotapi.NewMessage(msg.Chat.ID, greeting)
		msgReply.ParseMode = "Markdown"
		h.tg.Bot.Send(msgReply)
		return
	}

	if msg.Text == "/total" {
		monthlyTotal, _ := h.gs.GetMonthlyTotalByUser(msg.From.ID)

		reply := fmt.Sprintf(
			"💰 Total pengeluaran kamu bulan ini: Rp%s",
			humanize.Comma(int64(monthlyTotal)),
		)

		h.tg.Bot.Send(tgbotapi.NewMessage(msg.Chat.ID, reply))
		return
	}

	if msg.Text == "/export" {
		records, err := h.gs.GetMonthlyDataByUser(msg.From.ID)
		if err != nil {
			h.tg.Bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ Gagal ambil data dari Google Sheets"))
			return
		}

		exporter := service.NewExporter()

		username := msg.From.UserName
		if username == "" {
			username = msg.From.FirstName
		}

		filePath, err := exporter.ExportToCSV(username, records)
		if err != nil {
			h.tg.Bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ Gagal export ke CSV"))
			return
		}

		doc := tgbotapi.NewDocument(msg.Chat.ID, tgbotapi.FilePath(filePath))
		doc.Caption = fmt.Sprintf("📊 Export data bulan ini untuk %s", username)
		h.tg.Bot.Send(doc)

		os.Remove(filePath)
		return
	}

	amount, desc := parseExpense(msg.Text)
	if amount > 0 {
		reply := tgbotapi.NewMessage(msg.Chat.ID, "Pilih kategori:")
		reply.ReplyMarkup = buildCategoryKeyboard(amount, desc, msg.From.ID)
		h.tg.Bot.Send(reply)
	} else {
		h.tg.Bot.Send(tgbotapi.NewMessage(msg.Chat.ID,
			"Ketik pengeluaran, contoh: Beli kopi 10000\nAtau ketik /total untuk cek total bulan ini."))
	}
}

func (h *TelegramHandler) handleCallback(cb *tgbotapi.CallbackQuery) {
	parts := strings.SplitN(cb.Data, "|", 4)
	if len(parts) < 4 {
		return
	}

	category := parts[0]
	amount, _ := strconv.Atoi(parts[1])
	desc := parts[2]
	userID, _ := strconv.ParseInt(parts[3], 10, 64)

	_ = h.gs.Save(category, amount, desc, userID)

	monthlyTotal, _ := h.gs.GetMonthlyTotalByUser(userID)

	reply := fmt.Sprintf(
		"Pengeluaran dicatat:\n"+
			"- Deskripsi: %s\n"+
			"- Jumlah: Rp%s\n"+
			"- Kategori: %s\n"+
			"- Tanggal: %s\n"+
			"- Total bulan ini: Rp%s",
		desc,
		humanize.Comma(int64(amount)),
		category,
		time.Now().Format("2006-01-02"),
		humanize.Comma(int64(monthlyTotal)),
	)

	h.tg.Bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, reply))
}

func parseExpense(input string) (int, string) {
	re := regexp.MustCompile(`(\d+)`)
	match := re.FindString(input)
	if match == "" {
		return 0, ""
	}
	amount, _ := strconv.Atoi(match)
	desc := strings.TrimSpace(strings.Replace(input, match, "", 1))
	return amount, desc
}

func buildCategoryKeyboard(amount int, desc string, userID int64) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, cat := range categories {
		callback := fmt.Sprintf("%s|%d|%s|%d", cat, amount, desc, userID)
		btn := tgbotapi.NewInlineKeyboardButtonData(cat, callback)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
