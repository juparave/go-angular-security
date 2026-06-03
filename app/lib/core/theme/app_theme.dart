import 'package:flutter/material.dart';

import 'app_colors.dart';

/// Material 3 app theme.
///
/// Fonts: HankenGrotesk (body) + Piazzolla (display/serif headings) — bundled
/// in assets/fonts/, no network fetch required.
/// Corner radius: 12 dp — matches the web `--radius-box: 1rem`.
class AppTheme {
  const AppTheme._();

  static const double _radius = 12;
  static const String _body = 'HankenGrotesk';
  static const String _display = 'Piazzolla';

  static final ThemeData light = _build(Brightness.light);
  static final ThemeData dark = _build(Brightness.dark);

  static ThemeData _build(Brightness brightness) {
    final isLight = brightness == Brightness.light;

    final scheme =
        ColorScheme.fromSeed(
          seedColor: AppColors.primary,
          brightness: brightness,
        ).copyWith(
          primary: AppColors.azul,
          onPrimary: AppColors.azulContent,
          secondary: AppColors.teal,
          onSecondary: AppColors.tealContent,
          tertiary: AppColors.marigold,
          onTertiary: AppColors.marigoldContent,
          surface: isLight ? AppColors.base100 : AppColors.darkBase100,
          onSurface: isLight
              ? AppColors.baseContent
              : AppColors.darkBaseContent,
          surfaceContainerLow: isLight
              ? AppColors.base200
              : AppColors.darkBase200,
          surfaceContainer: isLight ? AppColors.base200 : AppColors.darkBase200,
          outline: isLight ? AppColors.base300 : AppColors.darkBase300,
          error: AppColors.error,
        );

    final textTheme = _buildTextTheme(scheme.onSurface);

    return ThemeData(
      useMaterial3: true,
      brightness: brightness,
      colorScheme: scheme,
      textTheme: textTheme,
      fontFamily: _body,
      scaffoldBackgroundColor: scheme.surface,
      extensions: [isLight ? TalaveraColors.light : TalaveraColors.dark],
      appBarTheme: AppBarTheme(
        centerTitle: false,
        elevation: 0,
        scrolledUnderElevation: 1,
        backgroundColor: scheme.surface,
        foregroundColor: scheme.onSurface,
        titleTextStyle: TextStyle(
          fontFamily: _display,
          fontWeight: FontWeight.w700,
          fontSize: 18,
          color: scheme.onSurface,
        ),
      ),
      cardTheme: CardThemeData(
        elevation: 0,
        color: isLight ? AppColors.base200 : AppColors.darkBase200,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(_radius),
          side: BorderSide(
            color: isLight ? AppColors.base300 : AppColors.darkBase300,
          ),
        ),
      ),
      dividerTheme: DividerThemeData(
        color: isLight ? AppColors.base300 : AppColors.darkBase300,
        thickness: 1,
        space: 1,
      ),
      inputDecorationTheme: InputDecorationTheme(
        isDense: true,
        filled: true,
        fillColor: isLight ? AppColors.base200 : AppColors.darkBase200,
        contentPadding: const EdgeInsets.symmetric(
          horizontal: 16,
          vertical: 14,
        ),
        hintStyle: TextStyle(
          fontFamily: _body,
          color: isLight
              ? AppColors.baseContent.withValues(alpha: 0.4)
              : AppColors.darkBaseContent.withValues(alpha: 0.4),
          fontWeight: FontWeight.w400,
        ),
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(_radius),
          borderSide: BorderSide(color: scheme.outline),
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(_radius),
          borderSide: BorderSide(color: scheme.outline),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(_radius),
          borderSide: const BorderSide(color: AppColors.azul, width: 2),
        ),
        errorBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(_radius),
          borderSide: const BorderSide(color: AppColors.error),
        ),
      ),
      filledButtonTheme: FilledButtonThemeData(
        style: FilledButton.styleFrom(
          backgroundColor: AppColors.azul,
          foregroundColor: AppColors.azulContent,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(_radius),
          ),
          padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 14),
          textStyle: const TextStyle(
            fontFamily: _body,
            fontWeight: FontWeight.w700,
            fontSize: 15,
            letterSpacing: 0.3,
          ),
        ),
      ),
      outlinedButtonTheme: OutlinedButtonThemeData(
        style: OutlinedButton.styleFrom(
          foregroundColor: scheme.onSurface,
          side: BorderSide(color: scheme.outline),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(_radius),
          ),
          padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 14),
          textStyle: const TextStyle(
            fontFamily: _body,
            fontWeight: FontWeight.w600,
            fontSize: 15,
          ),
        ),
      ),
      textButtonTheme: TextButtonThemeData(
        style: TextButton.styleFrom(
          foregroundColor: AppColors.azul,
          textStyle: const TextStyle(
            fontFamily: _body,
            fontWeight: FontWeight.w600,
            fontSize: 14,
          ),
        ),
      ),
      chipTheme: ChipThemeData(
        backgroundColor: isLight ? AppColors.base200 : AppColors.darkBase200,
        labelStyle: const TextStyle(
          fontFamily: _body,
          fontSize: 13,
          fontWeight: FontWeight.w500,
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
          side: BorderSide(color: scheme.outline),
        ),
      ),
      snackBarTheme: SnackBarThemeData(
        behavior: SnackBarBehavior.floating,
        backgroundColor: AppColors.navy,
        contentTextStyle: const TextStyle(
          fontFamily: _body,
          color: Colors.white,
          fontWeight: FontWeight.w500,
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(_radius),
        ),
      ),
      dialogTheme: DialogThemeData(
        backgroundColor: isLight ? AppColors.base100 : AppColors.darkBase100,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(_radius),
        ),
        titleTextStyle: TextStyle(
          fontFamily: _display,
          fontWeight: FontWeight.w700,
          fontSize: 18,
          color: scheme.onSurface,
        ),
        contentTextStyle: TextStyle(
          fontFamily: _body,
          fontSize: 14,
          fontWeight: FontWeight.w400,
          height: 1.5,
          color: scheme.onSurface.withValues(alpha: 0.7),
        ),
      ),
      floatingActionButtonTheme: FloatingActionButtonThemeData(
        backgroundColor: AppColors.marigold,
        foregroundColor: AppColors.marigoldContent,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(_radius),
        ),
      ),
      bottomNavigationBarTheme: BottomNavigationBarThemeData(
        backgroundColor: isLight ? AppColors.base100 : AppColors.darkBase200,
        selectedItemColor: AppColors.azul,
        unselectedItemColor: isLight
            ? AppColors.baseContent.withValues(alpha: 0.5)
            : AppColors.darkBaseContent.withValues(alpha: 0.5),
        type: BottomNavigationBarType.fixed,
        elevation: 0,
      ),
    );
  }

  static TextTheme _buildTextTheme(Color color) {
    TextStyle serif(double size, FontWeight weight, {double spacing = 0}) =>
        TextStyle(
          fontFamily: _display,
          fontSize: size,
          fontWeight: weight,
          letterSpacing: spacing,
          color: color,
          height: 1.2,
        );

    TextStyle sans(double size, FontWeight weight, {double spacing = 0}) =>
        TextStyle(
          fontFamily: _body,
          fontSize: size,
          fontWeight: weight,
          letterSpacing: spacing,
          color: color,
          height: 1.4,
        );

    return TextTheme(
      displayLarge: serif(32, FontWeight.w700),
      displayMedium: serif(26, FontWeight.w700),
      displaySmall: serif(22, FontWeight.w700),
      headlineLarge: serif(20, FontWeight.w700),
      headlineMedium: serif(18, FontWeight.w700),
      headlineSmall: serif(16, FontWeight.w600),
      titleLarge: sans(17, FontWeight.w700),
      titleMedium: sans(15, FontWeight.w600, spacing: 0.1),
      titleSmall: sans(13, FontWeight.w600, spacing: 0.1),
      bodyLarge: sans(16, FontWeight.w400),
      bodyMedium: sans(14, FontWeight.w400),
      bodySmall: sans(12, FontWeight.w400),
      labelLarge: sans(14, FontWeight.w600),
      labelMedium: sans(12, FontWeight.w500),
      labelSmall: sans(11, FontWeight.w500, spacing: 0.5),
    );
  }
}
