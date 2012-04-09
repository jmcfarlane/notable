/**
  * Encapsulate and namespace
  */
(function() {
  if ( typeof(NOTABLE) === "undefined") {
    NOTABLE = {};
  };

  // Constants
  var ENCRYPTED = '<ENCRYPTED>';
  var RE_ENCRYPTED = new RegExp(/^[^ ]{32,}$/);
  var RE_DECRYPTED = new RegExp(/^[\000-\177]*$/);

  /**
   * Main application
   */
  NOTABLE.Application = function(el) {
    var that = this;

    /**
     * Application setup
     */
    this.init = function() {

      var pkgs = ['corechart', 'table'];
      google.load('visualization', '1.0', {'packages':pkgs});
      google.setOnLoadCallback(that.perform_search);

      // Add event handlers
      $('#create').on('click', that.create);
      $('#persist').on('click', that.persist);
      $('#refresh').on('click', that.perform_search);
      $('#reset').on('click', that.reset_all);
      $('#search input').on('keypress', that.perform_search);
      $('#password-dialog form').on('submit', that.decrypt);
      $('#editor').on('click', that.launch_editor);
      $('#delete').on('click', that.delete);

      // Key bindings
      $(document).keydown(function(e) {
        switch (e.which) {
          case 27:
            that.reset_all();
            break;
          case 78:
            that.create();
            break;
          case 83:
            that.search_dialog();
            break;
        }
      });

      return this;
    };

    /**
     * Center an element
     */
    this.center = function(el) {
      el = $(el);
      el.css({
        position: 'absolute',
        left: '50%',
        'margin-left': 0 - (el.width() / 2)
      });
      return el;
    };

    /**
     * Post process the note content as it might be encrypted
     */
    this.content = function(data, i) {
      var idx = data.getColumnIndex('content');
      var content = data.getValue(i, idx).split('\n');
      content.shift();
      if (RE_ENCRYPTED.test(content)) {
        return ENCRYPTED;
      }
      return content.join('\n').substr(0, 100);
    };

    /**
     * Create a new note
     */
    this.create = function() {
      if (NOTABLE.any_visible(['#content textarea'])) {
        return;
      }
      $('#create').hide();
      $('#listing').hide();
      $('#refresh').hide();
      $('#content').show();
      $('#delete').show();
      $('#editor').show();
      $('#persist').show();
      $('#reset').show();

      setTimeout(foo = function() {
        if (! NOTABLE.any_visible(['#search input', '#password-dialog'])) {
          $('#content textarea').focus()
        }}, 100);
      this.reset_search();
    }

    /**
     * Decrypt a particular note's contents
     */
    this.decrypt = function() {
      var post = {
        password: $('#password-dialog input').val(),
        uid: $('#content #uid').val(),
      }
      $.post('/api/decrypt', post, function (response) {
        if (RE_DECRYPTED.test(response)) {
          $('#content textarea').val(response);
          that.reset_password();
          $('#delete').show();
          $('#password').val(post.password);
        } else {
          $('#password-dialog input').val('');
        };
      });
      return false;
    }

    /**
     * Delete a note
     */
    this.delete = function() {
      if (confirm('Are you sure you want to delete this note?')) {
        var post = {
          password: $('#content #password').val(),
          uid: $('#content #uid').val(),
        }
        $.post('/api/delete', post, function (response) {
          that.reset_all();
          that.perform_search();
        });
      }
    }

    /**
     * Edit an existing note
     */
    this.edit = function(data, row) {
      that.create();
      $('#content textarea').val(data.getValue(row, data.getColumnIndex('content')))
      $('#content #uid').val(data.getValue(row, data.getColumnIndex('uid')))
      $('#content #tags').val(data.getValue(row, data.getColumnIndex('tags')))

      var b = data.getValue(row, data.getColumnIndex('content')).split('\n')[1];
      if (RE_ENCRYPTED.test(b)) {
        that.password_dialog();
        $('#delete').hide();
      }
    }

    /**
     * Launch external editor
     */
    this.launch_editor = function() {
      var post = {
        content: $('#content textarea').val(),
        uid: $('#content #uid').val(),
      }
      $.post('/api/launch_editor', post, function (response) {
        that.poll_disk(response);
      });
    }

    /**
     * Spawn password entry dialog for an encrypted note
     */
    this.password_dialog = function() {
      this.center('#password-dialog').show();
      setTimeout("$('#password-dialog input').focus()", 100);
    }

    /**
     * Make network call to search for notes
     */
    this.perform_search = function() {
      var q = {s: $('#search input').val()};
      $.get('/api/list', q, function (response) {
        that.render_listing(response);
      });
    }

    /**
     * Persist a note to the backend
     */
    this.persist = function() {
      var post = {
        content: $('#content textarea').val(),
        password: $('#content #password').val(),
        tags: $('#content #tags').val(),
        uid: $('#content #uid').val(),
      }
      $.post('/api/persist', post, function (response) {
        that.reset_all();
        that.perform_search();
      });
    }

    /**
     * Poll disk for updates by an external text editor
     */
    this.poll_disk = function(uid) {
      $.get('/api/from_disk/' + uid, function (response) {
        if (response != 'missing') {
          $('#content textarea').val(response);
          setTimeout(function(){that.poll_disk(uid)}, 200);
        } else {
          that.persist();
        }
      });
    }

    /**
     * Render note listing
     */
    this.render_listing = function(json) {
      var div = document.getElementById('listing')
      var data = new google.visualization.DataTable(json);
      var table = new google.visualization.Table(div);
      var view = new google.visualization.DataView(data);
      var columns = [
        data.getColumnIndex('updated'),
        data.getColumnIndex('subject'),
        {calc: that.content, id:'body', type:'string', label:'Body'},
        data.getColumnIndex('tags'),
      ]
      var options = {
        height: '200px',
        sortAscending: false,
        sortColumn: 0,
        width: '100%',
      };

      view.setColumns(columns);
      $(div).height($(window).height() - $('#buttons').height() - 30);
      table.draw(view, options);

      // Add navigation listener
      google.visualization.events.addListener(table, 'select',
        function(e) {
          that.edit(data, table.getSelection()[0].row);
      });

      // Redraw if the window size changes
      $(window).resize(this.perform_search);
    }

    /**
     * Reset all dialogs and forms
     */
    this.reset_all = function() {
      $('#content').hide();
      $('#reset').hide();
      $('#persist').hide();
      $('#delete').hide();
      $('#editor').hide();
      $('#listing').show();
      $('#create').show();
      $('#refresh').show();
      $('#content #tags').val('');
      $('#content textarea').val('');
      $('#content #uid').val('');
      $('#content #password').val('');
      that.reset_password();
      that.reset_search();
    }

    /**
     * Reset password dialog
     */
    this.reset_password = function() {
      $('#password-dialog').hide();
      $('#password-dialog input').val('');
    }

    /**
     * Reset search dialog
     */
    this.reset_search = function() {
      $('#search').hide();
      $('#search input').val('');
      this.perform_search();
    }

    /**
     * Render search dialog
     */
    this.search_dialog = function() {
      var found = NOTABLE.any_visible(['#content textarea']);
      if (! found) {
        $('#search').show();
        setTimeout("$('#search input').focus()", 100);
      }
    }

  }; // eoc

  /**
   * Return boolean if any of the specified selectors are visible
   */
  NOTABLE.any_visible = function(selectors) {
    var guilty = false
    $.each(selectors, function(idx, value) {
      if ($(value).is(":visible")) {
        guilty = true;
        return false;
      };
    });
    return guilty;
  }

}());

