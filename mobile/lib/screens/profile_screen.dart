import 'package:flutter/material.dart';
import 'package:supabase_flutter/supabase_flutter.dart';
import '../services/api_service.dart';

class ProfileScreen extends StatefulWidget {
  const ProfileScreen({super.key});

  @override
  State<ProfileScreen> createState() => _ProfileScreenState();
}

class _ProfileScreenState extends State<ProfileScreen> {
  Map<String, dynamic>? _profile;
  List<dynamic> _badges = [];
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
        final companyId = m['company_id'];
        final memberId = m['id'];
        final profileData = await ApiService.get('/companies/$companyId/members/$memberId/profile');
        setState(() {
          _profile = profileData['membership'];
          _badges = profileData['badges'] as List? ?? [];
        });
      }
    } catch (_) {}
    if (mounted) setState(() => _loading = false);
  }

  @override
  Widget build(BuildContext context) {
    final email = Supabase.instance.client.auth.currentUser?.email ?? '';

    return SafeArea(
      child: _loading
          ? const Center(child: CircularProgressIndicator())
          : ListView(
              padding: const EdgeInsets.all(20),
              children: [
                const Text('Profile',
                    style: TextStyle(fontSize: 28, fontWeight: FontWeight.bold, color: Colors.white)),
                const SizedBox(height: 20),
                Card(
                  child: Padding(
                    padding: const EdgeInsets.all(20),
                    child: Column(
                      children: [
                        CircleAvatar(
                          radius: 36,
                          backgroundColor: const Color(0xFF8B5CF6).withValues(alpha: 0.2),
                          child: const Icon(Icons.person_rounded, size: 36, color: Color(0xFF8B5CF6)),
                        ),
                        const SizedBox(height: 12),
                        Text(email, style: TextStyle(color: Colors.grey[400], fontSize: 14)),
                        if (_profile != null) ...[
                          const SizedBox(height: 16),
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                            children: [
                              _ProfileStat('Level', '${_profile!['level'] ?? 1}', const Color(0xFF8B5CF6)),
                              _ProfileStat('XP', '${_profile!['xp'] ?? 0}', const Color(0xFF10B981)),
                              _ProfileStat('Coins', '${_profile!['coins'] ?? 0}', const Color(0xFFF59E0B)),
                            ],
                          ),
                        ],
                      ],
                    ),
                  ),
                ),
                if (_badges.isNotEmpty) ...[
                  const SizedBox(height: 24),
                  Text('Badges', style: TextStyle(color: Colors.grey[400], fontSize: 14, fontWeight: FontWeight.w600)),
                  const SizedBox(height: 12),
                  Wrap(
                    spacing: 8,
                    runSpacing: 8,
                    children: _badges.map<Widget>((b) {
                      return Container(
                        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 8),
                        decoration: BoxDecoration(
                          color: const Color(0xFF8B5CF6).withValues(alpha: 0.15),
                          borderRadius: BorderRadius.circular(20),
                          border: Border.all(color: const Color(0xFF8B5CF6).withValues(alpha: 0.3)),
                        ),
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            const Icon(Icons.military_tech_rounded, color: Color(0xFF8B5CF6), size: 16),
                            const SizedBox(width: 6),
                            Text(b['name'] ?? '', style: const TextStyle(color: Colors.white, fontSize: 13)),
                          ],
                        ),
                      );
                    }).toList(),
                  ),
                ],
                const SizedBox(height: 32),
                OutlinedButton.icon(
                  onPressed: () => Supabase.instance.client.auth.signOut(),
                  icon: const Icon(Icons.logout_rounded, color: Color(0xFFEF4444)),
                  label: const Text('Sign Out', style: TextStyle(color: Color(0xFFEF4444))),
                  style: OutlinedButton.styleFrom(
                    side: const BorderSide(color: Color(0xFF1F1F2E)),
                    padding: const EdgeInsets.symmetric(vertical: 14),
                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                  ),
                ),
              ],
            ),
    );
  }
}

class _ProfileStat extends StatelessWidget {
  final String label;
  final String value;
  final Color color;

  const _ProfileStat(this.label, this.value, this.color);

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Text(value, style: TextStyle(color: color, fontSize: 22, fontWeight: FontWeight.bold)),
        const SizedBox(height: 2),
        Text(label, style: TextStyle(color: Colors.grey[600], fontSize: 12)),
      ],
    );
  }
}
