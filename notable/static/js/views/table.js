/**
 * @fileoverview Describes a view of notes
 */
define([
  'backbone',
  'models/note',
  'views/row',
  'text!templates/not-saved.html',
  'text!templates/saved.html',
  'text!templates/table.html',
  'lib/jquery.hotkeys'
],
function(Backbone,
         NoteModel,
         notesTableRowView,
         notSavedTemplate,
         savedTemplate,
         notesTableTemplate) {
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
      options.searchModal.on('next', this.selectNextNote, this);
      options.searchModal.on('previous', this.selectPreviousNote, this);
      options.searchModal.on('open', this.openSelectedNote, this);
      options.searchModal.on('search', _.bind(function() {
        setTimeout(_.bind(this.defaultSelected, this), 100);
      }, this));
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

    defaultSelected: function() {
      this.collection.each(function(model) {
        model.set('selected', false);
      });
      var visible = this.visibleRows().at(0);
      if (visible) {
       visible.set('selected', true);
      }
    },

    openSelectedNote: function() {
      // I don't like this condition.  The enter binding should not be
      // firing when a modal input (or any input) has focus, hmm.
      if (this.options.passwordModal.$('input').is(':visible')) {
        return;
      }
      this.collection.find(function(model){
        return model.get('selected');
      }, this).view.onRowClick();
    },

    render: function(collection) {
      $('body').prepend(_.template(notSavedTemplate));
      $('body').prepend(_.template(savedTemplate));
      this.$el.append(_.template(notesTableTemplate)());
      this.collection.each(this.addRow, this);
      this.defaultSelected();
    },

    selectNextNote: function() {
      this.navigateNote(-1);
    },

    selectPreviousNote: function() {
      this.navigateNote(+1);
    },

    navigateNote: function(step) {
      var visible = this.visibleRows();
      visible.every(function(model, idx) {
        if (model.get('selected')) {
          var possible = visible.at(idx - step);
          if (possible) {
            visible.at(idx).set('selected', false);
            possible.set('selected', true);
            return false;
          }
        }
        return true;
      }, this)
    },

    visibleRows: function() {
      return new Backbone.Collection(this.collection.filter(function(model) {
        return model.view.$el.is(':visible');
      }));
    }


  });
});

