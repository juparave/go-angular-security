# Flutter client — Architecture

Mobile/desktop client for the `go-angular-security` Fiber backend. Feature-first
layout with three layers per feature. Riverpod (codegen) for state, GoRouter for
navigation, Dio for HTTP.

```
lib/
├── main.dart                  Entry point — restores the session, then boots Sentry + ProviderScope.
├── app.dart                   Root MaterialApp.router — theme, router, locale.
│
├── core/                      Cross-cutting infrastructure. No business rules.
│   ├── api/                   ApiClient (Dio façade) + ApiException hierarchy.
│   ├── config/                EnvironmentConfig — host selection + /api/v1 prefix.
│   ├── logging/               Shared Logger provider.
│   ├── router/                GoRouter + GoRouterRefreshStream.
│   ├── storage/               SecureStorage façade (Keychain / Encrypted SP).
│   └── theme/                 AppColors + AppTheme (light/dark).
│
├── features/<feature>/        Self-contained slice of the product.
│   ├── data/                  Repositories + DTO models. Talks to ApiClient only.
│   ├── application/           Riverpod controllers (AsyncNotifier, providers).
│   └── presentation/          Screens + feature-private widgets.
│
└── shared/                    Reusable across features.
    ├── widgets/               AsyncValueWidget, empty states, loading indicators.
    └── models/                Cross-feature value objects (User).
```

## Layer rules

1. **Presentation → application → data → core.** Never import in reverse.
2. **One repository per backend concern.** Repositories call `ApiClient` and
   translate `DioException` → `ApiException`.
3. **Controllers** are `AsyncNotifier`s. They orchestrate repositories and expose
   `AsyncValue<T>` to the UI.
4. **Screens are dumb.** `ref.watch`, render, dispatch intents. No `dio`, no
   `SharedPreferences`, no business logic.
5. **Shared is opt-in.** If exactly one feature uses something, keep it inside
   that feature.

## Environment

`core/config/environment_config.dart` picks the backend host at runtime from
compile-time defines, in this precedence order:

| Define | Host |
|---|---|
| `--dart-define=baseUrl=https://…` | that host (always wins) |
| `--dart-define=remoteDebug=true` | staging host |
| *(none)* | **production host** — the default, so a release build is never pointed at localhost |
| `--dart-define=development=true` | `http://localhost:5000`, or `http://10.0.2.2:5000` on Android (emulators cannot reach the host's `localhost`) |

Platform detection uses `defaultTargetPlatform`, not `dart:io`'s `Platform`,
which would break the web build.

`EnvironmentConfig.apiBaseUrl` appends `EnvironmentConfig.apiPrefix` (`/api/v1`)
— the group `server/cmd/api/main.go` mounts every route under. Repositories
therefore use bare paths (`/login`, `/team`, `/subscriptions/current`).

`.vscode/launch.json` wires these up: **app** → local backend, **Remote debug** →
staging, profile/release → production.

## State

- `@Riverpod(keepAlive: true)` for long-lived dependencies (config, api client,
  storage, repositories, router).
- `@riverpod` for screen-scoped derived state and controllers.
- Codegen: `dart run build_runner build` after touching an annotated provider,
  or `dart run build_runner watch` while developing. `riverpod_lint` runs through
  the analyzer plugin declared in `analysis_options.yaml`.

## Auth & navigation

`AuthRepository` owns the session:

- Holds the signed-in `User` in memory (`currentUser`) so GoRouter's `redirect`
  can read it **synchronously**. Reading a stream-backed provider there would
  rebuild a microtask later and leave a just-signed-in user on `/sign-in` for a
  frame.
- Broadcasts every change on `authStateChanges()`, which
  `GoRouterRefreshStream` turns into the router's `refreshListenable`.
- `restoreSession()` runs in `main()` before the first frame, exchanging a
  stored token for `GET /user` so a cold start doesn't flash the sign-in screen.

Everything except `/sign-in` sits behind a `StatefulShellRoute.indexedStack`
(Panel / Equipo / Perfil), each branch keeping its own navigation stack.

Tokens live in `SecureStorage` and are sent as the `jwt` cookie the Fiber
middleware expects. On a 401, `_RefreshInterceptor` calls `POST /refresh-token`
once and retries the original request; if that fails the tokens are cleared and
the router redirects.

## Error handling

- Repositories throw `ApiException` subtypes.
- Controllers wrap with `AsyncValue.guard(...)`.
- Screens render errors via `AsyncValueWidget` or a snackbar.
- Never `catch (_) {}` silently.

## Build

```bash
# local dev (backend: cd server && go run cmd/api/main.go)
flutter run --dart-define=development=true

# staging
flutter run --dart-define=remoteDebug=true

# production
flutter build apk --dart-define=sentryDsn=<dsn>

# one-off host override
flutter run --dart-define=baseUrl=https://api.staging.example.com
```

## Tests

```bash
flutter test                      # unit + widget
dart run build_runner build       # regenerate *.g.dart first if providers changed

# end-to-end against a running backend
flutter test integration_test/sign_in_flow_test.dart \
  --dart-define=development=true \
  --dart-define=e2eEmail=you@example.com \
  --dart-define=e2ePassword=secret
```

The integration test skips itself when the credential defines are absent.
