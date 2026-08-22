// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'auth_notifier.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning
/// The signed-in user, or null. Seeded with whatever `restoreSession()` already
/// put in the repository so the first build is never an empty loading state.

@ProviderFor(currentUser)
final currentUserProvider = CurrentUserProvider._();

/// The signed-in user, or null. Seeded with whatever `restoreSession()` already
/// put in the repository so the first build is never an empty loading state.

final class CurrentUserProvider
    extends $FunctionalProvider<AsyncValue<User?>, User?, Stream<User?>>
    with $FutureModifier<User?>, $StreamProvider<User?> {
  /// The signed-in user, or null. Seeded with whatever `restoreSession()` already
  /// put in the repository so the first build is never an empty loading state.
  CurrentUserProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'currentUserProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$currentUserHash();

  @$internal
  @override
  $StreamProviderElement<User?> $createElement($ProviderPointer pointer) =>
      $StreamProviderElement(pointer);

  @override
  Stream<User?> create(Ref ref) {
    return currentUser(ref);
  }
}

String _$currentUserHash() => r'3ae5b66d1a99a118dc8732e019543cd1eec2b1b8';

/// Convenience for widgets that only care whether someone is signed in.

@ProviderFor(isAuthenticated)
final isAuthenticatedProvider = IsAuthenticatedProvider._();

/// Convenience for widgets that only care whether someone is signed in.

final class IsAuthenticatedProvider
    extends $FunctionalProvider<bool, bool, bool>
    with $Provider<bool> {
  /// Convenience for widgets that only care whether someone is signed in.
  IsAuthenticatedProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'isAuthenticatedProvider',
        isAutoDispose: true,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$isAuthenticatedHash();

  @$internal
  @override
  $ProviderElement<bool> $createElement($ProviderPointer pointer) =>
      $ProviderElement(pointer);

  @override
  bool create(Ref ref) {
    return isAuthenticated(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(bool value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<bool>(value),
    );
  }
}

String _$isAuthenticatedHash() => r'7695b0d1903ee386f1ee989ed07de5f0fff7cee0';

/// Drives the sign-in / sign-out actions and exposes their loading and error
/// states to the UI. The user itself lives in [currentUserProvider].
///
/// keepAlive because logging out unmounts the screen that triggered it: an
/// auto-dispose notifier would be torn down while `logout()` is still in
/// flight, and the trailing `state =` would throw on a disposed notifier.

@ProviderFor(AuthController)
final authControllerProvider = AuthControllerProvider._();

/// Drives the sign-in / sign-out actions and exposes their loading and error
/// states to the UI. The user itself lives in [currentUserProvider].
///
/// keepAlive because logging out unmounts the screen that triggered it: an
/// auto-dispose notifier would be torn down while `logout()` is still in
/// flight, and the trailing `state =` would throw on a disposed notifier.
final class AuthControllerProvider
    extends $AsyncNotifierProvider<AuthController, void> {
  /// Drives the sign-in / sign-out actions and exposes their loading and error
  /// states to the UI. The user itself lives in [currentUserProvider].
  ///
  /// keepAlive because logging out unmounts the screen that triggered it: an
  /// auto-dispose notifier would be torn down while `logout()` is still in
  /// flight, and the trailing `state =` would throw on a disposed notifier.
  AuthControllerProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'authControllerProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$authControllerHash();

  @$internal
  @override
  AuthController create() => AuthController();
}

String _$authControllerHash() => r'14a7ebabbd9e37295ec19467fc302153281ba315';

/// Drives the sign-in / sign-out actions and exposes their loading and error
/// states to the UI. The user itself lives in [currentUserProvider].
///
/// keepAlive because logging out unmounts the screen that triggered it: an
/// auto-dispose notifier would be torn down while `logout()` is still in
/// flight, and the trailing `state =` would throw on a disposed notifier.

abstract class _$AuthController extends $AsyncNotifier<void> {
  FutureOr<void> build();
  @$mustCallSuper
  @override
  WhenComplete runBuild() {
    final ref = this.ref as $Ref<AsyncValue<void>, void>;
    final element =
        ref.element
            as $ClassProviderElement<
              AnyNotifier<AsyncValue<void>, void>,
              AsyncValue<void>,
              Object?,
              Object?
            >;
    return element.handleCreate(ref, build);
  }
}
