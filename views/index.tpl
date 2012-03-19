% setdefault('jq', 'http://ajax.googleapis.com/ajax/libs/jquery')
% setdefault('jqv', '1.7.1')
<html>
  <head>
  </head>
  <body>
    <div id='content'></div>
    <input id="refresh" type="button" value="Refresh" />
    <script type="text/javascript" src="/static/app.js"></script>
    <script type="text/javascript" src="{{jq}}/{{jqv}}/jquery.min.js"></script>
    <script type="text/javascript" src="http://www.google.com/jsapi"></script>
    <script type="text/javascript">
      // Load the gviz table
      var pkgs = ['corechart', 'table'];
      google.load('visualization', '1.0', {'packages':pkgs});
      google.setOnLoadCallback(notes.fetch_by_tags);

      // Add a refresh button
      $('#refresh').on('click', notes.fetch_by_tags);

    </script>
  </body>
</html>
