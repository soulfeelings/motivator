import 'package:flutter/material.dart';
import '../services/api_service.dart';

class ChallengesScreen extends StatefulWidget {
  const ChallengesScreen({super.key});

  @override
  State<ChallengesScreen> createState() => _ChallengesScreenState();
}

class _ChallengesScreenState extends State<ChallengesScreen> {
  List<dynamic> _challenges = [];
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
        final data = await ApiService.get('/companies/$_companyId/challenges');
        setState(() => _challenges = data as List? ?? []);
      }
    } catch (_) {}
    if (mounted) setState(() => _loading = false);
  }

  Color _statusColor(String status) {
    switch (status) {
      case 'active': return const Color(0xFF8B5CF6);
      case 'completed': return const Color(0xFF10B981);
      case 'pending': return const Color(0xFFF59E0B);
      default: return Colors.grey;
    }
  }

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Padding(
            padding: EdgeInsets.fromLTRB(20, 20, 20, 16),
            child: Text('Challenges',
                style: TextStyle(fontSize: 28, fontWeight: FontWeight.bold, color: Colors.white)),
          ),
          Expanded(
            child: _loading
                ? const Center(child: CircularProgressIndicator())
                : _challenges.isEmpty
                    ? Center(child: Text('No challenges yet', style: TextStyle(color: Colors.grey[600])))
                    : RefreshIndicator(
                        onRefresh: _load,
                        child: ListView.builder(
                          padding: const EdgeInsets.symmetric(horizontal: 20),
                          itemCount: _challenges.length,
                          itemBuilder: (ctx, i) {
                            final ch = _challenges[i];
                            return Card(
                              margin: const EdgeInsets.only(bottom: 12),
                              child: Padding(
                                padding: const EdgeInsets.all(16),
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Row(
                                      children: [
                                        Container(
                                          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                                          decoration: BoxDecoration(
                                            color: _statusColor(ch['status']).withValues(alpha: 0.15),
                                            borderRadius: BorderRadius.circular(8),
                                          ),
                                          child: Text(
                                            ch['status'].toString().toUpperCase(),
                                            style: TextStyle(color: _statusColor(ch['status']), fontSize: 11, fontWeight: FontWeight.bold),
                                          ),
                                        ),
                                        const SizedBox(width: 8),
                                        Text(ch['metric'] ?? '', style: const TextStyle(color: Color(0xFF8B5CF6), fontFamily: 'monospace')),
                                        const Spacer(),
                                        Text('Target: ${ch['target']}', style: TextStyle(color: Colors.grey[500], fontSize: 12)),
                                      ],
                                    ),
                                    const SizedBox(height: 16),
                                    Row(
                                      children: [
                                        Expanded(
                                          child: Column(
                                            children: [
                                              Text('Challenger', style: TextStyle(color: Colors.grey[600], fontSize: 11)),
                                              const SizedBox(height: 4),
                                              Text('${ch['challenger_score'] ?? 0}',
                                                  style: const TextStyle(color: Color(0xFF10B981), fontSize: 24, fontWeight: FontWeight.bold)),
                                            ],
                                          ),
                                        ),
                                        Text('VS', style: TextStyle(color: Colors.grey[700], fontWeight: FontWeight.bold)),
                                        Expanded(
                                          child: Column(
                                            children: [
                                              Text('Opponent', style: TextStyle(color: Colors.grey[600], fontSize: 11)),
                                              const SizedBox(height: 4),
                                              Text('${ch['opponent_score'] ?? 0}',
                                                  style: const TextStyle(color: Color(0xFF10B981), fontSize: 24, fontWeight: FontWeight.bold)),
                                            ],
                                          ),
                                        ),
                                      ],
                                    ),
                                    if ((ch['xp_reward'] ?? 0) > 0 || (ch['wager'] ?? 0) > 0)
                                      Padding(
                                        padding: const EdgeInsets.only(top: 12),
                                        child: Text(
                                          'XP: +${ch['xp_reward']} | Wager: ${ch['wager']} coins',
                                          style: TextStyle(color: Colors.grey[600], fontSize: 12),
                                        ),
                                      ),
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
