import 'package:flutter/material.dart';
import '../services/api_service.dart';

class QuestScreen extends StatefulWidget {
  const QuestScreen({super.key});

  @override
  State<QuestScreen> createState() => _QuestScreenState();
}

class _QuestScreenState extends State<QuestScreen> {
  String? _companyId;
  List<dynamic> _quests = [];
  Map<String, dynamic>? _activeTarget;
  List<dynamic> _receivedMessages = [];
  String? _activeQuestId;
  bool _loading = true;
  final _messageController = TextEditingController();

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
        final qs = await ApiService.get('/companies/$_companyId/quests') as List? ?? [];
        setState(() => _quests = qs);

        // Find active quest and load target
        for (final q in qs) {
          if (q['status'] == 'active' || q['status'] == 'voting' || q['status'] == 'revealed') {
            _activeQuestId = q['id'];
            await _loadQuestDetails(q['id']);
            break;
          }
        }
      }
    } catch (_) {}
    if (mounted) setState(() => _loading = false);
  }

  Future<void> _loadQuestDetails(String questId) async {
    try {
      final target = await ApiService.get('/companies/$_companyId/quests/$questId/my-target');
      final messages = await ApiService.get('/companies/$_companyId/quests/$questId/messages') as List? ?? [];
      setState(() {
        _activeTarget = target;
        _receivedMessages = messages;
      });
    } catch (_) {}
  }

  Future<void> _sendMessage() async {
    if (_activeQuestId == null || _messageController.text.isEmpty) return;
    try {
      await ApiService.post('/companies/$_companyId/quests/$_activeQuestId/send', {
        'message': _messageController.text,
      });
      _messageController.clear();
      await _loadQuestDetails(_activeQuestId!);
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Message sent!'), backgroundColor: Color(0xFF10B981)),
      );
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('$e'), backgroundColor: const Color(0xFFEF4444)),
      );
    }
  }

  Future<void> _vote(String pairId) async {
    if (_activeQuestId == null) return;
    try {
      await ApiService.post('/companies/$_companyId/quests/$_activeQuestId/vote', {
        'pair_id': pairId,
      });
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Vote cast!'), backgroundColor: Color(0xFF8B5CF6)),
      );
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('$e'), backgroundColor: const Color(0xFFEF4444)),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Secret Motivator', style: TextStyle(fontWeight: FontWeight.bold))),
      body: _loading
          ? const Center(child: CircularProgressIndicator())
          : _activeQuestId == null
              ? _buildNoQuest()
              : _buildActiveQuest(),
    );
  }

  Widget _buildNoQuest() {
    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(Icons.favorite_rounded, size: 64, color: const Color(0xFFEC4899).withValues(alpha: 0.3)),
          const SizedBox(height: 16),
          Text('No active quest', style: TextStyle(color: Colors.grey[500], fontSize: 16)),
          const SizedBox(height: 8),
          Text('Your admin will start one soon!', style: TextStyle(color: Colors.grey[700], fontSize: 14)),
        ],
      ),
    );
  }

  Widget _buildActiveQuest() {
    final quest = _quests.firstWhere((q) => q['id'] == _activeQuestId, orElse: () => null);
    final status = quest?['status'] ?? 'active';
    final alreadySent = _activeTarget?['sent'] == true;

    return RefreshIndicator(
      onRefresh: _load,
      child: ListView(
        padding: const EdgeInsets.all(20),
        children: [
          // Quest header
          Card(
            child: Padding(
              padding: const EdgeInsets.all(20),
              child: Column(
                children: [
                  const Icon(Icons.favorite_rounded, size: 40, color: Color(0xFFEC4899)),
                  const SizedBox(height: 12),
                  Text(quest?['name'] ?? 'Secret Motivator',
                      style: const TextStyle(color: Colors.white, fontSize: 20, fontWeight: FontWeight.bold)),
                  const SizedBox(height: 4),
                  Text(quest?['description'] ?? '', style: TextStyle(color: Colors.grey[500], fontSize: 13), textAlign: TextAlign.center),
                  const SizedBox(height: 12),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                    decoration: BoxDecoration(
                      color: const Color(0xFF8B5CF6).withValues(alpha: 0.15),
                      borderRadius: BorderRadius.circular(20),
                    ),
                    child: Text(status.toString().toUpperCase(),
                        style: const TextStyle(color: Color(0xFF8B5CF6), fontSize: 12, fontWeight: FontWeight.bold)),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 20),

          // Send message section
          if (status == 'active') ...[
            Text('Your Secret Target', style: TextStyle(color: Colors.grey[400], fontSize: 14, fontWeight: FontWeight.w600)),
            const SizedBox(height: 8),
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: alreadySent
                    ? Column(
                        children: [
                          const Icon(Icons.check_circle_rounded, color: Color(0xFF10B981), size: 32),
                          const SizedBox(height: 8),
                          const Text('Message sent!', style: TextStyle(color: Color(0xFF10B981), fontWeight: FontWeight.bold)),
                          const SizedBox(height: 4),
                          Text('"${_activeTarget?['message'] ?? ''}"',
                              style: TextStyle(color: Colors.grey[400], fontStyle: FontStyle.italic)),
                        ],
                      )
                    : Column(
                        children: [
                          Text('Send a positive message to your colleague',
                              style: TextStyle(color: Colors.grey[500], fontSize: 13)),
                          const SizedBox(height: 12),
                          TextField(
                            controller: _messageController,
                            maxLines: 3,
                            decoration: const InputDecoration(
                              hintText: 'Write something nice...',
                              hintStyle: TextStyle(color: Color(0xFF4B5563)),
                            ),
                            style: const TextStyle(color: Colors.white),
                          ),
                          const SizedBox(height: 12),
                          SizedBox(
                            width: double.infinity,
                            child: ElevatedButton.icon(
                              onPressed: _sendMessage,
                              icon: const Icon(Icons.send_rounded, size: 18),
                              label: const Text('Send Anonymously'),
                            ),
                          ),
                        ],
                      ),
              ),
            ),
          ],

          // Received messages
          if (_receivedMessages.isNotEmpty) ...[
            const SizedBox(height: 24),
            Text('Messages You Received', style: TextStyle(color: Colors.grey[400], fontSize: 14, fontWeight: FontWeight.w600)),
            const SizedBox(height: 8),
            ..._receivedMessages.map<Widget>((m) => Card(
              margin: const EdgeInsets.only(bottom: 8),
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        const Icon(Icons.favorite, color: Color(0xFFEC4899), size: 16),
                        const SizedBox(width: 8),
                        Text(m['sender_id'] != null ? 'From: ${(m['sender_id'] as String).substring(0, 8)}' : 'Anonymous',
                            style: TextStyle(color: Colors.grey[600], fontSize: 12)),
                      ],
                    ),
                    const SizedBox(height: 8),
                    Text('"${m['message']}"',
                        style: const TextStyle(color: Colors.white, fontSize: 15, fontStyle: FontStyle.italic)),
                    if (status == 'voting') ...[
                      const SizedBox(height: 12),
                      SizedBox(
                        width: double.infinity,
                        child: OutlinedButton.icon(
                          onPressed: () => _vote(m['pair_id']),
                          icon: const Icon(Icons.thumb_up_rounded, size: 16, color: Color(0xFF8B5CF6)),
                          label: const Text('Vote for this', style: TextStyle(color: Color(0xFF8B5CF6))),
                          style: OutlinedButton.styleFrom(side: const BorderSide(color: Color(0xFF1F1F2E))),
                        ),
                      ),
                    ],
                  ],
                ),
              ),
            )),
          ],
        ],
      ),
    );
  }
}
