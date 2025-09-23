# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

---
## [0.1.2] - 2025-09-22
### Added
- Export report using /export command
- Dismiss after choose category

## [0.1.1] - 2025-09-22
### Added
- Command to log expenses via Telegram bot.
- Category selection with inline keyboard.
- Google Sheets integration for saving expenses.
- Service & handler refactor for clean architecture.
- Auto-calculation of monthly total expenses (per user).
- Save data to Google Sheets (Timestamp, User ID, Description, Amount, Total Expense).
- Reply message includes monthly total.
- Number formatting with thousand separators.

---

## [0.1.0] - 2025-09-19
### Added
- Initial project setup.
- Basic expense logging (description + amount).
- Save data to Google Sheets (Timestamp, User ID, Description, Amount, Category).
- Project structure with `service/` and `handler/`.
- `.env` configuration support.
- `.gitignore` and `.env-example` files.

---
