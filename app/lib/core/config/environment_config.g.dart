// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'environment_config.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning

@ProviderFor(environmentConfig)
final environmentConfigProvider = EnvironmentConfigProvider._();

final class EnvironmentConfigProvider
    extends
        $FunctionalProvider<
          EnvironmentConfig,
          EnvironmentConfig,
          EnvironmentConfig
        >
    with $Provider<EnvironmentConfig> {
  EnvironmentConfigProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'environmentConfigProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$environmentConfigHash();

  @$internal
  @override
  $ProviderElement<EnvironmentConfig> $createElement(
    $ProviderPointer pointer,
  ) => $ProviderElement(pointer);

  @override
  EnvironmentConfig create(Ref ref) {
    return environmentConfig(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(EnvironmentConfig value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<EnvironmentConfig>(value),
    );
  }
}

String _$environmentConfigHash() => r'92eee0b9e35a35808e31bf06533a702ad2eee2a5';
