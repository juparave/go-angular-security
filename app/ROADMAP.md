# GuiaChida Flutter App — Roadmap

Mobile companion to the GuiaChida web app. Lets travelers discover and walk curated city trails, and gives authenticated users (curators, team) access to management surfaces.

Architecture mirrors the medynota Flutter app: feature-first layout, Riverpod for state, GoRouter for navigation, Dio for HTTP.

---

## Phase 0 — Foundation ✅

- [x] `pubspec.yaml` — all dependencies (Riverpod, GoRouter, Dio, flutter_map, google_sign_in, etc.)
- [x] Asset pipeline — fonts (HankenGrotesk, Piazzolla) bundled in `assets/fonts/`; brand images in `assets/brand/`
- [x] `core/theme/` — `AppColors` (Talavera palette), `AppTheme` (light + dark `ThemeData`), `TalaveraColors` extension
- [x] `core/config/` — `Env` (compile-time `API_BASE_URL`, `PRODUCTION`, `SENTRY_DSN`)
- [x] `core/api/` — `ApiClient` Dio façade with cookie-auth + error interceptors; typed `ApiException` hierarchy
- [x] `core/storage/` — `SecureStorage` façade (Keychain / EncryptedSharedPreferences)
- [x] `core/router/` — GoRouter wired to `isAuthenticatedProvider`; auth guard on `/app`
- [x] `core/logging/` — shared `Logger` provider
- [x] `main.dart` + `app.dart` — `ProviderScope`, `MaterialApp.router`, es-MX locale, Sentry bootstrap
- [x] `shared/widgets/` — `AsyncValueWidget`
- [x] `shared/models/` — `User` model
- [x] `ARCHITECTURE.md` inside `app/`
- [x] App renamed to "GuiaChida" across Android, iOS, macOS, Linux, Windows, Web

---

## Phase 1 — Auth ✅

- [x] `features/auth/data/AuthRepository` — Google Sign-In, email/password login, logout
- [x] `features/auth/application/AuthNotifier` — `AsyncNotifier<User?>`, `isAuthenticatedProvider`
- [x] `features/auth/presentation/LoginScreen` — Google Sign-In button (primary) + collapsible email/password form
- [x] Router guard — unauthenticated users redirected from `/app` to `/login`
- [x] Tokens stored as `jwt` / `refreshjwt` cookies via `SecureStorage`, injected as `Cookie:` header by `ApiClient`
- [x] Google ID token sent to `POST /api/glogin` for server-side verification

**Pending in Phase 1:**
- [x] Token refresh — `_RefreshInterceptor` in `ApiClient`: silent retry on 401, falls back to logout on failure
- [ ] Register screen (`POST /api/register-account`)
- [ ] Forgot password flow (`POST /api/request-password-reset`)

---

## Phase 2 — Public discovery (no auth required) ← NEXT

### 2a — Home
- [ ] `features/home/data/HomeRepository` — `GET /api/trails?featured=true`, `GET /api/cities?featured=true`
- [ ] `features/home/application/HomeNotifier`
- [ ] `HomeScreen` — hero banner, featured trails horizontal scroll, featured cities strip

### 2b — Trail list & search
- [ ] `features/trails/data/TrailRepository` — list with filters (city, category, page)
- [ ] `TrailListScreen` — paginated list with filter chips
- [ ] `features/search/` — `SearchScreen` (debounced query → trails + cities results)
- [ ] Shared `TrailCard` widget (cover image, title, city, checkpoint count)

### 2c — Trail detail & city detail
- [ ] `TrailDetailScreen` — cover image, description, checkpoint list, static map preview
- [ ] `features/cities/data/CityRepository`
- [ ] `CityDetailScreen` — city info + trails in city

---

## Phase 3 — Walking mode (map-first)

The most complex public feature. Requires location permissions.

- [ ] `features/trails/presentation/WalkingModeScreen`
  - Full-screen `FlutterMap` centered on the trail
  - Checkpoint markers with tap-to-reveal bottom sheet
  - User location dot (permission request flow)
  - Progress bar (checkpoints visited / total)
  - Bottom sheet: checkpoint description, image, next/prev navigation
- [ ] `core/location/LocationService` — geolocator wrapper + permission handling

> **Maps decision checkpoint**: after walking mode works with `flutter_map` (OSM tiles), evaluate `google_maps_flutter` with the existing Google Places (New) API key. Google Maps gives better satellite imagery and Places autocomplete — decide before Phase 5.

---

## Phase 4 — Authenticated app (curator & team)

Mirrors the `/app` route group in the web app.

- [ ] `features/dashboard/` — post-login home: quick stats, recent activity
- [ ] `features/trails_mgmt/` — trail CRUD for curators (list, create, edit, publish/unpublish, checkpoint ordering)
- [ ] `features/cities_mgmt/` — city CRUD
- [ ] `features/team/` — team member list (read-only in v1)

---

## Phase 5 — Subscription & account

- [ ] `features/subscription/` — plan info, upgrade prompt, payment status
- [ ] `features/account/` — profile, change password, sign out
- [ ] Deep-link handling for post-payment redirect

---

## Phase 6 — Polish & platform

- [ ] Splash screen (`flutter_native_splash`) — warm base-100 background + wordmark
- [ ] Launcher icons (`flutter_launcher_icons`) from `assets/brand/favicon.png`
- [ ] Reduced-motion support (`MediaQuery.disableAnimations`)
- [ ] Accessibility pass: semantic labels, minimum 48 dp tap targets, focus order
- [ ] Dark mode testing pass

---

## Phase 7 — Release prep

- [ ] Android: signing config, `minSdkVersion 21`, Play Store metadata
- [ ] iOS: App Store metadata, entitlements, `GoogleService-Info.plist`
- [ ] CI: GitHub Actions — `flutter analyze`, `flutter test`, build artifacts
- [ ] Version strategy: semver in `pubspec.yaml`, changelog

---

## Open questions

1. **Maps provider** — `flutter_map` (OSM, free) vs `google_maps_flutter` (key exists, better data). Decide at Phase 3 checkpoint.
2. **Offline tiles** — cache trail-area tiles for spotty outdoor connectivity?
3. **Token refresh** — Dio interceptor to silently refresh on 401 before Phase 2 hits protected endpoints.
4. **Push notifications** — trail updates, team mentions? (Phase 7+)
5. **Web target** — `flutter build web`? Or keep web on Angular only?
