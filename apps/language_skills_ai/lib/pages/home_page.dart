import 'dart:io';
import 'dart:convert' show json;

import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter/material.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import 'package:http/http.dart' as http;
import '../env.dart';
import 'settings_page.dart';

final FirebaseAuth _auth = FirebaseAuth.instance;

class HomePage extends StatefulWidget {
  const HomePage({Key? key}) : super(key: key);

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  bool _isButtonDisabled = true;
  final inputController = TextEditingController();
  final outputController = TextEditingController();

  @override
  void initState() {
    super.initState();
    inputController.addListener(() {
      if (inputController.text.length > 10) {
        if (_isButtonDisabled) {
          setState(() {
            _isButtonDisabled = false;
          });
        }
      } else {
        if (!_isButtonDisabled) {
          setState(() {
            _isButtonDisabled = true;
          });
        }
      }
    });
  }

  @override
  void dispose() {
    inputController.dispose();
    outputController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final AppLocalizations l10n = AppLocalizations.of(context)!;
    return Scaffold(
      appBar: AppBar(
        actions: _appBarActions(context, l10n),
        title: Text(l10n.appTitle),
      ),
      body: Padding(
        padding: const EdgeInsets.all(20),
        child: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: <Widget>[
              TextField(
                controller: inputController,
                minLines: 5,
                maxLines: 5,
                maxLength: 200,
                decoration: const InputDecoration(
                  border: OutlineInputBorder(),
                  labelText: 'Your input',
                ),
              ),
              // const SizedBox(height: 30),
              ElevatedButton(
                onPressed: _isButtonDisabled ? null : () => sendRequest(inputController.text),
                child: Text(_isButtonDisabled ? 'Loading' : 'Improve'),
              ),
              TextField(
                controller: outputController,
                minLines: 5,
                maxLines: 5,
                decoration: const InputDecoration(
                  labelText: 'Output',
                ),
              )
            ],
          ),
        ),
      ),
    );
  }
  
  void sendRequest(String input) async {
    final idToken = await _auth.currentUser?.getIdToken();
    if (idToken != null) {
      setState(() {
        _isButtonDisabled = true;
      });

      final response = await http.post(
          Uri.parse(correctionAIURL),
          headers: {
            HttpHeaders.contentTypeHeader: 'application/json',
            HttpHeaders.authorizationHeader: 'Bearer $idToken',
          },
        body: json.encode({
          'text': input
        }));
      if (response.statusCode == 200) {
        outputController.text = response.body;
      } else {
        // TODO: Show error (e.g. in toast)
      }
      setState(() {
        _isButtonDisabled = false;
      });
    } else {
        // TODO: Show error (e.g. in toast)
    }
  }

  List<Widget> _appBarActions(BuildContext context, AppLocalizations l10n) {
    return [
      Padding(
        padding: const EdgeInsets.only(right: 3.0),
        child: IconButton(
          icon: const Icon(
            Icons.settings,
          ),
          tooltip: l10n.settingsPage,
          onPressed: () {
            Navigator.push(context,
                MaterialPageRoute(builder: (context) => const SettingsPage()));
          },
        ),
      )
    ];
  }
}
