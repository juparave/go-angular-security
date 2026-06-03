import 'package:flutter/material.dart';

class CityDetailScreen extends StatelessWidget {
  const CityDetailScreen({super.key, required this.slug});
  final String slug;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(slug)),
      body: const Center(child: Text('City detail — Phase 2')),
    );
  }
}
