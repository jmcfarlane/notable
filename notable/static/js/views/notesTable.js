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
      $(document).bind('keydown', 'j', _.bind(this.selectNextNote, this));
      $(document).bind('keydown', 'k', _.bind(this.selectPreviousNote, this));
      $(document).bind('keydown', 'return', _.bind(this.openSelectedNote, this));
    },

    addRow: function(note, idx) {
      var row = new notesTableRowView({
        model: note,
        tabs: $('.nav-tabs'),
        tabContent: $('.tab-content'),
        passwordModal: this.options.passwordModal,
        searchModal: this.options.searchModal
      });
      this.$('tbody').append(row.render().el);
      note.set('selected', (idx === 0));
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

    openSelectedNote: function() {
      this.collection.find(function(model){
        return model.get('selected');
      }, this).view.onRowClick();
    },

    render: function(collection) {
      $('body').prepend(_.template(savedTemplate));
      this.$el.append(_.template(notesTableTemplate)());
      this.collection.each(this.addRow, this);
    },

    selectNextNote: function() {
      this.navigateNote(-1);
    },

    selectPreviousNote: function() {
      this.navigateNote(+1);
    },

    navigateNote: function(step) {
      this.collection.every(function(model, idx) {
        if (model.get('selected')) {
          var possible = this.collection.at(idx - step);
          if (possible) {
            this.collection.at(idx).set('selected', false);
            possible.set('selected', true);
            return false;
          }
        }
        return true;
      }, this)
    }


  });
});

