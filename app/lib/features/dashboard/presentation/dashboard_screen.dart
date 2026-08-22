import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../auth/application/auth_notifier.dart';

class DashboardScreen extends ConsumerWidget {
  const DashboardScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final user = ref.watch(currentUserProvider).value;

    return Scaffold(
      appBar: AppBar(title: const Text('Panel')),
      body: Center(
        child: Text(
          user == null ? 'Panel' : 'Hola, ${user.displayName}',
          style: Theme.of(context).textTheme.titleLarge,
        ),
      ),
    );
  }
}
