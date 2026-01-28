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
- [ ] Step 2: Subscriptions (download, decode, parse, normalize)
- [ ] Step 3: UI (Home, Subscriptions, Profiles/Settings, About/Logs) + tray
- [ ] Step 4: Integration tests (vless parsing, sing-box start/stop w/ mock)
- [ ] Step 5: Windows build instructions + portable release layout

## Layout
```
.
├── frontend/                 # React + TS (Wails v2)
├── internal/
│   ├── model/                # Data models
│   ├── paths/                # AppData paths
│   ├── singbox/              # sing-box process + config generation
│   └── storage/              # JSON storage + schema migrations
├── go.mod
└── README.md
```

## MVP scope
See the user story in the task description. Focus is: subscriptions → nodes → connect with sing-box, tray control, logs.

## License
MIT (see `LICENSE`).
