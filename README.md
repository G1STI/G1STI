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
- [x] Step 5: Windows build instructions + portable release layout

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
3. Build with PowerShell:
   ```powershell
   ./scripts/build_windows.ps1 -Arch amd64
   ```

### Manual build (PowerShell)
```powershell
cd frontend
npm install
npm run build
cd ..
wails build -platform windows/amd64
```

## Install, run, debug (Windows)
### 1) Install prerequisites
* Go 1.22+
* Node.js 18+ (npm included)
* Wails CLI:
  ```powershell
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```
  The command is usually silent on success. Verify with:
  ```powershell
  wails version
  ```

### 2) Clone and fetch dependencies
```powershell
git clone <YOUR_REPO_URL>
cd G1STI
```

### 3) Development mode (UI hot-reload)
```powershell
cd frontend
npm install
npm run dev
```
In another PowerShell:
```powershell
cd G1STI
wails dev
```
This runs the Go backend and the Vite dev server together, with live reload for the UI.

### 4) Production build
```powershell
./scripts/build_windows.ps1 -Arch amd64
```
The output binary will be in `./build/bin/`.

### 5) First run and data folder
On first run, the app creates:
```
%AppData%\G1STI\
├── data.json
└── logs\
    ├── app.log
    └── singbox.log
```

### 6) Configure sing-box
For MVP, place `sing-box.exe` and `wintun.dll` next to `G1STI.exe` (portable layout below). Then set the sing-box path in Settings (UI) once that screen is wired.

### 7) Debugging tips
* Backend logs: `%AppData%\G1STI\logs\app.log`
* Sing-box logs: `%AppData%\G1STI\logs\singbox.log`
* If `wails dev` fails, confirm `npm install` succeeded and `wails` is in PATH.
* If `npm run build` reports `Unexpected "\xff" in JSON`, there may be a UTF-16 `package.json` in a parent folder. Ensure you run the build inside `frontend/` and keep the repo root `package.json` in UTF-8 (this repo includes one).

### Portable release layout
```
G1STI/
├── G1STI.exe
├── sing-box.exe
├── wintun.dll
└── resources/
```

Keep `sing-box.exe` and `wintun.dll` next to the app executable for the MVP. Configure the sing-box path in Settings (UI) once bindings are wired.

### Release checklist (MVP)
1. Bundle `sing-box.exe` and `wintun.dll` into the portable folder above.
2. Run the app once to create `%AppData%\G1STI\` storage and logs.
3. Add a subscription URL and update nodes.
4. Create a profile, select a node, and connect.
5. Verify `singbox.log` and `app.log` in `%AppData%\G1STI\logs`.

## License
MIT (see `LICENSE`).
