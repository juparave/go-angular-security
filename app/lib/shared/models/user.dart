class User {
  const User({
    required this.id,
    required this.email,
    required this.firstName,
    required this.lastName,
    required this.roleId,
    required this.enabled,
    this.accountId,
    this.accessToken,
    this.refreshToken,
  });

  final String id;
  final String email;
  final String firstName;
  final String lastName;
  final int roleId;
  final bool enabled;
  final String? accountId;
  final String? accessToken;
  final String? refreshToken;

  String get displayName => '$firstName $lastName'.trim();

  factory User.fromJson(Map<String, dynamic> json) => User(
    id: json['id'] as String,
    email: json['email'] as String,
    firstName: json['firstName'] as String? ?? '',
    lastName: json['lastName'] as String? ?? '',
    roleId: (json['roleId'] as num?)?.toInt() ?? 1,
    enabled: json['enabled'] as bool? ?? true,
    accountId: json['accountId'] as String?,
    accessToken: json['accessToken'] as String?,
    refreshToken: json['refreshToken'] as String?,
  );
}
