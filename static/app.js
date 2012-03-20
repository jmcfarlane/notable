var notes = {};
notes.search = {};

var calculate_content = function (data, i) {
  var content = data.getValue(i, data.getColumnIndex('content')).split('\n');
  content.shift();
  return content.join('\n').substr(0, 100);
}

notes.render_notes = function(notes) {
  var div = document.getElementById('listing')
  var data = new google.visualization.DataTable(notes);
  var table = new google.visualization.Table(div);
  var view = new google.visualization.DataView(data);
  var columns = [
    data.getColumnIndex('updated'),
    data.getColumnIndex('subject'),
    {calc: calculate_content, type:'string', label:'Body'},
    data.getColumnIndex('tags'),
  ]
  var options = {
    height: '200px',
    sortAscending: false,
    sortColumn: 0,
    width: '100%',
  };

  view.setColumns(columns);
  table.draw(view, options);

  // Add navigation listener
  google.visualization.events.addListener(table, 'select',
    function(e) {
      edit(data, table.getSelection()[0].row);
  });
}

notes.create = function () {
  if (! $('#search input').is(":visible")) {
    $('#create').hide();
    $('#refresh').hide();
    $('#content').show();
    $('#reset').show();
    $('#persist').show();
    setTimeout("$('#content textarea').focus()", 100);
  }
}

edit = function (data, row) {
  notes.create();
  $('#content textarea').val(data.getValue(row, data.getColumnIndex('content')))
  $('#content #uid').val(data.getValue(row, data.getColumnIndex('uid')))
  $('#content #tags').val(data.getValue(row, data.getColumnIndex('tags')))
}

notes.search = function () {
  var q = {s: $('#search input').val()};
  $.get('/api/list', q, function (response) {
    notes.render_notes(response);
  });
}

notes.persist = function () {
  var post = {
    content: $('#content textarea').val(),
    tags: $('#content #tags').val(),
    uid: $('#content #uid').val(),
  }
  $.post('/api/persist', post, function (response) {
    notes.reset();
    notes.search();
  });
}

notes.reset = function () {
  $('#content').hide();
  $('#reset').hide();
  $('#persist').hide();
  $('#create').show();
  $('#refresh').show();
  $('#content #tags').val('');
  $('#content textarea').val('');
  $('#content #uid').val('');
}

notes.search.reset = function () {
  $('#search').hide();
  $('#search input').val('');
}

notes.search.perform = function () {
  if (! $('#content textarea').is(":visible")) {
    $('#search').show();
    setTimeout("$('#search input').focus()", 100);
  }
}
