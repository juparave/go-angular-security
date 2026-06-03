/// Typed error surfaced by all repositories.
/// Controllers wrap calls in `AsyncValue.guard(...)` — the UI never sees `DioException`.
sealed class ApiException implements Exception {
  const ApiException(this.message);
  final String message;

  @override
  String toString() => message;
}

/// HTTP 4xx client errors (validation, not-found, forbidden, etc.)
final class ClientException extends ApiException {
  const ClientException(super.message, {required this.statusCode});
  final int statusCode;
}

/// HTTP 401 — token expired or invalid credentials
final class UnauthorizedException extends ApiException {
  const UnauthorizedException([super.message = 'Sesión expirada']);
}

/// HTTP 5xx or network failures
final class ServerException extends ApiException {
  const ServerException([super.message = 'Error del servidor']);
}

/// No internet connectivity
final class NetworkException extends ApiException {
  const NetworkException([super.message = 'Sin conexión a internet']);
}
