/*jshint node:true*/
"use strict";

var http = require('http');

http.createServer(function (req, res) {
  console.log('request');
  req.on('data', function(chunk) {
    console.log(chunk.toString().length);
  });
  setTimeout(function(){
    res.writeHead(200, {'Content-Type': 'text/plain'});
    res.end('Hello World\n');
  }, 350);
}).listen(1337, '127.0.0.1');

console.log('Server running at http://127.0.0.1:1337/');
