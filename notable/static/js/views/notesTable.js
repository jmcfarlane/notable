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
      this.collection.on('reset', _.bind(this.render, this));
      this.collection.fetch();
    },

    render: function(collection) {
      this.$el.html(_.template(notesTableTemplate)());
      this.collection.each(function(note) {
        var row = new notesTableRowView({
          model: note
        });
        this.$('tbody').append(row.render().el);
      });
    }

  });
});

