import 'package:flutter/material.dart';
import 'dashboard_screen.dart';
import 'leaderboard_screen.dart';
import 'challenges_screen.dart';
import 'rewards_screen.dart';
import 'profile_screen.dart';
import 'game_screen.dart';

class ShellScreen extends StatefulWidget {
  const ShellScreen({super.key});

  @override
  State<ShellScreen> createState() => _ShellScreenState();
}

class _ShellScreenState extends State<ShellScreen> {
  int _index = 0;

  final _screens = const [
    DashboardScreen(),
    LeaderboardScreen(),
    ChallengesScreen(),
    RewardsScreen(),
    ProfileScreen(),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: _screens[_index],
      floatingActionButton: FloatingActionButton(
        onPressed: () => Navigator.of(context).push(
          MaterialPageRoute(builder: (_) => const GameScreen()),
        ),
        backgroundColor: const Color(0xFF8B5CF6),
        child: const Icon(Icons.sports_esports_rounded, color: Colors.white),
      ),
      bottomNavigationBar: NavigationBar(
        selectedIndex: _index,
        onDestinationSelected: (i) => setState(() => _index = i),
        backgroundColor: const Color(0xFF111118),
        indicatorColor: const Color(0xFF8B5CF6).withValues(alpha: 0.2),
        destinations: const [
          NavigationDestination(icon: Icon(Icons.dashboard_rounded), label: 'Home'),
          NavigationDestination(icon: Icon(Icons.leaderboard_rounded), label: 'Rank'),
          NavigationDestination(icon: Icon(Icons.sports_mma_rounded), label: 'Battle'),
          NavigationDestination(icon: Icon(Icons.card_giftcard_rounded), label: 'Store'),
          NavigationDestination(icon: Icon(Icons.person_rounded), label: 'Profile'),
        ],
      ),
    );
  }
}
