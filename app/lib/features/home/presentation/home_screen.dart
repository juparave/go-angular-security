import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(
        title: Image.asset('assets/brand/logotipo.png', height: 28),
        actions: [
          IconButton(
            icon: const Icon(Icons.search),
            onPressed: () => context.go('/search'),
          ),
          IconButton(
            icon: const Icon(Icons.person_outline),
            onPressed: () => context.go('/login'),
          ),
        ],
      ),
      body: Center(
        child: Text(
          'Descubre el México\nque los turistas nunca ven',
          textAlign: TextAlign.center,
          style: theme.textTheme.headlineMedium,
        ),
      ),
    );
  }
}
