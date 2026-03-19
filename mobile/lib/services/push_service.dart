import 'dart:io';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'api_service.dart';

class PushService {
  static final FirebaseMessaging _fcm = FirebaseMessaging.instance;

  static Future<void> init(String companyId) async {
    final settings = await _fcm.requestPermission(
      alert: true,
      badge: true,
      sound: true,
    );

    if (settings.authorizationStatus == AuthorizationStatus.denied) return;

    final token = await _fcm.getToken();
    if (token != null) {
      await _registerToken(companyId, token);
    }

    _fcm.onTokenRefresh.listen((newToken) {
      _registerToken(companyId, newToken);
    });

    FirebaseMessaging.onMessage.listen(_handleForeground);
    FirebaseMessaging.onMessageOpenedApp.listen(_handleBackground);
  }

  static Future<void> _registerToken(String companyId, String token) async {
    try {
      await ApiService.post('/companies/$companyId/notifications/register', {
        'token': token,
        'platform': Platform.isIOS ? 'ios' : 'android',
      });
    } catch (_) {}
  }

  static void _handleForeground(RemoteMessage message) {
    // Foreground notification received — can show in-app banner
  }

  static void _handleBackground(RemoteMessage message) {
    // User tapped notification — can navigate to relevant screen
  }
}
