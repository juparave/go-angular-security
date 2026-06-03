import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'core/api/api_exception.dart';
import 'core/router/app_router.dart';
import 'core/theme/app_theme.dart';
import 'features/auth/application/auth_notifier.dart';

class App extends ConsumerWidget {
  const App({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.watch(routerProvider);

    // Sign out automatically when a refresh attempt fails (unrecoverable 401).
    ref.listen<AsyncValue<bool>>(isAuthenticatedProvider, (_, next) {
      next.whenData((authenticated) {
        if (!authenticated) return;
      });
    });

    // Global guard: if any feature controller surfaces an UnauthorizedException
    // after a failed refresh, clear the session and let the router redirect.
    ref.listen<AsyncValue<bool>>(isAuthenticatedProvider, (_, next) {
      if (next.hasError && next.error is UnauthorizedException) {
        ref.read(authNotifierProvider.notifier).logout();
      }
    });

    return MaterialApp.router(
      title: 'app',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.light,
      darkTheme: AppTheme.dark,
      themeMode: ThemeMode.system,
      routerConfig: router,
      localizationsDelegates: const [
        GlobalMaterialLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
      ],
      supportedLocales: const [Locale('es', 'MX'), Locale('en', 'US')],
      locale: const Locale('es', 'MX'),
    );
  }
}
