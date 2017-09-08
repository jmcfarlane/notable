/**
 * @fileoverview Describes a table row view for a note
 */
define([
  'views/note',
  'text!templates/row.html',
  'backbone',
  'underscore'
],
function(noteView, notesTableRowTemplate) {
  return Backbone.View.extend({
    events: {
      'click': 'onRowClick'
    },

    className: 'note-row',
    tagName: 'tr',

    initialize: function(options) {
      var modal = this.options.passwordModal;
      this._index = {};
      this._tab = null;
      this.index();
      this.model.on('change:uid', this.displayContent, this);
      this.model.on('change:selected', this.toggleSelected, this);
      this.model.on('content:fetched', this.displayContent, this);
      this.model.on('change', this.index, this);
      this.model.on('destroy', this.onDestroy, this);
      this.model.on('decryption:error', modal.renderError, modal);
      this.options.searchModal.on('search', this.search, this);
      this.model.set('selected', false);
      this.model.view = this;
    },

    index: function() {
      this._index.tags = (this.model.get('tags') || '').split(' ');
    },

    render: function(collection) {
      this.$el.html(_.template(notesTableRowTemplate, {
        note: this.model.toJSON()
      }));
      if (this.model.get('ahead_of_primary') || (this.model.get('deleted'))) {
        this.$el.css('color', 'red');
      }
      if (this.model.get('deleted')) {
        this.$el.css('text-decoration', 'line-through');
      }
      return this;
    },

    /**
     * When a row is clicked call a method on the model to fetch the
     * content itself (prompting for password if necessary).
     * Displaying the note will happen once the content is
     * successfully fetched.
     */
    onRowClick: function() {
      if (this._tab) {
        this._tab.show();
        return;
      }

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
    displayContent: function() {
      this._tab = new noteView({
        _row: this,
        model: this.model,
        tabContent: this.options.tabContent,
        tabs: this.options.tabs
      }).render();
      this.options.passwordModal.hide();
      this._tab.on('destroy', this.onTabClose, this);
      return this._tab;
    },

    toggleSelected: function() {
      var td = this.$('td:first');
      td[this.model.get('selected') ? 'addClass' : 'removeClass']('selected');
    },

    onDestroy: function() {
      this.$el.detach().remove();
      delete this._tab;
    },

    onTabClose: function() {
      delete this._tab;
      this.model.set('content', null, {
        silent: true
      });
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
        return this.startsWith(tag, q.query);
      }, this);

      if (this.startsWith(subject, q.query) || match || _.isEmpty(q.query)) {
        this.show();
      } else {
        // Secondary search is full text index
        if (_.contains(q.uids, this.model.get('uid'))) {
          this.show();
        } else {
          // No soup for you
          this.hide();
        }
      }
    },

    startsWith: function(haystack, needle) {
      return haystack.indexOf(needle) == 0;
    }

  });
});
