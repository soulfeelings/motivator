import 'package:flutter/material.dart';
import '../services/api_service.dart';

class LeaderboardScreen extends StatefulWidget {
  const LeaderboardScreen({super.key});

  @override
  State<LeaderboardScreen> createState() => _LeaderboardScreenState();
}

class _LeaderboardScreenState extends State<LeaderboardScreen> {
  List<dynamic> _entries = [];
  String? _companyId;
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
        _companyId = memberships[0]['company_id'];
        final data = await ApiService.get('/companies/$_companyId/leaderboard?limit=50');
        setState(() => _entries = data as List? ?? []);
      }
    } catch (_) {}
    if (mounted) setState(() => _loading = false);
  }

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Padding(
            padding: EdgeInsets.fromLTRB(20, 20, 20, 16),
            child: Text('Leaderboard',
                style: TextStyle(fontSize: 28, fontWeight: FontWeight.bold, color: Colors.white)),
          ),
          Expanded(
            child: _loading
                ? const Center(child: CircularProgressIndicator())
                : RefreshIndicator(
                    onRefresh: _load,
                    child: ListView.builder(
                      padding: const EdgeInsets.symmetric(horizontal: 20),
                      itemCount: _entries.length,
                      itemBuilder: (ctx, i) {
                        final e = _entries[i];
                        final rank = e['rank'] ?? i + 1;
                        return Card(
                          margin: const EdgeInsets.only(bottom: 8),
                          child: ListTile(
                            leading: CircleAvatar(
                              backgroundColor: rank <= 3
                                  ? [const Color(0xFFF59E0B), const Color(0xFF9CA3AF), const Color(0xFFCD7F32)][rank - 1]
                                      .withValues(alpha: 0.2)
                                  : const Color(0xFF1F1F2E),
                              child: Text(
                                '$rank',
                                style: TextStyle(
                                  fontWeight: FontWeight.bold,
                                  color: rank <= 3
                                      ? [const Color(0xFFF59E0B), const Color(0xFF9CA3AF), const Color(0xFFCD7F32)][rank - 1]
                                      : Colors.grey,
                                ),
                              ),
                            ),
                            title: Text(
                              e['display_name'] ?? (e['user_id'] as String?)?.substring(0, 8) ?? '?',
                              style: const TextStyle(color: Colors.white, fontWeight: FontWeight.w600),
                            ),
                            subtitle: Text('Level ${e['level'] ?? 1}', style: TextStyle(color: Colors.grey[600])),
                            trailing: Row(
                              mainAxisSize: MainAxisSize.min,
                              children: [
                                const Icon(Icons.bolt_rounded, color: Color(0xFF10B981), size: 18),
                                const SizedBox(width: 4),
                                Text('${e['xp'] ?? 0}',
                                    style: const TextStyle(color: Color(0xFF10B981), fontWeight: FontWeight.bold)),
                              ],
                            ),
                          ),
                        );
                      },
                    ),
                  ),
          ),
        ],
      ),
    );
  }
}
