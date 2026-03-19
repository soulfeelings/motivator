import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:supabase_flutter/supabase_flutter.dart';

class ApiService {
  static const String _baseUrl = 'http://localhost:8080/api/v1';
  static const int _maxRetries = 2;
  static const Duration _retryDelay = Duration(seconds: 1);

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
    return _withRetry(() async {
      final res = await http.get(
        Uri.parse('$_baseUrl$path'),
        headers: await _headers(),
      );
      return _handleResponse(res);
    });
  }

  static Future<dynamic> post(String path, [Map<String, dynamic>? body]) async {
    return _withRetry(() async {
      final res = await http.post(
        Uri.parse('$_baseUrl$path'),
        headers: await _headers(),
        body: body != null ? jsonEncode(body) : null,
      );
      return _handleResponse(res);
    });
  }

  static dynamic _handleResponse(http.Response res) {
    if (res.statusCode == 401) {
      throw ApiException('Session expired. Please sign in again.', 401);
    }
    final json = jsonDecode(res.body);
    if (json['success'] != true) {
      throw ApiException(json['error'] ?? 'Request failed', res.statusCode);
    }
    return json['data'];
  }

  static Future<dynamic> _withRetry(Future<dynamic> Function() fn) async {
    for (int attempt = 0; attempt <= _maxRetries; attempt++) {
      try {
        return await fn();
      } on ApiException {
        rethrow;
      } catch (e) {
        if (attempt == _maxRetries) rethrow;
        await Future.delayed(_retryDelay * (attempt + 1));
      }
    }
  }
}

class ApiException implements Exception {
  final String message;
  final int statusCode;

  ApiException(this.message, this.statusCode);

  @override
  String toString() => message;
}
