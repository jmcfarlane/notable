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
    }
  },
	paths: {
    backbone: '../lib/backbone',
    jquery: '../lib/jquery',
    templates: '../templates',
    text: '../lib/text',
    underscore: '../lib/underscore'
	}
});

// The notable client side application
require(
	[
   'collections/notes',
   'views/notes',
   'underscore',
   'backbone'
  ],
	function(NotesCollection, NotesView) {
    var notesView = new NotesView({
      collection: new NotesCollection(),
      el: $('#notes')
    });
	}
);
