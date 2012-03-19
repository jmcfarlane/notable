notes = {};

notes.render_notes = function(notes) {
  var div = document.getElementById('content')
  var data = new google.visualization.DataTable(notes);
  var table = new google.visualization.Table(div);
  table.draw(data);
}

notes.fetch_by_tags = function () {
  $.ajax({url: '/api/list/', success: notes.render_notes });
}
