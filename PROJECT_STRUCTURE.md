```
mindkit-gitk/
├── cmd/
│   └── gitk/               # Command-line interface
│       └── main.go
├── internal/
│   ├── storage/           # BNB Greenfield storage implementation
│   │   ├── object.go
│   │   ├── reference.go
│   │   └── storage.go
│   ├── commands/          # Git command implementations
│   │   ├── init.go
│   │   ├── add.go
│   │   ├── commit.go
│   │   └── push.go
│   └── mindkit/          # MindKit integration
│       ├── ai.go
│       └── client.go
├── pkg/                   # Public packages
│   └── protocol/         # Protocol definitions
├── docs/                 # Documentation
├── go.mod
├── go.sum
└── README.md
```