import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../features/auth/application/auth_notifier.dart';
import '../../features/auth/presentation/login_screen.dart';
import '../../features/cities/presentation/city_detail_screen.dart';
import '../../features/dashboard/presentation/dashboard_screen.dart';
import '../../features/home/presentation/home_screen.dart';
import '../../features/search/presentation/search_screen.dart';
import '../../features/trails/presentation/trail_detail_screen.dart';
import '../../features/trails/presentation/trail_list_screen.dart';

final routerProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(isAuthenticatedProvider);

  return GoRouter(
    initialLocation: '/',
    redirect: (context, state) {
      if (authState.isLoading) return null;

      final isAuthenticated = authState.value ?? false;
      final isAuthRoute = state.matchedLocation.startsWith('/login');
      final isProtectedRoute = state.matchedLocation.startsWith('/app');

      if (!isAuthenticated && isProtectedRoute) return '/login';
      if (isAuthenticated && isAuthRoute) return '/app';
      return null;
    },
    routes: [
      GoRoute(path: '/', builder: (_, _) => const HomeScreen()),
      GoRoute(path: '/trails', builder: (_, _) => const TrailListScreen()),
      GoRoute(
        path: '/trails/:slug',
        builder: (_, s) => TrailDetailScreen(slug: s.pathParameters['slug']!),
      ),
      GoRoute(path: '/search', builder: (_, _) => const SearchScreen()),
      GoRoute(
        path: '/cities/:slug',
        builder: (_, s) => CityDetailScreen(slug: s.pathParameters['slug']!),
      ),
      GoRoute(path: '/login', builder: (_, _) => const LoginScreen()),
      GoRoute(path: '/app', builder: (_, _) => const DashboardScreen()),
    ],
    errorBuilder: (context, state) => Scaffold(
      body: Center(child: Text('Página no encontrada: ${state.uri}')),
    ),
  );
});
