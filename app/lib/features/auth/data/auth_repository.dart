import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:google_sign_in/google_sign_in.dart';

import '../../../core/api/api_client.dart';
import '../../../core/storage/secure_storage.dart';
import '../../../shared/models/user.dart';

final _googleSignIn = GoogleSignIn(scopes: ['email', 'profile']);

class AuthRepository {
  AuthRepository(this._api, this._storage);

  final ApiClient _api;
  final SecureStorage _storage;

  Future<bool> hasValidToken() async {
    final token = await _storage.readAccessToken();
    return token != null && token.isNotEmpty;
  }

  Future<User> googleSignIn() async {
    final account = await _googleSignIn.signIn();
    if (account == null) throw Exception('Google sign-in cancelled');

    final auth = await account.authentication;
    final idToken = auth.idToken;
    if (idToken == null) throw Exception('No ID token from Google');

    return _postToBackend('/api/glogin', {'credential': idToken});
  }

  Future<User> login({required String email, required String password}) async {
    return _postToBackend('/api/login', {
      'user': {'email': email, 'password': password},
    });
  }

  Future<void> logout() async {
    try {
      await _api.post<void>('/api/logout');
      await _googleSignIn.signOut();
    } finally {
      await _storage.clearTokens();
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

    return user;
  }
}

final authRepositoryProvider = Provider<AuthRepository>((ref) {
  return AuthRepository(
    ref.watch(apiClientProvider),
    ref.watch(secureStorageProvider),
  );
});
