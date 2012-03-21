var notes = {};
notes.RE_ENCRYPTED = new RegExp(/[^=]{12,}==$/);
notes.encrypted = '<ENCRYPTED>';

var calculate_content = function (data, i) {
  var content = data.getValue(i, data.getColumnIndex('content')).split('\n');
  content.shift();
  if (notes.RE_ENCRYPTED.test(content)) {
    return notes.encrypted;
  } else {
    return content.join('\n').substr(0, 100);
  }
}

notes.render_notes = function(notes) {
  var div = document.getElementById('listing')
  var data = new google.visualization.DataTable(notes);
  var table = new google.visualization.Table(div);
  var view = new google.visualization.DataView(data);
  var columns = [
    data.getColumnIndex('updated'),
    data.getColumnIndex('subject'),
    {calc: calculate_content, id:'body', type:'string', label:'Body'},
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
  var found = notes.any_visible(['#search input', '#content textarea']);
  if (! found) {
    $('#create').hide();
    $('#refresh').hide();
    $('#content').show();
    $('#reset').show();
    $('#persist').show();
    $('#editor').show();
    setTimeout("$('#content textarea').focus()", 100);
  }
}

notes.launch_editor = function () {
  var post = {
    content: $('#content textarea').val(),
    uid: $('#content #uid').val(),
  }
  $.post('/api/launch_editor', post, function (response) {
    notes.poll_disk(response);
  });
}

notes.poll_disk = function (uid) {
  $.get('/api/from_disk/' + uid, function (response) {
    if (response != 'missing') {
      $('#content textarea').val(response);
      setTimeout('notes.poll_disk("'+ uid +'")', 1000)
    }
  });
}

notes.pwd_prompt = function () {
  $('#password-dialog').show();
  setTimeout("$('#password-dialog input').focus()", 100);
}

notes.pwd_submit = function () {
  var post = {
    password: $('#password-dialog input').val(),
    uid: $('#content #uid').val(),
  }
  $.post('/api/decrypt', post, function (response) {
    if (! notes.RE_ENCRYPTED.test(response)) {
      $('#content textarea').val(response);
    };
    notes.password_reset();
  });
  return false;
}

edit = function (data, row) {
  notes.create();
  $('#content textarea').val(data.getValue(row, data.getColumnIndex('content')))
  $('#content #uid').val(data.getValue(row, data.getColumnIndex('uid')))
  $('#content #tags').val(data.getValue(row, data.getColumnIndex('tags')))

  var b = data.getValue(row, data.getColumnIndex('content')).split('\n')[1];
  if (notes.RE_ENCRYPTED.test(b)) {
    notes.pwd_prompt();
  }
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
    password: $('#content #password').val(),
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
  $('#editor').hide();
  $('#create').show();
  $('#refresh').show();
  $('#content #tags').val('');
  $('#content textarea').val('');
  $('#content #uid').val('');
  $('#content #password').val('');
  notes.password_reset();
}

notes.search_reset = function () {
  $('#search').hide();
  $('#search input').val('');
}

notes.password_reset = function () {
  $('#password-dialog').hide();
  $('#password-dialog input').val('');
}

notes.any_visible = function (selectors) {
  var guilty = false
  $.each(selectors, function (idx, value) {
    if ($(value).is(":visible")) {
      guilty = true;
      return false;
    };
  });
  return guilty;
}

notes.search_perform = function () {
  var found = notes.any_visible(['#content textarea']);
  if (! found) {
    $('#search').show();
    setTimeout("$('#search input').focus()", 100);
  }
}
