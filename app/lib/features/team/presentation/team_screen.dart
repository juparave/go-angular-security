import 'package:flutter/material.dart';

/// Placeholder for the team-management feature backed by `GET /api/v1/team`
/// and friends (see `server/internal/handlers/team.go`).
class TeamScreen extends StatelessWidget {
  const TeamScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Equipo')),
      body: const Center(child: Text('Gestión de equipo — pendiente')),
    );
  }
}
