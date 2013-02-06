/**
 * @fileoverview Describes a view of notes
 */
define([
  'backbone',
  'views/notesTableRow',
  'text!templates/notesTable.html'
],
function(Backbone, notesTableRowView, notesTableTemplate) {
  return Backbone.View.extend({

    initialize: function(options) {
      this.collection.on('reset', this.render, this);
      this.collection.fetch();
      $('body').append(this.options.passwordModal.render().el);
    },

    addRow: function(note) {
      var row = new notesTableRowView({
        model: note,
        tabs: $('.nav-tabs'),
        tabContent: $('.tab-content'),
        passwordModal: this.options.passwordModal
      });
      this.$('tbody').append(row.render().el);
    },

    render: function(collection) {
      this.$el.append(_.template(notesTableTemplate)());
      this.collection.each(this.addRow, this);
    }

  });
});

