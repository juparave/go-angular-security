# app

Flutter client for the [`go-angular-security`](../) SaaS starter — the same
Fiber + GORM backend the Angular frontend talks to.

## Quick start

```bash
# 1. backend
cd ../server && cp .env.example .env && go run cmd/api/main.go   # :5000

# 2. client
cd ../app
flutter pub get
dart run build_runner build      # generates the Riverpod *.g.dart files
flutter run --dart-define=development=true
```

In VS Code, pick the **app** launch configuration — it passes
`--dart-define=development=true` for you. **Remote debug** targets the staging
host; profile/release target production.

## Where things are

See [ARCHITECTURE.md](ARCHITECTURE.md) for the layer rules, the environment
precedence table, and how auth drives routing.

Before shipping, set the real hosts in
`lib/core/config/environment_config.dart` (`_prodHost`, `_stagingHost`) — they
ship as `example.com` placeholders.
