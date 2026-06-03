import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:sentry_flutter/sentry_flutter.dart';

import 'app.dart';
import 'core/config/env.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  // Fonts are bundled in assets/fonts/ — never fetch from the network.
  GoogleFonts.config.allowRuntimeFetching = false;

  if (Env.isProduction && Env.sentryDsn.isNotEmpty) {
    await SentryFlutter.init((options) {
      options.dsn = Env.sentryDsn;
      options.tracesSampleRate = 0.2;
    }, appRunner: () => runApp(const ProviderScope(child: App())));
  } else {
    runApp(const ProviderScope(child: App()));
  }
}
