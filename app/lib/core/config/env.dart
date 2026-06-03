/// Compile-time environment configuration.
///
/// Values are injected via `--dart-define` at build time:
///   flutter run --dart-define=API_BASE_URL=https://api.example.com
///
/// Local dev falls back to the default values below.
class Env {
  const Env._();

  static const String apiBaseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://localhost:8080',
  );

  static const bool isProduction = bool.fromEnvironment(
    'PRODUCTION',
    defaultValue: false,
  );

  static const String sentryDsn = String.fromEnvironment('SENTRY_DSN');
}
