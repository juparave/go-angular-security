import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:sentry_flutter/sentry_flutter.dart';

import 'app.dart';
import 'core/config/environment_config.dart';
import 'features/auth/data/auth_repository.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  // Fonts are bundled in assets/fonts/ — never fetch from the network.
  GoogleFonts.config.allowRuntimeFetching = false;

  const config = EnvironmentConfig();

  // The router's redirect reads AuthRepository.currentUser synchronously, so
  // the stored token has to be exchanged for a user before the first frame —
  // otherwise every cold start flashes the sign-in screen. Build the container
  // here so that await happens once, against the same providers the app uses.
  final container = ProviderContainer();
  await container.read(authRepositoryProvider).restoreSession();

  Widget app() =>
      UncontrolledProviderScope(container: container, child: const App());

  if (config.isProduction && config.sentryDsn.isNotEmpty) {
    await SentryFlutter.init((options) {
      options.dsn = config.sentryDsn;
      options.tracesSampleRate = 0.2;
    }, appRunner: () => runApp(app()));
  } else {
    runApp(app());
  }
}
