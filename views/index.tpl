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
        <input type="text" id="tags" name="tags" />
        <label for="tags">Tags</label>
        <br/>
        <input type="text" id="password" name="password" />
        <label for="tags">Password</label>
        <input type="hidden" id="uid" />
      </form>
    </div>

    <div id="buttons">
      <input id="refresh" type="button" value="Refresh" />
      <input id="create" type="button" value="New" />
      <input id="persist" type="button" value="Save" class="hidden" />
      <input id="reset" type="button" value="Cancel" class="hidden" />
      <input id="editor" type="button" value="Editor" class="hidden" />
    </div>

    <div id="search" class="hidden">
      <input type="text" id="q" />
    </div>

    <div id="password-dialog" class="hidden">
      <form>
        <input type="password" />
      </form>
    </div>

    <script type="text/javascript" src="/static/app.js"></script>
    <script type="text/javascript" src="{{jq}}/{{jqv}}/jquery.min.js"></script>
    <script type="text/javascript" src="http://www.google.com/jsapi"></script>
    <script type="text/javascript">
      // Load the gviz table
      var pkgs = ['corechart', 'table'];
      google.load('visualization', '1.0', {'packages':pkgs});
      google.setOnLoadCallback(notes.search);

      // Add event handlers
      $('#create').on('click', notes.create);
      $('#persist').on('click', notes.persist);
      $('#refresh').on('click', notes.search);
      $('#reset').on('click', notes.reset);
      $('#search input').on('keypress', notes.search);
      $('#password-dialog form').on('submit', notes.pwd_submit);
      $('#editor').on('click', notes.launch_editor);

      // Key bindings
      $(document).keydown(function(e){
        switch (e.which) {
          case 27:
            notes.reset();
            notes.search_reset();
            notes.password_reset();
            break;
          case 78:
            notes.create();
            break;
          case 83:
            notes.search_perform();
            break;
        }
      });

      // Dynamic placements
      var pwd = $('#password-dialog');
      pwd.css({
        position: 'absolute',
        left: '50%',
        'margin-left': 0 - (pwd.width() / 2)
      });
    </script>
  </body>
</html>
