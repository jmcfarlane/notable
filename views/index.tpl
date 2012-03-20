% setdefault('jq', 'http://ajax.googleapis.com/ajax/libs/jquery')
% setdefault('jqv', '1.7.1')
<html>
  <head>
    <link type="text/css" rel="stylesheet" href="/static/style.css" />
  </head>
  <body>
    <div id="listing-container">
      <div id="listing"></div>
    </div>

    <div id="content" name="content" class="hidden">
      <form method="post" target="/api/create">
        <textarea></textarea>
      </form>
    </div>

    <div id="buttons">
      <input id="refresh" type="button" value="Refresh" />
      <input id="create" type="button" value="New" />
      <input id="persist" type="button" value="Save" class="hidden" />
      <input id="reset" type="button" value="Cancel" class="hidden" />
    </div>

    <script type="text/javascript" src="/static/app.js"></script>
    <script type="text/javascript" src="{{jq}}/{{jqv}}/jquery.min.js"></script>
    <script type="text/javascript" src="http://www.google.com/jsapi"></script>
    <script type="text/javascript">
      // Load the gviz table
      var pkgs = ['corechart', 'table'];
      google.load('visualization', '1.0', {'packages':pkgs});
      google.setOnLoadCallback(notes.fetch_by_tags);

      // Add event handlers
      $('#create').on('click', notes.create);
      $('#persist').on('click', notes.persist);
      $('#refresh').on('click', notes.fetch_by_tags);
      $('#reset').on('click', notes.reset);

      // Key bindings
      $(document).keydown(function(e){
        switch (e.keyCode) {
          case 27:
            notes.reset();
            break;
          case 78:
            notes.create();
            break;
        }
      });
    </script>
  </body>
</html>
