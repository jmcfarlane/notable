/**
 * @fileoverview Describes a view of notes
 */
define([
  'backbone',
  'text!templates/notes.html'
],
function(Backbone, notesTemplate) {
  return Backbone.View.extend({

    initialize: function(options) {
      this.collection.on('reset', _.bind(this.render, this));
      this.collection.fetch();
    },

    render: function(collection) {
      this.$el.html(_.template(notesTemplate)({
        notes: this.collection.toJSON()
      }));
    }

  });
});

