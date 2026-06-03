import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../shared/models/user.dart';
import '../data/auth_repository.dart';

/// Null means unauthenticated; non-null is the signed-in user.
class AuthNotifier extends AsyncNotifier<User?> {
  @override
  Future<User?> build() async {
    final hasToken = await ref.watch(authRepositoryProvider).hasValidToken();
    return hasToken
        ? null
        : null; // user object loaded lazily on first protected request
  }

  Future<void> googleSignIn() async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(
      () => ref.read(authRepositoryProvider).googleSignIn(),
    );
  }

  Future<void> login({required String email, required String password}) async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(
      () => ref
          .read(authRepositoryProvider)
          .login(email: email, password: password),
    );
  }

  Future<void> logout() async {
    await ref.read(authRepositoryProvider).logout();
    state = const AsyncData(null);
  }
}

final authNotifierProvider = AsyncNotifierProvider<AuthNotifier, User?>(
  AuthNotifier.new,
);

/// Convenience: true when user is signed in (token present).
final isAuthenticatedProvider = FutureProvider<bool>((ref) async {
  final user = await ref.watch(authNotifierProvider.future);
  if (user != null) return true;
  return ref.watch(authRepositoryProvider).hasValidToken();
});
