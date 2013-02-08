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
      'click': 'onRowClick'
    },

    tagName: 'tr',

    initialize: function(options) {
      var modal = this.options.passwordModal;
      this._index = {};
      this.index();
      this.model.on('change', this.displayContent, this);
      this.model.on('change', this.index, this);
      this.model.on('decryption:error', modal.renderError, modal);
      this.options.searchModal.on('search', this.search, this);
    },

    index: function() {
      this._index.tags = this.model.get('tags').split(' ');
    },

    render: function(collection) {
      this.$el.html(_.template(notesTableRowTemplate, {
        note: this.model.toJSON()
      }));
      return this;
    },

    /**
     * When a row is clicked call a method on the model to fetch the
     * content itself (prompting for password if necessary).
     * Displaying the note will happen once the content is
     * successfully fetched.
     */
    onRowClick: function() {
      var modal = this.options.passwordModal,
        password = modal.getPassword().val();
      if (this.model.get('encrypted') && !password) {
        return modal.show(_.bind(this.onRowClick, this));
      }
      this.model.fetchContent(password);
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
      this.options.passwordModal.hide();
    },

    hide: function() {
      this.$el.hide();
    },

    show: function() {
      this.$el.show();
    },

    search: function(q) {
      var subject = this.model.get('subject').toLowerCase(),
        match = _.find(this._index.tags, function(tag) {
        return this.startsWith(tag, q);
      }, this);

      if (this.startsWith(subject, q) || match || _.isEmpty(q)) {
        this.show();
      } else {
        this.hide();
      }
    },

    startsWith: function(haystack, needle) {
      return haystack.indexOf(needle) == 0;
    }

  });
});
