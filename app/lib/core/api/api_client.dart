import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../config/env.dart';
import '../storage/secure_storage.dart';
import 'api_exception.dart';

class ApiClient {
  ApiClient(this._dio, this._storage) {
    _dio
      ..options.baseUrl = Env.apiBaseUrl
      ..options.connectTimeout = const Duration(seconds: 15)
      ..options.receiveTimeout = const Duration(seconds: 30)
      ..options.headers['Accept'] = 'application/json'
      ..interceptors.addAll([
        _CookieAuthInterceptor(_storage),
        _RefreshInterceptor(_storage, _dio),
        _ErrorInterceptor(),
      ]);
  }

  final Dio _dio;
  final SecureStorage _storage;

  Future<Response<T>> get<T>(
    String path, {
    Map<String, dynamic>? queryParameters,
  }) =>
      _dio.get<T>(path, queryParameters: queryParameters);

  Future<Response<T>> post<T>(String path, {Object? data}) =>
      _dio.post<T>(path, data: data);

  Future<Response<T>> put<T>(String path, {Object? data}) =>
      _dio.put<T>(path, data: data);

  Future<Response<T>> patch<T>(String path, {Object? data}) =>
      _dio.patch<T>(path, data: data);

  Future<Response<T>> delete<T>(String path) => _dio.delete<T>(path);
}

/// Injects the stored access token as a Cookie header on every request.
/// The backend reads `jwt` and `refreshjwt` as HttpOnly cookies; on mobile
/// we manage them ourselves via SecureStorage.
class _CookieAuthInterceptor extends Interceptor {
  _CookieAuthInterceptor(this._storage);
  final SecureStorage _storage;

  @override
  Future<void> onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    final token = await _storage.readAccessToken();
    if (token != null && token.isNotEmpty) {
      options.headers['Cookie'] = 'jwt=$token';
    }
    handler.next(options);
  }
}

/// On a 401, silently calls POST /api/refresh-token, saves the new tokens,
/// then retries the original request once. If the refresh itself fails the
/// original 401 is forwarded to _ErrorInterceptor so AuthNotifier can sign out.
///
/// Uses a separate Dio instance to avoid re-entering this interceptor.
class _RefreshInterceptor extends Interceptor {
  _RefreshInterceptor(this._storage, this._dio);

  final SecureStorage _storage;
  final Dio _dio;

  // Prevents concurrent refreshes when multiple requests 401 simultaneously.
  bool _refreshing = false;

  @override
  Future<void> onError(
    DioException err,
    ErrorInterceptorHandler handler,
  ) async {
    final response = err.response;

    // Only intercept 401s that are not already a refresh attempt.
    if (response?.statusCode != 401 ||
        err.requestOptions.path.contains('/refresh-token') ||
        err.requestOptions.path.contains('/login') ||
        err.requestOptions.path.contains('/glogin')) {
      handler.next(err);
      return;
    }

    if (_refreshing) {
      handler.next(err);
      return;
    }

    final refreshToken = await _storage.readRefreshToken();
    if (refreshToken == null || refreshToken.isEmpty) {
      handler.next(err);
      return;
    }

    _refreshing = true;
    try {
      final refreshDio = Dio()
        ..options.baseUrl = Env.apiBaseUrl
        ..options.connectTimeout = const Duration(seconds: 15)
        ..options.receiveTimeout = const Duration(seconds: 30)
        ..options.headers['Accept'] = 'application/json';

      final refreshResponse = await refreshDio.post<Map<String, dynamic>>(
        '/api/refresh-token',
        data: {'refreshToken': refreshToken},
      );

      final body = refreshResponse.data!;
      final user = body['user'] as Map<String, dynamic>?;
      final newAccess = user?['accessToken'] as String?;
      final newRefresh = user?['refreshToken'] as String?;

      if (newAccess == null || newAccess.isEmpty) {
        handler.next(err);
        return;
      }

      await _storage.saveTokens(
        accessToken: newAccess,
        refreshToken: newRefresh ?? refreshToken,
      );

      // Retry the original request with the new token.
      final retryOptions = err.requestOptions;
      retryOptions.headers['Cookie'] = 'jwt=$newAccess';

      final retryResponse = await _dio.fetch<dynamic>(retryOptions);
      handler.resolve(retryResponse);
    } catch (_) {
      // Refresh failed — forward the original 401 to _ErrorInterceptor.
      handler.next(err);
    } finally {
      _refreshing = false;
    }
  }
}

class _ErrorInterceptor extends Interceptor {
  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    final response = err.response;

    if (err.type == DioExceptionType.connectionError ||
        err.type == DioExceptionType.receiveTimeout ||
        err.type == DioExceptionType.sendTimeout) {
      handler.reject(
        DioException(
          requestOptions: err.requestOptions,
          error: const NetworkException(),
        ),
      );
      return;
    }

    if (response != null) {
      final message = _extractMessage(response);
      if (response.statusCode == 401) {
        handler.reject(
          DioException(
            requestOptions: err.requestOptions,
            error: UnauthorizedException(message),
          ),
        );
        return;
      }
      if (response.statusCode != null && response.statusCode! >= 500) {
        handler.reject(
          DioException(
            requestOptions: err.requestOptions,
            error: const ServerException(),
          ),
        );
        return;
      }
      handler.reject(
        DioException(
          requestOptions: err.requestOptions,
          error: ClientException(message, statusCode: response.statusCode ?? 0),
        ),
      );
      return;
    }

    handler.next(err);
  }

  String _extractMessage(Response<dynamic> response) {
    try {
      final data = response.data;
      if (data is Map<String, dynamic>) {
        return (data['message'] ?? data['error'] ?? 'Error desconocido')
            .toString();
      }
    } catch (_) {}
    return 'Error desconocido';
  }
}

final apiClientProvider = Provider<ApiClient>((ref) {
  return ApiClient(Dio(), ref.watch(secureStorageProvider));
});
