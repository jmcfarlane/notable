/**
 * @fileoverview Describes a view of notes
 */
define([
  'backbone',
  'models/note',
  'views/notesTableRow',
  'text!templates/saved.html',
  'text!templates/notesTable.html',
  'lib/jquery.hotkeys'
],
function(Backbone, NoteModel, notesTableRowView, savedTemplate, notesTableTemplate) {
  return Backbone.View.extend({

    initialize: function(options) {
      this.collection.on('reset', this.render, this);
      this.collection.fetch();
      $('body').append(this.options.passwordModal.render().el);
      $('body').append(this.options.searchModal.render().el);
      $('.create').on('click', _.bind(this.createNote, this));
      $(document).bind('keydown', 'ctrl+c', _.bind(this.createNote, this));
    },

    addRow: function(note) {
      var row = new notesTableRowView({
        model: note,
        tabs: $('.nav-tabs'),
        tabContent: $('.tab-content'),
        passwordModal: this.options.passwordModal,
        searchModal: this.options.searchModal
      });
      this.$('tbody').append(row.render().el);
      return row;
    },

    createNote: function() {
      var note = new NoteModel({
        content: '',
        subject: '',
        tags: ''
      });
      this.addRow(note);
      note.save();
      return false;
    },

    render: function(collection) {
      $('body').prepend(_.template(savedTemplate));
      this.$el.append(_.template(notesTableTemplate)());
      this.collection.each(this.addRow, this);
    }

  });
});

