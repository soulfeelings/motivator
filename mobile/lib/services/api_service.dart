import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:supabase_flutter/supabase_flutter.dart';

class ApiService {
  static const String _baseUrl = 'http://localhost:8080/api/v1';

  static Future<String> _getToken() async {
    final session = Supabase.instance.client.auth.currentSession;
    return session?.accessToken ?? '';
  }

  static Future<Map<String, String>> _headers() async {
    final token = await _getToken();
    return {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer $token',
    };
  }

  static Future<dynamic> get(String path) async {
    final res = await http.get(
      Uri.parse('$_baseUrl$path'),
      headers: await _headers(),
    );
    final json = jsonDecode(res.body);
    if (json['success'] != true) throw Exception(json['error'] ?? 'Request failed');
    return json['data'];
  }

  static Future<dynamic> post(String path, [Map<String, dynamic>? body]) async {
    final res = await http.post(
      Uri.parse('$_baseUrl$path'),
      headers: await _headers(),
      body: body != null ? jsonEncode(body) : null,
    );
    final json = jsonDecode(res.body);
    if (json['success'] != true) throw Exception(json['error'] ?? 'Request failed');
    return json['data'];
  }
}
