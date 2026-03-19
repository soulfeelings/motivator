import 'package:flutter/material.dart';
import '../services/api_service.dart';

class DashboardScreen extends StatefulWidget {
  const DashboardScreen({super.key});

  @override
  State<DashboardScreen> createState() => _DashboardScreenState();
}

class _DashboardScreenState extends State<DashboardScreen> {
  Map<String, dynamic>? _profile;
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final me = await ApiService.get('/me');
      final memberships = me['memberships'] as List? ?? [];
      if (memberships.isNotEmpty) {
        final m = memberships[0];
        setState(() => _profile = m);
      }
    } catch (_) {}
    if (mounted) setState(() => _loading = false);
  }

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: _loading
          ? const Center(child: CircularProgressIndicator())
          : RefreshIndicator(
              onRefresh: _load,
              child: ListView(
                padding: const EdgeInsets.all(20),
                children: [
                  const Text('Welcome back!',
                      style: TextStyle(fontSize: 28, fontWeight: FontWeight.bold, color: Colors.white)),
                  const SizedBox(height: 4),
                  Text('Keep pushing forward',
                      style: TextStyle(color: Colors.grey[500], fontSize: 14)),
                  const SizedBox(height: 24),
                  if (_profile != null) ...[
                    _StatCard(
                      label: 'Level',
                      value: '${_profile!['level'] ?? 1}',
                      icon: Icons.military_tech_rounded,
                      color: const Color(0xFF8B5CF6),
                    ),
                    const SizedBox(height: 12),
                    Row(
                      children: [
                        Expanded(
                          child: _StatCard(
                            label: 'XP',
                            value: '${_profile!['xp'] ?? 0}',
                            icon: Icons.bolt_rounded,
                            color: const Color(0xFF10B981),
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: _StatCard(
                            label: 'Coins',
                            value: '${_profile!['coins'] ?? 0}',
                            icon: Icons.monetization_on_rounded,
                            color: const Color(0xFFF59E0B),
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 12),
                    _StatCard(
                      label: 'Role',
                      value: (_profile!['role'] ?? 'employee').toString().toUpperCase(),
                      icon: Icons.shield_rounded,
                      color: const Color(0xFF6366F1),
                    ),
                  ] else
                    Card(
                      child: Padding(
                        padding: const EdgeInsets.all(24),
                        child: Text('No company membership yet.',
                            style: TextStyle(color: Colors.grey[500])),
                      ),
                    ),
                ],
              ),
            ),
    );
  }
}

class _StatCard extends StatelessWidget {
  final String label;
  final String value;
  final IconData icon;
  final Color color;

  const _StatCard({required this.label, required this.value, required this.icon, required this.color});

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Row(
          children: [
            Container(
              width: 44,
              height: 44,
              decoration: BoxDecoration(
                color: color.withValues(alpha: 0.15),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Icon(icon, color: color, size: 22),
            ),
            const SizedBox(width: 16),
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(label, style: TextStyle(color: Colors.grey[500], fontSize: 12)),
                const SizedBox(height: 2),
                Text(value, style: const TextStyle(color: Colors.white, fontSize: 24, fontWeight: FontWeight.bold)),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
