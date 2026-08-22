import 'package:flutter/foundation.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';

import 'package:app/main.dart' as app;

/// Drives the app against a real, already-running backend:
///
///   cd server && go run cmd/api/main.go     # listens on :5000
///   flutter test integration_test/sign_in_flow_test.dart \
///     --dart-define=development=true \
///     --dart-define=e2eEmail=you@example.com \
///     --dart-define=e2ePassword=secret
///
/// Skipped unless the credential defines are supplied, so `flutter test` on a
/// machine with no backend stays green.
void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  const email = String.fromEnvironment('e2eEmail');
  const password = String.fromEnvironment('e2ePassword');

  testWidgets('signing in lands on the dashboard', (tester) async {
    await app.main();
    await tester.pumpAndSettle();

    // Cold start with no stored token redirects to /sign-in.
    expect(find.byKey(const Key('showEmailFormButton')), findsOneWidget);

    await tester.tap(find.byKey(const Key('showEmailFormButton')));
    await tester.pumpAndSettle();

    await tester.enterText(find.byKey(const Key('signInEmailField')), email);
    await tester.enterText(
      find.byKey(const Key('signInPasswordField')),
      password,
    );
    await tester.pumpAndSettle();

    await tester.tap(find.byKey(const Key('signInSubmitButton')));

    // Hits the real backend (password verify + JWT issue), so give it real
    // time rather than a single pump.
    await tester.pumpAndSettle(const Duration(seconds: 10));

    expect(find.text('Panel'), findsWidgets);
  }, skip: email.isEmpty || password.isEmpty);
}
