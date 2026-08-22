import 'dart:convert';
import 'dart:typed_data';

import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:app/core/api/api_client.dart';
import 'package:app/core/api/api_exception.dart';
import 'package:app/core/config/environment_config.dart';
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

/// Minimal [HttpClientAdapter] that answers from a caller-supplied function and
/// records every request it saw. Unlike http_mock_adapter's route handlers, the
/// function runs per request, so a call can answer 401 then 200 on the retry.
class _StubAdapter implements HttpClientAdapter {
  _StubAdapter(this._respond);

  final ResponseBody Function(RequestOptions options, int callIndex) _respond;
  final requests = <RequestOptions>[];

  /// Header snapshots taken at call time. The retry mutates the original
  /// RequestOptions in place, so [requests] holds the same instance twice —
  /// copy the headers here to see what each individual call actually sent.
  final sentHeaders = <Map<String, dynamic>>[];

  @override
  Future<ResponseBody> fetch(
    RequestOptions options,
    Stream<Uint8List>? requestStream,
    Future<void>? cancelFuture,
  ) async {
    final index = requests.length;
    requests.add(options);
    sentHeaders.add(Map<String, dynamic>.from(options.headers));
    return _respond(options, index);
  }

  @override
  void close({bool force = false}) {}
}

ResponseBody _json(int status, Map<String, dynamic> body) =>
    ResponseBody.fromString(
      jsonEncode(body),
      status,
      headers: {
        Headers.contentTypeHeader: [Headers.jsonContentType],
      },
    );

class _Harness {
  _Harness(this.client, this.api, this.refresh, this.storage);

  final ApiClient client;

  /// Requests that went through the main Dio (business endpoints + retries).
  final _StubAdapter api;

  /// Requests that went through the refresh Dio (POST /refresh-token only).
  final _StubAdapter refresh;

  final _FakeSecureStorage storage;
}

