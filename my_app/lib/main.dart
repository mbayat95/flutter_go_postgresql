import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';


void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Go App',
      home: HomePage(),
    );
  }
}

class HomePage extends StatefulWidget {
  @override
  _HomePageState createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final TextEditingController _controller = TextEditingController();
  String _response = '';

  Future<void> _postData() async {
    var response = await http.post(
      Uri.parse('http://localhost:8080/post-endpoint'),
      headers: {"Content-Type": "application/json"},
      body: json.encode({'data': _controller.text}),
    );
    setState(() {
      _response = response.body;
    });
  }

  Future<void> _getData() async {
    var response = await http.get(
      Uri.parse('http://localhost:8080/get-endpoint?data=' + Uri.encodeComponent(_controller.text)),
    );
    setState(() {
      _response = response.body;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Flutter Go App'),
      ),
      body: Padding(
        padding: EdgeInsets.all(16.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            TextField(
              controller: _controller,
              decoration: InputDecoration(labelText: 'Enter data'),
            ),
            ElevatedButton(
              onPressed: _postData,
              child: Text('POST'),
            ),
            ElevatedButton(
              onPressed: _getData,
              child: Text('GET'),
            ),
            Text(_response),
          ],
        ),
      ),
    );
  }
}
