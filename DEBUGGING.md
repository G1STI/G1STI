# Debugging toolkit

This file describes a **quick debugging flow** so we can diagnose issues faster and iterate with fewer back-and-forth messages.

---

## 1) Quick checks (Windows)

Run these in PowerShell from the repo root:

```powershell
# Backend + frontend dev (two shells)
cd frontend
npm install
npm run dev
```

```powershell
cd G1STI
wails dev
```

If `wails dev` fails or the EXE closes immediately:

```powershell
# Show app log
Get-Content "$env:APPDATA\\G1STI\\logs\\app.log" -Tail 200

# Show sing-box log
Get-Content "$env:APPDATA\\G1STI\\logs\\singbox.log" -Tail 200
```

If `npm run build` fails with `Unexpected "\xff" in JSON`:

```powershell
Rename-Item "$env:USERPROFILE\\package.json" package.json.bak
```

---

## 2) In-app debug panel

Open the **Logs** tab and press **Refresh**.

The app will show the last ~200 lines from `%AppData%\\G1STI\\logs\\app.log` so you can paste the exact error message.

---

## 3) “Report a bug” template (for GitHub issues)

Please copy‑paste this into a GitHub issue or message:

```
### What I did
1)
2)
3)

### Expected

### Actual

### Logs (paste last 200 lines)
```

Also include:
* Windows version
* Go version (`go version`)
* Wails version (`wails version`)
* Node version (`node -v`)

---

## 4) Fast‑triage checklist (for the agent)

When reporting to the assistant, include:

```
1) Steps to reproduce
2) Exact error text
3) app.log tail
4) singbox.log tail (if connect issue)
5) Versions: Go / Wails / Node / Windows
```

This lets me fix issues quickly without multiple back‑and‑forth messages.
