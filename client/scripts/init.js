function init() {
  var ROOT = '//' + window.location.host + '/_ah/api';
  gapi.client.load('action', 'v1', function() {
    console.log('client loaded!');
  }, ROOT);
}

