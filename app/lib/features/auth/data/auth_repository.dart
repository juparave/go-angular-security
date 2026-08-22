import 'dart:async';

import 'package:google_sign_in/google_sign_in.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../../../core/api/api_client.dart';
import '../../../core/storage/secure_storage.dart';
import '../../../shared/models/user.dart';

part 'auth_repository.g.dart';

final _googleSignIn = GoogleSignIn(scopes: ['email', 'profile']);

/// Owns the session. Holds the signed-in [User] in memory so the router's
/// `redirect` can read it **synchronously**, and broadcasts every change so
/// `refreshListenable` fires on sign-in / sign-out.
///
/// Paths are relative to `EnvironmentConfig.apiBaseUrl`, which already
/// includes the `/api/v1` prefix the Fiber server mounts its routes under.
class AuthRepository {
  AuthRepository(this._api, this._storage);

  final ApiClient _api;
  final SecureStorage _storage;

  final _authStateController = StreamController<User?>.broadcast();

  User? _currentUser;

  /// Synchronously readable — see the note in `app_router.dart`.
  User? get currentUser => _currentUser;

  Stream<User?> authStateChanges() => _authStateController.stream;

  /// Called once from `main()` before the first frame: if a token survived the
  /// last run, exchange it for the current user so the router can start on a
  /// protected route instead of flashing the sign-in screen.
  Future<User?> restoreSession() async {
    final token = await _storage.readAccessToken();
    if (token == null || token.isEmpty) return _emit(null);

    try {
      final response = await _api.get<Map<String, dynamic>>('/user');
      return _emit(User.fromJson(response.data!));
    } catch (_) {
      // Token expired and refresh failed, or the server is unreachable.
      await _storage.clearTokens();
      return _emit(null);
    }
  }

  Future<User> googleSignIn() async {
    final account = await _googleSignIn.signIn();
    if (account == null) throw Exception('Google sign-in cancelled');

    final auth = await account.authentication;
    final idToken = auth.idToken;
    if (idToken == null) throw Exception('No ID token from Google');

    return _postToBackend('/glogin', {'credential': idToken});
  }

  Future<User> login({required String email, required String password}) async {
    return _postToBackend('/login', {
      'user': {'email': email, 'password': password},
    });
  }

  Future<void> logout() async {
    try {
      await _api.post<void>('/logout');
      await _googleSignIn.signOut();
    } finally {
      await _storage.clearTokens();
      _emit(null);
    }
  }

  Future<User> _postToBackend(String path, Object data) async {
    final response = await _api.post<Map<String, dynamic>>(path, data: data);
    final body = response.data!;
    final user = User.fromJson(body['user'] as Map<String, dynamic>);

    if (user.accessToken != null && user.accessToken!.isNotEmpty) {
      await _storage.saveTokens(
        accessToken: user.accessToken!,
        refreshToken: user.refreshToken ?? '',
      );
    }

    _emit(user);
    return user;
  }

  User? _emit(User? user) {
    _currentUser = user;
    if (!_authStateController.isClosed) _authStateController.add(user);
    return user;
  }

  void dispose() => _authStateController.close();
}

@Riverpod(keepAlive: true)
AuthRepository authRepository(Ref ref) {
  final repository = AuthRepository(
    ref.watch(apiClientProvider),
    ref.watch(secureStorageProvider),
  );
  ref.onDispose(repository.dispose);
  return repository;
}
