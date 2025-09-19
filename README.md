## 📄 README.md

```markdown
# Go Personal Expense Bot 💰🤖

A simple Telegram bot built with **Go** to track personal expenses and store them in **Google Sheets**.

---

## Features
- Add expenses directly via Telegram chat.
- Auto-parse amount and description.
- Choose expense category via inline keyboard.
- Save data to Google Sheets:
  - Timestamp
  - User ID
  - Description
  - Amount
  - Category
---

## Setup

1. **Clone the repo**
   ```bash
   git clone git@github.com:yasseryazid/go-personal-expense-bot.git
   cd go-personal-expense-bot
````

2. **Install dependencies**

   ```bash
   go mod tidy
   ```

3. **Create `.env`**

   ```env
   TELEGRAM_BOT_TOKEN=your_telegram_bot_token
   GOOGLE_CREDENTIALS=credentials.json
   GOOGLE_SHEET_ID=your_google_sheet_id
   ```

4. **Google Sheets**

   * Create a new Sheet and copy its ID.
   * Add headers: `Timestamp | UserID | Description | Amount | Category`.
   * Create a Service Account, download `credentials.json`, and share the Sheet with it (Editor access).

5. **Run the bot**

   ```bash
   go run main.go
   ```

---

## Usage

1. Send an expense in chat:

   ```
   Buy coffee 15000
   ```
2. Choose a category from inline buttons.
3. Bot confirms and saves it to Google Sheets.

---

## License

FREE © 2025 [Yasser Yazid](https://github.com/yasseryazid)
