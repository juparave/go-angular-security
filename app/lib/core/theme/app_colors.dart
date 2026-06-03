import 'package:flutter/material.dart';

/// Design tokens ported from the web app's Talavera DaisyUI theme.
/// Source of truth: `angular/src/tailwind.css` — `@plugin "daisyui/theme" { name: 'talavera' }`.
///
/// OKLCH values are approximated to sRGB hex. Never hardcode colors in widgets;
/// always reference these constants or the semantic [TalaveraColors] extension.
class AppColors {
  const AppColors._();

  // ─── Base surfaces ────────────────────────────────────────────────────────
  /// oklch(0.995 0.004 250) — near-white page background
  static const Color base100 = Color(0xFFFAFBFF);

  /// oklch(0.97 0.008 250) — card / sidebar background
  static const Color base200 = Color(0xFFF0F2FA);

  /// oklch(0.93 0.014 250) — divider / subtle background
  static const Color base300 = Color(0xFFE2E6F5);

  /// oklch(0.25 0.03 258) — body text / ink
  static const Color baseContent = Color(0xFF2B2D42);

  // ─── Primary — Azul Talavera ──────────────────────────────────────────────
  /// oklch(0.52 0.15 250)
  static const Color azul = Color(0xFF2D5BE3);
  static const Color azulContent = Color(0xFFFDFEFF);

  // ─── Secondary — Teal ────────────────────────────────────────────────────
  /// oklch(0.55 0.12 200)
  static const Color teal = Color(0xFF2A8C8C);
  static const Color tealContent = Color(0xFFFDFEFF);

  // ─── Accent — Marigold / Oro ─────────────────────────────────────────────
  /// oklch(0.82 0.15 80)
  static const Color marigold = Color(0xFFE8C14A);
  static const Color marigoldContent = Color(0xFF3D2E08);

  // ─── Neutral — Ink / Navy ────────────────────────────────────────────────
  /// oklch(0.25 0.03 258)
  static const Color navy = Color(0xFF2B2D42);
  static const Color navyContent = Color(0xFFF8F9FE);

  // ─── Semantic ─────────────────────────────────────────────────────────────
  /// oklch(0.6 0.12 230)
  static const Color info = Color(0xFF3B82C4);

  /// oklch(0.6 0.13 150)
  static const Color success = Color(0xFF22A869);

  /// oklch(0.76 0.15 82)
  static const Color warning = Color(0xFFD4920A);

  /// oklch(0.58 0.21 25)
  static const Color error = Color(0xFFD94430);

  // ─── Dark mode surfaces ───────────────────────────────────────────────────
  static const Color darkBase100 = Color(0xFF1A1F3A);
  static const Color darkBase200 = Color(0xFF141828);
  static const Color darkBase300 = Color(0xFF0F1220);
  static const Color darkBaseContent = Color(0xFFE8EAF6);

  // ─── Seed for ColorScheme.fromSeed ────────────────────────────────────────
  static const Color primary = azul;
}

/// Semantic tokens consumed by feature widgets — accessed via
/// `Theme.of(context).extension<TalaveraColors>()!` or the `context.talavera` shortcut.
@immutable
class TalaveraColors extends ThemeExtension<TalaveraColors> {
  const TalaveraColors({
    required this.success,
    required this.warning,
    required this.danger,
    required this.info,
    required this.microLabel,
    required this.accentBar,
    required this.navBg,
    required this.navBgDeep,
    required this.navText,
    required this.navTextMuted,
    required this.navBorder,
    required this.navActiveBg,
    required this.navActiveText,
    required this.navHoverBg,
  });

  final Color success;
  final Color warning;
  final Color danger;
  final Color info;
  final Color microLabel;
  final Color accentBar;
  final Color navBg;
  final Color navBgDeep;
  final Color navText;
  final Color navTextMuted;
  final Color navBorder;
  final Color navActiveBg;
  final Color navActiveText;
  final Color navHoverBg;

