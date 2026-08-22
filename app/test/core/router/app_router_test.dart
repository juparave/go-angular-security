import 'dart:convert';
import 'dart:typed_data';

import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:google_fonts/google_fonts.dart';

import 'package:app/app.dart';
import 'package:app/core/api/api_client.dart';
import 'package:app/core/config/environment_config.dart';
import 'package:app/core/storage/secure_storage.dart';
import 'package:app/features/auth/application/auth_notifier.dart';
import 'package:app/features/auth/data/auth_repository.dart';

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

class _StubAdapter implements HttpClientAdapter {
  _StubAdapter(this._respond);
  final ResponseBody Function(RequestOptions options) _respond;

  @override
  Future<ResponseBody> fetch(
    RequestOptions options,
    Stream<Uint8List>? requestStream,
    Future<void>? cancelFuture,
  ) async => _respond(options);

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

const _userJson = {
  'id': '1',
  'email': 'ana@example.com',
  'firstName': 'Ana',
  'lastName': 'García',
  'roleId': 1,
  'enabled': true,
  'accessToken': 'access-token',
  'refreshToken': 'refresh-token',
};

ProviderContainer _container({
  String? storedToken,
  required ResponseBody Function(RequestOptions options) respond,
}) {
  final storage = _FakeSecureStorage()..accessToken = storedToken;
  final dio = Dio()..httpClientAdapter = _StubAdapter(respond);

  return ProviderContainer(
    overrides: [
      secureStorageProvider.overrideWithValue(storage),
      apiClientProvider.overrideWithValue(
        ApiClient(dio, storage, const EnvironmentConfig(), refreshDio: Dio()),
      ),
    ],
  );
}

void main() {
  setUpAll(() => GoogleFonts.config.allowRuntimeFetching = false);

  // Every helper below awaits a real Dio round-trip through a stub adapter.
  // testWidgets runs its body under a fake clock that never advances the timers
  // Dio schedules, so those awaits must go through tester.runAsync.

  testWidgets('cold start with no stored token lands on sign-in', (
    tester,
  ) async {
    final container = _container(
      respond: (options) => _json(401, {'message': 'Unauthorized'}),
    );
    addTearDown(container.dispose);

    await tester.runAsync(
      () => container.read(authRepositoryProvider).restoreSession(),
    );
    await tester.pumpWidget(
      UncontrolledProviderScope(container: container, child: const App()),
    );
    await tester.pumpAndSettle();

    expect(find.byKey(const Key('showEmailFormButton')), findsOneWidget);
    expect(find.text('Panel'), findsNothing);
  });

  testWidgets('a stored token is restored before the first frame', (
    tester,
  ) async {
    final container = _container(
      storedToken: 'valid-token',
      respond: (options) => _json(200, _userJson),
    );
    addTearDown(container.dispose);

    await tester.runAsync(
      () => container.read(authRepositoryProvider).restoreSession(),
    );
    await tester.pumpWidget(
      UncontrolledProviderScope(container: container, child: const App()),
    );
    await tester.pumpAndSettle();

    // Never flashes /sign-in — the very first frame is already the dashboard.
    expect(find.byKey(const Key('showEmailFormButton')), findsNothing);
    expect(find.text('Hola, Ana García'), findsOneWidget);
  });

  testWidgets('signing in redirects to the dashboard', (tester) async {
    final container = _container(
      respond: (options) => options.path == '/login'
          ? _json(200, {'user': _userJson})
          : _json(401, {'message': 'Unauthorized'}),
    );
    addTearDown(container.dispose);

    await tester.runAsync(
      () => container.read(authRepositoryProvider).restoreSession(),
    );
    await tester.pumpWidget(
      UncontrolledProviderScope(container: container, child: const App()),
    );
    await tester.pumpAndSettle();
    expect(find.byKey(const Key('showEmailFormButton')), findsOneWidget);

    await tester.runAsync(
      () => container
          .read(authControllerProvider.notifier)
          .login(email: 'ana@example.com', password: 'secret123'),
    );
    await tester.pumpAndSettle();

    expect(find.text('Hola, Ana García'), findsOneWidget);
  });

  testWidgets('signing out redirects back to sign-in', (tester) async {
    final container = _container(
      storedToken: 'valid-token',
      respond: (options) => options.path == '/logout'
          ? _json(200, {'message': 'success'})
          : _json(200, _userJson),
    );
    addTearDown(container.dispose);

    await tester.runAsync(
      () => container.read(authRepositoryProvider).restoreSession(),
    );
    await tester.pumpWidget(
      UncontrolledProviderScope(container: container, child: const App()),
    );
    await tester.pumpAndSettle();
    expect(find.text('Hola, Ana García'), findsOneWidget);

    await tester.runAsync(
      () => container.read(authControllerProvider.notifier).logout(),
    );
    await tester.pumpAndSettle();

    expect(find.byKey(const Key('showEmailFormButton')), findsOneWidget);
  });
}
