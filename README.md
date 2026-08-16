# payslip-tool

Generates monthly payslips (xlsx + pdf) for teachers from a government
Excel export. Built for a single kendra (13 schools, 48 teachers).

## Requirements

- Go 1.22+
- LibreOffice (`soffice`) — used headless for xlsx → pdf conversion
  ```
  sudo pacman -S libreoffice-fresh
  ```

## Usage

```bash
go run main.go -f=<path to govt xlsx> -m=<month 1-12> -y=<year>

# example
go run main.go -f=Aug26.xlsx -m=8 -y=2026
```

Output lands in `output/` (xlsx) and `output/pdf/` (pdf), one file per
employee, named `{School}_{Employee Name}.xlsx` / `.pdf`.

Safe to re-run on the same file/month — inserts are upserts, so nothing
duplicates.

## Project structure

```
main.go                    — orchestrates the full pipeline
internal/
  db/                       — SQLite connection, schema, inserts, joined export query
  importer/                 — parses the govt xlsx into structs
  excelgen/                 — fills the payslip xlsx template per employee
  pdfgen/                   — xlsx → pdf via LibreOffice headless
  concurrency/               — worker pool for pdf conversion
  timer/                    — simple execution timing helper
template.xlsx               — blank payslip layout (safe to commit, no real data)
```

## Design notes

- **SQLite in WAL mode, `synchronous=FULL`** — small monthly write volume,
  so the durability cost of full fsync is negligible; worth it for payroll data.
- **Surrogate integer PKs, business keys (UDISE, Shalarth ID) as UNIQUE** —
  insulates the schema from the government ever changing those formats.
- **`ON CONFLICT ... DO UPDATE`** on every insert — re-running the same
  month's file updates existing rows instead of duplicating or erroring.
- **PDF conversion uses a 4-worker pool**, each with an isolated LibreOffice
  profile (`-env:UserInstallation=...`) to avoid lock collisions between
  concurrent `soffice` processes.
- **No mailer, no frontend** — dropped deliberately; recipient emails
  aren't available, and this is run as a CLI, not through a UI.

## .env

Only needed if a future mailer is added back. Currently unused by the
pipeline itself.