/// Builds a real [ApiClient] with both of its Dio instances stubbed, so the
/// production interceptors are the code actually under test.
_Harness _buildHarness({
  required ResponseBody Function(RequestOptions options, int callIndex) api,
  ResponseBody Function(RequestOptions options, int callIndex)? refresh,
}) {
  final storage = _FakeSecureStorage();
  final apiAdapter = _StubAdapter(api);
  final refreshAdapter = _StubAdapter(
    refresh ??
        (options, _) => throw StateError(
          'refresh endpoint should not have been called: ${options.path}',
        ),
  );

  final dio = Dio()..httpClientAdapter = apiAdapter;
  final refreshDio = Dio()..httpClientAdapter = refreshAdapter;

  final client = ApiClient(
    dio,
    storage,
    const EnvironmentConfig(),
    refreshDio: refreshDio,
  );

  return _Harness(client, apiAdapter, refreshAdapter, storage);
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

void main() {
  test('ApiClient targets the Fiber /api/v1 group', () {
    expect(
      const EnvironmentConfig().apiBaseUrl,
      endsWith(EnvironmentConfig.apiPrefix),
    );
  });

  group('_RefreshInterceptor', () {
    test('retries the original request after a successful refresh', () async {
      final h = _buildHarness(
        api: (options, i) => i == 0
            ? _json(401, {'message': 'Unauthorized'})
            : _json(200, {'ok': true}),
        refresh: (options, _) => _json(200, {
          'message': 'success',
          'user': {
            'id': '1',
            'accessToken': 'new-access-token',
            'refreshToken': 'new-refresh-token',
          },
        }),
      );
      h.storage.accessToken = 'expired-token';
      h.storage.refreshToken = 'valid-refresh';

      final response = await h.client.get<Map<String, dynamic>>('/data');

      expect(response.statusCode, 200);
      expect(h.api.requests, hasLength(2), reason: 'retried exactly once');
      expect(h.refresh.requests.single.path, '/refresh-token');
      expect(h.refresh.requests.single.data, {
        'refreshToken': 'valid-refresh',
      });
      expect(h.storage.accessToken, 'new-access-token');
      expect(h.storage.refreshToken, 'new-refresh-token');
    });

    test('sends the refreshed token as the jwt cookie on the retry', () async {
      final h = _buildHarness(
        api: (options, i) => i == 0
            ? _json(401, {'message': 'Unauthorized'})
            : _json(200, {'ok': true}),
        refresh: (options, _) => _json(200, {
          'user': {'id': '1', 'accessToken': 'fresh-token'},
        }),
      );
      h.storage.accessToken = 'expired-token';
      h.storage.refreshToken = 'valid-refresh';

      await h.client.get<Map<String, dynamic>>('/data');

      expect(h.api.sentHeaders.first['Cookie'], 'jwt=expired-token');
      expect(h.api.sentHeaders.last['Cookie'], 'jwt=fresh-token');
    });

    test('keeps the old refresh token when the response omits a new one',
        () async {
      final h = _buildHarness(
        api: (options, i) => i == 0
            ? _json(401, {'message': 'Unauthorized'})
            : _json(200, {'ok': true}),
        refresh: (options, _) => _json(200, {
          'user': {'id': '1', 'accessToken': 'fresh-token'},
        }),
      );
      h.storage.accessToken = 'expired-token';
      h.storage.refreshToken = 'valid-refresh';

      await h.client.get<Map<String, dynamic>>('/data');

      expect(h.storage.refreshToken, 'valid-refresh');
    });

    test('raises UnauthorizedException when no refresh token is stored',
        () async {
      final h = _buildHarness(
        api: (options, i) => _json(401, {'message': 'Unauthorized'}),
      );
      h.storage.accessToken = 'expired-token';
      h.storage.refreshToken = null;

      await expectLater(
        () => h.client.get<Map<String, dynamic>>('/data'),
        throwsA(
          isA<DioException>().having(
            (e) => e.error,
            'error',
            isA<UnauthorizedException>(),
          ),
        ),
      );
      expect(h.refresh.requests, isEmpty);
      expect(h.api.requests, hasLength(1), reason: 'no retry without a token');
    });

    test('raises UnauthorizedException when the refresh call itself fails',
        () async {
      final h = _buildHarness(
        api: (options, i) => _json(401, {'message': 'Unauthorized'}),
        refresh: (options, _) =>
            _json(401, {'message': 'Refresh token expired'}),
      );
      h.storage.accessToken = 'expired-token';
      h.storage.refreshToken = 'stale-refresh';

      await expectLater(
        () => h.client.get<Map<String, dynamic>>('/data'),
        throwsA(
          isA<DioException>().having(
            (e) => e.error,
            'error',
            isA<UnauthorizedException>(),
          ),
        ),
      );
    });

    for (final path in ['/login', '/glogin', '/refresh-token']) {
      test('does not attempt a refresh for a 401 on $path', () async {
        // The refresh adapter throws if touched.
        final h = _buildHarness(
          api: (options, i) => _json(401, {'message': 'Invalid credentials'}),
        );
        h.storage.accessToken = 'expired-token';
        h.storage.refreshToken = 'valid-refresh';

        await expectLater(
          () => h.client.post<Map<String, dynamic>>(path, data: const {}),
          throwsA(
            isA<DioException>().having(
              (e) => e.error,
              'error',
              isA<UnauthorizedException>(),
            ),
          ),
        );
        expect(h.refresh.requests, isEmpty);
        expect(h.api.requests, hasLength(1));
      });
    }
  });

  group('_ErrorInterceptor', () {
    test('maps 5xx to ServerException', () async {
      final h = _buildHarness(api: (options, i) => _json(500, {'m': 'boom'}));

      await expectLater(
        () => h.client.get<Map<String, dynamic>>('/data'),
        throwsA(
          isA<DioException>().having(
            (e) => e.error,
            'error',
            isA<ServerException>(),
          ),
        ),
      );
    });

    test('maps other 4xx to ClientException with the server message', () async {
      final h = _buildHarness(
        api: (options, i) => _json(422, {'message': 'Email ya registrado'}),
      );

      await expectLater(
        () => h.client.get<Map<String, dynamic>>('/data'),
        throwsA(
          isA<DioException>().having(
            (e) => e.error,
            'error',
            isA<ClientException>()
                .having((e) => e.statusCode, 'statusCode', 422)
                .having((e) => e.message, 'message', 'Email ya registrado'),
          ),
        ),
      );
    });
  });
}
