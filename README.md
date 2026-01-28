# G1STI Windows Client (MVP)

Minimal Windows 10/11 desktop client for VLESS (and other subscription URIs) powered by `sing-box`.

> **Disclaimer (legal/admin use only):** This software is intended strictly for lawful, authorized network administration and корпоративного/корпоративного применения. You are responsible for complying with local laws, policies, and provider terms. The authors do not endorse or support any misuse.

## Project goals
- Fast connect/disconnect
- Simple profile/subscription management
- Clear status and logs
- Minimal UI (no “100 tabs”)

## Checklist (plan of work)
- [x] Step 0: Repository structure + README checklist
- [x] Step 1: Backend core (models, storage, sing-box process manager, config generator)
- [x] Step 2: Subscriptions (download, decode, parse, normalize)
- [ ] Step 3: UI (Home, Subscriptions, Profiles/Settings, About/Logs) + tray
- [x] Step 4: Integration tests (vless parsing, sing-box start/stop w/ mock)
- [ ] Step 5: Windows build instructions + portable release layout

## Layout
```
.
├── frontend/                 # React + TS (Wails v2)
├── internal/
│   ├── model/                # Data models
│   ├── app/                  # Backend app services
│   ├── parser/               # URI parsing helpers
│   ├── paths/                # AppData paths
│   ├── singbox/              # sing-box process + config generation
│   └── storage/              # JSON storage + schema migrations
│   └── subscription/         # Subscription download + normalize
├── go.mod
└── README.md
```

## MVP scope
See the user story in the task description. Focus is: subscriptions → nodes → connect with sing-box, tray control, logs.

## Build (Windows, Wails v2)
1. Install Go 1.22+ and Node.js 18+.
2. Install Wails CLI:
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```
3. Install frontend dependencies:
   ```bash
   cd frontend
   npm install
   ```
4. Build Windows binary:
   ```bash
   wails build -platform windows/amd64
   ```

### Portable release layout
```
G1STI/
├── G1STI.exe
├── sing-box.exe
├── wintun.dll
└── resources/
```

Keep `sing-box.exe` and `wintun.dll` next to the app executable for the MVP. Configure the sing-box path in Settings (UI) once bindings are wired.

## License
MIT (see `LICENSE`).
