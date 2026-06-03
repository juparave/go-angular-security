import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http_mock_adapter/http_mock_adapter.dart';
import 'package:mocktail/mocktail.dart';

import 'package:app/core/api/api_client.dart';
import 'package:app/core/api/api_exception.dart';
import 'package:app/core/storage/secure_storage.dart';

// ---------------------------------------------------------------------------
// Fakes
// ---------------------------------------------------------------------------

class _FakeSecureStorage extends Fake implements SecureStorage {
  String? accessToken;
  String? refreshToken;

  @override
  Future<String?> readAccessToken() async => accessToken;

  @override
  Future<String?> readRefreshToken() async => refreshToken;

  @override
  Future<void> saveTokens({
    required String accessToken,
    required String refreshToken,
  }) async {
    this.accessToken = accessToken;
    this.refreshToken = refreshToken;
  }

  @override
  Future<void> clearTokens() async {
    accessToken = null;
    refreshToken = null;
  }
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

/// Builds a Dio + DioAdapter pair wired up with the same interceptors as
/// [ApiClient] but with the adapter injected so we can fake HTTP.
///
/// Returns the (Dio, DioAdapter, _FakeSecureStorage) triple.
(Dio, DioAdapter, _FakeSecureStorage) _buildTestClient() {
  final storage = _FakeSecureStorage();
  final dio = Dio(
    BaseOptions(
      baseUrl: 'http://localhost',
      validateStatus: (_) => true, // let interceptors handle status codes
    ),
  );
  final adapter = DioAdapter(dio: dio, matcher: const FullHttpRequestMatcher());

  dio.interceptors.addAll([
    _TestCookieInterceptor(storage),
    _TestRefreshInterceptor(storage, dio),
    _TestErrorInterceptor(),
  ]);

  return (dio, adapter, storage);
}

/// Minimal cookie interceptor mirroring production behaviour.
class _TestCookieInterceptor extends Interceptor {
  _TestCookieInterceptor(this._storage);
  final _FakeSecureStorage _storage;

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

/// Re-implements _RefreshInterceptor using the same logic but with an
/// injectable refresh Dio so we can mock both the original and refresh calls.
class _TestRefreshInterceptor extends Interceptor {
  _TestRefreshInterceptor(this._storage, this._dio);
  final _FakeSecureStorage _storage;
  final Dio _dio;
  bool _refreshing = false;

