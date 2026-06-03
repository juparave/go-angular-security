# GuiaChida Flutter — Architecture

Feature-first layout with three layers per feature. Riverpod for state, GoRouter for navigation, Dio for HTTP. Same pattern as the medynota companion app.

```
lib/
├── main.dart                  Entry point — bootstraps ProviderScope + Sentry.
├── app.dart                   Root MaterialApp.router — theme, router, locale.
│
├── core/                      Cross-cutting infrastructure. No business rules.
│   ├── api/                   ApiClient (Dio façade) + ApiException hierarchy.
│   ├── config/                Compile-time env (Env.apiBaseUrl, isProduction).
│   ├── logging/               Shared Logger provider.
│   ├── router/                GoRouter wired to auth state.
│   ├── storage/               SecureStorage façade (Keychain / Encrypted SP).
│   ├── theme/                 AppColors (Talavera palette) + AppTheme (light/dark).
│   └── location/              LocationService — geolocator + permission flow (Phase 3).
│
├── features/<feature>/        Self-contained slice of the product.
│   ├── data/                  Repositories + DTO models. Talks to ApiClient only.
│   ├── application/           Riverpod controllers (AsyncNotifier, providers).
│   └── presentation/          Screens + feature-private widgets.
│
└── shared/                    Reusable across features.
    ├── widgets/               AsyncValueWidget, empty states, loading indicators.
    ├── utils/                 Formatters (es-MX), validators.
    └── models/                Cross-feature value objects (rare).
```

## Layer rules

1. **Presentation → application → data → core.** Never import in reverse.
2. **One repository per backend concern.** Repositories call `ApiClient` and translate `DioException` → `ApiException`.
3. **Controllers** are `AsyncNotifier`s. They orchestrate repositories and expose `AsyncValue<T>` to the UI.
4. **Screens are dumb.** `ref.watch`, render, dispatch intents. No `dio`, no `SharedPreferences`, no business logic.
5. **Shared is opt-in.** If exactly one feature uses something, keep it inside that feature.

## Theme

Colors from the web app's Talavera DaisyUI theme — `AppColors` in `core/theme/app_colors.dart`.
Fonts: **Hanken Grotesk** (body) via `google_fonts`, **Piazzolla** (display/headings) via `google_fonts`.
Semantic tokens via `TalaveraColors` ThemeExtension — access with `context.talavera`.

## State

- `Provider` for pure dependencies (repositories, clients).
- `AsyncNotifierProvider` for async state (auth, lists, detail).
- Avoid `StateProvider` beyond trivial UI toggles.

## Navigation

GoRouter with auth redirect. Public routes accessible without token. `/app/**` redirects to `/login` when unauthenticated.

## Error handling

- Repositories throw `ApiException` subtypes.
- Controllers wrap with `AsyncValue.guard(...)`.
- Screens render errors via `AsyncValueWidget` or snackbar.
- Never `catch (_) {}` silently.

## Build

```bash
# local dev
flutter run --dart-define=API_BASE_URL=http://localhost:8080

# production
flutter build apk \
  --dart-define=API_BASE_URL=https://api.guiachida.mx \
  --dart-define=PRODUCTION=true \
  --dart-define=SENTRY_DSN=<dsn>
```
