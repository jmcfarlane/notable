/**
 * @fileoverview Describes a table row view for a note
 */
define([
  'views/noteTab',
  'text!templates/notesTableRow.html',
  'backbone',
  'underscore'
],
function(noteTab, notesTableRowTemplate) {
  return Backbone.View.extend({
    events: {
      'click': 'getContent'
    },

    tagName: 'tr',

    render: function(collection) {
      this.$el.html(_.template(notesTableRowTemplate, {
        note: this.model.toJSON()
      }));

      return this;
    },

    /**
     * Fetch content for a specific note, passing in the password in
     * the event it's an encrypted note.
     */
    getContent: function() {
      this.model.getContent('', _.bind(this.displayContent, this));
    },

    /**
     * Display the details of a particular note
     */
    displayContent: function(content) {
      var tab = new noteTab({
        tabs: this.options.tabs,
        tabContent: this.options.tabContent,
        model: this.model
      }).render();
    }

  });
});