  @override
  Future<void> onError(
    DioException err,
    ErrorInterceptorHandler handler,
  ) async {
    final response = err.response;

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
      // Use the same Dio (adapter already mocked) for the refresh call.
      final refreshResponse = await _dio.post<Map<String, dynamic>>(
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

      final retryOptions = err.requestOptions;
      retryOptions.headers['Cookie'] = 'jwt=$newAccess';
      final retryResponse = await _dio.fetch<dynamic>(retryOptions);
      handler.resolve(retryResponse);
    } catch (_) {
      handler.next(err);
    } finally {
      _refreshing = false;
    }
  }
}

class _TestErrorInterceptor extends Interceptor {
  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    final response = err.response;
    if (response?.statusCode == 401) {
      handler.reject(
        DioException(
          requestOptions: err.requestOptions,
          error: const UnauthorizedException(),
        ),
      );
      return;
    }
    handler.next(err);
  }
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

void main() {
  group('_RefreshInterceptor', () {
    test('retries original request after successful token refresh', () async {
      final (dio, adapter, storage) = _buildTestClient();

      storage.accessToken = 'expired-token';
      storage.refreshToken = 'valid-refresh';

      // First call to /api/data returns 401.
      adapter.onGet(
        '/api/data',
        (server) => server.reply(401, {'message': 'Unauthorized'}),
        headers: {'Cookie': 'jwt=expired-token'},
      );

      // Refresh endpoint returns new tokens.
      adapter.onPost(
        '/api/refresh-token',
        (server) => server.reply(200, {
          'message': 'success',
          'user': {
            'id': '1',
            'accessToken': 'new-access-token',
            'refreshToken': 'new-refresh-token',
          },
        }),
        data: {'refreshToken': 'valid-refresh'},
      );

      // Retry with new token succeeds.
      adapter.onGet(
        '/api/data',
        (server) => server.reply(200, {'result': 'ok'}),
        headers: {'Cookie': 'jwt=new-access-token'},
      );

      final response = await dio.get<Map<String, dynamic>>('/api/data');

      expect(response.statusCode, 200);
      expect(response.data, {'result': 'ok'});
      expect(storage.accessToken, 'new-access-token');
      expect(storage.refreshToken, 'new-refresh-token');
    });

    test('raises UnauthorizedException when no refresh token is stored',
        () async {
      final (dio, adapter, storage) = _buildTestClient();

      storage.accessToken = 'expired-token';
      storage.refreshToken = null; // nothing to refresh with

      adapter.onGet(
        '/api/data',
        (server) => server.reply(401, {'message': 'Unauthorized'}),
      );

      expect(
        () => dio.get<dynamic>('/api/data'),
        throwsA(
          isA<DioException>().having(
            (e) => e.error,
            'error',
            isA<UnauthorizedException>(),
          ),
        ),
      );
    });

    test('raises UnauthorizedException when refresh endpoint itself returns 401',
        () async {
      final (dio, adapter, storage) = _buildTestClient();

      storage.accessToken = 'expired-token';
      storage.refreshToken = 'also-expired-refresh';

      adapter.onGet(
        '/api/data',
        (server) => server.reply(401, {'message': 'Unauthorized'}),
      );

      adapter.onPost(
        '/api/refresh-token',
        (server) => server.reply(401, {'message': 'refresh token expired'}),
        data: {'refreshToken': 'also-expired-refresh'},
      );

      expect(
        () => dio.get<dynamic>('/api/data'),
        throwsA(
          isA<DioException>().having(
            (e) => e.error,
            'error',
            isA<UnauthorizedException>(),
          ),
        ),
      );
    });

    test('does not intercept 401 on /login', () async {
      final (dio, adapter, storage) = _buildTestClient();

      storage.refreshToken = 'some-refresh';

      adapter.onPost(
        '/api/login',
        (server) => server.reply(401, {'message': 'bad credentials'}),
      );

      expect(
        () => dio.post<dynamic>('/api/login', data: {}),
        throwsA(
          isA<DioException>().having(
            (e) => e.error,
            'error',
            isA<UnauthorizedException>(),
          ),
        ),
      );
    });

    test('does not intercept 401 on /glogin', () async {
      final (dio, adapter, storage) = _buildTestClient();

      storage.refreshToken = 'some-refresh';

      adapter.onPost(
        '/api/glogin',
        (server) => server.reply(401, {'message': 'invalid credential'}),
      );

      expect(
        () => dio.post<dynamic>('/api/glogin', data: {}),
        throwsA(
          isA<DioException>().having(
            (e) => e.error,
            'error',
            isA<UnauthorizedException>(),
          ),
        ),
      );
    });

    test('injects new access token as Cookie header on retry', () async {
      final (dio, adapter, storage) = _buildTestClient();

      storage.accessToken = 'old-token';
      storage.refreshToken = 'good-refresh';

      String? cookieOnRetry;

      adapter.onGet(
        '/api/data',
        (server) => server.reply(401, {'message': 'Unauthorized'}),
        headers: {'Cookie': 'jwt=old-token'},
      );

      adapter.onPost(
        '/api/refresh-token',
        (server) => server.reply(200, {
          'message': 'success',
          'user': {
            'id': '1',
            'accessToken': 'fresh-token',
            'refreshToken': 'fresh-refresh',
          },
        }),
        data: {'refreshToken': 'good-refresh'},
      );

      // Capture the cookie on the retry.
      dio.interceptors.add(
        InterceptorsWrapper(
          onRequest: (options, handler) {
            if (options.path == '/api/data' &&
                options.headers['Cookie'] != 'jwt=old-token') {
              cookieOnRetry = options.headers['Cookie'] as String?;
            }
            handler.next(options);
          },
        ),
      );

      adapter.onGet(
        '/api/data',
        (server) => server.reply(200, {'ok': true}),
        headers: {'Cookie': 'jwt=fresh-token'},
      );

      await dio.get<dynamic>('/api/data');

      expect(cookieOnRetry, 'jwt=fresh-token');
    });
  });
}
