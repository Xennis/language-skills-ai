import 'dart:io';

import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_signin_button/flutter_signin_button.dart';
import 'package:google_sign_in/google_sign_in.dart';

typedef OAuthSignIn = void Function();

final FirebaseAuth _auth = FirebaseAuth.instance;

/// Entrypoint example for various sign-in flows with Firebase.
class AuthGate extends StatefulWidget {
  // ignore: public_member_api_docs
  const AuthGate({Key? key}) : super(key: key);

  @override
  State<StatefulWidget> createState() => _AuthGateState();
}

class _AuthGateState extends State<AuthGate> {
  GlobalKey<FormState> formKey = GlobalKey<FormState>();
  String error = '';
  String verificationId = '';

  bool isLoading = false;

  void setIsLoading() {
    setState(() {
      isLoading = !isLoading;
    });
  }

  late Map<Buttons, OAuthSignIn> authButtons;

  @override
  void initState() {
    super.initState();
    authButtons = {
        Buttons.Google: _signInWithGoogle,
    };
  }

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: FocusScope.of(context).unfocus,
      child: Scaffold(
        body: Center(
          child: SingleChildScrollView(
            child: Center(
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20),
                child: SafeArea(
                  child: Form(
                    key: formKey,
                    autovalidateMode: AutovalidateMode.onUserInteraction,
                    child: ConstrainedBox(
                      constraints: const BoxConstraints(maxWidth: 400),
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Visibility(
                            visible: error.isNotEmpty,
                            child: MaterialBanner(
                              backgroundColor: Theme.of(context).errorColor,
                              content: Text(error),
                              actions: [
                                TextButton(
                                  onPressed: () {
                                    setState(() {
                                      error = '';
                                    });
                                  },
                                  child: const Text(
                                    'dismiss',
                                    style: TextStyle(color: Colors.white),
                                  ),
                                )
                              ],
                              contentTextStyle:
                              const TextStyle(color: Colors.white),
                              padding: const EdgeInsets.all(10),
                            ),
                          ),
                          const SizedBox(height: 20),
                          ...authButtons.keys
                              .map(
                                (button) => Padding(
                              padding:
                              const EdgeInsets.symmetric(vertical: 5),
                              child: AnimatedSwitcher(
                                duration: const Duration(milliseconds: 200),
                                child: isLoading
                                    ? Container(
                                  color: Colors.grey[200],
                                  height: 50,
                                  width: double.infinity,
                                )
                                    : SizedBox(
                                  width: double.infinity,
                                  height: 50,
                                  child: SignInButton(
                                    button,
                                    onPressed: authButtons[button]!,
                                  ),
                                ),
                              ),
                            ),
                          )
                              .toList(),
                        ],
                      ),
                    ),
                  ),
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }

  Future<void> _signInWithGoogle() async {
    setIsLoading();
    try {
      // Trigger the authentication flow
      final googleUser = await GoogleSignIn().signIn();

      // Obtain the auth details from the request
      final googleAuth = await googleUser?.authentication;

      if (googleAuth != null) {
        // Create a new credential
        final credential = GoogleAuthProvider.credential(
          accessToken: googleAuth.accessToken,
          idToken: googleAuth.idToken,
        );

        // Once signed in, return the UserCredential
        await _auth.signInWithCredential(credential);
      }
    } on FirebaseAuthException catch (e) {
      setState(() {
        error = '${e.message}';
      });
    } finally {
      setIsLoading();
    }
  }

}
