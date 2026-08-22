import 'package:flutter_test/flutter_test.dart';

import 'package:app/core/config/environment_config.dart';

void main() {
  group('EnvironmentConfig', () {
    // These run with no --dart-define, i.e. the release-build defaults.
    test('defaults to the production host, never localhost', () {
      const config = EnvironmentConfig();
      expect(config.isProduction, isTrue);
      expect(config.host, startsWith('https://'));
      expect(config.host, isNot(contains('localhost')));
      expect(config.host, isNot(contains('10.0.2.2')));
    });

    test('apiBaseUrl appends the Fiber /api/v1 group to the host', () {
      const config = EnvironmentConfig();
      expect(config.apiBaseUrl, '${config.host}${EnvironmentConfig.apiPrefix}');
      expect(EnvironmentConfig.apiPrefix, '/api/v1');
    });

    test('sentryDsn is empty unless injected at build time', () {
      expect(const EnvironmentConfig().sentryDsn, isEmpty);
    });
  });
}
