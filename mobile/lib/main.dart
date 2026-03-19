import 'package:flutter/material.dart';
import 'package:supabase_flutter/supabase_flutter.dart';
import 'app.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  await Supabase.initialize(
    url: const String.fromEnvironment('SUPABASE_URL',
        defaultValue: 'https://evfkxiphjhriwaozppsf.supabase.co'),
    anonKey: const String.fromEnvironment('SUPABASE_ANON_KEY',
        defaultValue:
            'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImV2Zmt4aXBoamhyaXdhb3pwcHNmIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjQwODgwMTcsImV4cCI6MjA3OTY2NDAxN30.Nqxq6azBpBq6qCq8aZ-DwEeH9E0eKlASpEbOf3Lgj9E'),
  );

  runApp(const MotivatorApp());
}
