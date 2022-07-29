import 'package:flutter/material.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import 'package:language_skills_ai/pages/auth_gate_widget.dart';


class AuthPage extends StatelessWidget {
  const AuthPage({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final AppLocalizations l10n = AppLocalizations.of(context)!;
    return Scaffold(
        body: LayoutBuilder(
          builder: (context, constraines) {
            return Row(
              children: [
                Visibility(
                  visible: constraines.maxWidth >= 1200,
                  child: Expanded(
                    child: Container(
                      height: double.infinity,
                      color: Theme.of(context).colorScheme.primary,
                      child: Center(
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Text(
                              l10n.appTitle,
                              style: Theme.of(context).textTheme.headline4,
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
                ),
                SizedBox(
                  width: constraines.maxWidth >= 1200
                      ? constraines.maxWidth / 2
                      : constraines.maxWidth,
                  child: AuthGate(),
                ),
              ],
            );
          },
        ),
      );
  }
}
