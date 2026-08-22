import 'package:go_router/go_router.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../../features/auth/data/auth_repository.dart';
import '../../features/auth/presentation/login_screen.dart';
import '../../features/auth/presentation/profile_screen.dart';
import '../../features/dashboard/presentation/dashboard_screen.dart';
import '../../features/home/presentation/app_shell.dart';
import '../../features/team/presentation/team_screen.dart';
import 'go_router_refresh_stream.dart';

part 'app_router.g.dart';

enum AppRoute { signIn, dashboard, team, profile }

@Riverpod(keepAlive: true)
GoRouter goRouter(Ref ref) {
  final authRepository = ref.watch(authRepositoryProvider);

  return GoRouter(
    initialLocation: '/dashboard',
    debugLogDiagnostics: true,
    redirect: (context, state) {
      // Read the repository directly rather than currentUserProvider:
      // refreshListenable fires synchronously off authStateChanges, while the
      // stream-backed provider rebuilds a microtask later, which would leave a
      // just-signed-in user stuck on the sign-in route for one frame.
      final isLoggedIn = authRepository.currentUser != null;
      final isAuthRoute = state.matchedLocation == '/sign-in';

      if (!isLoggedIn) return isAuthRoute ? null : '/sign-in';
      if (isAuthRoute) return '/dashboard';
      return null;
    },
    refreshListenable: GoRouterRefreshStream(authRepository.authStateChanges()),
    routes: [
      GoRoute(
        path: '/sign-in',
        name: AppRoute.signIn.name,
        builder: (context, state) => const LoginScreen(),
      ),
      StatefulShellRoute.indexedStack(
        builder: (context, state, navigationShell) =>
            AppShell(navigationShell: navigationShell),
        branches: [
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: '/dashboard',
                name: AppRoute.dashboard.name,
                builder: (context, state) => const DashboardScreen(),
              ),
            ],
          ),
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: '/team',
                name: AppRoute.team.name,
                builder: (context, state) => const TeamScreen(),
              ),
            ],
          ),
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: '/profile',
                name: AppRoute.profile.name,
                builder: (context, state) => const ProfileScreen(),
              ),
            ],
          ),
        ],
      ),
    ],
  );
}
