import 'dart:async';

import 'package:flutter/foundation.dart';

/// Bridges a Stream to a Listenable so GoRouter's `refreshListenable` can
/// react to auth-state changes (see the go_router migration guide).
class GoRouterRefreshStream extends ChangeNotifier {
  GoRouterRefreshStream(Stream<dynamic> stream) {
    notifyListeners();
    _subscription = stream.asBroadcastStream().listen(
      (dynamic _) => notifyListeners(),
    );
  }

  late final StreamSubscription<dynamic> _subscription;

  @override
  void dispose() {
    _subscription.cancel();
    super.dispose();
  }
}
