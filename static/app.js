notes = {};

notes.render_notes = function(notes) {
  var div = document.getElementById('listing')
  var data = new google.visualization.DataTable(notes);
  var table = new google.visualization.Table(div);
  var options = {
    height: '200px',
    showRowNumber: true,
    sortAscending: false,
    sortColumn: 0,
    width: '100%',
  };
  table.draw(data, options);
}

notes.create = function () {
  $('#create').hide();
  $('#refresh').hide();
  $('#content').show();
  $('#reset').show();
  $('#persist').show();
  setTimeout("$('#content textarea').focus()", 100);
}

notes.fetch_by_tags = function () {
  $.ajax({url: '/api/list/', success: notes.render_notes });
}

notes.persist = function () {
  var content = $('#content textarea').val();
  var subject = content.split('\n')[0];
  var note = {content: content, subject: subject};
  $.post('/api/persist', note, function (response) {
    notes.reset();
    notes.fetch_by_tags();
  });
}

notes.reset = function () {
  $('#content').hide();
  $('#reset').hide();
  $('#persist').hide();
  $('#create').show();
  $('#refresh').show();
  $('#content textarea').val('');
}
