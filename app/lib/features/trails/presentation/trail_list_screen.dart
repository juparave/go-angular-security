import 'package:flutter/material.dart';

class TrailListScreen extends StatelessWidget {
  const TrailListScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Trails')),
      body: const Center(child: Text('Trail list — Phase 2')),
    );
  }
}
