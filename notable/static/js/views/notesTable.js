/**
 * @fileoverview Describes a view of notes
 */
define([
  'backbone',
  'views/notesTableRow',
  'text!templates/saved.html',
  'text!templates/notesTable.html'
],
function(Backbone, notesTableRowView, savedTemplate, notesTableTemplate) {
  return Backbone.View.extend({

    initialize: function(options) {
      this.collection.on('reset', this.render, this);
      this.collection.fetch();
      $('body').append(this.options.passwordModal.render().el);
      $('body').append(this.options.searchModal.render().el);
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
    },

    render: function(collection) {
      $('body').prepend(_.template(savedTemplate));
      this.$el.append(_.template(notesTableTemplate)());
      this.collection.each(this.addRow, this);
    }

  });
});

