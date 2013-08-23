// Configure require itself
require.config({
  baseUrl: '/static/js',
  shim: {
    underscore: {
      exports: '_'
    },
    backbone: {
      deps: ['underscore', 'jquery'],
      exports: 'Backbone'
    },
    bootstrap: {
      deps: ['jquery']
    }
  },
	paths: {
    ace: '../lib/ace/src-min/ace',
    backbone: '../lib/backbone',
    bootstrap: '../lib/bootstrap/js/bootstrap.min',
    jquery: '../lib/jquery',
    lib: '../lib',
    templates: '../templates',
    text: '../lib/plugins/require/text',
    underscore: '../lib/underscore'
	}
});

// The notable client side application
require(
	[
   'collections/notes',
   'views/table',
   'views/password',
   'views/search',
   'jquery',
   'underscore',
   'backbone',
   'bootstrap',
   'lib/plugins/bootstrap/tab'
  ],
	function(NotesCollection, NotesTableView, PasswordModalView, SearchModalView) {
    var notesView = new NotesTableView({
      collection: new NotesCollection(),
      el: $('#notes'),
      passwordModal: new PasswordModalView(),
      searchModal: new SearchModalView()
    });
	}
);
