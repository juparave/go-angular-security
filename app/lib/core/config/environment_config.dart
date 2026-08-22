import 'package:flutter/foundation.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'environment_config.g.dart';

/// Base URL for the Go Fiber backend.
///
/// Defaults to the **production** host — a release build with no defines at
/// all still points somewhere real, never at localhost. Pass
/// `--dart-define=development=true` (see `.vscode/launch.json`, "app") to point
/// at a local `go run cmd/api/main.go` instead; platform-appropriate localhost
/// addresses are used for local dev, since Android emulators can't reach
/// `localhost` directly — hence `10.0.2.2`. Pass
/// `--dart-define=remoteDebug=true` (see "Remote debug") to point at the shared
/// staging host. An explicit `--dart-define=baseUrl=https://your-host` always
/// wins over all of the above.
///
/// Platform detection uses `defaultTargetPlatform` rather than
/// `dart:io`'s `Platform`, which is unavailable when compiling for web.
class EnvironmentConfig {
  const EnvironmentConfig();

  static const String _baseUrlOverride = String.fromEnvironment('baseUrl');
  static const bool _isDevelopment = bool.fromEnvironment('development');
  static const bool _isRemoteDebug = bool.fromEnvironment('remoteDebug');
  static const String _sentryDsn = String.fromEnvironment('sentryDsn');

  // Replace these with the real hosts for your deployment.
  static const String _prodHost = 'https://app.example.com';
  static const String _stagingHost = 'https://dev.example.com';

  /// Matches `APP_PORT` in `server/.env.example`.
  static const int _devPort = 5000;

  /// Every route in `server/internal/routes/routes.go` is mounted under this
  /// prefix by `main.go` (`server.Group("/api").Group("/v1")`).
  static const String apiPrefix = '/api/v1';

  String get host {
    if (_baseUrlOverride.isNotEmpty) return _baseUrlOverride;
    if (_isRemoteDebug) return _stagingHost;
    if (!_isDevelopment) return _prodHost;
    if (!kIsWeb && defaultTargetPlatform == TargetPlatform.android) {
      return 'http://10.0.2.2:$_devPort';
    }
    return 'http://localhost:$_devPort';
  }

  /// Dio base URL — the host with the API prefix already applied, so
  /// repositories use bare paths like `/login` and `/team`.
  String get apiBaseUrl => '$host$apiPrefix';

  /// True for builds that target neither a local nor the staging backend.
  bool get isProduction => !_isDevelopment && !_isRemoteDebug;

  /// Empty unless `--dart-define=sentryDsn=...` was passed at build time.
  String get sentryDsn => _sentryDsn;
}

@Riverpod(keepAlive: true)
EnvironmentConfig environmentConfig(Ref ref) => const EnvironmentConfig();
