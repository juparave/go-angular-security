import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:logger/logger.dart';

import '../config/env.dart';

final loggerProvider = Provider<Logger>((ref) {
  return Logger(
    level: Env.isProduction ? Level.warning : Level.debug,
    printer: PrettyPrinter(
      methodCount: 1,
      errorMethodCount: 8,
      lineLength: 100,
      colors: true,
      printEmojis: true,
    ),
  );
});