  static const TalaveraColors light = TalaveraColors(
    success: AppColors.success,
    warning: AppColors.warning,
    danger: AppColors.error,
    info: AppColors.info,
    microLabel: Color(0xFF8890B0),
    accentBar: AppColors.azul,
    navBg: AppColors.navy,
    navBgDeep: Color(0xFF1A1F3A),
    navText: AppColors.navyContent,
    navTextMuted: Color(0xFF8890B0),
    navBorder: Color(0x1F2B2D42),
    navActiveBg: Color(0x1A2D5BE3),
    navActiveText: AppColors.azul,
    navHoverBg: Color(0x0F2B2D42),
  );

  static const TalaveraColors dark = TalaveraColors(
    success: AppColors.success,
    warning: AppColors.warning,
    danger: AppColors.error,
    info: AppColors.info,
    microLabel: Color(0xFF6670A0),
    accentBar: AppColors.marigold,
    navBg: AppColors.darkBase200,
    navBgDeep: AppColors.darkBase300,
    navText: AppColors.darkBaseContent,
    navTextMuted: Color(0xFF6670A0),
    navBorder: Color(0x1FE8EAF6),
    navActiveBg: Color(0x1AE8C14A),
    navActiveText: AppColors.marigold,
    navHoverBg: Color(0x0FE8EAF6),
  );

  @override
  TalaveraColors copyWith({
    Color? success,
    Color? warning,
    Color? danger,
    Color? info,
    Color? microLabel,
    Color? accentBar,
    Color? navBg,
    Color? navBgDeep,
    Color? navText,
    Color? navTextMuted,
    Color? navBorder,
    Color? navActiveBg,
    Color? navActiveText,
    Color? navHoverBg,
  }) => TalaveraColors(
    success: success ?? this.success,
    warning: warning ?? this.warning,
    danger: danger ?? this.danger,
    info: info ?? this.info,
    microLabel: microLabel ?? this.microLabel,
    accentBar: accentBar ?? this.accentBar,
    navBg: navBg ?? this.navBg,
    navBgDeep: navBgDeep ?? this.navBgDeep,
    navText: navText ?? this.navText,
    navTextMuted: navTextMuted ?? this.navTextMuted,
    navBorder: navBorder ?? this.navBorder,
    navActiveBg: navActiveBg ?? this.navActiveBg,
    navActiveText: navActiveText ?? this.navActiveText,
    navHoverBg: navHoverBg ?? this.navHoverBg,
  );

  @override
  TalaveraColors lerp(ThemeExtension<TalaveraColors>? other, double t) {
    if (other is! TalaveraColors) return this;
    return TalaveraColors(
      success: Color.lerp(success, other.success, t)!,
      warning: Color.lerp(warning, other.warning, t)!,
      danger: Color.lerp(danger, other.danger, t)!,
      info: Color.lerp(info, other.info, t)!,
      microLabel: Color.lerp(microLabel, other.microLabel, t)!,
      accentBar: Color.lerp(accentBar, other.accentBar, t)!,
      navBg: Color.lerp(navBg, other.navBg, t)!,
      navBgDeep: Color.lerp(navBgDeep, other.navBgDeep, t)!,
      navText: Color.lerp(navText, other.navText, t)!,
      navTextMuted: Color.lerp(navTextMuted, other.navTextMuted, t)!,
      navBorder: Color.lerp(navBorder, other.navBorder, t)!,
      navActiveBg: Color.lerp(navActiveBg, other.navActiveBg, t)!,
      navActiveText: Color.lerp(navActiveText, other.navActiveText, t)!,
      navHoverBg: Color.lerp(navHoverBg, other.navHoverBg, t)!,
    );
  }
}

extension TalaveraColorsX on BuildContext {
  TalaveraColors get talavera => Theme.of(this).extension<TalaveraColors>()!;
}
