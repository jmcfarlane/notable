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
        <input type="password" id="password" name="password" />
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
      <input id="delete" type="button" value="Delete" class="hidden" />
    </div>

    <div id="search" class="hidden">
      <input type="text" id="q" />
    </div>

    <div id="password-dialog" class="hidden">
      <form>
        <input type="password" />
      </form>
    </div>

    <script type="text/javascript" src="{{jq}}/{{jqv}}/jquery.min.js"></script>
    <script type="text/javascript" src="http://www.google.com/jsapi"></script>
    <script type="text/javascript" src="/static/app.js"></script>
    <script type="text/javascript">
      // Load the application
      var notable = (new NOTABLE.Application('abc')).init();
    </script>
  </body>
</html>
