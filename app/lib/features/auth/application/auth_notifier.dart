import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../../../shared/models/user.dart';
import '../data/auth_repository.dart';

part 'auth_notifier.g.dart';

/// The signed-in user, or null. Seeded with whatever `restoreSession()` already
/// put in the repository so the first build is never an empty loading state.
@Riverpod(keepAlive: true)
Stream<User?> currentUser(Ref ref) {
  // ref.watch must run during build. An `async*` body would defer it until the
  // stream is first listened to, which re-enters the provider and loops
  // forever — hence the plain function plus the helper below.
  return _withCurrentValue(ref.watch(authRepositoryProvider));
}

/// authStateChanges is a broadcast stream with no replay, and restoreSession()
/// has already run by the time the tree builds — emit the current value first
/// so [currentUserProvider] never sits in AsyncLoading waiting for a change.
Stream<User?> _withCurrentValue(AuthRepository repository) async* {
  yield repository.currentUser;
  yield* repository.authStateChanges();
}

/// Convenience for widgets that only care whether someone is signed in.
@riverpod
bool isAuthenticated(Ref ref) => ref.watch(currentUserProvider).value != null;

/// Drives the sign-in / sign-out actions and exposes their loading and error
/// states to the UI. The user itself lives in [currentUserProvider].
///
/// keepAlive because logging out unmounts the screen that triggered it: an
/// auto-dispose notifier would be torn down while `logout()` is still in
/// flight, and the trailing `state =` would throw on a disposed notifier.
@Riverpod(keepAlive: true)
class AuthController extends _$AuthController {
  @override
  FutureOr<void> build() {}

  Future<void> login({required String email, required String password}) async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(
      () => ref
          .read(authRepositoryProvider)
          .login(email: email, password: password),
    );
  }

  Future<void> googleSignIn() async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(
      () => ref.read(authRepositoryProvider).googleSignIn(),
    );
  }

  Future<void> logout() async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(
      () => ref.read(authRepositoryProvider).logout(),
    );
  }
}
