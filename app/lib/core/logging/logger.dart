import 'package:logger/logger.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../config/environment_config.dart';

part 'logger.g.dart';

@Riverpod(keepAlive: true)
Logger logger(Ref ref) {
  return Logger(
    level: ref.watch(environmentConfigProvider).isProduction
        ? Level.warning
        : Level.debug,
    printer: PrettyPrinter(
      methodCount: 1,
      errorMethodCount: 8,
      lineLength: 100,
      colors: true,
      printEmojis: true,
    ),
  );
}
