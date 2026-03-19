import 'package:flutter/material.dart';
import '../services/api_service.dart';

class RewardsScreen extends StatefulWidget {
  const RewardsScreen({super.key});

  @override
  State<RewardsScreen> createState() => _RewardsScreenState();
}

class _RewardsScreenState extends State<RewardsScreen> {
  List<dynamic> _rewards = [];
  String? _companyId;
  int _coins = 0;
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
        _coins = memberships[0]['coins'] ?? 0;
        final data = await ApiService.get('/companies/$_companyId/rewards');
        setState(() => _rewards = data as List? ?? []);
      }
    } catch (_) {}
    if (mounted) setState(() => _loading = false);
  }

  Future<void> _redeem(String rewardId, int cost) async {
    if (_coins < cost) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Not enough coins'), backgroundColor: Color(0xFFEF4444)),
      );
      return;
    }
    try {
      await ApiService.post('/companies/$_companyId/rewards/redeem', {'reward_id': rewardId});
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Reward redeemed!'), backgroundColor: Color(0xFF10B981)),
      );
      _load();
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('$e'), backgroundColor: const Color(0xFFEF4444)),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(20, 20, 20, 4),
            child: Row(
              children: [
                const Text('Reward Store',
                    style: TextStyle(fontSize: 28, fontWeight: FontWeight.bold, color: Colors.white)),
                const Spacer(),
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                  decoration: BoxDecoration(
                    color: const Color(0xFFF59E0B).withValues(alpha: 0.15),
                    borderRadius: BorderRadius.circular(20),
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      const Icon(Icons.monetization_on_rounded, color: Color(0xFFF59E0B), size: 18),
                      const SizedBox(width: 4),
                      Text('$_coins', style: const TextStyle(color: Color(0xFFF59E0B), fontWeight: FontWeight.bold)),
                    ],
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: 12),
          Expanded(
            child: _loading
                ? const Center(child: CircularProgressIndicator())
                : _rewards.isEmpty
                    ? Center(child: Text('No rewards available', style: TextStyle(color: Colors.grey[600])))
                    : RefreshIndicator(
                        onRefresh: _load,
                        child: ListView.builder(
                          padding: const EdgeInsets.symmetric(horizontal: 20),
                          itemCount: _rewards.length,
                          itemBuilder: (ctx, i) {
                            final rw = _rewards[i];
                            final cost = rw['cost_coins'] ?? 0;
                            final canAfford = _coins >= cost;
                            return Card(
                              margin: const EdgeInsets.only(bottom: 12),
                              child: Padding(
                                padding: const EdgeInsets.all(16),
                                child: Row(
                                  children: [
                                    Container(
                                      width: 48,
                                      height: 48,
                                      decoration: BoxDecoration(
                                        color: const Color(0xFFF59E0B).withValues(alpha: 0.15),
                                        borderRadius: BorderRadius.circular(12),
                                      ),
                                      child: const Icon(Icons.card_giftcard_rounded, color: Color(0xFFF59E0B)),
                                    ),
                                    const SizedBox(width: 16),
                                    Expanded(
                                      child: Column(
                                        crossAxisAlignment: CrossAxisAlignment.start,
                                        children: [
                                          Text(rw['name'] ?? '',
                                              style: const TextStyle(color: Colors.white, fontWeight: FontWeight.w600)),
                                          if (rw['description'] != null)
                                            Text(rw['description'], style: TextStyle(color: Colors.grey[600], fontSize: 12)),
                                        ],
                                      ),
                                    ),
                                    Column(
                                      children: [
                                        Text('$cost', style: const TextStyle(color: Color(0xFFF59E0B), fontWeight: FontWeight.bold, fontSize: 18)),
                                        const SizedBox(height: 4),
                                        GestureDetector(
                                          onTap: canAfford ? () => _redeem(rw['id'], cost) : null,
                                          child: Container(
                                            padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 6),
                                            decoration: BoxDecoration(
                                              color: canAfford ? const Color(0xFF8B5CF6) : const Color(0xFF1F1F2E),
                                              borderRadius: BorderRadius.circular(8),
                                            ),
                                            child: Text(
                                              'Redeem',
                                              style: TextStyle(
                                                color: canAfford ? Colors.white : Colors.grey[600],
                                                fontSize: 12,
                                                fontWeight: FontWeight.w600,
                                              ),
                                            ),
                                          ),
                                        ),
                                      ],
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
