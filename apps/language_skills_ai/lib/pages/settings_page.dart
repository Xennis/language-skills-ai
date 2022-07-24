import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter/material.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import 'package:google_sign_in/google_sign_in.dart';

class SettingsPage extends StatelessWidget {
  const SettingsPage({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final AppLocalizations l10n = AppLocalizations.of(context)!;

    final User user = FirebaseAuth.instance.currentUser!;
    final String email = user.email!;

    return Scaffold(
        appBar: AppBar(
          title: Text(l10n.settingsPage),
        ),
        body: SingleChildScrollView(
              child: Padding(
                padding: const EdgeInsets.fromLTRB(10, 20, 10, 20),
                child: Column(
                  children: [
                    ListTile(
                      leading: const Icon(Icons.account_circle),
                      title: const Text('Logout'),
                      subtitle: Text('Account: $email'),
                      onTap: () {
                        _signOut();
                        // Close settings
                        Navigator.of(context).pop();
                      },
                    ),
                  ],
                ),
              ),
            )
          );
  }

  Future<void> _signOut() async {
    await FirebaseAuth.instance.signOut();
    await GoogleSignIn().signOut();
  }
}


